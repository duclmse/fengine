package grpc

import (
	"context"
	"fmt"
	"github.com/duclmse/fengine/cmd/fengine"
	. "github.com/duclmse/fengine/pb"
	logger "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/viot"
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
	execute, err := executorClient.Execute(context.Background(), &Script{
		Function: &Function{
			Input: []*Variable{
				{Name: "s", Value: &Variable_String_{String_: "string"}},
				{Name: "i", Value: &Variable_I32{I32: 100}},
			},
			Output: []*Variable{
				{Name: "i"},
			},
			Code: `
				me.test({s,i});
				Table('a').Select({and:[{a:{$gt:10,$lt:20}}]});
				me.i=0;
				return {i:i+me.i, s:s+me.s}
			`,
		},
		Attributes: []*Variable{
			{Name: "s", Value: &Variable_String_{String_: "string"}},
			{Name: "i", Value: &Variable_I32{I32: 100}},
		},
		Referee: map[string]*Function{
			"test": {
				Input: []*Variable{
					{Name: "str", Value: &Variable_String_{String_: "hello"}},
					{Name: "i32", Value: &Variable_I32{I32: 200}},
				},
				Output: []*Variable{{Name: "i"}},
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
