package config

import "github.com/bluele/gcache"

// Guard guard
type Guard struct {
	Cache gcache.Cache
}

// InitGCache init gcache
func InitGCache(config *GuardConfig) Guard {
	// gcache
	return Guard{
		Cache: gcache.New(config.Cache.Size).LRU().Build(),
	}
}
