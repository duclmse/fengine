package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/go-zoo/bone"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
)

//region Data structure
type entityRequest struct {
	thingId string
}

//endregion Data structure

func decodeEntityRequest(ctx context.Context, request *http.Request) (any, error) {
	thingId := bone.GetValue(request, "id")
	return entityRequest{thingId: thingId}, nil
}

func decodeUpsertEntityRequest(ctx context.Context, request *http.Request) (any, error) {
	uid, err := uuid.Parse(bone.GetValue(request, "id"))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	def := &sql.EntityDefinition{}
	err = json.Unmarshal(body, def)
	if err != nil {
		return nil, err
	}
	def.Id = uid
	return def, nil
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
