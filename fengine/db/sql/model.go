package sql

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/duclmse/fengine/pb"
	"github.com/google/uuid"
)

type EntityType uint8

const (
	Shape EntityType = iota
	Template
	Thing
	invalid
)

type VarType uint8

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
			return fmt.Errorf("scan source is not entity type, but %v", value)
		}
		return nil
	}
	return fmt.Errorf("scan source is not int, but %v", value)
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
	IsPrimaryKey bool    `json:"is_primary_key"`
	IsLogged     bool    `json:"is_logged"`
}

//#region SelectRequest

type OrderBy struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"`
}

type SelectRequest struct {
	Table   string    `json:"table"`
	Fields  []string  `json:"fields"`
	Filter  Filter    `json:"filter"`
	GroupBy []string  `json:"group_by"`
	Limit   int32     `json:"limit"`
	Offset  int32     `json:"offset"`
	OrderBy []OrderBy `json:"order_by"`
}

func (sr SelectRequest) ToSQL() (string, error) {
	logic, err := sr.Filter.ToSQL()
	if err != nil {
		fmt.Printf("err %s\n", err.Error())
		return "", err
	}
	defaultValue := func(prefix, a, b string) string {
		if a == "" {
			return b
		}
		return prefix + a
	}
	order := func(orderBy []OrderBy) string {
		length := len(orderBy)
		if length == 0 {
			return ""
		}
		var sb strings.Builder
		first := orderBy[0]
		asc := map[bool]string{
			true:  "",
			false: "DESC",
		}
		sb.WriteString(fmt.Sprintf("%s %s", first.Field, asc[first.Ascending]))
		for i := 1; i < length; i++ {
			sb.WriteString(fmt.Sprintf("%s %s", first.Field, asc[first.Ascending]))
		}
		return sb.String()
	}
	fields := defaultValue("", strings.Join(sr.Fields, ", "), "*")
	groupBy := defaultValue(" GROUP BY ", strings.Join(sr.GroupBy, ", "), "")
	orderBy := defaultValue(" ORDER BY ", order(sr.OrderBy), "")

	if sr.Limit == 0 || sr.Limit > 10000 {
		sr.Limit = 10000
	}
	if logic != "" {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s%s%s LIMIT %d OFFSET %d",
			fields, sr.Table, logic, groupBy, orderBy, sr.Limit, sr.Offset), nil
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s LIMIT %d OFFSET %d",
		fields, sr.Table, groupBy, orderBy, sr.Limit, sr.Offset), nil
}

type ResultSet struct {
	Columns []string
	Rows    [][]any
}

//#endregion SelectRequest

//#region InsertRequest

type InsertRequest struct {
	Table  string         `json:"table"`
	Values map[string]any `json:"values"`
}

type BatchInsertRequest struct {
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
	Data   [][]any  `json:"data"`
}

