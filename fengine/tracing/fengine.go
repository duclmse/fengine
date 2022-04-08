package tracing

import (
	"context"

	"github.com/duclmse/fengine/fengine"
	"github.com/opentracing/opentracing-go"
)

var (
	_ fengine.Repository = (*fengineRepositoryMiddleware)(nil)
	_ fengine.Cache      = (*fengineCacheMiddleware)(nil)
)

type fengineRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   fengine.Repository
}

type fengineCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  fengine.Cache
}

func FEngineRepositoryMiddleware(tracer opentracing.Tracer, repo fengine.Repository) fengine.Repository {
	return fengineRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func FEngineCacheMiddleware(tracer opentracing.Tracer, cache fengine.Cache) fengine.Cache {
	return fengineCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (frm fengineRepositoryMiddleware) Get(ctx context.Context, id string) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "Get")
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	return frm.repo.Get(ctx, id)
}

func (frm fengineCacheMiddleware) Get(ctx context.Context, id string) (interface{}, error) {
	span := createSpan(ctx, frm.tracer, "Get")
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	return frm.cache.Get(ctx, id)
}

func createSpan(ctx context.Context, tracer opentracing.Tracer, opName string) opentracing.Span {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return tracer.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	}
	return tracer.StartSpan(opName)
}
