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
// CREATE TABLE
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
		var dropcolumn Sqldroplocation
		err := json.Unmarshal([]byte(datas), &dropcolumn )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",dropcolumn )
	t.Logf("%v\n", Gendropcolumn (dropcolumn))
}

func Gendropcolumn (dropcolumn  Sqldroplocation) string{
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE ")
	sb.WriteString(dropcolumn.Name)
	sb.WriteString(" ADD \n")
	for _ ,column := range dropcolumn.Location{
		
		sb.WriteString(column.Name)
		sb.WriteString(" = ")
	
	}
	// sb.WriteString("PRIMARY KEY (")
	// for i, f := range dropcolumn.PrimaryKey{
	// 	if i > 0 {
	// 		sb.WriteString(", ")
	// 	}
	// 	sb.WriteString(f)
	// }
	// sb.WriteString(")\n)")
	return sb.String()
}



/// DROP TABLE


type SqlCo struct{
	Name string `json:"name"`
}
type Sqldropcolumn struct{
	Name string `json:"name"`
	Table []SqlCo `json:"table"`
	
}
func Test_dropcolumn(t*testing.T){
	datas := `{
		"name": "SQL",
		"table": [
			{
				"name": "table_name_1"		
	
			},
			{ 
				"name": "table_name_2"
				
			}
		]
		}`
		var dropcolumn Sqldropcolumn 
		err := json.Unmarshal([]byte(datas), &dropcolumn )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",dropcolumn )
	t.Logf("%v\n", Gendrop(dropcolumn))
}

func Gendrop (dropcolumn  Sqldropcolumn) string{
	var sb strings.Builder
	sb.WriteString("DROP ")
	sb.WriteString(dropcolumn.Name)
	sb.WriteString(" DROP \n")
	for _ ,table := range dropcolumn.Table{
		
		sb.WriteString(table.Name)
		sb.WriteString("\n")
	}
	return sb.String()
}

//ALTER DROP COLUM
type Columndata struct{
	Name string `json:"name"`
}
type SqlDropcolumn struct{
	Name string `json:"name"`
	Column []Columndata `json:"column"`
	
}
func Test_Dropcolum(t *testing.T){
	data :=`{
		"name": "SQL",
		"column":[
			{
				"name": "column_name"
			}
			]
		}`

var dropcolumn SqlDropcolumn
		err := json.Unmarshal([]byte(data), &dropcolumn )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",dropcolumn )
	t.Logf("%v\n", GenDropcolumn(dropcolumn))
}


func GenDropcolumn (dropcolumn  SqlDropcolumn) string{
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE  ")
	sb.WriteString(dropcolumn.Name)
	sb.WriteString("\n DROP COLUMN")
	for _ ,table := range dropcolumn.Column{
		sb.WriteString("   ")
		sb.WriteString(table.Name)
		sb.WriteString("\n")
	}
	return sb.String()
}


// ALTER SET COMPRESSION
type comp struct{
	Name string `json : "name"`
}
type Altercom struct{
	Name string `json : "name"`
	Method []comp `json:"method"`
}
func Test_Setcp(t *testing.T){
	com:=`{
		"name": "column_name",
		"method": [
			{
				"name": "compression_method"
			}
		]
	}`
	var cpn Altercom
	err:=json.Unmarshal([]byte(com),&cpn)
	if err!=nil {
		fmt.Printf("err passing json %v\n",err)
		return
	}
	fmt.Printf("%v\n",cpn)
	t.Logf("%v\n",GenCom(cpn))
}

func GenCom (cpn Altercom) string{
	var sb strings.Builder
	sb.WriteString(" ALTER  ")
	sb.WriteString(cpn.Name)
	sb.WriteString("\n SET COMPRESSION")
	for _,method := range cpn.Method{
       sb.WriteString("  ")
	   sb.WriteString( method.Name)
	}
	return sb.String()
}

// ALTER SET STATISTICS
type Static struct{
	Number string `json: "number"`
}
 type Clstt struct{
	 Name string `json: "name"`
	 Stts []Static `json: "stts"`
 }
func Test_Statistics(t *testing.T){
	dtas:=`{
		"name": "column_name",
		"stts":[
			{
				"number": "iteger_"
			}
		]
	}`
var cstt Clstt
	err:=json.Unmarshal([]byte(dtas),&cstt)
	 if err != nil {
		 fmt.Printf("err passing json %v\n",err)
		 return
	 }
	 fmt.Printf("%v\n",cstt)
	 t.Logf("%v\n",Genstt(cstt))
}
func Genstt(cstt Clstt) string{
	var sb strings.Builder
	sb.WriteString(" ALTER")
	sb.WriteString(cstt.Name)
	sb.WriteString(" \n SET STATISTICS ")
	for _ ,Stts := range cstt.Stts{
		sb.WriteString("  ")
		sb.WriteString(Stts.Number)
	}
	return sb.String()
}

