package router_test

import (
	"github.com/harrisonju123/mcp-agent-poc/router"
	"testing"
)

func TestHotSwap(t *testing.T) {
	r := router.New()
	a := router.Tool{
		Name: "a",
		Handler: func(b []byte) ([]byte, error) {
			return b, nil
		},
	}
	r.Replace([]router.Tool{a})
	if len(r.List()) != 1 {
		t.Fatal("expected 1 tool")
	}

	b := router.Tool{Name: "b", Handler: a.Handler}
	r.Replace([]router.Tool{b})
	if _, err := r.Call("a", nil); err == nil {
		t.Fatal("expected error")
	}
	if _, err := r.Call("b", nil); err != nil {
		t.Fatalf("new tool not callable: %v", err)
	}
}
