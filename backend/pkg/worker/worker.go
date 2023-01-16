package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func NewRedisWorker(ctx interface{}, opts ...Option) *work.WorkerPool {
	options := NewOptions(opts...)

	o := []redis.DialOption{
		redis.DialReadTimeout(options.RedisReadTimeout),
		redis.DialWriteTimeout(options.RedisWriteTimeout),
		redis.DialUsername(options.RedisUsername),
		redis.DialPassword(options.RedisPassword),
		redis.DialDatabase(options.RedisDatabase),
		redis.DialTLSSkipVerify(false),
	}

	pool := &redis.Pool{
		MaxActive: options.MaxActive,
		MaxIdle:   options.MaxIdle,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", options.RedisAddress, o...)
		},
	}

	return work.NewWorkerPool(ctx, options.MaxConcurrency, fmt.Sprintf("{%s}", options.RedisNamespace), pool)
}