// RENAME TABLE

type Re struct{
	Name string `json:"name"`
}
type SqlRetb struct{
	Name string `json:"name"`
	Tbname []Re `json:"tbname"`
	
}
func Test_Rename(t *testing.T){
	data :=`{
		"name": "SQL",
		"tbname":[
			{
				"name": "new_table_name"
			}
			]
		}`

var retb SqlRetb
		err := json.Unmarshal([]byte(data), &retb )
	if err != nil {
		fmt.Printf("err parsing json %v\n", err)
		return
	}
	fmt.Printf("%v\n",retb )
	t.Logf("%v\n", GenRname(retb))
}


func GenRname(retb  SqlRetb) string{
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE  ")
	sb.WriteString(retb.Name)
	sb.WriteString("\n RENAME TO")
	for _ ,tbname := range retb.Tbname{
		sb.WriteString("   ")
		sb.WriteString(tbname.Name)
		sb.WriteString("\n")
	}
	return sb.String()
}

// SET SCHEMA


type Sch struct{
	Name string `json: "name"`
}

type SqlRes struct{
	Name string `json: "name"`
	Tbname []Sch `json: "tbname"`
}
func Test_Setschema(t *testing.T){
	dts :=`{
		"name": "name",
		"tbname": [
			{
				"name": "new_schema"
			}
		]
		}`
		var  resch SqlRes
		err:= json.Unmarshal([]byte(dts), &resch)
		if err!=nil {
			fmt.Printf("err passing json %v\n",err)
			return
			}
			fmt.Printf("%v\n",resch)
			t.Logf("%v\n",GenSCH(resch))
}
func GenSCH (resch SqlRes) string {
	var sb strings.Builder
	sb.WriteString("ALTER_TABLE  ")
	sb.WriteString(resch.Name)
	sb.WriteString("\n SET SCHEMA  ")
	for _,tbname :=range resch.Tbname{
		sb.WriteString("  ")
		sb.WriteString(tbname.Name)
		sb.WriteString("\n")
	}
	return sb.String()

}

// SET TABLESPACE


type Space struct{
	Name string `json: "name"`
}

type Tbspace struct{
	Name string `json: "name"`
	Tbname []Space `json: ""tbname`
}
func Test_TbSpace(t *testing.T){
	tbs:=`{
		"name": "name",
		"tbname" :[
			{
				"name": "new_tablespace"
			}

		]
	}`
 var tbsp Tbspace
 err:=json.Unmarshal([]byte(tbs), &tbsp)
 if err!=nil {
	 fmt.Printf("err passing json %v\n", err)
	 return
 }
 fmt.Printf( "%v\n",tbsp)
 t.Logf("%v\n",GenTbsp(tbsp))
}
func GenTbsp( tbsp Tbspace) string{
	var sb strings.Builder
	sb.WriteString(" ALTER_TABLE ALL IN TABLE SPACE ")
	sb.WriteString(tbsp.Name)
	sb.WriteString("\n SET TABLE SPACE  ")

	for _,tbname := range tbsp.Tbname{
     sb.WriteString("  ")
	 sb.WriteString(tbname.Name)
	}
	return sb.String()
}
 
// DETACH PARTITION

type Par struct{
	Iname string `json : "iname"`
}
type Tbpar struct{
	Id string `json : "id"`
	Partition []Par `json: "partition"`
}
func Test_Detach(t *testing.T){
	describe:= `{
		"id" : "name",
		"partition": [
			{
				"iname" : "partition_name"
			}
		]
	}`
	var depar Tbpar
	err:=json.Unmarshal([]byte(describe), &depar)
	if err!=nil {
		fmt.Printf("err passing json %v\n",err)
		return
	}
	fmt.Printf("%v\n",depar)
	t.Logf("%v\n",GenDpar(depar))
}
func GenDpar (depar Tbpar) string{
	var sb strings.Builder
	sb.WriteString(" ALTER TABLE  ")
	sb.WriteString(depar.Id)
	sb.WriteString("\n DETACH PARTITION ")
	for _,partition := range depar.Partition{
		sb.WriteString("  ")
		sb.WriteString(partition.Iname)
		sb.WriteString("\n")
	}
	return sb.String()
}
