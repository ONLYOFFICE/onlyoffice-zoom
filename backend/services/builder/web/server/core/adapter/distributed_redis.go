package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	bigcache "github.com/allegro/bigcache/v3"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

// TODO: Bugfix (WIP)
type distributedRedisSessionAdapter struct {
	cache               *bigcache.BigCache
	redisClient         goredislib.UniversalClient
	redisLock           *redsync.Mutex
	redisPersistChannel string
	redisRemoveChannel  string
	persistBuffer       chan string
	removeBuffer        chan string
	initDone            chan bool
	bufferSize          int64
	logger              log.Logger
}

func NewDistributedRedisSessionAdapter(opts ...Option) (port.SessionServiceAdapter, error) {
	options := NewOptions(opts...)

	if len(options.RedisAddresses) < 1 {
		options.Logger.Fatal("could not create a new redis session adapter. Invalid address")
	}

	var client goredislib.UniversalClient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(options.RedisAddresses) == 1 {
		roptions := &goredislib.Options{
			Addr:        options.RedisAddresses[0],
			Username:    options.RedisUsername,
			Password:    options.RedisPassword,
			DB:          0,
			MaxRetries:  3,
			ReadTimeout: 2 * time.Second,
		}

		client = goredislib.NewClient(roptions)

		if err := client.Ping(ctx).Err(); err != nil {
			return nil, err
		}
	} else {
		roptions := &goredislib.ClusterOptions{
			Addrs:       options.RedisAddresses,
			ReadOnly:    true,
			Username:    options.RedisUsername,
			Password:    options.RedisPassword,
			MaxRetries:  3,
			ReadTimeout: 2 * time.Second,
		}

		rdb := goredislib.NewClusterClient(roptions)

		if err := rdb.ForEachShard(ctx, func(ctx context.Context, client *goredislib.Client) error {
			return client.Ping(ctx).Err()
		}); err != nil {
			return nil, err
		}

		client = rdb
	}

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	mutex := rs.NewMutex("session-mutex")

	bcache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	if err != nil {
		return nil, err
	}

	adapter := distributedRedisSessionAdapter{
		cache:               bcache,
		redisClient:         client,
		redisLock:           mutex,
		redisPersistChannel: "persist-session-channel",
		redisRemoveChannel:  "remove-session-channel",
		persistBuffer:       make(chan string, options.BufferSize),
		removeBuffer:        make(chan string, options.BufferSize),
		initDone:            make(chan bool),
		bufferSize:          int64(options.BufferSize),
		logger:              options.Logger,
	}

	go adapter.loadData()
	go adapter.subscribePersist()
	go adapter.subscribeRemove()

	return adapter, nil
}

func (s distributedRedisSessionAdapter) loadData() {
	var cursor uint64

	for {
		var keys []string
		var err error

		keys, cursor, err = s.redisClient.Scan(context.Background(), cursor, "*", s.bufferSize).Result()
		if err != nil {
			s.logger.Errorf("failed to retrieve data during cold start. Reason: %s\n", err.Error())
		}

		for _, key := range keys {
			val, _ := s.redisClient.Get(context.Background(), key).Result()
			if err := s.cache.Set(key, []byte(val)); err != nil {
				s.logger.Errorf("could not set a pair with key=%s during cold start. Reason: %s\n", key, err.Error())
			}
		}

		if cursor == 0 {
			break
		}
	}

	s.initDone <- true
	close(s.initDone)
}

func (s distributedRedisSessionAdapter) subscribePersist() {
	pubsub := s.redisClient.Subscribe(context.Background(), s.redisPersistChannel)

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				s.logger.Errorf("could not receive a persist message. Reason: %s", err.Error())
			} else {
				s.persistBuffer <- msg.Payload
				s.logger.Debugf("received a persist message: %v", msg.Payload)
			}
		}
	}()

	<-s.initDone

	for {
		msg := <-s.persistBuffer
		payload := strings.Split(msg, ";")
		key, value := payload[0], payload[1]
		s.cache.Set(key, []byte(value))
		s.logger.Debugf("persisting %s", key)
	}
}

func (s distributedRedisSessionAdapter) subscribeRemove() {
	pubsub := s.redisClient.Subscribe(context.Background(), s.redisRemoveChannel)

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				s.logger.Errorf("could not receive a remove message. Reason: %s", err.Error())
			} else {
				s.removeBuffer <- msg.Payload
				s.logger.Debugf("received a remove message: %v", msg.Payload)
			}
		}
	}()

	<-s.initDone

	for {
		msg := <-s.removeBuffer
		s.cache.Delete(msg)
		s.logger.Debugf("removing %s", msg)
	}
}

func (s distributedRedisSessionAdapter) broadcastAndPersist(key, value string) error {
	publish_message := key + ";" + value
	if err := s.redisLock.Lock(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer s.redisLock.Unlock()
	defer cancel()

	if err := s.redisClient.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
		return err
	}

	s.redisClient.Publish(ctx, s.redisPersistChannel, publish_message)
	return nil
}

func (s distributedRedisSessionAdapter) broadcastAndRemove(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Second)
	defer cancel()
	if err := s.redisClient.Publish(ctx, s.redisRemoveChannel, key).Err(); err != nil {
		return err
	}

	if err := s.redisLock.Lock(); err != nil {
		return err
	}

	defer s.redisLock.Unlock()
	return s.redisClient.Del(ctx, key).Err()
}

func (s distributedRedisSessionAdapter) InsertSession(ctx context.Context, mid string, session domain.Session, expiresIn time.Duration) (domain.Session, error) {
	if sess, err := s.SelectSessionByMettingID(ctx, mid); err == nil {
		return sess, ErrSessionAlreadyExists
	} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return session, err
	}

	buf, err := json.Marshal(session)
	if err != nil {
		return session, err
	}

	if err := s.broadcastAndPersist(mid, string(buf)); err != nil {
		return session, err
	}

	return session, nil
}

func (s distributedRedisSessionAdapter) SelectSessionByMettingID(ctx context.Context, mid string) (domain.Session, error) {
	var session domain.Session
	buf, err := s.cache.Get(mid)

	if err == nil {
		if uerr := json.Unmarshal(buf, &session); uerr != nil {
			return session, uerr
		}
		return session, nil
	}

	if sess, err := s.redisClient.Get(ctx, mid).Result(); err != nil {
		return session, err
	} else {
		if uerr := json.Unmarshal([]byte(sess), &session); uerr != nil {
			return session, uerr
		}

		s.cache.Set(mid, []byte(sess))

		return session, nil
	}
}

func (s distributedRedisSessionAdapter) DeleteSessionByMeetingID(ctx context.Context, mid string) error {
	if err := s.broadcastAndRemove(mid); err != nil {
		return err
	}

	return nil
}
