package tinycache

import (
	"github.com/cshum/gopkg/errw"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
}

var NotFound = errw.NotFound("tinycache: not found")
