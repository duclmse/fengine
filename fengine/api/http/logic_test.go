package http

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	//"github.com/goccy/go-json"
	"encoding/json"
)

func Test_DynamicUnmarshall(t *testing.T) {
	logic := Relations{}
	jsonb := `{
		"$and": [
			{"a": {"$gt": 10, "$lt": 20}},
			{"$or": [
				{"b": {"$gt": 50}},
				{"b": {"$lt": 20}},
				{"c": {"$in": ["abc", "def", 123]}}
			]}
		]
	}`
	if err := json.Unmarshal([]byte(jsonb), &logic); err != nil {
		fmt.Printf("err> %v\n", err)
		return
	}
	sb := new(strings.Builder)
	err := buildLogic(logic, sb)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", sb.String())
}

var LogicOperator = map[string]string{
	"$and": "and",
	"$or":  "or",
}

var ComparisonOperator = map[string]string{
	"$gt": ">",
	"$ge": ">=",
	"$eq": "=",
	"$ne": "!=",
	"$le": "<=",
	"$lt": "<",
	"$in": "in",
}

type Relations map[string]interface{}

func buildLogic(logic Relations, sb *strings.Builder) error {
	if len(logic) > 1 {
		return errors.New("cannot have more than one relation")
	}
	for k, v := range logic {
		err := buildRelations(k, v, sb)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildRelations(op string, relations interface{}, sb *strings.Builder) error {
	switch r := relations.(type) {
	case map[string]interface{}:
		//fmt.Printf("map %v\n", r)
	case []interface{}:
		sb.WriteString("(")
		for i, v := range r {
			if i > 0 {
				sb.WriteString(fmt.Sprintf(" %s ", LogicOperator[op]))
			}
			_, err := buildCondition(v, sb)
			if err != nil {
				return err
			}
		}
		sb.WriteString(")")
	}
	return nil
}

func buildCondition(condition interface{}, sb *strings.Builder) (interface{}, error) {
	switch c := condition.(type) {
	case map[string]interface{}:
		if len(c) > 1 {
			return nil, errors.New("condition cannot have more than one key")
		}
		for k, v := range c {
			if strings.HasPrefix(k, "$") {
				err := buildRelations(k, v, sb)
				if err != nil {
					return nil, err
				}
			} else {
				_, err := buildComparison(k, v, sb)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}

// build expression fragment (a < b) from {a: {$lt: b}}
func buildComparison(field string, comparison interface{}, sb *strings.Builder) (interface{}, error) {
	switch c := comparison.(type) {
	case map[string]interface{}:
		i := 0
		for k, v := range c {
			if i > 0 {
				sb.WriteString(fmt.Sprintf(" and "))
			}
			sb.WriteString(fmt.Sprintf("%s %s ", field, ComparisonOperator[k]))
			checkLogicOperand(v, sb)
			i++
		}
	}
	return nil, nil
}

// build
func checkLogicOperand(value interface{}, sb *strings.Builder) {
	switch o := value.(type) {
	case string:
		sb.WriteString(fmt.Sprintf("'%s'", o))
	case int8, int16, int32, int64, float32, float64:
		sb.WriteString(fmt.Sprintf("%v", o))
	case []interface{}:
		sb.WriteString("(")
		for i, v := range o {
			if i > 0 {
				sb.WriteString(",")
			}
			checkLogicOperand(v, sb)
		}
		sb.WriteString(")")
	}
}
