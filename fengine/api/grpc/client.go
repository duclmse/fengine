package grpc

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	pb "github.com/duclmse/fengine/pb"
)

var _ pb.FEngineDataClient = (*grpcDataClient)(nil)

type grpcDataClient struct {
	timeout    time.Duration
	grpcSelect endpoint.Endpoint
	grpcInsert endpoint.Endpoint
	grpcUpdate endpoint.Endpoint
	grpcDelete endpoint.Endpoint
}

func NewClient(conn *grpc.ClientConn, tracer opentracing.Tracer, timeout time.Duration) pb.FEngineDataClient {
	svcName := "pb.FEngineData"

	return &grpcDataClient{
		timeout: timeout,
		grpcSelect: kitot.TraceClient(tracer, "fengine_select")(kitgrpc.NewClient(
			conn, svcName, "grpcSelect", encodeGetRequest, decodeGetResponse, pb.Result{}).Endpoint()),
		grpcInsert: kitot.TraceClient(tracer, "fengine_insert")(kitgrpc.NewClient(
			conn, svcName, "grpcInsert", encodeGetRequest, decodeGetResponse, pb.Result{}).Endpoint()),
		grpcUpdate: kitot.TraceClient(tracer, "fengine_update")(kitgrpc.NewClient(
			conn, svcName, "grpcUpdate", encodeGetRequest, decodeGetResponse, pb.Result{}).Endpoint()),
		grpcDelete: kitot.TraceClient(tracer, "fengine_delete")(kitgrpc.NewClient(
			conn, svcName, "grpcDelete", encodeGetRequest, decodeGetResponse, pb.Result{}).Endpoint()),
	}
}

func (client grpcDataClient) Select(ctx context.Context, in *pb.SelectRequest, opts ...grpc.CallOption) (*pb.ResultSet, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Insert(ctx context.Context, in *pb.InsertRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Update(ctx context.Context, in *pb.UpdateRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	//TODO implement me
	panic("implement me")
}

func encodeGetRequest(ctx context.Context, grpcReq any) (request any, err error) {
	return nil, nil
}

func decodeGetResponse(ctx context.Context, grpcRes any) (response any, err error) {
	return nil, nil
}
