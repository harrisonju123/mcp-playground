package health

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/harrisonju123/mcp-agent-poc/internal/router"
)

type Breaker struct {
	r         *Recorder
	tripUntil atomic.Int64 //unix nano timestamp; 0 == closed
}

var (
	latencyThreshold = 1 * time.Second
	successThreshold = 0.95
	probeInterval    = 2 * time.Second
)

func (b *Breaker) Call(ctx context.Context, t router.Tool, in []byte) ([]byte, error) {
	now := time.Now()
	if trip := b.tripUntil.Load(); trip != 0 && now.UnixNano() < trip {
		return nil, fmt.Errorf("breaker open for %s", t.Name)
	}

	out, err ;= t.Call(ctx, in)

}
