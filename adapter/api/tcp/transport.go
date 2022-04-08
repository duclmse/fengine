package tcp

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func StartTCP(port string, errChan chan error) {
	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		fmt.Printf("[FATAL] Cannot listen to TCP port :%s > %v\n", port, err)
		errChan <- err
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Printf("Cannot close TCP adapter: %s\n", err.Error())
			errChan <- err
		}
	}(listener)

	rand.Seed(time.Now().Unix())
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("FATAL: %s\n", err.Error())
			return
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	client := connection.RemoteAddr().String()
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			fmt.Printf("[FATAL] Cannot close TCP connection > %v\n", err.Error())
		}
	}(connection)

	fmt.Printf("Serving %s\n", client)
	for {
		netData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(netData)
		if temp == "STOP" {
			break
		}

		fmt.Printf("%s: %s\n", client, temp)
		result := strconv.Itoa(rand.Intn(100)) + "\n"
		connection.Write([]byte(result))
	}
}
