package http

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Marshall(t *testing.T) {
	data := []byte(`{
		"function": {
			"input": [
				{"name": "str", "type": "String", "value": "Hello"},
				{"name": "i32", "type": "Int32", "value": 31071996}
			],
			"output": [
				{"name": "str", "type": "String"},
				{"name": "i32", "type": "Int32"}
			],
			"code": "return {str, i32}"
		}
	}`)
	var execution Execution
	err := json.Unmarshal(data, &execution)
	if err != nil {
		t.Errorf("Err parsing JSON: %s", err)
		return
	}
	fmt.Printf("%v\n", execution)
	function := execution.Function
	fmt.Printf("Input ---\n")
	printArgs(function.Input)
	fmt.Printf("Output ---\n")
	printArgs(function.Output)
}

func Test_Unmarshall(t *testing.T) {
	execution := Execution{
		Function: Function{
			Input: []Variable{
				{Name: "str", Type: String, Value: "hello"},
				{Name: "i32", Type: Int32, Value: 3212312},
			},
			Output: []Variable{
				{Name: "str", Type: String},
				{Name: "i32", Type: Int32},
			},
			Code: "return {str, i32}",
		},
	}
	data, err := json.Marshal(execution)
	if err != nil {
		fmt.Printf("Err marshalling: %v", err)
		return
	}
	fmt.Printf("%s\n", string(data))
}

func printArgs(variables []Variable) {
	for i, arg := range variables {
		fmt.Printf("%d %s %s(%v): %v\n", i, arg.Name, TypeString[arg.Type], arg.Type, arg.Value)
	}
}

func Assert(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected != actual {
		t.Errorf(`%s: Expected "%v" but got "%v"`, message, expected, actual)
	}
}
