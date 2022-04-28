package sql

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	pb "github.com/duclmse/fengine/pb"
	"github.com/google/uuid"
	"strings"
	"time"
)

type EntityType uint8
type VarType uint8

const (
	Shape EntityType = iota
	Template
	Thing
	invalid
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

func EntityTypeValue(s string) EntityType {
	switch s {
	case "shape":
		return Shape
	case "template":
		return Template
	case "thing":
		return Thing
	default:
		return invalid
	}
}

type JsonString string

type Entity struct {
	Id          *uuid.UUID `sql:"id,type:uuid"`
	Name        string     `sql:"name"`
	Type        EntityType `sql:"type"`
	Description *string    `sql:"description"`
	ProjectId   *uuid.UUID `sql:"project_id,type:uuid"`
}

type Attribute struct {
	EntityId    *uuid.UUID `sql:"entity_id,type:uuid"`
	Name        string     `sql:"name"`
	Type        pb.Type    `sql:"type"`
	From        *uuid.UUID `sql:"from,type:uuid"`
	ValueI32    int32      `sql:"value_i32"`
	ValueI64    int32      `sql:"value_i64"`
	ValueF32    float32    `sql:"value_f32"`
	ValueF64    float64    `sql:"value_f64"`
	ValueBool   bool       `sql:"value_bool"`
	ValueJson   *string    `sql:"value_json"`
	ValueString *string    `sql:"value_string"`
	ValueBinary []byte     `sql:"value_binary"`
}

type EntityService struct {
	Id          *uuid.UUID  `json:"id" sql:"id,type:uuid"`
	Name        string      `json:"name" sql:"name"`
	Type        EntityType  `json:"type" sql:"type"`
	Description *string     `json:"description" sql:"description"`
	ProjectId   *uuid.UUID  `json:"project_id" db:"project_id" sql:",type:uuid"`
	Input       *JsonString `json:"input,omitempty" sql:"entity_type"`
	Output      *string     `json:"output" sql:"output"`
	From        *uuid.UUID  `json:"from,omitempty" sql:"from,type:uuid"`
	Code        *string     `json:"code" sql:"code,type:uuid"`
	CreateTs    uuid.Time   `json:"create_ts" sql:"create_ts"`
	UpdateTs    uuid.Time   `json:"update_ts" sql:"update_ts"`
}

type EntitySubscription struct {
	Id          *uuid.UUID  `json:"id" sql:"id,type:uuid"`
	Name        string      `json:"name" sql:"name"`
	Type        EntityType  `json:"type" sql:"type"`
	Description *string     `json:"description" sql:"description"`
	ProjectId   *uuid.UUID  `json:"project_id" db:"project_id" sql:",type:uuid"`
	Input       *JsonString `json:"input,omitempty" sql:"entity_type"`
	Output      *string     `json:"output" sql:"output"`
	From        *uuid.UUID  `json:"from,omitempty" sql:"from,type:uuid"`
	Code        *string     `json:"code" sql:"code,type:uuid"`
	CreateTs    uuid.Time   `json:"create_ts" sql:"create_ts"`
	UpdateTs    uuid.Time   `json:"update_ts" sql:"update_ts"`
}

func (et *EntityType) Scan(value any) error {
	if i, ok := value.([]byte); ok {
		switch string(i) {
		case "shape", "Shape":
			*et = Shape
		case "template", "Template":
			*et = Template
		case "thing", "Thing":
			*et = Thing
		default:
			return errors.New(fmt.Sprintf("scan source is not entity type, but %v", value))
		}
		return nil
	}
	return errors.New(fmt.Sprintf("scan source is not int, but %v", value))
}

func (et EntityType) Value() (driver.Value, error) {
	return int(et), nil
}

func (et *EntityType) MarshalJSON() ([]byte, error) {
	if et == nil {
		return []byte(`"null"`), nil
	}
	switch *et {
	case Shape:
		return []byte(`"Shape"`), nil
	case Template:
		return []byte(`"Template"`), nil
	case Thing:
		return []byte(`"Thing"`), nil
	}

	return nil, errors.New("invalid entity type")
}

func (vt *VarType) Scan(value any) error {
	if i, ok := value.([]byte); ok {
		switch string(i) {
		case "i32":
			*vt = I32
		case "i64":
			*vt = I64
		case "f32":
			*vt = F32
		case "f64":
			*vt = F64
		case "bool":
			*vt = Bool
		case "json":
			*vt = Json
		case "string":
			*vt = String
		case "binary":
			*vt = Binary
		}
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

type UuidArray []uuid.UUID

type UuidString string

func (us UuidString) ToUuidArray() (*UuidArray, error) {
	l := len(us)
	uuids := strings.Split(string(us)[1:l-1], ",")
	var ids UuidArray
	for _, sid := range uuids {
		uid, e := uuid.Parse(sid)
		if e != nil {
			return nil, e
		}
		ids = append(ids, uid)
	}
	return &ids, nil
}

func (ua *UuidArray) MarshalJSON() ([]byte, error) {
	if ua == nil || len(*ua) == 0 {
		return []byte(`[]`), nil
	}
	bb := bytes.Buffer{}
	bb.Write([]byte(`["`))
	for i, v := range *ua {
		if i > 0 {
			bb.Write([]byte(`","`))
		}
		text, err := v.MarshalText()
		if err != nil {
			return nil, err
		}
		bb.Write(text)

	}
	bb.Write([]byte(`"]`))
	return bb.Bytes(), nil
}

type EntityDefinition struct {
	Id            uuid.UUID      `json:"id" db:"id,type:uuid"`
	Name          string         `json:"name"`
	Type          EntityType     `json:"type"`
	Description   *string        `json:"description"`
	ProjectId     *uuid.UUID     `json:"project_id" db:"project_id,type:uuid"`
	BaseTemplate  *uuid.UUID     `json:"base_template" db:"base_template,type:uuid"`
	BaseShapesStr *UuidString    `json:"-" db:"base_shapes"`
	BaseShapes    *UuidArray     `json:"base_shapes"`
	Attributes    []Variable     `json:"attributes,omitempty"`
	Services      []Function     `json:"services,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`
	CreateTs      *time.Time     `json:"create_ts" db:"create_ts"`
	UpdateTs      *time.Time     `json:"update_ts" db:"update_ts"`
}

func (d EntityDefinition) ToThingServices() ([]ThingService, error) {
	services := make([]ThingService, len(d.Services))
	for _, svc := range d.Services {
		services = append(services, ThingService{
			EntityId: d.Id,
			Name:     svc.Name,
			Input:    svc.Input,
			Output:   svc.Output,
			Code:     svc.Code,
		})
	}

	return services, nil
}

func (d EntityDefinition) ToThingSubscriptions() ([]ThingSubscription, error) {
	subs := make([]ThingSubscription, len(d.Subscriptions))
	for _, sub := range d.Subscriptions {
		subs = append(subs, ThingSubscription{
			EntityId:  d.Id,
			Name:      sub.Name,
			Enabled:   sub.Enabled,
			Event:     sub.Event,
			Attribute: sub.Attribute,
		})
	}

	return subs, nil
}

type Args []Variable

type Variable struct {
	Name  string      `json:"name,omitempty"`
	Type  pb.Type     `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

type Params []Parameter

type Parameter struct {
	Name string  `json:"name,omitempty"`
	Type pb.Type `json:"type"`
}

type ThingService struct {
	EntityId uuid.UUID  `json:"entity_id"`
	Name     string     `json:"name"`
	Input    Params     `json:"input,omitempty"`
	Output   pb.Type    `json:"output,omitempty"`
	Code     string     `json:"code,omitempty"`
	From     *uuid.UUID `json:"from"`
}

type ThingServiceId struct {
	EntityId uuid.UUID `json:"entity_id"`
	Name     string    `json:"name"`
}

type Function struct {
	Name   string  `json:"name"`
	Input  Params  `json:"input,omitempty"`
	Output pb.Type `json:"output,omitempty"`
	Code   string  `json:"code,omitempty"`
}

type ThingSubscription struct {
	EntityId  uuid.UUID  `json:"entity_id"`
	Name      string     `json:"name"`
	Enabled   bool       `json:"enabled"`
	Event     string     `json:"event"`
	Attribute Variable   `json:"attribute"`
	From      *uuid.UUID `json:"from"`
}

type ThingSubscriptionId struct {
	EntityId uuid.UUID `json:"entity_id"`
	Name     string    `json:"name"`
}

type Subscription struct {
	Name      string   `json:"name"`
	Enabled   bool     `json:"enabled"`
	Event     string   `json:"event"`
	Attribute Variable `json:"attribute"`
}

type ServiceRequest struct {
	ThingId     uuid.UUID `json:"thing_id"`
	ServiceName string    `json:"service_name"`
	Input       Args      `json:"input"`
}

type TableDefinition struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name         string  `json:"name"`
	Type         pb.Type `json:"type" sql:"type"`
	IsPrimaryKey bool    `json:"is_primary_key" json:"is_primary_key"`
	IsLogged     bool    `json:"is_logged" json:"is_logged"`
}

type SelectRequest struct {
	Table   string   `json:"table"`
	Fields  []string `json:"fields"`
	Filter  Filter   `json:"filter"`
	GroupBy []string `json:"group_by"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
	OrderBy []string `json:"order_by"`
}

func (sr SelectRequest) ToSQL() (string, error) {
	defaultValue := func(prefix, a, b string) string {
		if a == "" {
			return b
		}
		return prefix + a
	}
	fields := defaultValue("", strings.Join(sr.Fields, ", "), "*")
	groupBy := defaultValue(" GROUP BY ", strings.Join(sr.GroupBy, ", "), "")
	orderBy := defaultValue(" ORDER BY ", strings.Join(sr.OrderBy, ", "), "")
	logic, err := sr.Filter.BuildLogic()
	if err != nil {
		fmt.Printf("err %s\n", err.Error())
		return "", err
	}
	if sr.Limit == 0 || sr.Limit > 10000 {
		sr.Limit = 10000
	}
	if filter := logic.String(); filter != "" {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s%s%s LIMIT %d OFFSET %d",
			fields, sr.Table, filter, groupBy, orderBy, sr.Limit, sr.Offset), nil
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s LIMIT %d OFFSET %d",
		fields, sr.Table, groupBy, orderBy, sr.Limit, sr.Offset), nil
}

type InsertRequest struct {
	Table  string           `json:"table"`
	Values []map[string]any `json:"values"`
}

func (r InsertRequest) ToSQL() (sql string, e error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`INSERT INTO %s (`, r.Table))
	sb.WriteString(` (`)
	sb.WriteString(`) VALUES (`)

	return sb.String(), nil
}

