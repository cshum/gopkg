package util

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// ParamInt parse int from chi URL param
func ParamInt(r *http.Request, key string) (int, error) {
	val, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

// PanicLogger capture errors to raven, optional logger
func PanicLogger(logger *zap.Logger) {
	if r := recover(); r != nil {
		raven.CaptureError(r.(error), nil)
		logger.Panic("panic", zap.Error(r.(error)))
	}
}

// Bool2Int8 convert bool to int8
func Bool2Int8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// CurrentDir current source file dir
func CurrentDir() (string, error) {
	_, callerFileName, _, _ := runtime.Caller(1)
	return filepath.Abs(filepath.Dir(callerFileName))
}

// AnyOfString any of string
func AnyOfString(list ...string) string {
	for _, str := range list {
		str = strings.TrimSpace(str)
		if str != "" {
			return str
		}
	}
	return ""
}
