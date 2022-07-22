package sql_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/duclmse/fengine/viot"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	l "log"
	"os"
	"testing"

	. "github.com/duclmse/fengine/fengine/db/sql"
	"github.com/goccy/go-json"

	"github.com/google/uuid"

	"github.com/duclmse/fengine/pkg/logger"
)

func connect(t *testing.T) (logger.Logger, *pgxpool.Pool) {
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
		fmt.Printf("err %s\n", err)
	} else {
		fmt.Printf("uuid = %s\n", id)
	}
}

func TestUuid(t *testing.T) {
	log, db := connect(t)
	rows, err := db.Query(context.Background(), `SELECT * FROM entity WHERE id = $1::UUID`, `21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab`)
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	defer rows.Close()

	things := make([]Entity, 0)
	for rows.Next() {
		e := &Entity{}
		//FIXME
		//if err := rows.StructScan(e); err != nil {
		//	t.Errorf("err %s", err.Error())
		//}
		things = append(things, *e)
	}
	t.Logf("result len %d", len(things))
	log.Struct(things)
}

//func TestFengineRepository_GetThingService(t *testing.T) {
//	log, db := connect(t)
//	repository := NewFEngineRepository(db, log)
//	service, err := repository.GetThingService(context.Background(), ThingServiceId{
//		EntityId: uuid.MustParse("1d6d5123-3fb8-4ab1-956f-c6f96847471d"),
//		Name:     "templ_method",
//	})
//	if err != nil {
//		t.Errorf("err %s\n", err.Error())
//		return
//	}
//	log.Struct(service)
//}

//func TestFengineRepository_GetAllThingService(t *testing.T) {
//	log, db := connect(t)
//	repository := NewFEngineRepository(db, log)
//	service, err := repository.GetThingAllServices(context.Background(), uuid.MustParse("21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab"))
//	if err != nil {
//		t.Errorf("err %s\n", err.Error())
//		return
//	}
//	fmt.Printf("%+v\n", service)
//}

func TestPanic(t *testing.T) {
	a := func(msg string) {
		defer func() {
			if err := recover(); err != nil {
				l.Printf("err %s\n", err)
			}
		}()
		panic(msg)
	}
	a("1")
	a("2")
	a("3")
}

func TestGeneratedQuery(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"table":   "tbl_test",
		"fields":  ["id", "name as n", "description d", "a", "b", "j", "t"],
		"filter":  {
			"$and": [
				{"a": {"$gt": 10, "$lt": 20}},
				{"$or": [
					{"b": {"$gt": 50}},
					{"b": {"$lt": 20}},
					{"a": {"$in": [123, 234]}}
				]}
			]
		},
		"limit":   1000,
		"offset":  1,
		"order_by": [{"field": "name"}]
	}`)
	req := new(SelectRequest)
	if err := json.Unmarshal(jsonb, &req); err != nil {
		t.Logf("error unmarshalling req: %s", err.Error())
		t.FailNow()
	}

	sql, err := req.ToSQL()
	if err != nil {
		t.Fatalf("cannot generate query command %s", err)
	}
	t.Logf("%s\n", sql)

	_, db := connect(t)
	rows, err := db.Query(context.Background(), sql)
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	cols := Columns(rows)

	fmt.Printf("%+v\n", cols)
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			t.Fatalf("err getting values %s\n", err)
		}
		fmt.Printf("%v\n", vals)
	}
}

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
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
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
		t.Fatalf("QueryRow failed: %v\n", err)
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
		t.Fatalf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("%v\n", tags)
}

func TestConnectionPool(t *testing.T) {
	bg := context.Background()
	pool, err := pgxpool.ConnectConfig(bg, getPoolConfig())
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
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
	// language=postgresql
	tags, err := conn.Query(bg, `SELECT * FROM tbl_test WHERE a > $1`, 0)
	if err != nil {
		t.Fatalf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("%+v\n", tags)
	for tags.Next() {
		values, err := tags.Values()
		if err != nil {
			return
		}
		fmt.Printf("%v\n", values)
	}
}

func TestQueryIn(t *testing.T) {
	conn := getPgxConnection(bg)
	defer viot.CloseCtx(conn, bg)
	// language=postgresql
	tags, err := conn.Query(bg, `SELECT * FROM tbl_test WHERE a = ANY($1)`, []any{17, 18, 19, 20, 21})
	if err != nil {
		t.Fatalf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("%+v\n", tags)
	for tags.Next() {
		values, err := tags.Values()
		if err != nil {
			return
		}
		fmt.Printf("%v\n", values)
	}
}

func Test_DynamicUnmarshall(t *testing.T) {
	checkUnmarshall(t, `{
		"$and": [
			{"a": {"$gt": 10, "$lt": 20}},
			{"$or": [{"b": {"$gt": 50}},{"b": {"$lt": 20}},{"c": {"$in": ["abc", "def", 123]}}]}
		]
	}`)

	checkUnmarshall(t, `{"a": {"$gt": 10}}`)
}

func checkUnmarshall(t *testing.T, s string) {
	filter := Filter{}
	if err := json.Unmarshal([]byte(s), &filter); err != nil {
		t.Errorf("err> %v\n", err)
		return
	}

	if logic, err := filter.ToSQL(); err == nil {
		t.Logf("%s\n", logic)
	}
}

func TestTableDefinition_ToSQL(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"name": "test",
		"fields": [
			{"name": "id", "type": "i32", "is_primary_key": true, "is_logged": false},
			{"name": "name", "type": "string", "is_primary_key": false, "is_logged": false}
		]
	}`)

	def := TableDefinition{}
	err := json.Unmarshal(jsonb, &def)
	sql, err := def.ToSQL()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	t.Logf("%s\n", sql)
}

func TestSelectRequest_ToSQL(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"table":   "tbl_test",
		"fields":  ["id", "name as n", "description"],
		"filter":  {
			"$and": [
				{"a": {"$gt": 10, "$lt": 20}},
				{"$or": [
					{"b": {"$gt": 50}},
					{"b": {"$lt": 20}},
					{"c": {"$in": ["abc", "def", 123]}}
				]}
			]
		},
		"group_by": ["name"],
		"limit":   1000,
		"offset":  10,
		"order_by": [{"field":"name", "ascending": false}]
	}`)
	req := SelectRequest{}
	if err := json.Unmarshal(jsonb, &req); err != nil {
		t.Fatalf("error unmarshalling req: %s", err.Error())
	}

	sql, err := req.ToSQL()
	if err != nil {
		t.FailNow()
		return
	}
	t.Logf("%s\n", sql)
}

func TestUpdateRequest_ToSQL(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"table": "tbl_test",
		"values": {"b": 25, "c": 1},
		"filter": {"c": {"$eq": 3}}
	}`)
	req := UpdateRequest{}
	if err := json.Unmarshal(jsonb, &req); err != nil {
		t.Fatalf("error unmarshalling req: %s", err.Error())
	}

	sql, values, err := req.ToSQL()
	if err != nil {
		t.Logf("error = %v", err)
		t.FailNow()
		return
	}
	t.Logf("%s\n", sql)
	t.Logf("%v\n", values)
}
