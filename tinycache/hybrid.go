package tinycache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

type Hybrid struct {
	Redis       *redis.Client
	Local       *cache.Cache
	maxLocalTTL time.Duration
}

func NewHybrid(client *redis.Client, maxLocalTTL, cleanupInterval time.Duration) *Hybrid {
	return &Hybrid{
		Redis:       client,
		Local:       cache.New(maxLocalTTL, cleanupInterval),
		maxLocalTTL: maxLocalTTL,
	}
}

func (c *Hybrid) Get(key string) ([]byte, error) {
	if res, ok := c.Local.Get(key); ok {
		return res.([]byte), nil
	}
	if c.Redis != nil {
		res, err := c.Redis.Get(key).Result()
		if err == redis.Nil {
			return nil, NotFound
		}
		if err != nil {
			return nil, err
		}
		ttl, err := c.Redis.TTL(key).Result()
		if err != nil {
			return nil, err
		}
		value := []byte(res)
		if ttl > c.maxLocalTTL {
			c.Local.Set(key, value, c.maxLocalTTL)
		} else {
			c.Local.Set(key, value, ttl)
		}
		return value, err
	}
	return nil, NotFound
}

func (c *Hybrid) Set(key string, value []byte, ttl time.Duration) error {
	if ttl > c.maxLocalTTL {
		c.Local.Set(key, value, c.maxLocalTTL)
	} else {
		c.Local.Set(key, value, ttl)
	}
	if c.Redis != nil {
		if err := c.Redis.Set(key, value, ttl).Err(); err != nil {
			return err
		}
	}
	return nil
}
