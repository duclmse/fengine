package sql

import (
	"fmt"
	"os"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
)

func TestGenUUID(t *testing.T) {
	if id, err := uuid.NewV4(); err != nil {
		fmt.Printf("err")
	} else {
		fmt.Printf("uuid:= %s\n", id)
	}
}

func TestUuid(t *testing.T) {
	type Thing struct {
		ID          uuid.UUID `json:"id" sql:",type:uuid"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}
	things := make([]Thing, 0)

	config := Config{Host: "localhost", Port: "5432", User: "postgres", Pass: "1", Name: "postgres", SSLMode: "disable"}
	log, err := logger.New(os.Stdout, "debug")
	if err != nil {
		return
	}
	db, err := Connect(config, log)
	if err != nil {
		log.Error("cannot connect %s", err)
		return
	}
	rows, err := db.Query(`SELECT id, name FROM entity`)
	if err != nil {
		log.Error("err %s", err.Error())
		return
	}
	defer viot.Close(log, "")(rows)

	for rows.Next() {
		t := &Thing{}
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return
		}
		things = append(things, *t)
	}
	log.Struct(things)
}
