package cache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	Redis       *redis.Client
	Local       *cache.Cache
	maxLocalTTL time.Duration
}

func New(client *redis.Client, maxLocalTTL, cleanupInterval time.Duration) *Cache {
	return &Cache{
		Redis:       client,
		Local:       cache.New(maxLocalTTL, cleanupInterval),
		maxLocalTTL: maxLocalTTL,
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	if res, ok := c.Local.Get(key); ok {
		return res.([]byte), nil
	}
	if res, err := c.Redis.Get(key).Result(); res != "" {
		value := []byte(res)
		c.Local.Set(key, value, c.maxLocalTTL)
		return value, err
	}
	return []byte{}, nil
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	if ttl > c.maxLocalTTL {
		c.Local.Set(key, value, c.maxLocalTTL)
	} else {
		c.Local.Set(key, value, ttl)
	}
	if _, err := c.Redis.Set(key, value, ttl).Result(); err != nil {
		return err
	}
	return nil
}
