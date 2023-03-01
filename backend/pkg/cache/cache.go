package cache

import (
	"context"
	"encoding/json"
	"time"

	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"go-micro.dev/v4/cache"
)

type CustomCache struct {
	store *gocache.Cache[string]
	name  string
}

func (c *CustomCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	res, err := c.store.Get(ctx, key)
	if err != nil {
		return nil, time.Now(), err
	}

	return []byte(res), time.Now(), nil
}

func (c *CustomCache) Put(ctx context.Context, key string, val interface{}, d time.Duration) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.store.Set(ctx, key, string(buf), store.WithExpiration(d))
}

func (c *CustomCache) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}

func (c *CustomCache) String() string {
	return c.name
}

func NewCache(opts ...Option) cache.Cache {
	options := NewOptions(opts...)

	switch options.CacheType {
	case Memory:
		return &CustomCache{
			store: newMemory(options.Size),
			name:  "freecache",
		}
	case Redis:
		return &CustomCache{
			store: newRedis(options.Address, options.Username, options.Password, options.DB),
			name:  "redis",
		}
	default:
		return &CustomCache{
			store: newMemory(options.Size),
			name:  "freecache",
		}
	}
}
