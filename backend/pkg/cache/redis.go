package cache

import (
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/go-redis/redis/v8"
)

func newRedis(address, password string, db int) *marshaler.Marshaler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	redisStore := redis_store.NewRedis(redisClient)
	cacheManager := cache.New[string](redisStore)
	marshaller := marshaler.New(cacheManager.GetCodec().GetStore())
	return marshaller
}
