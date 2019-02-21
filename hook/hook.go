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
	result := make(chan context.Context)
	errors := make(chan error)
	count := len(fns)
	for _, fn := range fns {
		go func(
			fn func(ctx context.Context) (context.Context, error),
			result chan<- context.Context,
			errors chan<- error,
		) {
			c, err := fn(ctx)
			if err != nil {
				errors <- err
			} else {
				result <- c
			}
		}(fn, result, errors)
	}
	var clist []context.Context
	var elist []error
	for i := 0; i < count; i++ {
		select {
		case c := <-result:
			clist = append(clist, c)
		case e := <-errors:
			elist = append(elist, e)
		}
	}
	return clist, elist
}
