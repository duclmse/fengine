package grpc

import (
	"context"
	. "github.com/duclmse/fengine/pb"
)

func (g grpcDataServer) Select(ctx context.Context, request *SelectRequest) (*Result, error) {
	_, resp, err := g.selectEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Result), nil
}

func (g grpcDataServer) Insert(ctx context.Context, request *InsertRequest) (*Result, error) {
	_, resp, err := g.insertEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Result), nil
}

func (g grpcDataServer) Update(ctx context.Context, request *UpdateRequest) (*Result, error) {
	_, resp, err := g.updateEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Result), nil
}

func (g grpcDataServer) Delete(ctx context.Context, request *DeleteRequest) (*Result, error) {
	_, resp, err := g.deleteEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Result), nil
}

func (g grpcThingServer) ResolveService(ctx context.Context, request *ScriptRequest) (*Result, error) {
	_, resp, err := g.resolveEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Result), nil
}
