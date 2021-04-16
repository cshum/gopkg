package tinycache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
)

type Hybrid struct {
	Pool     *redis.Pool
	Prefix   string
	Local    *cache.Cache
	RedisTTL time.Duration
	LocalTTL time.Duration
}

func NewHybrid(redis *redis.Pool, redisTTL, localTTL time.Duration) *Hybrid {
	return &Hybrid{
		Pool:     redis,
		Local:    cache.New(localTTL, localTTL),
		RedisTTL: redisTTL,
		LocalTTL: localTTL,
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
		if ttl > c.LocalTTL {
			c.Local.Set(key, value, c.LocalTTL)
		} else {
			c.Local.Set(key, value, ttl)
		}
		return
	}
	err = NotFound
	return
}

func (c *Hybrid) Set(key string, value []byte) error {
	c.Local.Set(key, value, c.LocalTTL)
	if c.Pool != nil {
		conn := c.Pool.Get()
		defer conn.Close()
		if _, err := conn.Do(
			"PSETEX", c.Prefix+key, int64(c.RedisTTL/time.Millisecond), value); err != nil {
			return err
		}
	}
	return nil
}
