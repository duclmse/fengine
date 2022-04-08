package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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

//region Data structure

type Variable struct {
	Name  string      `json:"name,omitempty"`
	Type  Type        `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

type Function struct {
	Input  []Variable `json:"input,omitempty"`
	Output []Variable `json:"output,omitempty"`
	Code   string     `json:"code,omitempty"`
}

type Execution struct {
	Function Function   `json:"function,omitempty"`
	Referee  []Function `json:"referee,omitempty"`
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

func decodeRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var execution Execution
	err = json.Unmarshal(body, &execution)
	if err != nil {
		return body, nil
	}

	return nil, nil
}
