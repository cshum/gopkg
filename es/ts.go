package es

import (
	"time"
)

func Timestamp() int64 {
	return ToTimestamp(time.Now())
}

func ToTimestamp(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
