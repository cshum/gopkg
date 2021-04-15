package tinycache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Local struct {
	Cache *cache.Cache
	TTL   time.Duration
}

func NewLocal(ttl time.Duration) *Local {
	return &Local{
		Cache: cache.New(ttl, ttl),
		TTL:   ttl,
	}
}

func (c *Local) Get(key string) ([]byte, error) {
	if res, ok := c.Cache.Get(key); ok {
		return res.([]byte), nil
	}
	return nil, NotFound
}

func (c *Local) Set(key string, value []byte) error {
	c.Cache.Set(key, value, c.TTL)
	return nil
}
