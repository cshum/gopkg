package simplecache

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func DoTestCache(t *testing.T, c Cache) {
	// not found
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	// set and found
	if err := c.Set("a", []byte{'b'}, time.Millisecond*100); err != nil {
		t.Error(err)
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err, "should value and no error")
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err, "should value and no error")
	}
	time.Sleep(time.Millisecond * 500)
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(v, err, "should value nil and err not found")
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

	// test JSON JSONMarshal JSONUnmarshal
	var v []int
	if err := JSONUnmarshal(c, "a", v); err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	if err := JSONUnmarshal(c, "a", v); err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	// set and found
	if err := JSONMarshal(c, "a", []int{1, 2, 3}, time.Millisecond*100); err != nil {
		t.Error(err)
	}
	if raw, err := c.Get("a"); string(raw) != "[1,2,3]" || err != nil {
		t.Error(err, "should value and no error")
	}
	if err := JSONUnmarshal(c, "a", &v); len(v) != 3 || v[2] != 3 || err != nil {
		t.Error(err, "should value and no error")
	}
	time.Sleep(time.Millisecond * 500)
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(v, err, "should value nil and err not found")
	}
}

func TestLocal(t *testing.T) {
	DoTestCache(t, NewLocal(10, int64(10<<20)))
}

func TestRedis(t *testing.T) {
	DoTestCache(t, NewRedis(&redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", ":6379")
		},
	}))
}

func TestHybrid(t *testing.T) {
	DoTestCache(t, NewHybrid(&redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", ":6379")
		},
	}, NewLocal(10, int64(10<<20)), time.Minute*1))
}

func TestHybridRedis(t *testing.T) {
	DoTestCache(t, NewHybrid(&redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", ":6379")
		},
	}, NewLocal(10, int64(10<<20)), time.Nanosecond))
}
