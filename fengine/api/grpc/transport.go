package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/pb"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.FEngineServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	identityNameInfo kitgrpc.Handler
}

func NewServer(tracer opentracing.Tracer, svc fengine.Service) pb.FEngineServiceServer {
	return &grpcServer{
		identityNameInfo: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Get")(grpcGet(svc)),
			decodeIdentifyNameRequest,
			encodeIdentifyNameResponse,
		),
	}
}

func (gs grpcServer) GrpcGet(ctx context.Context, id *pb.ID) (*pb.Info, error) {
	fmt.Printf("Executing GrpcGet...")
	return nil, encodeError(errors.New("not implemented"))
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
