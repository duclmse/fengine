package api

import (
	"context"
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

func (mm metricsMiddleware) CreateEntity(ctx context.Context, entityDef *fengine.EntityDefinition) (r *fengine.Result, e error) {
	defer mm.count("UpsertEntity")(time.Now())
	return mm.svc.CreateEntity(ctx, entityDef)
}

func (mm metricsMiddleware) ExecuteService(ctx context.Context, script *fengine.ServiceRequest) (*fengine.Result, error) {
	defer mm.count("ExecuteService")(time.Now())
	return mm.svc.ExecuteService(ctx, script)
}

func (mm metricsMiddleware) CreateTable(ctx context.Context, definition *fengine.TableDefinition) (*fengine.Result, error) {
	defer mm.count("CreateTable")(time.Now())
	return mm.svc.CreateTable(ctx, definition)
}

func (mm metricsMiddleware) Select(ctx context.Context, request *fengine.SelectRequest) (*fengine.Result, error) {
	defer mm.count("Select")(time.Now())
	return mm.svc.Select(ctx, request)
}

func (mm metricsMiddleware) Insert(ctx context.Context, request *fengine.InsertRequest) (*fengine.Result, error) {
	defer mm.count("Insert")(time.Now())
	return mm.svc.Insert(ctx, request)
}

func (mm metricsMiddleware) Update(ctx context.Context, request *fengine.UpdateRequest) (*fengine.Result, error) {
	defer mm.count("Update")(time.Now())
	return mm.svc.Update(ctx, request)
}

func (mm metricsMiddleware) Delete(ctx context.Context, request *fengine.DeleteRequest) (*fengine.Result, error) {
	defer mm.count("Delete")(time.Now())
	return mm.svc.Delete(ctx, request)
}

func (mm metricsMiddleware) GetThingAllServices(ctx context.Context, uuid string) (*fengine.Result, error) {
	defer mm.count("GetThingAllServices")(time.Now())
	return mm.svc.GetThingAllServices(ctx, uuid)
}

func (mm metricsMiddleware) GetThingService(ctx context.Context, uuid string, s string) (*fengine.Result, error) {
	defer mm.count("GetThingService")(time.Now())
	return mm.svc.GetThingService(ctx, uuid, s)
}

func (mm metricsMiddleware) count(name string) func(begin time.Time) {
	return func(begin time.Time) {
		mm.counter.With("method", name).Add(1)
		mm.latency.With("method", name).Observe(time.Since(begin).Seconds())
	}
}
