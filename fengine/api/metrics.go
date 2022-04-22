package api

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/duclmse/fengine/fengine"
	"github.com/go-kit/kit/metrics"
)

var _ fengine.Service = (*metricsMiddleware)(nil)

func MetricsMiddleware(svc fengine.Service, counter metrics.Counter, latency metrics.Histogram) fengine.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fengine.Service
}

func (mm metricsMiddleware) ExecuteService(ctx context.Context, script *fengine.JsonScript) (*fengine.Result, error) {
	defer mm.count("ExecuteService")(time.Now())
	return mm.svc.ExecuteService(ctx, script)
}

func (mm metricsMiddleware) Select(ctx context.Context, request *fengine.JsonSelectRequest) (*fengine.Result, error) {
	defer mm.count("Select")(time.Now())
	return mm.svc.Select(ctx, request)
}

func (mm metricsMiddleware) Insert(ctx context.Context, request *fengine.JsonInsertRequest) (*fengine.Result, error) {
	defer mm.count("Insert")(time.Now())
	return mm.svc.Insert(ctx, request)
}

func (mm metricsMiddleware) Update(ctx context.Context, request *fengine.JsonUpdateRequest) (*fengine.Result, error) {
	defer mm.count("Update")(time.Now())
	return mm.svc.Update(ctx, request)
}

func (mm metricsMiddleware) Delete(ctx context.Context, request *fengine.JsonDeleteRequest) (*fengine.Result, error) {
	defer mm.count("Delete")(time.Now())
	return mm.svc.Delete(ctx, request)
}

func (mm metricsMiddleware) GetThingAllServices(ctx context.Context, uuid uuid.UUID) (*fengine.Result, error) {
	defer mm.count("GetThingAllServices")(time.Now())
	return mm.svc.GetThingAllServices(ctx, uuid)
}

func (mm metricsMiddleware) GetThingService(ctx context.Context, uuid uuid.UUID, s string) (*fengine.Result, error) {
	defer mm.count("GetThingService")(time.Now())
	return mm.svc.GetThingService(ctx, uuid, s)
}

func (mm metricsMiddleware) count(name string) func(begin time.Time) {
	return func(begin time.Time) {
		mm.counter.With("method", name).Add(1)
		mm.latency.With("method", name).Observe(time.Since(begin).Seconds())
	}
}
