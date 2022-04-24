package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	ot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	pb "github.com/duclmse/fengine/pb"
)

var _ pb.FEngineExecutorClient = (*grpcExecutorClient)(nil)

func NewExecutorClient(conn *grpc.ClientConn, tracer opentracing.Tracer, config GrpcService) pb.FEngineExecutorClient {
	svcName := "viot.FEngineExecutor"
	return &grpcExecutorClient{
		timeout: time.Duration(config.Timeout) * time.Second,
		insertService: ot.TraceClient(tracer, "insert")(kitgrpc.NewClient(
			conn, svcName, "insertService", encodeExecuteRequest, decodeExecuteResponse, pb.Result{}).Endpoint()),
		updateService: ot.TraceClient(tracer, "update")(kitgrpc.NewClient(
			conn, svcName, "updateService", encodeExecuteRequest, decodeExecuteResponse, pb.Result{}).Endpoint()),
		deleteService: ot.TraceClient(tracer, "delete")(kitgrpc.NewClient(
			conn, svcName, "deleteService", encodeExecuteRequest, decodeExecuteResponse, pb.Result{}).Endpoint()),
		executeService: ot.TraceClient(tracer, "execute")(kitgrpc.NewClient(
			conn, svcName, "executeService", encodeExecuteRequest, decodeExecuteResponse, pb.Result{}).Endpoint()),
	}
}

type grpcExecutorClient struct {
	timeout        time.Duration
	insertService  endpoint.Endpoint
	updateService  endpoint.Endpoint
	deleteService  endpoint.Endpoint
	executeService endpoint.Endpoint
}

func (client grpcExecutorClient) AddService(ctx context.Context, in *pb.ThingMethod, opts ...grpc.CallOption) (*pb.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	res, err := client.insertService(ctx, *in)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	result := res.(*pb.Result)
	return result, nil
}

func (client grpcExecutorClient) UpdateService(ctx context.Context, in *pb.ThingMethod, opts ...grpc.CallOption) (*pb.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	res, err := client.updateService(ctx, *in)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	result := res.(*pb.Result)
	return result, nil
}

func (client grpcExecutorClient) DeleteService(ctx context.Context, in *pb.ThingMethod, opts ...grpc.CallOption) (*pb.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	res, err := client.deleteService(ctx, *in)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	result := res.(*pb.Result)
	return result, nil
}

func (client grpcExecutorClient) Execute(ctx context.Context, in *pb.Script, opts ...grpc.CallOption) (*pb.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	res, err := client.executeService(ctx, *in)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	result := res.(*pb.Result)
	return result, nil
}

func encodeExecuteRequest(ctx context.Context, grpcReq any) (request any, err error) {
	//script := grpcReq.(pb.Script)
	return &pb.Script{}, nil
}

func decodeExecuteResponse(ctx context.Context, grpcRes any) (response any, err error) {
	return grpcRes, nil
}
