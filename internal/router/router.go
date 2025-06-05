package router

import (
	"sync"
)

type Router struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

type ErrNotFound struct{ Tool string }

func (e ErrNotFound) Error() string {
	return "tool not found: " + e.Tool
}
func New() *Router { return &Router{tools: make(map[string]Tool)} }

// Replace Hot swap in a single swap, new map and overwrite pointer
// no partial updates, no in place deletions
func (r *Router) Replace(tools []Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools = make(map[string]Tool, len(tools))
	for _, t := range tools {
		r.tools[t.Name] = t
	}
}

// List reads concurrently and copies the map into a slice to iterate without holding lock
func (r *Router) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t)
	}

	return out
}

// Call calls a tool by getting read lock, fetch tool, release lock before handler execution.
// hot swap still valid since Tool values are copied.
func (r *Router) Call(name string, args []byte) ([]byte, error) {
	r.mu.RLock()
	h, ok := r.tools[name]
	r.mu.RUnlock()

	if !ok {
		return nil, ErrNotFound{name}
	}

	return h.Handler(args)
}
