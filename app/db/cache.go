package db

import (
	_ "embed"
	"sync"

	"github.com/coocood/freecache"
)

var (
	cacheOnce sync.Once
	cache     *freecache.Cache
)

func SetupCache() {
	syncCache()
}

func Cache() *freecache.Cache {
	return cache
}

func syncCache() *freecache.Cache {
	if cache != nil {
		return cache
	}
	cacheOnce.Do(func() {
		cacheSize := 32 * 1024 * 1024 // 32 MB
		cache := freecache.NewCache(cacheSize)
		configureCache(cache)
	})

	return cache
}

func configureCache(fc *freecache.Cache) {

}
