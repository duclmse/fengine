package fengine

import (
	"context"
)

type Cache interface {
	Get(ctx context.Context, id string) (interface{}, error)
}
