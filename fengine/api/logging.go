package api

import (
	"context"
	"time"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/duclmse/fengine/pkg/logger"
)

var _ fengine.Service = (*loggingMiddleware)(nil)

func LoggingMiddleware(svc fengine.Service, log logger.Logger) fengine.Service {
	return &loggingMiddleware{log: log, svc: svc}
}

type loggingMiddleware struct {
	log logger.Logger
	svc fengine.Service
}

func (l loggingMiddleware) GetEntity(ctx context.Context, id string) (r fengine.Result, e error) {
	defer l.log.Elapse("Get entity id=%s", id)(time.Now(), &e)
	return l.svc.GetEntity(ctx, id)
}

func (l loggingMiddleware) UpsertEntity(ctx context.Context, entityDef sql.EntityDefinition) (r fengine.Result, e error) {
	defer l.log.Elapse("Create entity name=%s", entityDef.Name)(time.Now(), &e)
	return l.svc.UpsertEntity(ctx, entityDef)
}

func (l loggingMiddleware) DeleteEntity(ctx context.Context, id string) (r fengine.Result, e error) {
	defer l.log.Elapse("Delete entity name=%s", id)(time.Now(), &e)
	return l.svc.DeleteEntity(ctx, id)
}

func (l loggingMiddleware) ExecuteService(ctx context.Context, s sql.ServiceRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("Execute service")(time.Now(), &e)
	return l.svc.ExecuteService(ctx, s)
}

func (l loggingMiddleware) CreateTable(ctx context.Context, definition sql.TableDefinition) (r fengine.Result, e error) {
	defer l.log.Elapse("CreateTable name=%s", definition.Name)(time.Now(), &e)
	return l.svc.CreateTable(ctx, definition)
}

func (l loggingMiddleware) Select(ctx context.Context, request sql.SelectRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("Select")(time.Now(), &e)
	return l.svc.Select(ctx, request)
}

func (l loggingMiddleware) Insert(ctx context.Context, request sql.InsertRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("Insert")(time.Now(), &e)
	return l.svc.Insert(ctx, request)
}

func (l loggingMiddleware) Update(ctx context.Context, request sql.UpdateRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("Update")(time.Now(), &e)
	return l.svc.Update(ctx, request)
}

func (l loggingMiddleware) Delete(ctx context.Context, request sql.DeleteRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("Delete")(time.Now(), &e)
	return l.svc.Delete(ctx, request)
}

func (l loggingMiddleware) GetThingAllServices(ctx context.Context, thingId string) (r fengine.Result, e error) {
	defer l.log.Elapse("Get all service of thing (%v)", thingId)(time.Now(), &e)
	return l.svc.GetThingAllServices(ctx, thingId)
}

func (l loggingMiddleware) GetThingService(ctx context.Context, id sql.ThingServiceId) (result fengine.Result, e error) {
	defer l.log.Elapse("Get thing (%v) service %s", id.EntityId, id.Name)(time.Now(), &e)
	return l.svc.GetThingService(ctx, id)
}

func (l loggingMiddleware) Execute(ctx context.Context, script sql.ServiceRequest) (r fengine.Result, e error) {
	defer l.log.Elapse("ExecuteService with %+v", script)(time.Now(), &e)
	return l.svc.ExecuteService(ctx, script)
}
