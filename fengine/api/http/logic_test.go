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
	buildLogic(logic)
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

func buildLogic(logic Relations) error {
	if len(logic) > 1 {
		return errors.New("cannot have more than one relation")
	}
	for k, v := range logic {
		err := buildRelations(k, v)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\n")
	return nil
}

func buildRelations(op string, relations interface{}) error {
	switch r := relations.(type) {
	case map[string]interface{}:
		//fmt.Printf("map %v\n", r)
	case []interface{}:
		fmt.Printf("(")
		for i, v := range r {
			if i > 0 {
				fmt.Printf(" %s ", LogicOperator[op])
			}
			_, err := buildCondition(v)
			if err != nil {
				return err
			}
		}
		fmt.Printf(")")
	}
	return nil
}

func buildCondition(condition interface{}) (interface{}, error) {
	switch c := condition.(type) {
	case map[string]interface{}:
		if len(c) > 1 {
			return nil, errors.New("condition cannot have more than one key")
		}
		for k, v := range c {
			if strings.HasPrefix(k, "$") {
				err := buildRelations(k, v)
				if err != nil {
					return nil, err
				}
			} else {
				_, err := buildComparison(k, v)
				if err != nil {
					return nil, err
				}

			}
		}
	}
	return nil, nil
}

func buildComparison(field string, comparison interface{}) (interface{}, error) {
	switch c := comparison.(type) {
	case map[string]interface{}:
		i := 0
		for k, v := range c {
			if i > 0 {
				fmt.Printf(" and ")
			}
			fmt.Printf("%s %s %v", field, ComparisonOperator[k], checkLogicOperand(v))
			i++
		}
		//fmt.Printf("\n")
	}
	return nil, nil
}

func checkLogicOperand(value interface{}) string {
	switch o := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", o)
	case int8, int16, int32, int64, float32, float64:
		return fmt.Sprintf("%v", o)
	case []interface{}:
		var sb strings.Builder
		sb.WriteString("(")
		for i, v := range o {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(checkLogicOperand(v))
		}
		sb.WriteString(")")
		return sb.String()
	}
	return ""
}
