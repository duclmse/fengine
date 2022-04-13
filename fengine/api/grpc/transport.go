package grpc

import (
	"context"
	"google.golang.org/grpc"

	"github.com/duclmse/fengine/fengine"
	. "github.com/duclmse/fengine/pb"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ FEngineExecutorClient = (*grpcServer)(nil)

type grpcServer struct {
	identityNameInfo kitgrpc.Handler
}

func NewServer(tracer opentracing.Tracer, svc fengine.Service) FEngineExecutorClient {
	return &grpcServer{
		identityNameInfo: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Execute")(grpcGet(svc)),
			decodeIdentifyNameRequest,
			encodeIdentifyNameResponse,
		),
	}
}

func (g grpcServer) Execute(ctx context.Context, in *Script, opts ...grpc.CallOption) (*Result, error) {
	//TODO implement me
	panic("implement me")
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
