package http

import (
	"context"
	"github.com/google/uuid"

	"github.com/go-kit/kit/endpoint"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/pkg/errors"
)

func getAllServicesEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func getServiceEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//c.DB.
		req, ok := request.(serviceRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		id, err := uuid.Parse(req.thingId)
		if err != nil {
			return nil, errors.New("invalid input")
		}
		service, err := svc.GetThingService(ctx, id, req.serviceName)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func execServiceEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		script, ok := request.(fengine.JsonScript)
		if !ok {
			return nil, errors.New("invalid input")
		}
		c.Log.Info("Received script %v", script)
		result, err := svc.ExecuteService(context.Background(), &script)
		if err != nil {
			c.Log.Error("Error in executing %s\n", err.Error())
			return nil, err
		}
		c.Log.Info("Done: %v", result)
		return nil, nil
	}
}
