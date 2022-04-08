package main

import (
	"fmt"
	"os"

	"github.com/duclmse/fengine/adapter/api/tcp"
	"github.com/duclmse/fengine/adapter/api/udp"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	fmt.Printf("args: %v\n", arguments[1:])

	err := make(chan error, 2)
	go tcp.StartTCP(arguments[1], err)
	go udp.StartUDP(arguments[2], err)
	oops := <-err
	if oops != nil {
		fmt.Printf("> Terminating adapter service due to error! %v", oops)
	} else {
		fmt.Printf("> Terminating adapter service!")
	}
}
