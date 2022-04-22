package grpc

import (
	"github.com/duclmse/fengine/fengine"
	. "github.com/duclmse/fengine/pb"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
)

var (
	_ FEngineDataServer  = (*grpcDataServer)(nil)
	_ FEngineThingServer = (*grpcThingServer)(nil)
)

type grpcDataServer struct {
	selectEndpoint kitgrpc.Handler
	insertEndpoint kitgrpc.Handler
	updateEndpoint kitgrpc.Handler
	deleteEndpoint kitgrpc.Handler
}

type grpcThingServer struct {
	resolveEndpoint *kitgrpc.Server
}

type GrpcService struct {
	URL     string
	Timeout int
}

func NewDataServer(tracer opentracing.Tracer, svc fengine.Service) FEngineDataServer {
	return &grpcDataServer{
		selectEndpoint: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Select")(grpcSelect(svc)),
			decodeSelectRequest, encodeSelectResponse),
		insertEndpoint: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Insert")(grpcInsert(svc)),
			decodeInsertRequest, encodeInsertResponse),
		updateEndpoint: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Update")(grpcUpdate(svc)),
			decodeUpdateRequest, encodeUpdateResponse),
		deleteEndpoint: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Delete")(grpcDelete(svc)),
			decodeDeleteRequest, encodeDeleteResponse),
	}
}

func NewThingServer(tracer opentracing.Tracer, svc fengine.Service) FEngineThingServer {
	return &grpcThingServer{
		resolveEndpoint: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "Resolve")(grpcResolve(svc)),
			decodeResolveRequest,
			encodeResolveResponse,
		),
	}
}
