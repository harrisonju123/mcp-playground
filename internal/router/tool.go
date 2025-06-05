package router

import (
	"context"
)

type Handler func(context.Context, []byte) ([]byte, error)

type Tool struct {
	Name        string
	Description string
	Handler     Handler
}

func (t Tool) ID() string {
	return t.Name
}

func (t Tool) Call(ctx context.Context, in []byte) ([]byte, error) { return t.Handler(ctx, in) }
