package http_test

import (
	"encoding/json"
	"fmt"
	. "github.com/duclmse/fengine/fengine"
	pb "github.com/duclmse/fengine/pb"
	"testing"
)

func Test_Unmarshall(t *testing.T) {
	data := []byte(`{
		"function": {
			"input": [
				{"name": "str", "type": "String", "value": "Hello"},
				{"name": "i32", "type": "I32", "value": 31071996}
			],
			"output": [
				{"name": "str", "type": "String"},
				{"name": "i32", "type": "I32"}
			],
			"code": "return {str, i32}"
		}
	}`)
	var execution ServiceRequest
	err := json.Unmarshal(data, &execution)
	if err != nil {
		t.Errorf("Err parsing JSON: %s", err)
		return
	}
	fmt.Printf("%v\n", execution)
	function := execution.Function
	fmt.Printf("Input ---\n")
	printParams(function.Input)
	fmt.Printf("Output ---\n")
	//printArgs(function.Output)
}

func Test_Marshall(t *testing.T) {
	execution := ServiceRequest{
		Function: Function{
			Input: []Parameter{
				{Name: "str", Type: pb.Type_string},
				{Name: "i32", Type: pb.Type_i32},
			},
			Output: pb.Type_json,
			Code:   "return {str, i32}",
		},
	}
	data, err := json.Marshal(execution)
	if err != nil {
		fmt.Printf("Err marshalling: %v", err)
		return
	}
	fmt.Printf("%s\n", string(data))
}

func Test_Filter(t *testing.T) {
	type a struct {
		A map[string]string `json:"-"`
	}
	data := []byte(`{
		"key": "value"
	}`)
	va := a{}
	err := json.Unmarshal(data, &va)
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Printf("%v\n", va.A)
		return
	}
	fmt.Printf("%v\n", va.A)
}

func printParams(args Params) {
	for i, arg := range args {
		fmt.Printf("%d %s %s(%v)\n", i, arg.Name, pb.Type_name[int32(arg.Type)], arg.Type)
	}
}

func printArgs(args []Variable) {
	for i, arg := range args {
		fmt.Printf("%d %s %s(%v): %v\n", i, arg.Name, pb.Type_name[int32(arg.Type)], arg.Type, arg.Value)
	}
}
