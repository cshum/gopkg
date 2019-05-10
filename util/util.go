package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Bool2Int8 convert bool to int8
func Bool2Int8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// CallerDir current caller abs file dir
func CallerDir() (string, error) {
	_, callerFileName, _, _ := runtime.Caller(1)
	return filepath.Abs(filepath.Dir(callerFileName))
}

// AnyOfString any of string
func AnyOfString(list ...string) string {
	for _, str := range list {
		if str = strings.TrimSpace(str); str != "" {
			return str
		}
	}
	return ""
}

func PrintJSON(vs ...interface{}) {
	for _, v := range vs {
		bytes, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Println(string(bytes))
		}
	}
}

func PrintJSONIndent(vs ...interface{}) {
	for _, v := range vs {
		bytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Println(string(bytes))
		}
	}
}

func ResolveURL(prefix, path string) (string, error) {
	u, err := url.Parse(prefix)
	if err != nil {
		return "", err
	}
	p, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return u.ResolveReference(p).String(), nil
}

func Timestamp() int64 {
	return ToTimestamp(time.Now())
}

func ToTimestamp(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func ToMD5Hash(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:]), nil
}
