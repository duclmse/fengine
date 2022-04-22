package api

import (
	"context"
	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/pkg/logger"
	"github.com/google/uuid"
	"time"
)

var _ fengine.Service = (*loggingMiddleware)(nil)

func LoggingMiddleware(svc fengine.Service, log logger.Logger) fengine.Service {
	return &loggingMiddleware{log: log, svc: svc}
}

type loggingMiddleware struct {
	log logger.Logger
	svc fengine.Service
}

func (l loggingMiddleware) ExecuteService(ctx context.Context, s *fengine.JsonScript) (res *fengine.Result, err error) {
	defer l.log.Elapse("Execute service")(time.Now(), &err)
	return l.svc.ExecuteService(ctx, s)
}

func (l loggingMiddleware) Select(ctx context.Context, request *fengine.JsonSelectRequest) (res *fengine.Result, err error) {
	defer l.log.Elapse("Select")(time.Now(), &err)
	return l.svc.Select(ctx, request)
}

func (l loggingMiddleware) Insert(ctx context.Context, request *fengine.JsonInsertRequest) (res *fengine.Result, err error) {
	defer l.log.Elapse("Insert")(time.Now(), &err)
	return l.svc.Insert(ctx, request)
}

func (l loggingMiddleware) Update(ctx context.Context, request *fengine.JsonUpdateRequest) (res *fengine.Result, err error) {
	defer l.log.Elapse("Update")(time.Now(), &err)
	return l.svc.Update(ctx, request)
}

func (l loggingMiddleware) Delete(ctx context.Context, request *fengine.JsonDeleteRequest) (res *fengine.Result, err error) {
	defer l.log.Elapse("Delete")(time.Now(), &err)
	return l.svc.Delete(ctx, request)
}

func (l loggingMiddleware) GetThingAllServices(ctx context.Context, thingId uuid.UUID) (result *fengine.Result, err error) {
	defer l.log.Elapse("Get all service of thing (%v)", thingId)(time.Now(), &err)
	return l.svc.GetThingAllServices(ctx, thingId)
}

func (l loggingMiddleware) GetThingService(ctx context.Context, thingId uuid.UUID, service string) (result *fengine.Result, err error) {
	defer l.log.Elapse("Get thing (%v) service %s", thingId, service)(time.Now(), &err)
	return l.svc.GetThingService(ctx, thingId, service)
}

func (l loggingMiddleware) Execute(ctx context.Context, script *fengine.JsonScript) (result *fengine.Result, err error) {
	defer l.log.Elapse("ExecuteService with %+v", script)(time.Now(), &err)
	return l.svc.ExecuteService(ctx, script)
}
