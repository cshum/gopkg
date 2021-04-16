package simplecache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Hybrid struct {
	Pool   *redis.Pool
	TTL    time.Duration
	Prefix string
	Local  *Local
}

func NewHybrid(redis *redis.Pool, ttl, localTTL time.Duration) *Hybrid {
	return &Hybrid{
		Pool:  redis,
		TTL:   ttl,
		Local: NewLocal(localTTL),
	}
}

func (c *Hybrid) Get(key string) (value []byte, err error) {
	if val, err2 := c.Local.Get(key); err2 == nil || err2 != NotFound {
		value = val
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
		if ttl > c.Local.TTL {
			// if redis still more ttl than local, re-cache at local
			if err = c.Local.Set(key, value); err != nil {
				return
			}
		}
		return
	}
	err = NotFound
	return
}

func (c *Hybrid) Set(key string, value []byte) error {
	if err := c.Local.Set(key, value); err != nil {
		return err
	}
	if c.Pool != nil {
		conn := c.Pool.Get()
		defer conn.Close()
		if _, err := conn.Do(
			"PSETEX", c.Prefix+key, int64(c.TTL/time.Millisecond), value); err != nil {
			return err
		}
	}
	return nil
}
