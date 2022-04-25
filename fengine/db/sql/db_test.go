package sql_test

import (
	"context"
	"fmt"
	. "github.com/duclmse/fengine/fengine/db/sql"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
)

func connect(t *testing.T) (logger.Logger, *sqlx.DB) {
	config := Config{Host: "localhost", Port: "5432", User: "postgres", Pass: "1", Name: "postgres", SSLMode: "disable"}
	log, err := logger.New(os.Stdout, "debug")
	if err != nil {
		t.Fatalf("cannot create logger %s", err)
		return nil, nil
	}
	db, err := Connect(config, log)
	if err != nil {
		log.Fatalf("cannot connect %s", err)
	}
	return log, db
}

func TestGenUUID(t *testing.T) {
	if id, err := uuid.NewRandom(); err != nil {
		fmt.Printf("err")
	} else {
		fmt.Printf("uuid:= %s\n", id)
	}
}

func TestUuid(t *testing.T) {
	log, db := connect(t)
	rows, err := db.QueryxContext(context.Background(), `SELECT * FROM entity WHERE id = $1::UUID`, `21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab`)
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	defer viot.Close(log, "")(rows)

	things := make([]Entity, 0)
	for rows.Next() {
		e := &Entity{}
		if err := rows.StructScan(e); err != nil {
			t.Errorf("err %s", err.Error())
		}
		things = append(things, *e)
	}
	t.Logf("result len %d", len(things))
	log.Struct(things)
}

func TestFengineRepository_GetThingService(t *testing.T) {
	log, db := connect(t)
	repository := NewFEngineRepository(db, log)
	service, err := repository.GetThingService(context.Background(), ThingServiceId{
		EntityId: uuid.MustParse("1d6d5123-3fb8-4ab1-956f-c6f96847471d"),
		Name:     "templ_method",
	})
	if err != nil {
		t.Errorf("err %s\n", err.Error())
		return
	}
	log.Struct(service)
}

func TestFengineRepository_GetAllThingService(t *testing.T) {
	log, db := connect(t)
	repository := NewFEngineRepository(db, log)
	service, err := repository.GetThingAllServices(context.Background(), uuid.MustParse("21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab"))
	if err != nil {
		t.Errorf("err %s\n", err.Error())
		return
	}
	fmt.Printf("%+v\n", service)
}

func TestJsonString_Scan(t *testing.T) {
}
