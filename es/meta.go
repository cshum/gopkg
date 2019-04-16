package es

import (
	"context"
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/olivere/elastic"
	"time"
)

type Meta struct {
	es      *elastic.Client
	m       map[string]string
	index   string
	doctype string
}

func NewMeta(es *elastic.Client, index, doctype string) *Meta {
	return &Meta{
		es:      es,
		m:       map[string]string{},
		index:   index,
		doctype: doctype,
	}
}

func (m *Meta) Load(ctx context.Context) error {
	ret, err := m.es.GetMapping().
		Index(m.index).Type(m.doctype).Do(ctx)
	if err != nil {
		return err
	}
	for _, data := range ret {
		parsed, err := gabs.Consume(data)
		if err != nil {
			return err
		}
		data = parsed.Search("mappings", m.doctype, "_meta").Data()
		if mp, ok := data.(map[string]interface{}); ok {
			for key, val := range mp {
				if str, ok := val.(string); ok {
					m.m[key] = str
				} else if bytes, err := json.Marshal(val); err == nil {
					m.m[key] = string(bytes)
				}
			}
		}
		return nil
	}
	return nil
}

func (m *Meta) Map() map[string]string {
	return m.m
}

func (m *Meta) Set(key, val string) {
	m.m[key] = val
}

func (m *Meta) Get(key string) (string, bool) {
	val, ok := m.m[key]
	return val, ok
}

func (m *Meta) SetTime(key string, t time.Time) {
	m.Set(key, t.Format(time.RFC3339Nano))
}

func (m *Meta) GetTime(key string) (time.Time, bool) {
	if val, ok := m.Get(key); ok {
		if t, err := time.Parse(time.RFC3339Nano, val); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func (m *Meta) Del(key string) {
	delete(m.m, key)
}

func (m *Meta) Save(ctx context.Context) error {
	_, err := m.es.PutMapping().
		Index(m.index).Type(m.doctype).
		BodyJson(map[string]interface{}{"_meta": m.m}).
		Do(ctx)
	return err
}
