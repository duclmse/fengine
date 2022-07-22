package tracing

import (
	"context"
	. "github.com/google/uuid"
	. "github.com/opentracing/opentracing-go"

	. "github.com/duclmse/fengine/fengine/db/cache"
	db "github.com/duclmse/fengine/fengine/db/sql"
)

var (
	_ db.Repository = (*fengineRepositoryMiddleware)(nil)
	_ Cache         = (*fengineCacheMiddleware)(nil)
)

//#region FEngineRepositoryMiddleware

func FEngineRepositoryMiddleware(tracer Tracer, repo db.Repository) db.Repository {
	return fengineRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

type fengineRepositoryMiddleware struct {
	tracer Tracer
	repo   db.Repository
}

func (frm fengineRepositoryMiddleware) GetEntity(ctx context.Context, id UUID) (*db.EntityDefinition, error) {
	span := createSpan(ctx, frm.tracer, "GetEntity")
	defer span.Finish()
	return frm.repo.GetEntity(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) UpsertEntity(ctx context.Context, def db.EntityDefinition) (int64, error) {
	span := createSpan(ctx, frm.tracer, "UpsertEntity")
	defer span.Finish()
	return frm.repo.UpsertEntity(ContextWithSpan(ctx, span), def)
}

func (frm fengineRepositoryMiddleware) DeleteEntity(ctx context.Context, thingId UUID) (int64, error) {
	span := createSpan(ctx, frm.tracer, "DeleteEntity")
	defer span.Finish()
	return frm.repo.DeleteEntity(ContextWithSpan(ctx, span), thingId)
}

func (frm fengineRepositoryMiddleware) GetThingAllServices(ctx context.Context, thingId UUID) ([]db.EntityService, error) {
	span := createSpan(ctx, frm.tracer, "GetThingAllServices")
	defer span.Finish()
	return frm.repo.GetThingAllServices(ContextWithSpan(ctx, span), thingId)
}

func (frm fengineRepositoryMiddleware) UpsertThingSubscription(ctx context.Context, sub ...db.ThingSubscription) (int64, error) {
	span := createSpan(ctx, frm.tracer, "UpsertThingSubscription")
	defer span.Finish()
	return frm.repo.UpsertThingSubscription(ContextWithSpan(ctx, span), sub...)
}

func (frm fengineRepositoryMiddleware) DeleteThingSubscription(ctx context.Context, id db.ThingSubscriptionId) (int64, error) {
	span := createSpan(ctx, frm.tracer, "DeleteThingSubscription")
	defer span.Finish()
	return frm.repo.DeleteThingSubscription(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) GetThingAttributes(ctx context.Context, attrs ...string) ([]db.Variable, error) {
	span := createSpan(ctx, frm.tracer, "GetThingAttributes")
	defer span.Finish()
	return frm.repo.GetThingAttributes(ContextWithSpan(ctx, span), attrs...)
}

func (frm fengineRepositoryMiddleware) SetThingAttributes(ctx context.Context, attrs []db.Variable) (int64, error) {
	span := createSpan(ctx, frm.tracer, "SetThingAttributes")
	defer span.Finish()
	return frm.repo.SetThingAttributes(ContextWithSpan(ctx, span), attrs)
}

func (frm fengineRepositoryMiddleware) GetThingService(ctx context.Context, id db.ThingServiceId) (*db.EntityService, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.GetThingService(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) UpsertThingService(ctx context.Context, svc ...db.ThingService) (int, error) {
	span := createSpan(ctx, frm.tracer, "UpsertThingService")
	defer span.Finish()
	return frm.repo.UpsertThingService(ContextWithSpan(ctx, span), svc...)
}

func (frm fengineRepositoryMiddleware) DeleteThingService(ctx context.Context, id db.ThingServiceId) (int, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.DeleteThingService(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) GetThingAllSubscriptions(ctx context.Context, thingId UUID) ([]db.EntitySubscription, error) {
	span := createSpan(ctx, frm.tracer, "GetThingAllSubscriptions")
	defer span.Finish()
	return frm.repo.GetThingAllSubscriptions(ContextWithSpan(ctx, span), thingId)
}

func (frm fengineRepositoryMiddleware) GetThingSubscriptions(ctx context.Context, id db.ThingSubscriptionId) (*db.EntitySubscription, error) {
	span := createSpan(ctx, frm.tracer, "GetThingAllSubscriptions")
	defer span.Finish()
	return frm.repo.GetThingSubscriptions(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) GetAttributeHistory(ctx context.Context, attrs db.AttributeHistoryRequest) ([]db.Variable, error) {
	span := createSpan(ctx, frm.tracer, "GetAttributeHistory")
	defer span.Finish()
	return frm.repo.GetAttributeHistory(ContextWithSpan(ctx, span), attrs)
}

func (frm fengineRepositoryMiddleware) Select(ctx context.Context, sql string, params ...any) (r *db.ResultSet, err error) {
	span := createSpan(ctx, frm.tracer, "Select")
	defer span.Finish()
	return frm.repo.Select(ContextWithSpan(ctx, span), sql, params...)
}

func (frm fengineRepositoryMiddleware) Insert(ctx context.Context, sql string, params []any) (r int64, e error) {
	span := createSpan(ctx, frm.tracer, "Insert")
	defer span.Finish()
	return frm.repo.Insert(ContextWithSpan(ctx, span), sql, params)
}

func (frm fengineRepositoryMiddleware) BatchInsert(ctx context.Context, sql string, fields []string, params [][]any) (r int64, e error) {
	span := createSpan(ctx, frm.tracer, "Insert")
	defer span.Finish()
	return frm.repo.BatchInsert(ContextWithSpan(ctx, span), sql, fields, params)
}

func (frm fengineRepositoryMiddleware) Update(ctx context.Context, sql string, params []any) (r int64, e error) {
	span := createSpan(ctx, frm.tracer, "Update")
	defer span.Finish()
	return frm.repo.Update(ContextWithSpan(ctx, span), sql, params)
}

func (frm fengineRepositoryMiddleware) Delete(ctx context.Context, sql string, params ...any) (r int64, e error) {
	span := createSpan(ctx, frm.tracer, "Update")
	defer span.Finish()
	return frm.repo.Update(ContextWithSpan(ctx, span), sql, params)
}

//#endregion FEngineRepositoryMiddleware

//#region FEngineCacheMiddleware
func FEngineCacheMiddleware(tracer Tracer, cache Cache) Cache {
	return fengineCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

type fengineCacheMiddleware struct {
	tracer Tracer
	cache  Cache
}

func (frm fengineCacheMiddleware) Get(ctx context.Context, id string) (any, error) {
	span := createSpan(ctx, frm.tracer, "Get")
	defer span.Finish()

	ctx = ContextWithSpan(ctx, span)
	return frm.cache.Get(ctx, id)
}

//#endregion FEngineCacheMiddleware

func createSpan(ctx context.Context, tracer Tracer, opName string) Span {
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		return tracer.StartSpan(opName, ChildOf(parentSpan.Context()))
	}
	return tracer.StartSpan(opName)
}
