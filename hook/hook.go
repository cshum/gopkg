package hook

import (
	"context"
	"sync"
)

var lock = &sync.RWMutex{}
var hooks = map[interface{}][]Handler{}

type Handler func(ctx context.Context) (context.Context, error)

// Add add hook
func Add(hookType interface{}, hook Handler) {
	lock.Lock()
	defer lock.Unlock()
	hooks[hookType] = append(hooks[hookType], hook)
}

// Clear clear hook func by type
func Clear(hookType interface{}) {
	lock.Lock()
	defer lock.Unlock()
	hooks[hookType] = []Handler{}
}

// Reset reset all hooks
func Reset() {
	lock.Lock()
	defer lock.Unlock()
	hooks = map[interface{}][]Handler{}
}

func getByType(hookType interface{}) ([]Handler, int) {
	lock.RLock()
	defer lock.RUnlock()
	cnt := len(hooks[hookType])
	fns := make([]Handler, cnt)
	copy(fns, hooks[hookType])
	return fns, cnt
}

// Invoke invoke hook
func Invoke(ctx context.Context, hookType interface{}) (context.Context, error) {
	fns, _ := getByType(hookType)
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
func Parallel(ctx context.Context, hookType interface{}) ([]context.Context, []error) {
	fns, cnt := getByType(hookType)
	ctxs := make(chan context.Context, cnt)
	errors := make(chan error, cnt)
	for _, fn := range fns {
		go func(
			fn Handler,
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
	for i := 0; i < cnt; i++ {
		select {
		case c := <-ctxs:
			clist = append(clist, c)
		case e := <-errors:
			elist = append(elist, e)
		}
	}
	close(ctxs)
	close(errors)
	return clist, elist
}
