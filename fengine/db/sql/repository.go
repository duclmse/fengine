package sql

import (
	. "context"
	"database/sql"
	"fmt"

	. "github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
)

var _ Repository = (*fengineRepository)(nil)

type Repository interface {
	UpsertEntity(ctx Context, def EntityDefinition) (int, error)
	DeleteEntity(ctx Context, thingId UUID) (int, error)

	GetThingAllServices(ctx Context, thingId UUID) ([]EntityService, error)
	GetThingService(ctx Context, id ThingServiceId) (*EntityService, error)
	UpsertThingService(ctx Context, service ThingService) (int, error)
	DeleteThingService(ctx Context, id ThingServiceId) (int, error)

	GetThingAllSubscriptions(ctx Context, thingId UUID) ([]EntitySubscription, error)
	GetThingSubscriptions(ctx Context, id ThingSubscriptionId) (*EntitySubscription, error)
	UpsertThingSubscription(ctx Context, sub ThingSubscription) (int, error)
	DeleteThingSubscription(ctx Context, id ThingSubscriptionId) (int, error)

	GetThingAttributes(ctx Context, attrs ...string) ([]Variable, error)
	SetThingAttributes(ctx Context, attrs []Variable) (int, error)
	GetAttributeHistory(cts Context, attrs AttributeHistoryRequest)

	Select(ctx Context, sql string) (r map[string]Variable, e error)
	Insert(ctx Context, sql string) (r *sql.Result, e error)
	Update(ctx Context, sql string) (r *sql.Result, e error)
	Delete(ctx Context, sql string) (r *sql.Result, e error)
}

// NewFEngineRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db *sqlx.DB, log logger.Logger) Repository {
	return &fengineRepository{
		db:  NewDatabase(db),
		log: log,
	}
}

type fengineRepository struct {
	db  Database
	log logger.Logger
}

func (fer fengineRepository) GetEntity(ctx Context, thingId UUID) (EntityDefinition, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) UpsertEntity(ctx Context, def EntityDefinition) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) DeleteEntity(ctx Context, thingId UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingAllServices(ctx Context, thingId UUID) ([]EntityService, error) {
	// language=postgresql
	query := `SELECT m1.entity_id AS id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "service" m1 LEFT JOIN "service" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = $1::UUID`
	entities, err := fer.db.QueryxContext(ctx, query, thingId)
	if err != nil {
		return nil, err
	}
	defer viot.Close(nil, "db rows")(entities)

	result := []EntityService{}
	for entities.Next() {
		entity := EntityService{}
		if err := entities.StructScan(&entity); err != nil {
			return nil, err
		}
		result = append(result, entity)
	}
	return result, nil
}

func (fer fengineRepository) GetThingService(ctx Context, id ThingServiceId) (*EntityService, error) {
	// language=postgresql
	query := `SELECT m1.entity_id AS id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "service" m1 LEFT JOIN "service" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = $1::UUID AND m1.name = $2`
	entities, err := fer.db.QueryxContext(ctx, query, id.EntityId, id.Name)
	if err != nil {
		fmt.Printf("err selecting %s", err.Error())
		return nil, err
	}
	defer viot.Close(nil, "db rows")(entities)

	result := new(EntityService)
	count := 0
	for entities.Next() {
		entity := new(EntityService)
		if err := entities.StructScan(entity); err != nil {
			fer.log.Debug("err selecting %s", err.Error())
			return entity, err
		}
		result = entity
		count += 1
	}
	fmt.Printf("selected %d\n", count)
	return result, nil
}

func (fer fengineRepository) UpsertThingService(ctx Context, service ThingService) (int, error) {
	// language=postgresql
	query := `INSERT INTO service("entity_id", "name", "input", "output", "code") VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO UPDATE SET "input" = $3, "output" = $4, "code" = $5`
	result, err := fer.db.ExecContext(ctx, query, service.EntityId, service.Name, service.Input, service.Output, service.Code)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (fer fengineRepository) DeleteThingService(ctx Context, id ThingServiceId) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingAllSubscriptions(ctx Context, thingId UUID) ([]EntitySubscription, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingSubscriptions(ctx Context, id ThingSubscriptionId) (*EntitySubscription, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) UpsertThingSubscription(ctx Context, sub ThingSubscription) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) DeleteThingSubscription(ctx Context, id ThingSubscriptionId) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetThingAttributes(ctx Context, attrs ...string) ([]Variable, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) SetThingAttributes(ctx Context, attrs []Variable) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetAttributeHistory(cts Context, attrs AttributeHistoryRequest) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) Select(ctx Context, sql string) (r map[string]Variable, e error) {
	return nil, nil
}

func (fer fengineRepository) Insert(ctx Context, sql string) (r *sql.Result, e error) {
	return nil, nil
}

func (fer fengineRepository) Update(ctx Context, sql string) (r *sql.Result, e error) {
	return nil, nil
}

func (fer fengineRepository) Delete(ctx Context, sql string) (r *sql.Result, e error) {
	return nil, nil
}
