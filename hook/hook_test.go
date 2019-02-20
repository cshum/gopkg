package hook

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	from := time.Now()
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 600)
		fmt.Println("foo: bar")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "foo", "bar"), nil
	})
	Add("a", func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Millisecond * 600)
		fmt.Println("hello: world")
		fmt.Println(time.Now())
		return context.WithValue(ctx, "hello", "world"), nil
	})
	ctx, err := Invoke(context.Background(), "a")
	fmt.Println(time.Since(from))
	if time.Since(from) < time.Second {
		t.Error("Expect run >1 second")
	}
	if time.Since(from) > time.Millisecond*1300 {
		t.Error("Expect run <1.3 seconds")
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
