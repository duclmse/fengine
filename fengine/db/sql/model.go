package sql

import (
	"database/sql/driver"
	"errors"
	"fmt"
	. "github.com/google/uuid"
)

type EntityType uint8
type MethodType uint8
type VarType uint8

const (
	Shape EntityType = iota
	Template
	Thing
)

const (
	Service MethodType = iota
	Subscription
)

const (
	I32 VarType = iota
	I64
	F32
	F64
	Bool
	Json
	String
	Binary
)

type JsonString string

type Entity struct {
	Id          *UUID      `sql:"id,type:uuid"`
	Name        string     `sql:"name"`
	Type        EntityType `sql:"type"`
	Description *string    `sql:"description"`
	ProjectId   *UUID      `sql:"project_id,type:uuid"`
}

type Attribute struct {
	EntityId    *UUID   `sql:"entity_id,type:uuid"`
	Name        string  `sql:"name"`
	Type        VarType `sql:"var_type"`
	From        *UUID   `sql:"from,type:uuid"`
	ValueI32    int32   `sql:"value_i32"`
	ValueI64    int32   `sql:"value_i64"`
	ValueF32    float32 `sql:"value_f32"`
	ValueF64    float64 `sql:"value_f64"`
	ValueBool   bool    `sql:"value_bool"`
	ValueJson   *string `sql:"value_json"`
	ValueString *string `sql:"value_string"`
	ValueBinary []byte  `sql:"value_binary"`
}

type Method struct {
	EntityId *UUID       `sql:"entity_id,type:uuid"`
	Name     string      `sql:"name"`
	Input    *JsonString `sql:"input"`
	Output   *string     `sql:"output"`
	From     UUID        `sql:"from,type:uuid"`
	Code     *string     `sql:"code,type:uuid"`
}

type EntityMethod struct {
	Id          *UUID       `json:"id" sql:"id,type:uuid"`
	Name        string      `json:"name" sql:"name"`
	Type        EntityType  `json:"type" sql:"type"`
	Description *string     `json:"description" sql:"description"`
	ProjectId   *UUID       `json:"project_id" db:"project_id" sql:",type:uuid"`
	MethodName  string      `json:"method_name" sql:"name"`
	Input       *JsonString `json:"input,omitempty" sql:"entity_type"`
	Output      *string     `json:"output" sql:"output"`
	From        *UUID       `json:"from,omitempty" sql:"from,type:uuid"`
	Code        *string     `json:"code" sql:"code,type:uuid"`
}

type Field struct {
}

type Table struct {
	Fields map[string]Field
}

func (et *EntityType) Scan(value interface{}) error {
	if i, ok := value.([]byte); ok {
		*et = map[string]EntityType{
			"shape":    Shape,
			"template": Template,
			"thing":    Thing,
		}[string(i)]
		return nil
	}
	return errors.New(fmt.Sprintf("scan source is not int, but %v", value))
}

func (et EntityType) Value() (driver.Value, error) {
	return int(et), nil
}

func (mt *MethodType) Scan(value interface{}) error {
	if i, ok := value.([]byte); ok {
		*mt = map[string]MethodType{
			"service":      Service,
			"subscription": Subscription,
		}[string(i)]
		return nil
	}
	return errors.New("scan source is not int")
}

func (mt MethodType) Value() (driver.Value, error) {
	return int(mt), nil
}

func (vt *VarType) Scan(value interface{}) error {
	if i, ok := value.([]byte); ok {
		*vt = map[string]VarType{
			"i32":    I32,
			"i64":    I64,
			"f32":    F32,
			"f64":    F64,
			"bool":   Bool,
			"json":   Json,
			"string": String,
			"binary": Binary,
		}[string(i)]
		return nil
	}
	return errors.New("scan source is not int")
}

func (vt VarType) Value() (driver.Value, error) {
	return int(vt), nil
}

func (js *JsonString) MarshalJSON() ([]byte, error) {
	if js == nil {
		return []byte(`"null"`), nil
	}
	return []byte(*js), nil
}
