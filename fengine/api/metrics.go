package api

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/db/sql"
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

func (mm metricsMiddleware) GetEntity(ctx context.Context, id string) (r fengine.Result, e error) {
	defer mm.count("GetEntity")(time.Now())
	return mm.svc.GetEntity(ctx, id)
}

func (mm metricsMiddleware) UpsertEntity(ctx context.Context, entityDef sql.EntityDefinition) (r fengine.Result, e error) {
	defer mm.count("UpsertEntity")(time.Now())
	return mm.svc.UpsertEntity(ctx, entityDef)
}

func (mm metricsMiddleware) DeleteEntity(ctx context.Context, id string) (r fengine.Result, e error) {
	defer mm.count("DeleteEntity")(time.Now())
	return mm.svc.DeleteEntity(ctx, id)
}

func (mm metricsMiddleware) ExecuteService(ctx context.Context, script sql.ServiceRequest) (fengine.Result, error) {
	defer mm.count("ExecuteService")(time.Now())
	return mm.svc.ExecuteService(ctx, script)
}

func (mm metricsMiddleware) CreateTable(ctx context.Context, definition sql.TableDefinition) (fengine.Result, error) {
	defer mm.count("CreateTable")(time.Now())
	return mm.svc.CreateTable(ctx, definition)
}

func (mm metricsMiddleware) Select(ctx context.Context, req sql.SelectRequest) (*sql.ResultSet, error) {
	defer mm.count("Select")(time.Now())
	return mm.svc.Select(ctx, req)
}

func (mm metricsMiddleware) Insert(ctx context.Context, request sql.InsertRequest) (fengine.Result, error) {
	defer mm.count("Insert")(time.Now())
	return mm.svc.Insert(ctx, request)
}

func (mm metricsMiddleware) Update(ctx context.Context, request sql.UpdateRequest) (fengine.Result, error) {
	defer mm.count("Update")(time.Now())
	return mm.svc.Update(ctx, request)
}

func (mm metricsMiddleware) Delete(ctx context.Context, request sql.DeleteRequest) (fengine.Result, error) {
	defer mm.count("Delete")(time.Now())
	return mm.svc.Delete(ctx, request)
}

func (mm metricsMiddleware) GetThingAllServices(ctx context.Context, uuid uuid.UUID) (fengine.Result, error) {
	defer mm.count("GetThingAllServices")(time.Now())
	return mm.svc.GetThingAllServices(ctx, uuid)
}

func (mm metricsMiddleware) GetThingService(ctx context.Context, id sql.ThingServiceId) (fengine.Result, error) {
	defer mm.count("GetThingService")(time.Now())
	return mm.svc.GetThingService(ctx, id)
}

func (mm metricsMiddleware) count(name string) func(begin time.Time) {
	return func(begin time.Time) {
		mm.counter.With("method", name).Add(1)
		mm.latency.With("method", name).Observe(time.Since(begin).Seconds())
	}
}