func (r InsertRequest) ToSQL() (sql string, params []any, e error) {
	var nameSB strings.Builder
	var varSB strings.Builder

	vars := make([]any, len(r.Values))
	i := 0
	for k, v := range r.Values {
		if i == 1 {
			nameSB.WriteString(fmt.Sprintf("%s", k))
			varSB.WriteString("$1")
		}
		i += 1
		vars[i] = v
		nameSB.WriteString(fmt.Sprintf(", %s", k))
		varSB.WriteString(fmt.Sprintf(", $%d", i+1))
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s);`, r.Table, nameSB.String(), varSB.String()), vars, nil
}

//#endregion InsertRequest

//#region UpdateRequest

type UpdateRequest struct {
	Table  string         `json:"table"`
	Values map[string]any `json:"values"`
	Filter Filter         `json:"filter"`
}

func (r UpdateRequest) ToSQL() (sql string, values []any, e error) {
	logic, e := r.Filter.ToSQL()
	if e != nil {
		fmt.Printf("logic = %v\n", logic)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`UPDATE %s SET `, r.Table))
	after := false
	for k, v := range r.Values {
		if after {
			sb.WriteString(", ")
		} else {
			after = true
		}
		sb.WriteString(fmt.Sprintf("%s = %v", k, v))
	}

	if logic != "" {
		sb.WriteString(fmt.Sprintf(" WHERE %s", logic))
	}
	return sb.String(), values, nil
}

//#endregion UpdateRequest

//#region DeleteRequest

type DeleteRequest struct {
	Table  string `json:"table"`
	Filter Filter `json:"filter"`
}

func (r DeleteRequest) ToSQL() (sql string, values []any, e error) {
	logic, e := r.Filter.ToSQL()
	if e != nil {
		fmt.Printf("logic = %v\n", logic)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`DELETE FROM %s`, r.Table))

	if logic != "" {
		sb.WriteString(fmt.Sprintf(" WHERE %s", logic))
	}
	return sb.String(), values, nil
}

//#endregion DeleteRequest

type AttributeHistoryRequest struct {
}

//#region Function
type Function struct {
	Name   string  `json:"name"`
	Input  Params  `json:"input,omitempty"`
	Output pb.Type `json:"output,omitempty"`
	Code   string  `json:"code,omitempty"`
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
		return &pb.Variable{Name: name, Value: &pb.Variable_I32{I32: int32(a.Value.(float64))}}
	case pb.Type_i64:
		return &pb.Variable{Name: name, Value: &pb.Variable_I64{I64: int64(a.Value.(float64))}}
	case pb.Type_f32:
		return &pb.Variable{Name: name, Value: &pb.Variable_F32{F32: float32(a.Value.(float64))}}
	case pb.Type_f64:
		return &pb.Variable{Name: name, Value: &pb.Variable_F64{F64: a.Value.(float64)}}
	case pb.Type_bool:
		return &pb.Variable{Name: name, Value: &pb.Variable_Bool{Bool: a.Value.(bool)}}
	case pb.Type_json:
		return &pb.Variable{Name: name, Value: &pb.Variable_Json{Json: a.Value.(string)}}
	case pb.Type_str:
		return &pb.Variable{Name: name, Value: &pb.Variable_Str{Str: a.Value.(string)}}
	case pb.Type_bin:
		s := a.Value.(string)
		binary, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			fmt.Printf("cannot decode base64 with %s\n", s)
			return &pb.Variable{}
		}
		return &pb.Variable{Name: name, Value: &pb.Variable_Bin{Bin: binary}}
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

//#endregion Function

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
	case pb.Type_str:
		return "VARCHAR(5000)", nil
	case pb.Type_bin:
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

func comparisonOperator(op string) string {
	switch op {
	case "$gt":
		return ">"
	case "$ge":
		return ">="
	case "$eq":
		return "="
	case "$ne":
		return "!="
	case "$le":
		return "<="
	case "$lt":
		return "<"
	case "$in":
		return "IN"
	default:
		return ""
	}
}

type Filter map[string]interface{}

func (logic Filter) ToSQL() (string, error) {
	//fmt.Printf("filter to sql %v\n", logic)
	logicLength := len(logic)
	sb := &strings.Builder{}
	for k, v := range logic {
		if strings.HasPrefix(k, "$") {
			if logicLength > 1 {
				return "", errors.New("cannot have more than one relation")
			}
			if err := buildRelations(k, v, sb); err != nil {
				fmt.Printf("error building logic %s\n", err.Error())
				return "", err
			}
		} else if err := buildComparison(k, v, sb); err != nil {
			fmt.Printf("error building condition %s\n", err.Error())
			return "", err
		}
	}
	return sb.String(), nil
}

func buildRelations(op string, relations interface{}, sb *strings.Builder) error {
	var opr string
	switch op {
	case "$and":
		opr = "AND"
	case "$or":
		opr = "OR"
	}
	if opr == "" {
		return fmt.Errorf("invalid logic operator '%s'", op)
	}
	switch r := relations.(type) {
	case map[string]interface{}:
		//fmt.Printf("map %v\n", r)
	case []interface{}:
		sb.WriteString("(")
		for i, v := range r {
			if i > 0 {
				sb.WriteString(" ")
				sb.WriteString(strings.ToUpper(opr))
				sb.WriteString(" ")
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
		return nil
	default:
		return fmt.Errorf("build condition with %v", c)
	}
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
				sb.WriteString(" AND ")
			}
			sb.WriteString(fmt.Sprintf("%s %s ", field, comparisonOperator(k)))
			checkLogicOperand(v, sb)
			i++
		}
	case int, int8, int16, int32, int64, float32, float64, string:
		sb.WriteString(fmt.Sprintf("%s = ", field))
		checkLogicOperand(comparison, sb)
	default:
		fmt.Printf("build comp with %v\n", c)
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
