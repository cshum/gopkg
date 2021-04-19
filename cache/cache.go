package cache

import (
	"time"

	"github.com/cshum/gopkg/errw"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
}

var NotFound = errw.NotFound("cache: not found")

func toMillis(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

func fromMillis(pTTL int64) time.Duration {
	return time.Duration(pTTL) * time.Millisecond
}
