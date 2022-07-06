package sql_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/duclmse/fengine/viot"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var bg = context.Background()

type TblTest struct {
	id          int    `sql:"id"`
	name        string `sql:"name"`
	description string `sql:"description"`
	a           int
	b           int
	c           int
}

func getPgxConnection(ctx context.Context) *pgx.Conn {
	config, err := pgx.ParseConfig("postgres://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	config.User = "postgres"
	config.Password = "1"

	connection, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Panicf("Unable to connect to database: %v\n", err)
	}
	return connection
}

func getPoolConfig() (config *pgxpool.Config) {
	config, err := pgxpool.ParseConfig("postgres://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	connConfig := config.ConnConfig
	connConfig.User = "postgres"
	connConfig.Password = "1"
	return
}

func TestConnect(t *testing.T) {

	conn := getPgxConnection(bg)
	defer viot.CloseCtx(conn, bg)
	// language=postgresql
	rows, err := conn.Query(bg, `SELECT "id", "name", "base_shapes" FROM entity`)
	if err != nil {
		log.Panicf("QueryRow failed: %v\n", err)
	}
	for rows.Next() {
		var id uuid.UUID
		var name string
		var shapes []uuid.UUID
		if err := rows.Scan(&id, &name, &shapes); err != nil {
			errorHandling(err)
		}

		fmt.Printf("%s - %s - %v\n", id.String(), name, shapes)
	}
}

func TestCreateTable(t *testing.T) {
	conn := getPgxConnection(bg)
	defer viot.CloseCtx(conn, bg)

	// language=postgresql
	tags, err := conn.Exec(bg, `CREATE TABLE a(id INT)`)
	if err != nil {
		log.Panicf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("%v\n", tags)
}

func TestConnectionPool(t *testing.T) {
	bg := context.Background()
	pool, err := pgxpool.ConnectConfig(bg, getPoolConfig())
	if err != nil {
		log.Panicf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	var greeting string
	if err = pool.QueryRow(bg, "select 'Hello, world!'").Scan(&greeting); err != nil {
		t.Fatalf("QueryRow failed: %v\n", err)
	}

	fmt.Println(greeting)
}

func errorHandling(err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Printf("pg err %s: %s\n", pgErr.Code, pgErr.Message) // => 42601: syntax error at end of input
	} else {
		fmt.Printf("err %s\n", err.Error())
	}
}

func TestErrorHandling(t *testing.T) {
	conn := getPgxConnection(bg)
	var greeting string
	if err := conn.QueryRow(context.Background(), "select 1 +").Scan(&greeting); err != nil {
		errorHandling(err)
	}
}

func TestQueryFunc(t *testing.T) {
	conn := getPgxConnection(bg)
	defer viot.CloseCtx(conn, bg)

	var n int
	// language=postgresql
	tags, err := conn.QueryFunc(bg, `SELECT * FROM tbl_test WHERE a > $1`, []any{0}, []any{&n}, func(row pgx.QueryFuncRow) error {
		row.RawValues()
		return nil
	})
	if err != nil {
		log.Panicf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("%v\n", tags)
}
