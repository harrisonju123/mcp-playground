package health

import (
	"sync/atomic"
	"time"
)

type Recorder struct {
	latencyEwma uint64 //nano seconds 2^16
	successEwma uint32 //0-1

	totalCalls  uint64
	totalErrors uint64
}

const (
	latAlphaShift = 4 //EWMA a = 1/16
	sucAlphaShift = 3 // EWMA a = 1/8
)

func (r *Recorder) Observe(latency time.Duration, err error) {
	// << 16 is the same thing as multiplying by 2^16
	latNs := uint64(latency.Nanoseconds()) << 16
	oldLat := atomic.LoadUint64(&r.latencyEwma)
	newLat := oldLat - (oldLat >> latAlphaShift) + (latNs >> latAlphaShift)
	atomic.StoreUint64(&r.latencyEwma, newLat)

	var suc uint32
	if err == nil {
		suc = 0xFFFF
	}
	oldSuc := atomic.LoadUint32(&r.successEwma)
	newSuc := oldSuc - (oldSuc >> sucAlphaShift) + (suc >> sucAlphaShift)
	atomic.StoreUint32(&r.successEwma, newSuc)

	atomic.AddUint64(&r.totalCalls, 1)
	if err != nil {
		atomic.AddUint64(&r.totalErrors, 1)
	}
}

func (r *Recorder) P95Latency() time.Duration {
	ns := atomic.LoadUint64(&r.latencyEwma) >> 16
	return time.Duration(ns)
}

func (r *Recorder) SuccessRate() float64 {
	return float64(atomic.LoadUint32(&r.successEwma)) / 65535
}
