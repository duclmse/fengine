package http

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/duclmse/fengine/pkg/errors"
)

func getEntityEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(entityRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		service, err := svc.GetEntity(ctx, req.thingId)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func upsertEntityEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(sql.EntityDefinition)
		if !ok {
			return nil, errors.New("invalid input")
		}
		service, err := svc.UpsertEntity(ctx, req)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func deleteEntityEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(entityRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		service, err := svc.DeleteEntity(ctx, req.thingId)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
}

func getAllServicesEndpoint(svc fengine.Service, c fengine.ServiceComponent) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(entityRequest)
		if !ok {
			return nil, errors.New("invalid input")
		}
		id, err := uuid.Parse(req.thingId)
		if err != nil {
			return nil, errors.New("thing id is not a valid uuid")
		}
		service, err := svc.GetThingAllServices(ctx, id)
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
