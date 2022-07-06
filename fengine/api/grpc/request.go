package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duclmse/fengine/fengine/db/sql"
	viot "github.com/duclmse/fengine/pb"
)

func decodeSelectRequest(ctx context.Context, r any) (request any, err error) {
	req, ok := r.(*viot.SelectRequest)
	fmt.Printf("decodeSelectRequest %v %+v\n", ok, r)
	if !ok {
		return nil, errors.New("request is not a select request")
	}
	filter := sql.Filter{}
	if err = json.Unmarshal([]byte(req.Filter), &filter); err != nil {
		return nil, err
	}
	request = sql.SelectRequest{
		Table:   req.Table,
		Fields:  req.FieldNames,
		Filter:  filter,
		GroupBy: req.GroupBy,
		Limit:   req.Limit,
		Offset:  req.Offset,
		OrderBy: req.OrderBy,
	}
	return
}

func decodeInsertRequest(ctx context.Context, r any) (request any, err error) {
	request, ok := r.(*viot.InsertRequest)
	if !ok {
		return nil, errors.New("request is not a select request")
	}
	return
}

func decodeUpdateRequest(ctx context.Context, r any) (request any, err error) {
	request, ok := r.(*viot.UpdateRequest)
	if !ok {
		return nil, errors.New("request is not a select request")
	}
	return
}

func decodeDeleteRequest(ctx context.Context, r any) (request any, err error) {
	request, ok := r.(*viot.DeleteRequest)
	if !ok {
		return nil, errors.New("request is not a select request")
	}
	return
}

func decodeResolveRequest(ctx context.Context, r any) (request any, err error) {
	return nil, nil
}
