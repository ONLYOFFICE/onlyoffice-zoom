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
	store *gocache.Cache[[]byte]
}

func (c *CustomCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	buf, err := c.store.Get(ctx, key)
	return buf, time.Now(), err
}

func (c *CustomCache) Put(ctx context.Context, key string, val interface{}, d time.Duration) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.store.Set(ctx, key, buf, store.WithExpiration(d))
}

func (c *CustomCache) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}

func (c *CustomCache) String() string {
	return "freecache"
}

func NewCache(opts ...Option) cache.Cache {
	options := NewOptions(opts...)

	switch options.CacheType {
	case Memory:
		return &CustomCache{
			store: newMemory(options.Size),
		}
	default:
		return &CustomCache{
			store: newMemory(options.Size),
		}
	}
}
