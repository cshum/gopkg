package tinycache

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func DoTestCache(t *testing.T, c Cache) {
	// not found
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error("should value nil and err not found")
	}
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error("should value nil and err not found")
	}
	// set and found
	if err := c.Set("a", []byte{'b'}, time.Minute*1); err != nil {
		t.Error(err)
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error("should value and no error")
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error("should value and no error")
	}
	// set nil and found nil
	/*
		if err := c.Set("n", nil, time.Minute*1); err != nil {
			t.Error("should nil value and no error")
		}
		if v, err := c.Get("n"); v != nil || err != nil {
			t.Error("should nil value and no error")
		}
		if v, err := c.Get("n"); v != nil || err != nil {
			t.Error("should nil value and no error")
		}
	*/
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
