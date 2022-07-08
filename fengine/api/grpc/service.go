package grpc

import (
	"context"
	"fmt"
	pb "github.com/duclmse/fengine/pb"
)

func (g grpcDataServer) Select(ctx context.Context, request *pb.SelectRequest) (*pb.SelectResult, error) {
	_, resp, err := g.selectEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	set, ok := resp.(*pb.SelectResult)
	if !ok {
		fmt.Printf("GRPC Select: result is not result set")
		return nil, err
	}
	return set, nil
}

func (g grpcDataServer) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResult, error) {
	_, resp, err := g.insertEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.InsertResult), nil
}

func (g grpcDataServer) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResult, error) {
	_, resp, err := g.updateEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UpdateResult), nil
}

func (g grpcDataServer) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResult, error) {
	_, resp, err := g.deleteEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteResult), nil
}

//func (g grpcThingServer) ResolveService(ctx context.Context, request *pb.ScriptRequest) (*pb.Result, error) {
//	_, resp, err := g.resolveEndpoint.ServeGRPC(ctx, request)
//	if err != nil {
//		return nil, err
//	}
//	return resp.(*pb.Result), nil
//}
