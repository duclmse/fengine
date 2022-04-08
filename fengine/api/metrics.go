package api

import (
	"context"
	"time"

	"github.com/duclmse/fengine/fengine"
	"github.com/go-kit/kit/metrics"
)

var _ fengine.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fengine.Service
}

func MetricsMiddleware(svc fengine.Service, counter metrics.Counter, latency metrics.Histogram) fengine.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

// Get implements fengine.Service
func (mm metricsMiddleware) Get(ctx context.Context, id string) (interface{}, error) {
	defer mm.count("Get")(time.Now())
	return mm.svc.Get(ctx, id)
}

func (mm metricsMiddleware) count(name string) func(begin time.Time) {
	return func(begin time.Time) {
		mm.counter.With("method", name).Add(1)
		mm.latency.With("method", name).Observe(time.Since(begin).Seconds())
	}
}
