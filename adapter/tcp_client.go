package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	c, err := net.Dial("tcp", arguments[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		_, err := fmt.Fprintf(c, text+"\n")
		if err != nil {
			fmt.Printf("err handling request %v\n", err)
			return
		}

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Printf(">: %s\n", message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
