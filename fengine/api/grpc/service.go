package grpc

import (
	"context"
	"errors"
	pb "github.com/duclmse/fengine/pb"
)

func (g grpcDataServer) Select(ctx context.Context, request *pb.SelectRequest) (*pb.SelectResult, error) {
	_, resp, err := g.selectEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*pb.SelectResult)
	if !ok {
		return nil, errors.New("GRPC Select: result is not valid result set")
	}
	return result, nil
}

func (g grpcDataServer) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResult, error) {
	_, resp, err := g.insertEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*pb.InsertResult)
	if !ok {
		return nil, errors.New("GRPC Insert: result is not valid")
	}
	return result, nil
}

func (g grpcDataServer) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResult, error) {
	_, resp, err := g.updateEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*pb.UpdateResult)
	if !ok {
		return nil, errors.New("GRPC Update: result is not valid")
	}
	return result, nil
}

func (g grpcDataServer) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResult, error) {
	_, resp, err := g.deleteEndpoint.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*pb.DeleteResult)
	if !ok {
		return nil, errors.New("GRPC Delete: result is not valid")
	}
	return result, nil
}

//func (g grpcThingServer) ResolveService(ctx context.Context, request *pb.ScriptRequest) (*pb.Result, error) {
//	_, resp, err := g.resolveEndpoint.ServeGRPC(ctx, request)
//	if err != nil {
//		return nil, err
//	}
//	return resp.(*pb.Result), nil
//}
