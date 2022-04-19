package fengine

import (
	"context"
	"github.com/duclmse/fengine/fengine/db/sql"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	. "github.com/duclmse/fengine/pb"
	"github.com/duclmse/fengine/pkg/logger"
)

var _ Service = (*FengineService)(nil)

type FengineService struct {
	Repository sql.Repository
	Cache      Cache
	ExecClient FEngineExecutorClient
	Log        logger.Logger
}

type Service interface {
	Get(ctx context.Context, id string) (interface{}, error)
}

type ServiceComponent struct {
	Tracer      opentracing.Tracer
	Cache       *redis.Client
	CacheTracer opentracing.Tracer
	DB          *sqlx.DB
	Log         logger.Logger
	ExeClient   FEngineExecutorClient
}

func (s FengineService) New() Service {
	return &s
}

func (s FengineService) Get(ctx context.Context, id string) (interface{}, error) {
	s.Log.Info("service.go Get")
	return nil, nil
}
