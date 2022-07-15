package fengine

import (
	ctx "context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"

	"github.com/duclmse/fengine/fengine/db/cache"
	"github.com/duclmse/fengine/fengine/db/sql"
	pb "github.com/duclmse/fengine/pb"
	"github.com/duclmse/fengine/pkg/logger"
)

var _ Service = (*FengineService)(nil)

type ServiceComponent struct {
	Tracer      opentracing.Tracer
	Cache       *redis.Client
	CacheTracer opentracing.Tracer
	DB          *pgxpool.Pool
	Log         logger.Logger
	ExeClient   pb.FEngineExecutorClient
}

type Service interface {
	GetEntity(ctx ctx.Context, id string) (r Result, e error)
	UpsertEntity(ctx ctx.Context, entityDef sql.EntityDefinition) (r Result, e error)
	DeleteEntity(ctx ctx.Context, id string) (r Result, e error)

	GetThingAllServices(ctx ctx.Context, thing uuid.UUID) (r Result, e error)
	GetThingService(ctx ctx.Context, req sql.ThingServiceId) (r Result, e error)
	ExecuteService(ctx ctx.Context, req sql.ServiceRequest) (r Result, e error)

	CreateTable(ctx ctx.Context, req sql.TableDefinition) (r Result, e error)
	Select(ctx ctx.Context, req sql.SelectRequest) (res *sql.ResultSet, err error)
	Insert(ctx ctx.Context, req sql.InsertRequest) (r Result, e error)
	InsertBatch(ctx ctx.Context, req sql.BatchInsertRequest) (Result, error)
	Update(ctx ctx.Context, req sql.UpdateRequest) (r Result, e error)
	Delete(ctx ctx.Context, req sql.DeleteRequest) (r Result, e error)
}

func (s FengineService) New() Service {
	return &s
}

type FengineService struct {
	Repository sql.Repository
	Cache      cache.Cache
	ExecClient pb.FEngineExecutorClient
	Log        logger.Logger
}

func (s FengineService) GetEntity(ctx ctx.Context, id string) (r Result, e error) {
	uid, e := uuid.Parse(id)
	if e != nil {
		return result(e)
	}

	entity, e := s.Repository.GetEntity(ctx, uid)
	if e != nil {
		return result(e)
	}

	return Result{Data: entity}, e
}

func (s FengineService) UpsertEntity(ctx ctx.Context, entityDef sql.EntityDefinition) (r Result, err error) {
	upserted, err := s.Repository.UpsertEntity(ctx, entityDef)
	if err != nil {
		return result(err)
	}
	return Result{Msg: fmt.Sprintf("upserted %d", upserted)}, nil
}

func (s FengineService) DeleteEntity(ctx ctx.Context, id string) (r Result, e error) {
	uid, e := uuid.Parse(id)
	if e != nil {
		return result(e)
	}

	entity, e := s.Repository.DeleteEntity(ctx, uid)
	if e != nil {
		return result(e)
	}

	return Result{Data: entity}, e
}

func (s FengineService) GetThingAllServices(ctx ctx.Context, thingId uuid.UUID) (r Result, err error) {
	services, err := s.Repository.GetThingAllServices(ctx, thingId)
	if err != nil {
		return result(err)
	}
	return Result{Data: services}, nil
}

func (s FengineService) GetThingService(ctx ctx.Context, id sql.ThingServiceId) (Result, error) {
	services, err := s.Repository.GetThingService(ctx, id)
	if err != nil {
		return Result{Code: 1, Msg: err.Error()}, err
	}
	return Result{Data: services}, nil
}

func (s FengineService) ExecuteService(ctx ctx.Context, script sql.ServiceRequest) (Result, error) {
	return Result{}, nil
}

func (s FengineService) CreateTable(ctx ctx.Context, table sql.TableDefinition) (Result, error) {
	return Result{}, nil
}

func (s FengineService) Select(ctx ctx.Context, req sql.SelectRequest) (res *sql.ResultSet, err error) {
	_sql, err := req.ToSQL()
	if err != nil {
		return
	}

	return s.Repository.Select(ctx, _sql)
}

func (s FengineService) InsertBatch(ctx ctx.Context, req sql.InsertRequest) (Result, error) {
	rs, err := s.Repository.InsertBatch(ctx, req)
	if err != nil {
		return result(err)
	}

	return Result{Data: rs}, nil
}

func (s FengineService) Insert(ctx ctx.Context, req sql.InsertRequest) (Result, error) {
	_sql, params, err := req.ToSQL()
	if err != nil {
		return result(err)
	}

	rs, err := s.Repository.Insert(ctx, _sql)
	if err != nil {
		return result(err)
	}

	return Result{Data: rs}, nil
}

func (s FengineService) Update(ctx ctx.Context, req sql.UpdateRequest) (Result, error) {
	_sql, params, err := req.ToSQL()
	if err != nil {
		return result(err)
	}

	rs, err := s.Repository.Update(ctx, _sql, params)
	if err != nil {
		return result(err)
	}

	return Result{Data: rs}, nil
}

func (s FengineService) Delete(ctx ctx.Context, req sql.DeleteRequest) (Result, error) {
	_sql, params, err := req.ToSQL()
	if err != nil {
		return result(err)
	}

	rs, err := s.Repository.Update(ctx, _sql, params)
	if err != nil {
		return result(err)
	}

	return Result{Data: rs}, nil
}

func result(e error) (Result, error) {
	return Result{Code: 1, Msg: e.Error()}, e
}
