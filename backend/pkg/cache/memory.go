package cache

import (
	"github.com/coocood/freecache"
	"github.com/eko/gocache/lib/v4/cache"
	freecache_store "github.com/eko/gocache/store/freecache/v4"
)

func newMemory(size int) *cache.Cache[[]byte] {
	freecacheStore := freecache_store.NewFreecache(freecache.NewCache(size * 1024 * 1024))
	cacheManager := cache.New[[]byte](freecacheStore)
	return cacheManager
}
