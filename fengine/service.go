package fengine

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	. "github.com/duclmse/fengine/pb"
	"github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
)

var _ Service = (*fengineService)(nil)

type fengineService struct {
	repository Repository
	cache      Cache
	idp        viot.UUIDProvider
	execClient FEngineExecutorClient
	log        logger.Logger
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

func New(idp viot.UUIDProvider, repository Repository, cache Cache, client FEngineExecutorClient, log logger.Logger) Service {
	return &fengineService{
		repository: repository,
		cache:      cache,
		idp:        idp,
		execClient: client,
		log:        log,
	}
}

func (s fengineService) Get(ctx context.Context, id string) (interface{}, error) {
	s.log.Info("service.go Get")
	return nil, nil
}
