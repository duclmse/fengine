package sql

import (
	"context"
	"github.com/duclmse/fengine/viot"
	"github.com/google/uuid"
)

var _ Repository = (*fengineRepository)(nil)

type Repository interface {
	GetAllThingServices(ctx context.Context, id string) (interface{}, error)
	GetThingService(ctx context.Context, thingId uuid.UUID, serviceName string) ([]Entity, error)
	InsertThingService(ctx context.Context, id string) (interface{}, error)
	UpdateThingService(ctx context.Context, id string) (interface{}, error)
	DeleteThingService(ctx context.Context, id string) (interface{}, error)
}

// NewFEngineRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db Database) Repository {
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

func (fer fengineRepository) GetThingService(ctx context.Context, thingId uuid.UUID, serviceName string) ([]Entity, error) {
	params := map[string]interface{}{
		"id":   thingId,
		"name": serviceName,
	}
	// language=postgresql
	entities, err := fer.db.NamedQueryContext(ctx,
		`SELECT m1.entity_id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "method" m1 LEFT JOIN "method" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = :id::UUID AND m1.name = :name`, params)
	if err != nil {
		return nil, err
	}
	defer viot.Close(nil, "db rows")(entities)

	result := []Entity{}
	for entities.Next() {
		entity := Entity{}
		if err := entities.StructScan(&entity); err != nil {
			return result, err
		}
		result = append(result, entity)
	}
	return result, nil
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
