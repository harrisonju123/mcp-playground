package router

import (
	"context"
	"sync/atomic"
	"unsafe"
)

type ErrNotFound struct{ Tool string }

func (e ErrNotFound) Error() string {
	return "tool not found: " + e.Tool
}

type mapPtr = unsafe.Pointer

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
	for _, t := range tools {
		m[t.Name] = t
	}
	atomic.StorePointer(&r.table, unsafe.Pointer(&m))
}

// List returns a snapshot slice ( read only to callers)
func (r *Router) List() []Tool {
	tab := *(*map[string]Tool)(atomic.LoadPointer(&r.table))

	out := make([]Tool, 0, len(tab))
	for _, t := range tab {
		out = append(out, t)
	}

	return out
}

// Call looks up the tool without locks; safe becauase map is immutable
func (r *Router) Call(ctx context.Context, name string, in []byte) ([]byte, error) {
	tab := *(*map[string]Tool)(atomic.LoadPointer(&r.table))
	t, ok := tab[name]
	if !ok {
		return nil, ErrNotFound{name}
	}

	return t.Call(ctx, in)
}
