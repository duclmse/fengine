package sql

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

type Entity struct {
	Id          *uuid.UUID `sql:"id,type:uuid"`
	Name        string     `sql:"name"`
	Type        EntityType `sql:"type"`
	Description string     `sql:"description"`
	ProjectId   *uuid.UUID `db:"project_id" sql:",type:uuid"`
}

type Attribute struct {
	EntityId    *uuid.UUID `sql:"entity_id,type:uuid"`
	Name        string     `sql:"name"`
	Type        VarType    `sql:"var_type"`
	From        *uuid.UUID `sql:"from,type:uuid"`
	ValueI32    int32      `sql:"value_i32"`
	ValueI64    int32      `sql:"value_i64"`
	ValueF32    float32    `sql:"value_f32"`
	ValueF64    float64    `sql:"value_f64"`
	ValueBool   bool       `sql:"value_bool"`
	ValueJson   string     `sql:"value_json"`
	ValueString string     `sql:"value_string"`
	ValueBinary []byte     `sql:"value_binary"`
}

type Method struct {
	EntityId    *uuid.UUID `sql:"entity_id,type:uuid"`
	Name        string     `sql:"name"`
	EntityType  EntityType `sql:"entity_type"`
	Description string     `sql:"description"`
	ProjectId   *uuid.UUID `sql:"project_id,type:uuid"`
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
