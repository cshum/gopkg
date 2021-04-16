package simplecache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type Local struct {
	Cache *ristretto.Cache
}

func NewLocal(maxItems, maxSize int64) *Local {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxItems * 10,
		MaxCost:     maxSize,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return &Local{c}
}

func (c *Local) Get(key string) ([]byte, error) {
	if res, ok := c.Cache.Get(key); ok {
		return res.([]byte), nil
	}
	return nil, NotFound
}

func (c *Local) Set(key string, value []byte, ttl time.Duration) error {
	c.Cache.SetWithTTL(key, value, int64(len(value)), ttl)
	return nil
}
