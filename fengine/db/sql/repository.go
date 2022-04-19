package sql

import (
	"context"

	"github.com/duclmse/fengine/fengine"
)

var _ fengine.Repository = (*fengineRepository)(nil)

// NewFEngineRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db Database) fengine.Repository {
	return &fengineRepository{
		db: db,
	}
}

type fengineRepository struct {
	db Database
}

func (fer fengineRepository) GetAllThingServices(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingService(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) InsertThingService(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) UpdateThingService(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) DeleteThingService(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) Get(ctx context.Context, id string) (interface{}, error) {

	return "done!", nil
}
