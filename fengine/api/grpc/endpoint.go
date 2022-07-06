package grpc

import (
	. "context"
	"errors"
	. "github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/db/sql"
	. "github.com/go-kit/kit/endpoint"
)

func grpcSelect(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		selectRequest, ok := request.(sql.SelectRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		return svc.Select(ctx, selectRequest)
	}
}

func grpcInsert(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		return nil, nil
	}
}

func grpcUpdate(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		return nil, nil
	}
}

func grpcDelete(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		return nil, nil
	}
}

func grpcResolve(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		return nil, nil
	}
}
