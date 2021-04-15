package tinycache

import (
	"encoding/json"
)

func Marshal(c Cache, key string, v interface{}) (err error) {
	var raw []byte
	if raw, err = json.Marshal(v); err != nil {
		return
	}
	if err = c.Set(key, raw); err != nil {
		return
	}
	return
}

func Unmarshal(c Cache, key string, v interface{}) (err error) {
	var raw []byte
	if raw, err = c.Get(key); err != nil {
		return
	}
	if err = json.Unmarshal(raw, v); err != nil {
		return
	}
	return
}
