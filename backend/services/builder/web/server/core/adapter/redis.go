package adapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type redisSessionAdapter struct {
	redisClient goredislib.UniversalClient
	redisLock   *redsync.Mutex
	bufferSize  int64
	logger      log.Logger
}

func NewRedisSessionAdapter(opts ...Option) (port.SessionServiceAdapter, error) {
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
			DB:          options.RedisDatabase,
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

	adapter := redisSessionAdapter{
		redisClient: client,
		redisLock:   mutex,
		bufferSize:  int64(options.BufferSize),
		logger:      options.Logger,
	}

	return adapter, nil
}

func (s redisSessionAdapter) broadcastAndPersist(ctx context.Context, key, value string, expiresIn time.Duration) error {
	if err := s.redisLock.Lock(); err != nil {
		return err
	}
	defer s.redisLock.Unlock()

	if err := s.redisClient.Set(ctx, key, value, expiresIn).Err(); err != nil {
		return err
	}

	return nil
}

func (s redisSessionAdapter) broadcastAndRemove(ctx context.Context, key string) error {
	if err := s.redisLock.Lock(); err != nil {
		return err
	}

	defer s.redisLock.Unlock()
	return s.redisClient.Del(ctx, key).Err()
}

func (s redisSessionAdapter) InsertSession(ctx context.Context, mid string, session domain.Session, expiresAt time.Duration) (domain.Session, error) {
	buf, err := json.Marshal(session)
	if err != nil {
		return session, err
	}

	if err := s.broadcastAndPersist(ctx, mid, string(buf), expiresAt); err != nil {
		return session, err
	}

	return session, nil
}

func (s redisSessionAdapter) SelectSessionByMettingID(ctx context.Context, mid string) (domain.Session, error) {
	var session domain.Session

	if sess, err := s.redisClient.Get(ctx, mid).Result(); err != nil {
		return session, err
	} else {
		if uerr := json.Unmarshal([]byte(sess), &session); uerr != nil {
			return session, uerr
		}

		return session, nil
	}
}

func (s redisSessionAdapter) DeleteSessionByMeetingID(ctx context.Context, mid string) error {
	if err := s.broadcastAndRemove(ctx, mid); err != nil {
		return err
	}

	return nil
}
