package fengine

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/duclmse/fengine/fengine/db/cache"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/duclmse/fengine/pb"
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
	Execute(ctx context.Context, s *pb.Script) (*pb.Result, error)
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

func (s FengineService) Execute(ctx context.Context, script *pb.Script) (*pb.Result, error) {
	//TODO implement me
	panic("implement me")
}
