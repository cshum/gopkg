package tinycache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
)

type Hybrid struct {
	Pool        *redis.Pool
	Prefix      string
	Local       *cache.Cache
	MaxLocalTTL time.Duration
}

func NewHybrid(redis *redis.Pool, maxLocalTTL, cleanupInterval time.Duration) *Hybrid {
	return &Hybrid{
		Pool:        redis,
		Local:       cache.New(maxLocalTTL, cleanupInterval),
		MaxLocalTTL: maxLocalTTL,
	}
}

func (c *Hybrid) Get(key string) (value []byte, err error) {
	if res, ok := c.Local.Get(key); ok {
		value = res.([]byte)
		return
	}
	if c.Pool != nil {
		conn := c.Pool.Get()
		defer conn.Close()
		if err = conn.Send("GET", c.Prefix+key); err != nil {
			return
		}
		if err = conn.Send("PTTL", c.Prefix+key); err != nil {
			return
		}
		if err = conn.Flush(); err != nil {
			return
		}
		if value, err = redis.Bytes(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = NotFound
			}
			return
		}
		var pttl int64
		if pttl, err = redis.Int64(conn.Receive()); err != nil {
			return
		}
		ttl := time.Duration(pttl) * time.Millisecond
		if ttl > c.MaxLocalTTL {
			c.Local.Set(key, value, c.MaxLocalTTL)
		} else {
			c.Local.Set(key, value, ttl)
		}
		return
	}
	err = NotFound
	return
}

func (c *Hybrid) Set(key string, value []byte, ttl time.Duration) error {
	if ttl > c.MaxLocalTTL {
		c.Local.Set(key, value, c.MaxLocalTTL)
	} else {
		c.Local.Set(key, value, ttl)
	}
	if c.Pool != nil {
		conn := c.Pool.Get()
		defer conn.Close()
		if _, err := conn.Do(
			"PSETEX", c.Prefix+key, int64(ttl/time.Millisecond), value); err != nil {
			return err
		}
	}
	return nil
}
