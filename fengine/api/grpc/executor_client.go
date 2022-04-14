package grpc

import (
	"context"
	"time"

	pb "github.com/duclmse/fengine/pb"
	"github.com/go-kit/kit/endpoint"
	ot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var _ pb.FEngineExecutorClient = (*grpcExecutorClient)(nil)

type grpcExecutorClient struct {
	timeout time.Duration
	execute endpoint.Endpoint
}

func NewExecutorClient(conn *grpc.ClientConn, tracer opentracing.Tracer, timeout time.Duration) pb.FEngineExecutorClient {
	svcName := "pb.PricingService"

	return &grpcExecutorClient{
		timeout: timeout,
		execute: ot.TraceClient(tracer, "execute")(kitgrpc.NewClient(
			conn, svcName, "Execute", encodeRequest, decodeResponse, pb.Result{},
		).Endpoint()),
	}
}

func (client grpcExecutorClient) Execute(ctx context.Context, in *pb.Script, opts ...grpc.CallOption) (*pb.Result, error) {

	panic("implement me")
}

func encodeRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	return grpcReq.(*pb.Script), nil
}

func decodeResponse(ctx context.Context, grpcRes interface{}) (response interface{}, err error) {
	return grpcRes.(*pb.Result), nil
}
