package fengine

import (
	. "context"
	"errors"
	"fmt"
	. "github.com/google/uuid"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/duclmse/fengine/fengine/db/cache"
	. "github.com/duclmse/fengine/fengine/db/sql"
	pb "github.com/duclmse/fengine/pb"
	"github.com/duclmse/fengine/pkg/logger"
)

var _ Service = (*FengineService)(nil)

type ServiceComponent struct {
	Tracer      opentracing.Tracer
	Cache       *redis.Client
	CacheTracer opentracing.Tracer
	DB          *sqlx.DB
	Log         logger.Logger
	ExeClient   pb.FEngineExecutorClient
}

type Service interface {
	GetEntity(ctx Context, id string) (r Result, e error)
	UpsertEntity(ctx Context, entityDef EntityDefinition) (r Result, e error)
	DeleteEntity(ctx Context, id string) (r Result, e error)

	GetThingAllServices(ctx Context, thing string) (r Result, e error)
	GetThingService(ctx Context, req ThingServiceId) (r Result, e error)
	ExecuteService(ctx Context, req ServiceRequest) (r Result, e error)

	CreateTable(ctx Context, req TableDefinition) (r Result, e error)
	Select(ctx Context, req SelectRequest) (r Result, e error)
	Insert(ctx Context, req InsertRequest) (r Result, e error)
	Update(ctx Context, req UpdateRequest) (r Result, e error)
	Delete(ctx Context, req DeleteRequest) (r Result, e error)
}

func (s FengineService) New() Service {
	return &s
}

type FengineService struct {
	Repository Repository
	Cache      cache.Cache
	ExecClient pb.FEngineExecutorClient
	Log        logger.Logger
}

func (s FengineService) GetEntity(ctx Context, id string) (r Result, e error) {
	uid, e := Parse(id)
	if e != nil {
		return Result{Code: 1, Msg: e.Error()}, e
	}

	entity, e := s.Repository.GetEntity(ctx, uid)
	if e != nil {
		return Result{Code: 1, Msg: e.Error()}, e
	}

	return Result{Data: entity}, e
}
func (s FengineService) UpsertEntity(ctx Context, entityDef EntityDefinition) (r Result, e error) {
	upserted, err := s.Repository.UpsertEntity(ctx, entityDef)
	if err != nil {
		return Result{Code: 1, Msg: err.Error()}, err
	}
	return Result{Msg: fmt.Sprintf("upserted %d", upserted)}, err
}

func (s FengineService) DeleteEntity(ctx Context, id string) (r Result, e error) {
	uid, e := Parse(id)
	if e != nil {
		return Result{Code: 1, Msg: e.Error()}, e
	}

	entity, e := s.Repository.DeleteEntity(ctx, uid)
	if e != nil {
		return Result{Code: 1, Msg: e.Error()}, e
	}

	return Result{Data: entity}, e
}

func (s FengineService) GetThingAllServices(ctx Context, thingId string) (Result, error) {
	id, err := Parse(thingId)
	if err != nil {
		return Result{Code: 1}, errors.New("thing id is not a valid uuid")
	}
	services, err := s.Repository.GetThingAllServices(ctx, id)
	if err != nil {
		return Result{Code: 1, Msg: err.Error()}, err
	}
	return Result{Data: services}, nil
}

func (s FengineService) GetThingService(ctx Context, id ThingServiceId) (Result, error) {
	services, err := s.Repository.GetThingService(ctx, id)
	if err != nil {
		return Result{Code: 1, Msg: err.Error()}, err
	}
	return Result{Data: services}, nil
}

func (s FengineService) ExecuteService(ctx Context, script ServiceRequest) (Result, error) {
	return Result{}, nil
}

func (s FengineService) CreateTable(ctx Context, table TableDefinition) (Result, error) {
	return Result{}, nil
}

func (s FengineService) Select(ctx Context, req SelectRequest) (Result, error) {
	return Result{}, nil
}

func (s FengineService) Insert(ctx Context, req InsertRequest) (Result, error) {
	return Result{}, nil
}

func (s FengineService) Update(ctx Context, req UpdateRequest) (Result, error) {
	return Result{}, nil
}

func (s FengineService) Delete(ctx Context, req DeleteRequest) (Result, error) {
	return Result{}, nil
}
