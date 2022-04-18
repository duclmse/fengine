package grpc

import (
	"context"
	"fmt"
	"time"

	. "github.com/duclmse/fengine/pb"
	"github.com/go-kit/kit/endpoint"
	ot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var _ FEngineExecutorClient = (*grpcExecutorClient)(nil)

type grpcExecutorClient struct {
	timeout time.Duration
	execute endpoint.Endpoint
}

func NewExecutorClient(conn *grpc.ClientConn, tracer opentracing.Tracer, config GrpcService) FEngineExecutorClient {
	svcName := "viot.FEngineExecutor"
	return &grpcExecutorClient{
		timeout: time.Duration(config.Timeout) * time.Second,
		execute: ot.TraceClient(tracer, "execute")(kitgrpc.NewClient(
			conn, svcName, "Execute", encodeExecuteRequest, decodeExecuteResponse, Result{},
		).Endpoint()),
	}
}

func (client grpcExecutorClient) Execute(ctx context.Context, in *Script, opts ...grpc.CallOption) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	res, err := client.execute(ctx, *in)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	result := res.(*Result)
	return result, nil
}

func encodeExecuteRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	script := grpcReq.(Script)
	return &Script{
		Function:   script.Function,
		Attributes: script.Attributes,
		Referee:    script.Referee,
	}, nil
}

func decodeExecuteResponse(ctx context.Context, grpcRes interface{}) (response interface{}, err error) {
	return grpcRes, nil
}
