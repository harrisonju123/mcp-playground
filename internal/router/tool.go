package router

import "context"

type Handler func([]byte) ([]byte, error)

type Tool struct {
	Name        string
	Description string
	Handler     Handler
}

func (t Tool) Call(_ context.Context, in []byte) ([]byte, error) {
	return t.Handler(in)
}
