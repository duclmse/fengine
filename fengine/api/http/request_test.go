package http_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-json"

	pb "github.com/duclmse/fengine/pb"
	. "github.com/duclmse/fengine/fengine/db/sql"
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
	// function := execution.Function
	// fmt.Printf("Input ---\n")
	// printParams(function.Input)
	// fmt.Printf("Output ---\n")
	//printArgs(function.Output)
}

func Test_Marshall(t *testing.T) {
	// execution := ServiceRequest{
	// 	Function: Function{
	// 		Input: []Parameter{
	// 			{Name: "str", Type: pb.Type_string},
	// 			{Name: "i32", Type: pb.Type_i32},
	// 		},
	// 		Output: pb.Type_json,
	// 		Code:   "return {str, i32}",
	// 	},
	// }
	// data, err := json.Marshal(execution)
	// if err != nil {
	// 	fmt.Printf("Err marshalling: %v", err)
	// 	return
	// }
	// fmt.Printf("%s\n", string(data))
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

type SqlField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Constraint string `json:"constraint"`
}

type SqlTable struct {
	Name string `json:"name"`
	Fields []SqlField `json:"fields"`
	PrimaryKey []string `json:"primary_key"`
} 

func Test_GenCreateTable(t*testing.T) {
	schema := `{
		"name": "test",
		"fields": [
			{
				"name": "id",
				"type": "serial"
			},
			{
				"name": "name",
				"type": "varchar(50)",
				"constraint": "not null"
			}
		],
		"primary_key": ["id"]
	}`
	
	var table SqlTable
	err := json.Unmarshal([]byte(schema), &table)
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n", table)
	t.Logf("%v\n", GenCreateTable(table))
}

func GenCreateTable(table SqlTable) string {
	var sb strings.Builder
	sb.WriteString("CREATE TABLE ")
	sb.WriteString(table.Name)
	sb.WriteString(" (\n")
	for _, field := range table.Fields {
		sb.WriteString(field.Name)
		sb.WriteString(" ")
		sb.WriteString(field.Type)
		sb.WriteString(" ")
		sb.WriteString(field.Constraint)
		sb.WriteString(",\n")
	}
	sb.WriteString("PRIMARY KEY (")
	for i, f := range table.PrimaryKey {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(f)
	}
	sb.WriteString(")\n)")
	return sb.String()
}
// UPDATE DATA
type SqlNew struct{
	Name string `json:"name"`
	Value string `json:"value"`
}
type SqldataUpdate struct{
	Name string `json:"name"`
	Include []SqlNew `json:"include"`
}
func Test_Updatedata(t *testing.T){
	data := `{
		"name": "SQL",
		"include": [
			{
				"name": "HUONG",
				
				"value": "1"
			},
			{ 
				"name": "XUAN",
			   
			    "value": "2"
			},
			{ 
				"name": "THU",
	
				"value": "3"
			}
		]
		}`
		var Updatedata SqldataUpdate
		err := json.Unmarshal([]byte(data), &Updatedata)
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",Updatedata)
	t.Logf("%v\n", GenUpdatedata(Updatedata))
}

func GenUpdatedata(Updatedata SqldataUpdate) string{
	var sb strings.Builder
	sb.WriteString("Update ")

	sb.WriteString(Updatedata.Name)
	sb.WriteString(" SET \n")
	for i, include := range Updatedata.Include{
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(include.Name)
		sb.WriteString(" = ")
		sb.WriteString(include.Value)
	}

	return sb.String()
}


// ADD COLUMN


type Sqllocation struct{
	Name string `json:"name"`
}
type Sqldroplocation struct{
	Name string `json:"name"`
	Location[]Sqllocation `json:"location"`
	//PrimaryKey []string `json:"primary_key"`
}
func Test_Gendropcolumn(t*testing.T){
	datas := `{
		"name": "SQL",
		"location": [
			{
				"name": "column_name1"				
			}
		]
		}`
		var columndrop Sqldroplocation
		err := json.Unmarshal([]byte(datas), &columndrop )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",columndrop )
	t.Logf("%v\n", Gencolumndrop (columndrop))
}

func Gencolumndrop (columndrop  Sqldroplocation) string{
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE ")
	sb.WriteString(columndrop.Name)
	sb.WriteString(" ADD \n")
	for _ ,column := range columndrop.Column{
		
		sb.WriteString(column.Name)
		sb.WriteString(" = ")
		sb.WriteString(column.Type)
		sb.WriteString("\n")
	}
	// sb.WriteString("PRIMARY KEY (")
	// for i, f := range columndrop.PrimaryKey{
	// 	if i > 0 {
	// 		sb.WriteString(", ")
	// 	}
	// 	sb.WriteString(f)
	// }
	// sb.WriteString(")\n)")
	return sb.String()
}



/// DROP TABLE


type SqlColumn struct{
	Name string `json:"name"`
	Type string `json:"type"`
}
type SqlColumadd struct{
	Name string `json:"name"`
	Column []SqlColumn `json:"column"`
	//PrimaryKey []string `json:"primary_key"`
}
func Test_Genaddcolumn(t*testing.T){
	datas := `{
		"name": "SQL",
		"column": [
			{
				"name": "column_name1",				
				"type": "varchar(20)"
			},
			{ 
				"name": "column_name2",				
				"type": "serial"
			}
		]
		}`
		var columndrop SqlColumadd
		err := json.Unmarshal([]byte(datas), &columndrop )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",columndrop )
	t.Logf("%v\n", Gencolumndrop (columndrop))
}

func Gencolumndrop (columndrop  SqlColumadd) string{
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE ")
	sb.WriteString(columndrop.Name)
	sb.WriteString(" ADD \n")
	for _ ,column := range columndrop.Column{
		
		sb.WriteString(column.Name)
		sb.WriteString(" = ")
		sb.WriteString(column.Type)
		sb.WriteString("\n")
	}
	return sb.String()
}