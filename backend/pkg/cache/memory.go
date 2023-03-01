package cache

import (
	"time"

	"github.com/coocood/freecache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	freecache_store "github.com/eko/gocache/store/freecache/v4"
)

func newMemory(size int) *marshaler.Marshaler {
	freecacheStore := freecache_store.NewFreecache(freecache.NewCache(size*1024*1024), store.WithExpiration(10*time.Second))
	cacheManage := cache.New[[]byte](freecacheStore)
	return marshaler.New(cacheManage.GetCodec().GetStore())
}
