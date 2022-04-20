package tracing

import (
	"context"

	. "github.com/google/uuid"
	. "github.com/opentracing/opentracing-go"

	. "github.com/duclmse/fengine/fengine/db/cache"
	"github.com/duclmse/fengine/fengine/db/sql"
)

var (
	_ sql.Repository = (*fengineRepositoryMiddleware)(nil)
	_ Cache          = (*fengineCacheMiddleware)(nil)
)

//#region FEngineRepositoryMiddleware

func FEngineRepositoryMiddleware(tracer Tracer, repo sql.Repository) sql.Repository {
	return fengineRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

type fengineRepositoryMiddleware struct {
	tracer Tracer
	repo   sql.Repository
}

func (frm fengineRepositoryMiddleware) GetAllThingServices(ctx context.Context, thingId UUID) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "GetAllThingServices")
	defer span.Finish()
	return frm.repo.GetAllThingServices(ContextWithSpan(ctx, span), thingId)
}

func (frm fengineRepositoryMiddleware) GetThingService(ctx context.Context, thingId UUID, service string) (sql.EntityMethod, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.GetThingService(ContextWithSpan(ctx, span), thingId, service)
}

func (frm fengineRepositoryMiddleware) InsertThingService(ctx context.Context, id string) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.InsertThingService(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) UpdateThingService(ctx context.Context, id string) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.InsertThingService(ContextWithSpan(ctx, span), id)
}

func (frm fengineRepositoryMiddleware) DeleteThingService(ctx context.Context, id string) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "GetThingService")
	defer span.Finish()
	return frm.repo.InsertThingService(ContextWithSpan(ctx, span), id)
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

func (frm fengineCacheMiddleware) Get(ctx context.Context, id string) (interface{}, error) {
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
