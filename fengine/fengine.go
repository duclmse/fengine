package fengine

import (
	"context"
)

type Repository interface {
	GetAllThingServices(ctx context.Context, id string) (interface{}, error)
	GetThingService(ctx context.Context, id string) (interface{}, error)
	InsertThingService(ctx context.Context, id string) (interface{}, error)
	UpdateThingService(ctx context.Context, id string) (interface{}, error)
	DeleteThingService(ctx context.Context, id string) (interface{}, error)
}

type Cache interface {
	Get(ctx context.Context, id string) (interface{}, error)
}
