package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/go-zoo/bone"
	"github.com/google/uuid"
	"io"
	"net/http"
)

//region Data structure

// contextKeyType is a private struct that is used for storing bone values in net.Context
type contextKeyType struct{}

// contextKey is the key that is used to store bone values in the net.Context for each request
var contextKey = contextKeyType{}

//endregion Data structure

type allServiceRequest struct {
	thingId string
}

func decodeAllServiceRequest(ctx context.Context, request *http.Request) (any, error) {
	values := bone.GetAllValues(request)
	thingId := values["id"]
	return allServiceRequest{thingId: thingId}, nil
}

func decodeServiceRequest(ctx context.Context, request *http.Request) (any, error) {
	values := bone.GetAllValues(request)
	thingId, err := uuid.Parse(values["id"])
	if err != nil {
		return nil, err
	}
	serviceName := values["service"]
	return sql.ThingServiceId{EntityId: thingId, Name: serviceName}, nil
}

func decodeExecRequest(ctx context.Context, request *http.Request) (any, error) {
	values := bone.GetAllValues(request)
	thingId := values["id"]
	serviceName := values["service"]
	fmt.Printf("/thing/%s/service/%s\n", thingId, serviceName)
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var execution sql.ServiceRequest
	err = json.Unmarshal(body, &execution)
	if err != nil {
		fmt.Printf("error in decode exec: %s\n", err)
		return nil, err
	}

	return execution, nil
}
