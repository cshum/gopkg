package tinycache

import (
	"time"

	"github.com/go-redis/redis"
)

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		Client: client,
	}
}

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Get(key string) ([]byte, error) {
	res, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, NotFound
	}
	return []byte(res), err
}

func (r *Redis) Set(key string, value []byte, ttl time.Duration) error {
	return r.Client.Set(key, value, ttl).Err()
}
