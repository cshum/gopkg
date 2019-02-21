package hook

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func DoInvoke(t *testing.T, hookType string) {
	from := time.Now()
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "hello", "world"), nil
	})
	ctx, err := Invoke(context.Background(), hookType)
	fmt.Println(time.Since(from))
	if time.Since(from) < time.Millisecond*600 {
		t.Error("Expect run >=0.6 second")
	}
	if time.Since(from) > time.Millisecond*700 {
		t.Error("Expect run <0.7 seconds")
	}
	if err != nil {
		t.Error("Expect err nil")
	}
	if ctx.Value("foo") != "bar" {
		t.Error("Expect ctx value foo: bar")
	}
	if ctx.Value("hello") != "world" {
		t.Error("Expect ctx value hello: world")
	}
}

func DoInvokeWithError(t *testing.T, hookType string) {
	from := time.Now()
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return ctx, errors.New("err")
	})
	ctx, err := Invoke(context.Background(), hookType)
	fmt.Println(time.Since(from))
	if time.Since(from) < time.Millisecond*600 {
		t.Error("Expect run >=0.6 second")
	}
	if time.Since(from) > time.Millisecond*700 {
		t.Error("Expect run <0.7 seconds")
	}
	if err.Error() != "err" {
		t.Error("Expect err")
	}
	if ctx.Value("foo") != "bar" {
		t.Error("Expect ctx value foo: bar")
	}
}

func DoInvokeParallel(t *testing.T, hookType string) {
	from := time.Now()
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "hello", "world"), nil
	})
	ctxs, errs := Parallel(context.Background(), hookType)
	fmt.Println(time.Since(from))
	if time.Since(from) < time.Millisecond*300 {
		t.Error("Expect run >=0.3 second")
	}
	if time.Since(from) > time.Millisecond*400 {
		t.Error("Expect run <0.4 second")
	}
	if len(errs) != 0 {
		t.Error("Expect no errors")
	}
	if len(ctxs) != 2 {
		t.Error("Expect 2 result")
	}
	if !((ctxs[0].Value("foo") == "bar" && ctxs[1].Value("hello") == "world") ||
		(ctxs[1].Value("foo") == "bar" && ctxs[0].Value("hello") == "world")) {
		t.Error("Invalid result")
	}
}

func DoInvokeWithErrorParallel(t *testing.T, hookType string) {
	from := time.Now()
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add(hookType, func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return ctx, errors.New("err")
	})
	ctxs, errs := Parallel(context.Background(), hookType)
	fmt.Println(time.Since(from))
	if time.Since(from) < time.Millisecond*300 {
		t.Error("Expect run >=0.3 second")
	}
	if time.Since(from) > time.Millisecond*400 {
		t.Error("Expect run <0.4 second")
	}
	if len(errs) != 1 {
		t.Error("Expect 1 error")
	}
	if len(ctxs) != 1 {
		t.Error("Expect 1 result")
	}
	if ctxs[0].Value("foo") != "bar" {
		t.Error("Invalid result")
	}
}

func TestHook(t *testing.T) {
	t.Parallel()
	DoInvoke(t, "a")
	Clear("a")
}

func TestHookWithError(t *testing.T) {
	t.Parallel()
	DoInvokeWithError(t, "b")
	Clear("b")
}

func TestHookParallel(t *testing.T) {
	t.Parallel()
	DoInvokeParallel(t, "c")
	Clear("c")
}

func TestHookWithErrorParallel(t *testing.T) {
	t.Parallel()
	DoInvokeWithErrorParallel(t, "d")
	Clear("d")
}

func TestHookMixed(t *testing.T) {
	DoInvoke(t, "e")
	DoInvokeWithError(t, "f")
	DoInvokeParallel(t, "g")
	DoInvokeWithErrorParallel(t, "h")
	Reset()
}
