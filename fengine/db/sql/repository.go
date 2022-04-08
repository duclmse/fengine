package sql

import (
	"context"

	"github.com/duclmse/fengine/fengine"
)

var _ fengine.Repository = (*fengineRepository)(nil)

type fengineRepository struct {
	db Database
}

// NewPricingRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db Database) fengine.Repository {
	return &fengineRepository{
		db: db,
	}
}

func (fer fengineRepository) Get(ctx context.Context, id string) (interface{}, error) {

	return "done!", nil
}
