package router_test

import (
	"context"
	"testing"

	"github.com/harrisonju123/mcp-agent-poc/internal/router"
)

func TestHotSwap(t *testing.T) {
	r := router.NewRouter(nil)

	// Tool "a"
	a := router.Tool{
		Name: "a",
		Handler: func(ctx context.Context, b []byte) ([]byte, error) {
			return b, nil
		},
	}
	r.Replace([]router.Tool{a})

	// Assert tool "a" is registered
	if tools := r.List(); len(tools) != 1 || tools[0].Name != "a" {
		t.Fatalf("expected tool 'a', got: %+v", tools)
	}

	// Swap to tool "b"
	b := router.Tool{
		Name:    "b",
		Handler: a.Handler,
	}
	r.Replace([]router.Tool{b})

	// Call tool "a" should fail
	if _, err := r.Call(context.Background(), "a", nil); err == nil {
		t.Fatal("expected error when calling removed tool 'a'")
	}

	// Call tool "b" should succeed
	if _, err := r.Call(context.Background(), "b", nil); err != nil {
		t.Fatalf("expected tool 'b' to be callable, got error: %v", err)
	}
}

func BenchmarkCall(b *testing.B) {
	r := router.NewRouter([]router.Tool{{
		Name:    "echo",
		Handler: func(ctx context.Context, in []byte) ([]byte, error) { return in, nil },
	}})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = r.Call(context.Background(), "echo", nil)
		}
	})
}
