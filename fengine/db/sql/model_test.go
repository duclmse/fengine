package sql_test

import (
	"encoding/json"
	"log"
	"testing"

	. "github.com/duclmse/fengine/fengine/db/sql"
)

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
		log.Printf("%s\n", err)
		return
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
		log.Printf("error unmarshalling req: %s", err.Error())
		t.FailNow()
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
		log.Printf("error unmarshalling req: %s", err.Error())
		t.FailNow()
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
