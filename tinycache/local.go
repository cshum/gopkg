package tinycache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Local struct {
	Cache *cache.Cache
}

func NewLocal(defaultTTL, cleanupInterval time.Duration) *Local {
	return &Local{
		Cache: cache.New(defaultTTL, cleanupInterval),
	}
}

func (c *Local) Get(key string) ([]byte, error) {
	if res, ok := c.Cache.Get(key); ok {
		return res.([]byte), nil
	}
	return nil, NotFound
}

func (c *Local) Set(key string, value []byte, ttl time.Duration) error {
	c.Cache.Set(key, value, ttl)
	return nil
}
