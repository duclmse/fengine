package grpc

import (
	. "context"
	"github.com/duclmse/fengine/fengine"
	"github.com/go-kit/kit/endpoint"
)

func grpcGet(svc fengine.Service) endpoint.Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
