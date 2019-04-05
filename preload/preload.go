package preload

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Handler func(ctx context.Context) error

var before []Handler
var preloads []Handler
var after []Handler
var once sync.Once
var err error

func Before(fns ...Handler) {
	before = append(before, fns...)
}

func Add(fns ...Handler) {
	preloads = append(preloads, fns...)
}

func After(fns ...Handler) {
	after = append(after, fns...)
}

func invoke(ctx context.Context) error {
	if preloads == nil {
		return nil
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, fn := range preloads {
		func(fn Handler) {
			g.Go(func() error {
				return fn(ctx)
			})
		}(fn)
	}
	return g.Wait()
}

// Preload app Preload hook trigger
func Do() error {
	once.Do(func() {
		// Preload hook
		ctx := context.Background()
		if before != nil {
			for _, fn := range before {
				if err = fn(ctx); err != nil {
					return
				}
			}
		}
		if err = invoke(ctx); err != nil {
			return
		}
		if after != nil {
			for _, fn := range after {
				if err = fn(ctx); err != nil {
					return
				}
			}
		}
	})
	return err
}

func MustDo() {
	if Do() != nil {
		panic(err)
	}
}
