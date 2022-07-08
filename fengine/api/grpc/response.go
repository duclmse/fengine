package grpc

import (
	"context"
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

	rows := []*viot.ResultRow{}
	for _, row := range res.Rows {
		//fmt.Printf("%d:", i)
		a := make([]*viot.Value, len(row))
		for j, value := range row {
			//fmt.Printf(" %d v=%v %t\n", j, value, value)
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
			}
		}
		//fmt.Printf("row = %v\n", a)
		rows = append(rows, &viot.ResultRow{Value: a})
	}

	return &viot.SelectResult{Column: res.Columns, Row: rows}, nil
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
