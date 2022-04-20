package sql

import (
	"context"
	"github.com/duclmse/fengine/viot"
	. "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ Repository = (*fengineRepository)(nil)

type Repository interface {
	GetAllThingServices(ctx context.Context, thingId UUID) (interface{}, error)
	GetThingService(ctx context.Context, thingId UUID, serviceName string) (EntityMethod, error)
	InsertThingService(ctx context.Context, id string) (interface{}, error)
	UpdateThingService(ctx context.Context, id string) (interface{}, error)
	DeleteThingService(ctx context.Context, id string) (interface{}, error)
}

// NewFEngineRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db *sqlx.DB) Repository {
	return &fengineRepository{
		db: NewDatabase(db),
	}
}

type fengineRepository struct {
	db Database
}

func (fer fengineRepository) GetAllThingServices(ctx context.Context, thingId UUID) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingService(ctx context.Context, thingId UUID, serviceName string) (EntityMethod, error) {
	// language=postgresql
	query := `SELECT m1.entity_id AS id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "method" m1 LEFT JOIN "method" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = $1::UUID AND m1.name = $2`
	entities, err := fer.db.QueryxContext(ctx, query, thingId, serviceName)
	if err != nil {
		return EntityMethod{}, err
	}
	defer viot.Close(nil, "db rows")(entities)

	result := EntityMethod{}
	for entities.Next() {
		entity := EntityMethod{}
		if err := entities.StructScan(&entity); err != nil {
			return entity, err
		}
		result = entity
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