type UpdateRequest struct {
	Table  string        `json:"table"`
	Values []pb.Variable `json:"values"`
	Filter Filter        `json:"filter"`
}

type DeleteRequest struct {
	Table  string        `json:"table"`
	Values []pb.Variable `json:"values"`
	Filter Filter        `json:"filter"`
}

type AttributeHistoryRequest struct {
}

func (f Function) ToFunction() *pb.Function {
	return &pb.Function{
		Input: ReadParams(f.Input),
		//Output: f.Output,
		Code: f.Code,
	}
}

func ReadParams(params Params) []*pb.Parameter {
	variables := make([]*pb.Parameter, len(params))
	for _, v := range params {
		variables = append(variables, v.ToParameter())
	}
	return variables
}

func ReadArgs(args Args) []*pb.Variable {
	variables := make([]*pb.Variable, len(args))
	for _, v := range args {
		variables = append(variables, v.ToArgument())
	}
	return variables
}

func (a Parameter) ToParameter() *pb.Parameter {
	return &pb.Parameter{
		Name: a.Name,
		Type: a.Type,
	}
}

func (a Variable) ToArgument() *pb.Variable {
	name := a.Name
	if a.Value == nil {
		return &pb.Variable{Name: name}
	}
	switch a.Type {
	case pb.Type_i32:
		return &pb.Variable{Name: name, Type: pb.Type_i32, Value: &pb.Variable_I32{I32: int32(a.Value.(float64))}}
	case pb.Type_i64:
		return &pb.Variable{Name: name, Type: pb.Type_i64, Value: &pb.Variable_I64{I64: int64(a.Value.(float64))}}
	case pb.Type_f32:
		return &pb.Variable{Name: name, Type: pb.Type_f32, Value: &pb.Variable_F32{F32: float32(a.Value.(float64))}}
	case pb.Type_f64:
		return &pb.Variable{Name: name, Type: pb.Type_f64, Value: &pb.Variable_F64{F64: a.Value.(float64)}}
	case pb.Type_bool:
		return &pb.Variable{Name: name, Type: pb.Type_bool, Value: &pb.Variable_Bool{Bool: a.Value.(bool)}}
	case pb.Type_json:
		return &pb.Variable{Name: name, Type: pb.Type_json, Value: &pb.Variable_Json{Json: a.Value.(string)}}
	case pb.Type_string:
		return &pb.Variable{Name: name, Type: pb.Type_string, Value: &pb.Variable_String_{String_: a.Value.(string)}}
	case pb.Type_binary:
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

func ReadMethods(referee map[string]Function) map[string]*pb.Function {
	m := make(map[string]*pb.Function, len(referee))
	for k, v := range referee {
		m[k] = v.ToFunction()
	}
	return m
}

func SqlType(t pb.Type) (string, error) {
	switch t {
	case pb.Type_i32:
		return "INT", nil
	case pb.Type_i64:
		return "INT", nil
	case pb.Type_f32:
		return "BIGINT", nil
	case pb.Type_f64:
		return "FLOAT(4)", nil
	case pb.Type_bool:
		return "FLOAT(8)", nil
	case pb.Type_json:
		return "JSONB", nil
	case pb.Type_string:
		return "VARCHAR(5000)", nil
	case pb.Type_binary:
		return "BYTEA", nil
	default:
		return "", errors.New("invalid db type")
	}
}

func (td TableDefinition) ToSQL() (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (", td.Name))
	keys := []string{}
	for i, v := range td.Fields {
		if i > 0 {
			sb.WriteString(", ")
		}
		t, err := SqlType(v.Type)
		if err != nil {
			return "", err
		}
		s := fmt.Sprintf("%s %s", v.Name, t)
		if v.IsPrimaryKey {
			keys = append(keys, v.Name)
		}
		sb.WriteString(s)
	}
	if len(keys) == 0 {
		sb.WriteString(");")
		return sb.String(), nil
	}

	sb.WriteString(", PRIMARY KEY (")
	for i, v := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v)
	}
	sb.WriteString("));")
	return sb.String(), nil
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

type Filter map[string]interface{}

func (logic Filter) BuildLogic() (*strings.Builder, error) {
	if len(logic) > 1 {
		return nil, errors.New("cannot have more than one relation")
	}
	sb := new(strings.Builder)
	for k, v := range logic {
		if strings.HasPrefix(k, "$") {
			if err := buildRelations(k, v, sb); err != nil {
				fmt.Printf("error building logic %s\n", err.Error())
				return nil, err
			}
		} else if err := buildComparison(k, v, sb); err != nil {
			fmt.Printf("error building condition %s\n", err.Error())
			return nil, err
		}

	}
	return sb, nil
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
			if err := buildCondition(v, sb); err != nil {
				return err
			}
		}
		sb.WriteString(")")
	}
	return nil
}

