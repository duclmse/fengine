package api

import (
	"context"
	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/pkg/logger"
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

func (l loggingMiddleware) CreateEntity(ctx context.Context, entityDef *fengine.EntityDefinition) (r *fengine.Result, e error) {
	defer l.log.Elapse("Execute service")(time.Now(), &e)
	return l.svc.CreateEntity(ctx, entityDef)
}

func (l loggingMiddleware) ExecuteService(ctx context.Context, s *fengine.ServiceRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("Execute service")(time.Now(), &e)
	return l.svc.ExecuteService(ctx, s)
}

func (l loggingMiddleware) CreateTable(ctx context.Context, definition *fengine.TableDefinition) (r *fengine.Result, e error) {
	defer l.log.Elapse("CreateTable")(time.Now(), &e)
	return l.svc.CreateTable(ctx, definition)
}

func (l loggingMiddleware) Select(ctx context.Context, request *fengine.SelectRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("Select")(time.Now(), &e)
	return l.svc.Select(ctx, request)
}

func (l loggingMiddleware) Insert(ctx context.Context, request *fengine.InsertRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("Insert")(time.Now(), &e)
	return l.svc.Insert(ctx, request)
}

func (l loggingMiddleware) Update(ctx context.Context, request *fengine.UpdateRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("Update")(time.Now(), &e)
	return l.svc.Update(ctx, request)
}

func (l loggingMiddleware) Delete(ctx context.Context, request *fengine.DeleteRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("Delete")(time.Now(), &e)
	return l.svc.Delete(ctx, request)
}

func (l loggingMiddleware) GetThingAllServices(ctx context.Context, thingId string) (r *fengine.Result, e error) {
	defer l.log.Elapse("Get all service of thing (%v)", thingId)(time.Now(), &e)
	return l.svc.GetThingAllServices(ctx, thingId)
}

func (l loggingMiddleware) GetThingService(ctx context.Context, thingId string, service string) (result *fengine.Result, e error) {
	defer l.log.Elapse("Get thing (%v) service %s", thingId, service)(time.Now(), &e)
	return l.svc.GetThingService(ctx, thingId, service)
}

func (l loggingMiddleware) Execute(ctx context.Context, script *fengine.ServiceRequest) (r *fengine.Result, e error) {
	defer l.log.Elapse("ExecuteService with %+v", script)(time.Now(), &e)
	return l.svc.ExecuteService(ctx, script)
}
