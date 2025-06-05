package health

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// Breaker Half-open with throttled probes
// Fast fail when open. No request queuing
// Wrapped tool owns its own breaker; routing picks first healthy tool.
type Breaker struct {
	r         *Recorder
	tripUntil atomic.Int64 //unix nano timestamp; 0 == closed
}

func NewBreaker() *Breaker {
	return &Breaker{r: &Recorder{}}
}

var (
	latencyThreshold = 1 * time.Second
	successThreshold = 0.95
	probeInterval    = 2 * time.Second
)

func (b *Breaker) Call(ctx context.Context, c Callable, in []byte) ([]byte, error) {
	now := time.Now()
	if trip := b.tripUntil.Load(); trip != 0 && now.UnixNano() < trip {
		return nil, fmt.Errorf("breaker open for %s", c.ID())
	}

	out, err := c.Call(ctx, in)
	b.r.Observe(time.Since(now), err)

	if b.shouldTrip() {
		b.tripUntil.Store(now.Add(probeInterval).UnixNano())
		go b.probe(c)
	}
	return out, err
}

func (b *Breaker) probe(c Callable) {
	time.Sleep(probeInterval)
	_, err := c.Call(context.Background(), nil)
	if err != nil {
		b.tripUntil.Store(0) // closes breaker
	} else {
		b.tripUntil.Store(time.Now().Add(probeInterval).UnixNano())
	}
}

func (b *Breaker) shouldTrip() bool {
	return b.r.P95Latency() > latencyThreshold || b.r.SuccessRate() < successThreshold
}
