package grpc

import (
	"context"
	"fmt"
	"github.com/duclmse/fengine/fengine/db/sql"
	viot "github.com/duclmse/fengine/pb"
)

func encodeSelectResponse(ctx context.Context, r any) (response any, err error) {
	fmt.Printf("encodeSelectResponse %t\n", r)
	res, ok := r.([]map[string]sql.Variable)
	if !ok {
		return &viot.ResultSet{}, err
	}
	//res.Data.()
	rows := []*viot.ResultRow{}
	for _, v := range res {
		a := []*viot.Value{}
		for _, value := range v {
			switch vl := value.Value.(type) {
			case int32:
				a = append(a, &viot.Value{Value: &viot.Value_I32{I32: vl}})
			case int64:
				a = append(a, &viot.Value{Value: &viot.Value_I64{I64: vl}})
			case float32:
				a = append(a, &viot.Value{Value: &viot.Value_F32{F32: vl}})
			case float64:
				a = append(a, &viot.Value{Value: &viot.Value_F64{F64: vl}})
			case string:
				a = append(a, &viot.Value{Value: &viot.Value_String_{String_: vl}})
				//case json:
				//	a = append(a, &viot.Value{Value: &viot.Value_I64{I64: vl}})
			}
		}
		rows = append(rows, &viot.ResultRow{Values: a})
	}

	response = &viot.ResultSet{
		ColumnNames: nil,
		Rows:        rows,
	}
	return response, nil
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
