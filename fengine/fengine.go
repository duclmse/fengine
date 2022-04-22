package fengine

import (
	"bytes"
	"encoding/base64"
	"fmt"
	pb "github.com/duclmse/fengine/pb"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type Type int32

const (
	I32 Type = iota
	I64
	F32
	F64
	Bool
	Json
	String
	Binary
)

type JsonArgument struct {
	Name  string `json:"name,omitempty"`
	Type  Type   `json:"type"`
	Value any    `json:"value,omitempty"`
}

type JsonParameter struct {
	Name string  `json:"name,omitempty"`
	Type pb.Type `json:"type"`
}

type JsonArgs []JsonArgument

type JsonParams []JsonParameter

type JsonFunction struct {
	Input  JsonParams `json:"input,omitempty"`
	Output Type       `json:"output,omitempty"`
	Code   string     `json:"code,omitempty"`
}

type JsonScript struct {
	Function      JsonFunction            `json:"function,omitempty"`
	Attributes    []JsonArgument          `json:"attributes"`
	Services      map[string]JsonFunction `json:"services,omitempty"`
	Subscriptions map[string]JsonFunction `json:"subscriptions,omitempty"`
}

type Result struct {
	Value any `json:"value"`
}

type JsonSelectRequest struct {
	Table  string         `json:"table"`
	Fields []string       `json:"fields"`
	Filter map[string]any `json:"filter"`
}

type JsonInsertRequest struct {
	Table  string         `json:"table"`
	Values map[string]any `json:"values"`
}

type JsonUpdateRequest struct {
	Table  string         `json:"table"`
	Values map[string]any `json:"values"`
	Filter map[string]any `json:"filter"`
}

type JsonDeleteRequest struct {
	Table  string         `json:"table"`
	Values map[string]any `json:"values"`
	Filter map[string]any `json:"filter"`
}

var TypeString = map[Type]string{
	I32:    "I32",
	I64:    "I64",
	F32:    "F32",
	F64:    "F64",
	Bool:   "Bool",
	Json:   "Json",
	String: "String",
	Binary: "Binary",
}

var TypeID = map[string]Type{
	"I32":    I32,
	"I64":    I64,
	"F32":    F32,
	"F64":    F64,
	"Bool":   Bool,
	"Json":   Json,
	"String": String,
	"Binary": Binary,
}

func (s *Type) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(TypeString[*s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *Type) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = TypeID[j]
	return nil
}

func (s JsonScript) ToScript() *pb.Script {
	id, _ := uuid.MustParse(``).MarshalBinary()
	return &pb.Script{
		Function: &pb.MethodId{
			ThingID:    id,
			MethodName: "",
			Type:       0,
		},
		Services:     nil,
		Subscription: nil,
	}
}

func (f JsonFunction) ToFunction() *pb.Function {
	return &pb.Function{
		Input: ReadParams(f.Input),
		//Output: f.Output,
		Code: f.Code,
	}
}

func ReadParams(params JsonParams) []*pb.Parameter {
	variables := make([]*pb.Parameter, len(params))
	for _, v := range params {
		variables = append(variables, v.ToParameter())
	}
	return variables
}

func ReadArgs(args JsonArgs) []*pb.Variable {
	variables := make([]*pb.Variable, len(args))
	for _, v := range args {
		variables = append(variables, v.ToArgument())
	}
	return variables
}

func (a JsonParameter) ToParameter() *pb.Parameter {
	return &pb.Parameter{
		Name: a.Name,
		Type: a.Type,
	}
}

func (a JsonArgument) ToArgument() *pb.Variable {
	name := a.Name
	if a.Value == nil {
		return &pb.Variable{Name: name}
	}
	switch a.Type {
	case I32:
		return &pb.Variable{Name: name, Type: pb.Type_i32, Value: &pb.Variable_I32{I32: int32(a.Value.(float64))}}
	case I64:
		return &pb.Variable{Name: name, Type: pb.Type_i64, Value: &pb.Variable_I64{I64: int64(a.Value.(float64))}}
	case F32:
		return &pb.Variable{Name: name, Type: pb.Type_f32, Value: &pb.Variable_F32{F32: float32(a.Value.(float64))}}
	case F64:
		return &pb.Variable{Name: name, Type: pb.Type_f64, Value: &pb.Variable_F64{F64: a.Value.(float64)}}
	case Bool:
		return &pb.Variable{Name: name, Type: pb.Type_bool, Value: &pb.Variable_Bool{Bool: a.Value.(bool)}}
	case Json:
		return &pb.Variable{Name: name, Type: pb.Type_json, Value: &pb.Variable_Json{Json: a.Value.(string)}}
	case String:
		return &pb.Variable{Name: name, Type: pb.Type_string, Value: &pb.Variable_String_{String_: a.Value.(string)}}
	case Binary:
		s := a.Value.(string)
		binary, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			fmt.Printf("cannot decode base64 with %s\n", s)
			return &pb.Variable{}
		}
		return &pb.Variable{Name: name, Type: pb.Type_binary, Value: &pb.Variable_Binary{Binary: binary}}
	}
	return nil
}

func ReadMethods(referee map[string]JsonFunction) map[string]*pb.Function {
	m := make(map[string]*pb.Function, len(referee))
	for k, v := range referee {
		m[k] = v.ToFunction()
	}
	return m
}
