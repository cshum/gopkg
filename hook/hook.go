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

// Clear clear hook func by type
func Clear(hookType string) {
	lock.Lock()
	defer lock.Unlock()
	hooks[hookType] = []func(context.Context) (context.Context, error){}
}

// Reset reset all hooks
func Reset() {
	lock.Lock()
	defer lock.Unlock()
	hooks = map[string][]func(context.Context) (context.Context, error){}
}

func getByType(hookType string) []func(context.Context) (context.Context, error) {
	lock.RLock()
	defer lock.RUnlock()
	fns := make([]func(context.Context) (context.Context, error), len(hooks[hookType]))
	copy(fns, hooks[hookType])
	return fns
}

// Invoke invoke hook
func Invoke(ctx context.Context, hookType string) (context.Context, error) {
	fns := getByType(hookType)
	var err error
	for _, fn := range fns {
		ctx, err = fn(ctx)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

// Parallel invoke hook in parallel
func Parallel(ctx context.Context, hookType string) ([]context.Context, []error) {
	fns := getByType(hookType)
	ctxs := make(chan context.Context)
	errors := make(chan error)
	count := len(fns)
	for _, fn := range fns {
		go func(
			fn func(ctx context.Context) (context.Context, error),
			ctxs chan<- context.Context,
			errors chan<- error,
		) {
			c, err := fn(ctx)
			if err != nil {
				errors <- err
			} else {
				ctxs <- c
			}
		}(fn, ctxs, errors)
	}
	var clist []context.Context
	var elist []error
	for i := 0; i < count; i++ {
		select {
		case c := <-ctxs:
			clist = append(clist, c)
		case e := <-errors:
			elist = append(elist, e)
		}
	}
	return clist, elist
}
