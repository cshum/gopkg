package hook

import (
	"context"
	"sync"
)

var lock = &sync.RWMutex{}
var hooks = map[string][]func(context.Context) (context.Context, error){}

// Add add hook
func Add(hookType string, hook func(ctx context.Context) (context.Context, error)) {
	lock.Lock()
	defer lock.Unlock()
	hooks[hookType] = append(hooks[hookType], hook)
}

// Invoke invoke hook
func Invoke(ctx context.Context, hookType string) (context.Context, error) {
	lock.RLock()
	defer lock.RUnlock()
	var err error
	for _, fn := range hooks[hookType] {
		ctx, err = fn(ctx)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
