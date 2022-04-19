package sql

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
)

func TestGenUUID(t *testing.T) {
	if id, err := uuid.NewRandom(); err != nil {
		fmt.Printf("err")
	} else {
		fmt.Printf("uuid:= %s\n", id)
	}
}

func TestUuid(t *testing.T) {
	config := Config{Host: "localhost", Port: "5432", User: "postgres", Pass: "1", Name: "postgres", SSLMode: "disable"}
	log, err := logger.New(os.Stdout, "debug")
	if err != nil {
		log.Error("cannot create logger %s", err)
		return
	}
	db, err := Connect(config, log)
	if err != nil {
		log.Error("cannot connect %s", err)
		return
	}
	rows, err := db.QueryxContext(context.Background(), `SELECT * FROM entity WHERE id = $1::UUID`, `21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab`)
	if err != nil {
		log.Error("err %s", err.Error())
		return
	}
	defer viot.Close(log, "")(rows)

	things := make([]Entity, 0)
	for rows.Next() {
		t := &Entity{}
		if err := rows.StructScan(t); err != nil {
			log.Error("err %s", err.Error())
			return
		}
		things = append(things, *t)
	}
	log.Info("%d", len(things))
	log.Struct(things)
}

func TestFengineRepository_GetThingService(t *testing.T) {

}
