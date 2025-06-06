package router

import (
	"context"
	"github.com/harrisonju123/mcp-agent-poc/pkg/health"
	"sync/atomic"
	"unsafe"
)

type ErrNotFound struct{ Tool string }

func (e ErrNotFound) Error() string {
	return "tool not found: " + e.Tool
}

type mapPtr = unsafe.Pointer

// tableContents holds both the tools map and the breakers map side by side.
// Once built, both amps are immutable, so we can swap the entire struct atomically.
type tableContents struct {
	tools    map[string]Tool
	breakers map[string]*health.Breaker
}

type Router struct {
	table mapPtr
}

func NewRouter(tools []Tool) *Router {
	r := &Router{}
	r.Replace(tools)
	return r
}

// Replace atomically swaps in a brand new map.
// No partial updates, callers will never observe an intermediate state.
func (r *Router) Replace(tools []Tool) {
	m := make(map[string]Tool, len(tools))
	b := make(map[string]*health.Breaker, len(tools))

	for _, t := range tools {
		m[t.Name] = t
		b[t.Name] = health.NewBreaker()
	}

	newContents := &tableContents{
		tools:    m,
		breakers: b,
	}
	atomic.StorePointer(&r.table, unsafe.Pointer(newContents))
}

// List returns a snapshot slice ( read only to callers)
func (r *Router) List() []Tool {
	ptr := (*tableContents)(atomic.LoadPointer(&r.table))
	out := make([]Tool, 0, len(ptr.tools))
	for _, t := range ptr.tools {
		out = append(out, t)
	}

	return out
}

// Call looks up the tool without locks; safe becauase map is immutable
func (r *Router) Call(ctx context.Context, name string, in []byte) ([]byte, error) {
	ptr := (*tableContents)(atomic.LoadPointer(&r.table))

	t, ok := ptr.tools[name]
	if !ok {
		return nil, ErrNotFound{name}
	}

	br := ptr.breakers[name]
	return br.Call(ctx, &t, in)
}
