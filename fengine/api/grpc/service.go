package grpc

import (
	"context"
	pb "github.com/duclmse/fengine/pb"
)

func (g grpcDataServer) Select(ctx context.Context, request *pb.SelectRequest) (*pb.ResultSet, error) {
	_, resp, err := g.selectEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ResultSet), nil
}

func (g grpcDataServer) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.Result, error) {
	_, resp, err := g.insertEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Result), nil
}

func (g grpcDataServer) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.Result, error) {
	_, resp, err := g.updateEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Result), nil
}

func (g grpcDataServer) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.Result, error) {
	_, resp, err := g.deleteEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Result), nil
}

func (g grpcThingServer) ResolveService(ctx context.Context, request *pb.ScriptRequest) (*pb.Result, error) {
	_, resp, err := g.resolveEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Result), nil
}
