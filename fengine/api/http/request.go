package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pb "github.com/duclmse/fengine/pb"
)

//region Data structure

type Type int

const (
	Int32 Type = iota
	Int64
	Float
	Double
	Bool
	String
	Bytes
)

type JsonVariable struct {
	Name  string      `json:"name,omitempty"`
	Type  Type        `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

type JsonVars []JsonVariable

type JsonFunction struct {
	Input  JsonVars `json:"input,omitempty"`
	Output JsonVars `json:"output,omitempty"`
	Code   string   `json:"code,omitempty"`
}

type JsonScript struct {
	Attributes []JsonVariable          `json:"attributes"`
	Function   JsonFunction            `json:"function,omitempty"`
	Referee    map[string]JsonFunction `json:"referee,omitempty"`
}

var TypeString = map[Type]string{
	Int32:  "Int32",
	Int64:  "Int64",
	Float:  "Float",
	Double: "Double",
	Bool:   "Bool",
	String: "String",
	Bytes:  "Bytes",
}

var TypeID = map[string]Type{
	"Int32":  Int32,
	"Int64":  Int64,
	"Float":  Float,
	"Double": Double,
	"Bool":   Bool,
	"String": String,
	"Bytes":  Bytes,
}

//endregion Data structure

func decodeAllServiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var execution JsonScript
	err = json.Unmarshal(body, &execution)
	if err != nil {
		fmt.Printf("decode exec: %s\n", err)
		return nil, err
	}

	return execution.ToScript(), nil
}

func decodeServiceRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeExecRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var execution JsonScript
	err = json.Unmarshal(body, &execution)
	if err != nil {
		fmt.Printf("decode exec: %s\n", err)
		return nil, err
	}

	return execution.ToScript(), nil
}

//#region handle data structure

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

func (json JsonScript) ToScript() pb.Script {
	return pb.Script{
		Function:   json.Function.ToFunction(),
		Attributes: ReadVars(json.Attributes),
		Referee:    ReadReferee(json.Referee),
	}
}

func (f JsonFunction) ToFunction() *pb.Function {
	return &pb.Function{
		Input:  ReadVars(f.Input),
		Output: ReadVars(f.Output),
		Code:   f.Code,
	}
}

func (json JsonVariable) ToVariable() *pb.Variable {
	name := json.Name
	if json.Value == nil {
		return &pb.Variable{Name: name}
	}
	switch json.Type {
	case Int32:
		return &pb.Variable{Name: name, Value: &pb.Variable_I32{I32: int32(json.Value.(float64))}}
	case Int64:
		return &pb.Variable{Name: name, Value: &pb.Variable_I64{I64: int64(json.Value.(float64))}}
	case Float:
		return &pb.Variable{Name: name, Value: &pb.Variable_F32{F32: float32(json.Value.(float64))}}
	case Double:
		return &pb.Variable{Name: name, Value: &pb.Variable_F64{F64: json.Value.(float64)}}
	case Bool:
		return &pb.Variable{Name: name, Value: &pb.Variable_Bool{Bool: json.Value.(bool)}}
	case String:
		return &pb.Variable{Name: name, Value: &pb.Variable_String_{String_: json.Value.(string)}}
	case Bytes:
		s := json.Value.(string)
		binary, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			fmt.Printf("cannot decode base64 with %s\n", s)
			return &pb.Variable{}
		}
		return &pb.Variable{Name: name, Value: &pb.Variable_Binary{Binary: binary}}
	}
	return &pb.Variable{Name: name}
}

func ReadReferee(referee map[string]JsonFunction) map[string]*pb.Function {
	m := make(map[string]*pb.Function, len(referee))
	for k, v := range referee {
		m[k] = v.ToFunction()
	}
	return m
}

func ReadVars(attrs JsonVars) []*pb.Variable {
	variables := make([]*pb.Variable, len(attrs))
	for _, v := range attrs {
		variables = append(variables, v.ToVariable())
	}
	return variables
}

//#endregion handle data structure
