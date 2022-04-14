package grpc

import (
	"context"
	. "github.com/duclmse/fengine/pb"
)

func (g grpcDataServer) Select(ctx context.Context, request *SelectRequest) (*Script, error) {
	_, resp, err := g.selectEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Script), nil
}

func (g grpcDataServer) Insert(ctx context.Context, request *InsertRequest) (*Script, error) {
	_, resp, err := g.insertEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Script), nil
}

func (g grpcDataServer) Update(ctx context.Context, request *UpdateRequest) (*Script, error) {
	_, resp, err := g.updateEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Script), nil
}

func (g grpcDataServer) Delete(ctx context.Context, request *DeleteRequest) (*Script, error) {
	_, resp, err := g.deleteEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Script), nil
}

func (g grpcThingServer) ResolveService(ctx context.Context, request *ScriptRequest) (*Script, error) {
	_, resp, err := g.resolveEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Script), nil
}
