package health

import "context"

type Callable interface {
	Call(ctx context.Context, in []byte) ([]byte, error)
	ID() string // use this for metrics and logs
}
