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
		req, ok := request.(sql.InsertRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		return svc.Insert(ctx, req)
	}
}

func grpcUpdate(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		req, ok := request.(sql.UpdateRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		return svc.Update(ctx, req)
	}
}

func grpcDelete(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		req, ok := request.(sql.DeleteRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		return svc.Delete(ctx, req)
	}
}

func grpcResolve(svc Service) Endpoint {
	return func(ctx Context, request any) (any, error) {
		selectRequest, ok := request.(sql.SelectRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		return svc.Select(ctx, selectRequest)
	}
}
