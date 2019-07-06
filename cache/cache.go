package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	Redis *redis.Client
	Local *cache.Cache
}

func New(client *redis.Client, defaultTTL, cleanupInterval time.Duration) *Cache {
	return &Cache{
		Redis: client,
		Local: cache.New(defaultTTL, cleanupInterval),
	}
}

func getbytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	// check local first
	if res, ok := c.Local.Get(key); ok {
		return getbytes(res)
	}
	// then check redis
	res, err := c.Redis.Get(key).Result()
	return []byte(res), err
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	// set redis first
	if _, err := c.Redis.Set(key, value, ttl).Result(); err != nil {
		return err
	}
	// then set local
	c.Local.Set(key, value, ttl)
	return nil
}
