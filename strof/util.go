package strof

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

func Hash(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:]), nil
}

// NonEmpty any of string
func NonEmpty(strs ...string) string {
	for _, str := range strs {
		if str = strings.TrimSpace(str); str != "" {
			return str
		}
	}
	return ""
}

// CallerDir current caller abs file dir
func CallerDir() (string, error) {
	_, callerFileName, _, _ := runtime.Caller(1)
	return filepath.Abs(filepath.Dir(callerFileName))
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

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
