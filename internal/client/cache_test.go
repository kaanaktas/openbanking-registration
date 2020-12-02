package client

import (
	"github.com/patrickmn/go-cache"
	"testing"
)

func TestLoadCache(t *testing.T) {
	tests := []struct {
		name           string
		usedCacheId    string
		queriedCacheId string
		want           bool
	}{
		{"retrieve_cache_success", "cacheId", "cacheId", true},
		{"retrieve_cache_fail", "cacheId", "wrongCacheId", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadCache(); got != nil {
				got.Set(tt.usedCacheId, "value", cache.DefaultExpiration)
				if _, found := got.Get(tt.queriedCacheId); found != tt.want {
					t.Errorf("LoadCache() = %v, want %v", found, tt.want)
				}
			} else {
				t.Errorf("Couldn't LoadCache()")
			}
		})
	}
}
