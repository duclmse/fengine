package fengine

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id string) (interface{}, error)
}

type Cache interface {
	Get(ctx context.Context, id string) (interface{}, error)
}
