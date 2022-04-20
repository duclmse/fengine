package http

import (
	"context"
	"github.com/duclmse/fengine/pkg/errors"

	"github.com/duclmse/fengine/fengine"
	. "github.com/duclmse/fengine/pb"
	"github.com/go-kit/kit/endpoint"
)

func getServiceEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func execEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if script, ok := request.(Script); ok {
			c.Log.Info("Received script %v", script)
			result, err := svc.Execute(context.Background(), &script)
			if err != nil {
				c.Log.Error("Error in executing %s\n", err.Error())
				return nil, err
			}
			c.Log.Info("Done: %v", result.Output)
			return map[string]interface{}{
				"in":  script,
				"out": result.Output.GetValue(),
			}, nil
		}
		c.Log.Error("invalid input")
		return nil, errors.New("invalid input")
	}
}
