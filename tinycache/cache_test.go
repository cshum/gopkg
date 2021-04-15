package tinycache

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func DoTestCache(t *testing.T, c Cache, sleepTime time.Duration) {
	// not found
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	// set and found
	if err := c.Set("a", []byte{'b'}); err != nil {
		t.Error(err)
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err, "should value and no error")
	}
	if v, err := c.Get("a"); string(v) != "b" || err != nil {
		t.Error(err, "should value and no error")
	}
	time.Sleep(sleepTime)
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

	// test JSON Marshal Unmarshal
	var v []int
	if err := Unmarshal(c, "a", v); err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	if err := Unmarshal(c, "a", v); err != NotFound {
		t.Error(err, "should value nil and err not found")
	}
	// set and found
	if err := Marshal(c, "a", []int{1, 2, 3}); err != nil {
		t.Error(err)
	}
	if raw, err := c.Get("a"); string(raw) != "[1,2,3]" || err != nil {
		t.Error(err, "should value and no error")
	}
	if err := Unmarshal(c, "a", &v); len(v) != 3 || v[2] != 3 || err != nil {
		t.Error(err, "should value and no error")
	}
	time.Sleep(sleepTime)
	if v, err := c.Get("a"); v != nil || err != NotFound {
		t.Error(v, err, "should value nil and err not found")
	}
}

func TestLocal(t *testing.T) {
	DoTestCache(t, NewLocal(time.Millisecond*100), time.Millisecond*500)
}

func TestRedis(t *testing.T) {
	DoTestCache(t, NewRedis(&redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", ":6379")
		},
	}, time.Millisecond*100), time.Millisecond*500)
}

//func TestHybrid(t *testing.T) {
//	DoTestCache(t, NewHybrid(&redis.Pool{
//		Dial: func() (conn redis.Conn, err error) {
//			return redis.Dial("tcp", ":6379")
//		},
//	}, time.Minute*1))
//}
//
//func TestHybridRedis(t *testing.T) {
//	DoTestCache(t, NewHybrid(&redis.Pool{
//		Dial: func() (conn redis.Conn, err error) {
//			return redis.Dial("tcp", ":6379")
//		},
//	}, time.Nanosecond))
//}
