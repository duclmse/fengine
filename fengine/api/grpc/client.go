package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/duclmse/fengine/pb"
)

var _ pb.FEngineServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	timeout time.Duration
	grpcGet endpoint.Endpoint
}

func NewClient(conn *grpc.ClientConn, tracer opentracing.Tracer, timeout time.Duration) pb.FEngineServiceClient {
	svcName := "pb.PricingService"

	return &grpcClient{
		timeout: timeout,
		grpcGet: kitot.TraceClient(tracer, "identify_name")(kitgrpc.NewClient(
			conn, svcName, "GrpcGet", encodeGetRequest, decodeGetResponse, pb.Result{},
		).Endpoint()),
	}
}

func encodeGetRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	return grpcReq.(*pb.ID), nil
}

func decodeGetResponse(ctx context.Context, grpcRes interface{}) (response interface{}, err error) {
	return grpcRes.(*pb.Info), nil
}

func (client grpcClient) GrpcGet(ctx context.Context, id *pb.ID, opts ...grpc.CallOption) (*pb.Info, error) {
	fmt.Printf("id=%v, opts=%v\n", *id, opts)
	return nil, nil
}
