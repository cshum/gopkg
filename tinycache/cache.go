package tinycache

import (
	"github.com/cshum/gopkg/errw"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

var NotFound = errw.NotFound("tinycache: not found")
