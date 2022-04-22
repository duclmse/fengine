package fengine

import (
	. "context"
	. "github.com/google/uuid"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
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
	DB          *sqlx.DB
	Log         logger.Logger
	ExeClient   pb.FEngineExecutorClient
}

type Service interface {
	GetThingAllServices(Context, UUID) (*Result, error)
	GetThingService(Context, UUID, string) (*Result, error)
	ExecuteService(Context, *JsonScript) (*Result, error)

	Select(Context, *JsonSelectRequest) (*Result, error)
	Insert(Context, *JsonInsertRequest) (*Result, error)
	Update(Context, *JsonUpdateRequest) (*Result, error)
	Delete(Context, *JsonDeleteRequest) (*Result, error)
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

func (s FengineService) GetThingAllServices(ctx Context, thingId UUID) (*Result, error) {
	services, err := s.Repository.GetThingAllServices(ctx, thingId)
	if err != nil {
		return nil, err
	}
	return &Result{Value: services}, nil
}

func (s FengineService) GetThingService(ctx Context, thingId UUID, serviceName string) (*Result, error) {
	services, err := s.Repository.GetThingService(ctx, thingId, serviceName)
	if err != nil {
		return nil, err
	}
	return &Result{Value: services}, nil
}

func (s FengineService) ExecuteService(ctx Context, script *JsonScript) (*Result, error) {
	return nil, nil
}

func (s FengineService) Select(ctx Context, req *JsonSelectRequest) (*Result, error) {
	return nil, nil
}

func (s FengineService) Insert(ctx Context, req *JsonInsertRequest) (*Result, error) {
	return nil, nil
}

func (s FengineService) Update(ctx Context, req *JsonUpdateRequest) (*Result, error) {
	return nil, nil
}

func (s FengineService) Delete(ctx Context, req *JsonDeleteRequest) (*Result, error) {
	return nil, nil
}
