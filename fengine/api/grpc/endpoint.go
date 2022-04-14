package grpc

import (
	. "context"
	. "github.com/duclmse/fengine/fengine"
	. "github.com/go-kit/kit/endpoint"
)

func grpcSelect(svc Service) Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func grpcInsert(svc Service) Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func grpcUpdate(svc Service) Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func grpcDelete(svc Service) Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func grpcResolve(svc Service) Endpoint {
	return func(ctx Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
