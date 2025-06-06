package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	ToolLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tool_latencyseconds",
			Help:    "Latency of individual tool calls",
			Buckets: prometheus.DefBuckets, //p50 or p95 downstream via recording rules
		},
		[]string{"tool", "code"},
	)

	ToolTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tool_calls_total",
			Help: "Total tool invocations partitioned by tool & status",
		},
		[]string{"tool", "code"},
	)
)

func Register() {
	prometheus.MustRegister(ToolLatency, ToolTotal, collectors.NewBuildInfoCollector())
}
