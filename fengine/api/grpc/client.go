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
	timeout time.Duration
	grpcGet endpoint.Endpoint
}

func NewClient(conn *grpc.ClientConn, tracer opentracing.Tracer, timeout time.Duration) pb.FEngineDataClient {
	svcName := "pb.PricingService"

	return &grpcDataClient{
		timeout: timeout,
		grpcGet: kitot.TraceClient(tracer, "identify_name")(kitgrpc.NewClient(
			conn, svcName, "GrpcGet", encodeGetRequest, decodeGetResponse, pb.Result{},
		).Endpoint()),
	}
}

func (client grpcDataClient) Select(ctx context.Context, in *pb.SelectRequest, opts ...grpc.CallOption) (*pb.Script, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Insert(ctx context.Context, in *pb.InsertRequest, opts ...grpc.CallOption) (*pb.Script, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Update(ctx context.Context, in *pb.UpdateRequest, opts ...grpc.CallOption) (*pb.Script, error) {
	//TODO implement me
	panic("implement me")
}

func (client grpcDataClient) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.Script, error) {
	//TODO implement me
	panic("implement me")
}

func encodeGetRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	return nil, nil
}

func decodeGetResponse(ctx context.Context, grpcRes interface{}) (response interface{}, err error) {
	return nil, nil
}
