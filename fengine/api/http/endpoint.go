package http

import (
	"context"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/go-kit/kit/endpoint"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/pkg/errors"
)

func getAllServicesEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(allServiceRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		service, err := svc.GetThingAllServices(ctx, req.thingId)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func getServiceEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(sql.ThingServiceId)
		if !ok {
			return nil, errors.New("invalid input")
		}
		service, err := svc.GetThingService(ctx, req)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func execServiceEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		script, ok := request.(sql.ServiceRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		c.Log.Info("Received script %v", script)
		result, err := svc.ExecuteService(ctx, script)
		if err != nil {
			c.Log.Error("Error in executing %s\n", err.Error())
			return nil, err
		}
		c.Log.Info("Done: %v", result)
		return nil, nil
	}
}
