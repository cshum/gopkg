package env

import (
	"sync"

	"github.com/gobuffalo/envy"
	"go.uber.org/zap"
)

// Env env loader based on envy, and show only loaded values
type Env struct {
	// only loaded on startup, should no need bother concurrent
	m map[string]string
	l *sync.RWMutex
}

// New env loader
func New(files ...string) *Env {
	_ = envy.Load(files...) // load .env file
	return &Env{
		m: make(map[string]string),
		l: &sync.RWMutex{},
	}
}

// Get env loader
func (m *Env) Get(key string, value string) string {
	m.l.Lock()
	defer m.l.Unlock()
	val := envy.Get(key, value)
	m.m[key] = val
	return val
}

// Env env loader, show only loaded values
func (m *Env) Map() map[string]string {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m
}

// ZapFields return env map in for of zap.Field slices
func (m *Env) ZapFields() []zap.Field {
	m.l.RLock()
	defer m.l.RUnlock()
	fields := make([]zap.Field, 0, len(m.m))
	for k, v := range m.m {
		fields = append(fields, zap.String(k, v))
	}
	return fields
}
