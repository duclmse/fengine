package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	//_ fengine.Service = (*eventStore)(nil)
	_ Cache = (*fengineCache)(nil)
)

type fengineCache struct {
	client *redis.Client
}

type eventStore struct {
	client *redis.Client
}

type Cache interface {
	Get(ctx context.Context, id string) (interface{}, error)
}

func NewFEngineCache(client *redis.Client) Cache {
	return &fengineCache{
		client: client,
	}
}

func (fec fengineCache) Get(ctx context.Context, id string) (interface{}, error) {
	return fec.client.Get(ctx, id), nil
}
