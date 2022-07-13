package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duclmse/fengine/fengine/db/sql"
	viot "github.com/duclmse/fengine/pb"
)

func encodeSelectResponse(ctx context.Context, r any) (response any, err error) {
	//fmt.Printf("encodeSelectResponse %t\n", r)
	res, ok := r.(*sql.ResultSet)
	if !ok {
		fmt.Printf("cannot convert %t\n", r)
		return &viot.SelectResult{}, err
	}

	resultRows := res.Rows
	rows := make([]*viot.ResultRow, len(resultRows))
	for i, row := range resultRows {
		a := make([]*viot.Value, len(row))
		for j, value := range row {
			switch vl := value.(type) {
			case int32:
				a[j] = &viot.Value{Value: &viot.Value_I32{I32: vl}}
			case int64:
				a[j] = &viot.Value{Value: &viot.Value_I64{I64: vl}}
			case float32:
				a[j] = &viot.Value{Value: &viot.Value_F32{F32: vl}}
			case float64:
				a[j] = &viot.Value{Value: &viot.Value_F64{F64: vl}}
			case string:
				a[j] = &viot.Value{Value: &viot.Value_String_{String_: vl}}
			case bool:
				a[j] = &viot.Value{Value: &viot.Value_Bool{Bool: vl}}
			case []byte:
				a[j] = &viot.Value{Value: &viot.Value_Binary{Binary: vl}}
			default:
				bytes, err := json.Marshal(vl)
				if err == nil {
					fmt.Printf("%t -> %s\n", vl, string(bytes))
					a[j] = &viot.Value{Value: &viot.Value_Json{Json: string(bytes)}}
				} else {
					fmt.Printf("encodeSelectResponse err = %s\n", err)
				}
			}
		}
		//fmt.Printf("row = %v\n", a)
		rows[i] = &viot.ResultRow{Value: a}
	}

	return &viot.SelectResult{Code: 0, Column: res.Columns, Row: rows}, nil
}

func encodeDeleteResponse(ctx context.Context, r any) (response any, err error) {
	return r, nil
}

func encodeUpdateResponse(ctx context.Context, r any) (response any, err error) {
	return r, nil
}

func encodeInsertResponse(ctx context.Context, r any) (response any, err error) {
	return r, nil
}

func encodeResolveResponse(ctx context.Context, r any) (response any, err error) {
	return r, nil
}
