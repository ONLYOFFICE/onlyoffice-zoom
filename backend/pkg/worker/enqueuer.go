package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func NewRedisEnqueuer(opts ...Option) *work.Enqueuer {
	options := NewOptions(opts...)

	o := []redis.DialOption{
		redis.DialReadTimeout(options.RedisReadTimeout),
		redis.DialWriteTimeout(options.RedisWriteTimeout),
		redis.DialUsername(options.RedisUsername),
		redis.DialPassword(options.RedisPassword),
		redis.DialTLSConfig(options.TLSConfig),
		redis.DialDatabase(options.RedisDatabase),
	}

	if options.TLSConfig != nil {
		o = append(o, redis.DialTLSConfig(options.TLSConfig))
	}

	pool := &redis.Pool{
		MaxActive: options.MaxActive,
		MaxIdle:   options.MaxIdle,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", options.RedisAddress, o...)
		},
	}

	return work.NewEnqueuer(fmt.Sprintf("{%s}", options.RedisNamespace), pool)
}
