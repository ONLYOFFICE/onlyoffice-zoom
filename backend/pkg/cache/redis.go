package cache

import (
	"log"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/rueidis/v4"
	"github.com/rueian/rueidis"
)

func newRedis(address, username, password string, db int) *cache.Cache[string] {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{address},
		Username:    username,
		Password:    password,
		SelectDB:    db,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	cacheManager := cache.New[string](redis_store.NewRueidis(
		client,
		store.WithExpiration(15*time.Second),
		store.WithClientSideCaching(15*time.Second)),
	)

	return cacheManager
}