// build expression fragment (a < b) from {a: {$lt: b}}
func buildCondition(condition interface{}, sb *strings.Builder) error {
	switch c := condition.(type) {
	case map[string]interface{}:
		if len(c) > 1 {
			return errors.New("condition cannot have more than one key")
		}
		for k, v := range c {
			if strings.HasPrefix(k, "$") {
				if err := buildRelations(k, v, sb); err != nil {
					return err
				}
			} else if err := buildComparison(k, v, sb); err != nil {
				return err
			}
		}
	default:
		fmt.Printf("build cond with %v\n", c)
	}
	return nil
}

// build expression fragment (a < b) from {$lt: b}
func buildComparison(field string, comparison interface{}, sb *strings.Builder) error {
	//fmt.Printf("build comp with %t\n", comparison)
	switch c := comparison.(type) {
	case map[string]interface{}:
		i := 0
		for k, v := range c {
			if !strings.HasPrefix(k, "$") {
				return errors.New("expect comparison operator (start with '$')")
			}
			if i > 0 {
				sb.WriteString(fmt.Sprintf(" and "))
			}
			sb.WriteString(fmt.Sprintf("%s %s ", field, ComparisonOperator[k]))
			checkLogicOperand(v, sb)
			i++
		}
	default:
		//fmt.Printf("build comp with %v\n", c)
	}
	return nil
}

// build
func checkLogicOperand(value interface{}, sb *strings.Builder) {
	switch o := value.(type) {
	case string:
		sb.WriteString(fmt.Sprintf("'%s'", o))
	case int, int8, int16, int32, int64, float32, float64:
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
