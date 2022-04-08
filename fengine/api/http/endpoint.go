package http

import (
	"context"

	"github.com/duclmse/fengine/fengine"
	"github.com/go-kit/kit/endpoint"
)

func getEndpoint(svc fengine.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return Response{Code: 0, Msg: "Okie", Data: "Done!"}, nil
	}
}
