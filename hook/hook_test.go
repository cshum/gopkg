package hook

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	from := time.Now()
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "hello", "world"), nil
	})
	ctx, err := Invoke(context.Background(), "a")
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
	Clear("a")
}

func TestHookWithError(t *testing.T) {
	from := time.Now()
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return ctx, errors.New("err")
	})
	ctx, err := Invoke(context.Background(), "a")
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
	Clear("a")
}

func TestHookParallel(t *testing.T) {
	from := time.Now()
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "hello", "world"), nil
	})
	ctxs, errs := Parallel(context.Background(), "a")
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
	Reset()
}

func TestHookWithErrorParallel(t *testing.T) {
	from := time.Now()
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 300)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return ctx, errors.New("err")
	})
	ctxs, errs := Parallel(context.Background(), "a")
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
	Reset()
}
