package tinycache

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func DoTestCache(t *testing.T, c Cache) {
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error("should value nil and err not found")
	}
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error("should value nil and err not found")
	}
	if err := c.Set("a", []byte{'b'}, time.Minute*1); err != nil {
		t.Error(err)
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err)
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err)
	}
}

func TestLocal(t *testing.T) {
	DoTestCache(t, NewLocal(time.Minute*1, time.Minute*1))
}

func TestRedis(t *testing.T) {
	DoTestCache(t, NewRedis(redis.NewClient(&redis.Options{})))
}

func TestHybrid(t *testing.T) {
	DoTestCache(t, NewHybrid(redis.NewClient(
		&redis.Options{}), time.Minute*1, time.Minute*1))
}
