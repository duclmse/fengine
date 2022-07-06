package grpc_test

import (
	"context"
	"fmt"
	"github.com/duclmse/fengine/cmd/fengine"
	. "github.com/duclmse/fengine/pb"
	logger "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestExecutorClient(t *testing.T) {
	config := main.LoadConfig("./.env", "executor")
	log, err := logger.New(os.Stdout, config.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	client := main.ConnectToGrpcService("executor", config, log)
	//fmt.Printf("%+v\n", client)
	execConfig := config.GrpcServices["executor"]
	//fmt.Printf("%+v\n")
	serviceTracer, dbCloser := main.InitJaeger("fengine", config.JaegerURL, log)
	defer viot.Close(log, "vtfengine_db")(dbCloser)

	executorClient := NewExecutorClient(client, serviceTracer, execConfig)
	id, _ := uuid.MustParse(``).MarshalBinary()
	execute, err := executorClient.Execute(context.Background(), &Script{
		Function: &MethodId{
			ThingID:    id,
			MethodName: "",
			Type:       MethodType_service,
		},
		Attributes: []*Variable{
			{Name: "s", Value: &Variable_String_{String_: "string"}},
			{Name: "i", Value: &Variable_I32{I32: 100}},
		},
		Services: map[string]*Function{
			"test": {
				Input: []*Parameter{
					{Name: "str", Type: Type_string},
					{Name: "i32", Type: Type_i32},
				},
				Output: Type_json,
				Code:   `return {i32: i32+1, str: str+'!'}`,
			},
		},
	})
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}

	value := execute.Output.Value
	fmt.Printf("value %v\n", value)
}
