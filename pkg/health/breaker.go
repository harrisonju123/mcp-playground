package health

import (
	"context"
	"fmt"
	"github.com/harrisonju123/mcp-agent-poc/internal/metrics"
	"log"
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
	start := time.Now()
	log.Printf("call start...")
	if trip := b.tripUntil.Load(); trip != 0 && start.UnixNano() < trip {
		return nil, fmt.Errorf("breaker open for %s", c.ID())
	}

	code := "ok"
	out, err := c.Call(ctx, in)
	if err != nil {
		code = "error"
	}
	lat := time.Since(start).Seconds()
	metrics.ToolLatency.WithLabelValues(c.ID(), code).Observe(lat)
	metrics.ToolTotal.WithLabelValues(c.ID(), code).Inc()

	// feed results to EWMA recorder
	b.r.Observe(time.Since(start), err)

	// if our latency is too slow, trip the circuit breaker
	if b.shouldTrip() {
		// only set a new trip window if we're not already tripped
		if old := b.tripUntil.Load(); old == 0 || start.UnixNano() >= old {
			b.tripUntil.Store(start.Add(probeInterval).UnixNano())
			go b.probe(c)
		}
	}
	return out, err
}

func (b *Breaker) probe(c Callable) {
	// wait for interval to elapse
	time.Sleep(probeInterval)

	// check health, if no error, assume it's healthy
	_, err := c.Call(context.Background(), nil)
	if err == nil {
		// success so clear tripUntil to close the breaker
		b.tripUntil.Store(0) // closes breaker
	} else {
		// failture so keep it open for another probeInterval
		b.tripUntil.Store(time.Now().Add(probeInterval).UnixNano())
		go b.probe(c)
	}
}

func (b *Breaker) shouldTrip() bool {
	return b.r.P95Latency() > latencyThreshold || b.r.SuccessRate() < successThreshold
}
