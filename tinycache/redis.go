package tinycache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool   *redis.Pool
	Prefix string
}

func NewRedis(pool *redis.Pool) *Redis {
	return &Redis{
		Pool: pool,
	}
}

func (r *Redis) Get(key string) ([]byte, error) {
	c := r.Pool.Get()
	defer c.Close()
	res, err := redis.Bytes(c.Do("GET", r.Prefix+key))
	if err == redis.ErrNil {
		return nil, NotFound
	}
	return res, err
}

func (r *Redis) Set(key string, value []byte, ttl time.Duration) error {
	c := r.Pool.Get()
	defer c.Close()
	_, err := c.Do("PSETEX", r.Prefix+key, int64(ttl/time.Millisecond), value)
	return err
}
