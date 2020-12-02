package client

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var cacheMem *cache.Cache
var onceCache sync.Once

func initiateCache() {
	onceCache.Do(func() {
		cacheMem = cache.New(5*time.Minute, 10*time.Minute)
	})
}

func LoadCache() *cache.Cache {
	if cacheMem == nil {
		initiateCache()
	}

	return cacheMem
}
