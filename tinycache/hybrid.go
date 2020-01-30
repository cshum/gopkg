package tinycache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
)

type Hybrid struct {
	Redis       *redis.Pool
	RedisPrefix string
	Local       *cache.Cache
	MaxLocalTTL time.Duration
}

func NewHybrid(redis *redis.Pool, maxLocalTTL, cleanupInterval time.Duration) *Hybrid {
	return &Hybrid{
		Redis:       redis,
		Local:       cache.New(maxLocalTTL, cleanupInterval),
		MaxLocalTTL: maxLocalTTL,
	}
}

func (c *Hybrid) Get(key string) ([]byte, error) {
	if res, ok := c.Local.Get(key); ok {
		return res.([]byte), nil
	}
	if c.Redis != nil {
		conn := c.Redis.Get()
		defer conn.Close()
		value, err := redis.Bytes(conn.Do("GET", c.RedisPrefix+key))
		if err == redis.ErrNil {
			return nil, NotFound
		}
		if err != nil {
			return nil, err
		}
		pttl, err := redis.Int64(conn.Do("PTTL", c.RedisPrefix+key))
		if err != nil {
			return nil, err
		}
		ttl := time.Duration(pttl) * time.Millisecond
		if ttl > c.MaxLocalTTL {
			c.Local.Set(key, value, c.MaxLocalTTL)
		} else {
			c.Local.Set(key, value, ttl)
		}
		return value, err
	}
	return nil, NotFound
}

func (c *Hybrid) Set(key string, value []byte, ttl time.Duration) error {
	if ttl > c.MaxLocalTTL {
		c.Local.Set(key, value, c.MaxLocalTTL)
	} else {
		c.Local.Set(key, value, ttl)
	}
	if c.Redis != nil {
		conn := c.Redis.Get()
		defer conn.Close()
		if _, err := conn.Do(
			"PSETEX", c.RedisPrefix+key, int64(ttl/time.Millisecond), value); err != nil {
			return err
		}
	}
	return nil
}
