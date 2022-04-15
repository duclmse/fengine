package fengine

import (
	"context"

	"github.com/duclmse/fengine/viot"
)

var _ Service = (*fengineService)(nil)

type fengineService struct {
	repository Repository
	cache      Cache
	idp        viot.UUIDProvider
}

type Service interface {
	Get(ctx context.Context, id string) (interface{}, error)
}

func New(idp viot.UUIDProvider, repository Repository, cache Cache) Service {
	return &fengineService{
		repository: repository,
		cache:      cache,
		idp:        idp,
	}
}

func (s fengineService) Get(ctx context.Context, id string) (interface{}, error) {
	return nil, nil
}
