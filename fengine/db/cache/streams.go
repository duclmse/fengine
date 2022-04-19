package cache

import (
	"context"

	"github.com/duclmse/fengine/fengine"
	"github.com/go-redis/redis/v8"
)

var (
	_ fengine.Service = (*eventStore)(nil)
	_ fengine.Cache   = (*fengineCache)(nil)
)

type fengineCache struct {
	client *redis.Client
}

type eventStore struct {
	svc    fengine.Service
	client *redis.Client
}

func NewFEngineCache(client *redis.Client) fengine.Cache {
	return &fengineCache{
		client: client,
	}
}

func NewEventStoreMiddleware(svc fengine.Service, client *redis.Client) fengine.Service {
	return &eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Get(ctx context.Context, id string) (interface{}, error) {
	return es.svc.Get(ctx, id)
}

func (fec fengineCache) Get(ctx context.Context, id string) (interface{}, error) {
	return fec.client.Get(ctx, id), nil
}
