package udp

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func StartUDP(port string, errChan chan error) {
	addr, err := net.ResolveUDPAddr("udp4", ":"+port)
	if err != nil {
		fmt.Printf("[FATAL] Cannot resolve UDP port :%s > %v\n", port, err)
		errChan <- err
		return
	}

	connection, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Printf("[FATAL] Cannot listen to UDP port :%s > %v\n", port, err)
		errChan <- err
		return
	}

	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			fmt.Printf("Cannot close TCP adapter: %s\n", err.Error())
			errChan <- err
		}
	}(connection)
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		data := strconv.Itoa(rand.Intn(1000))
		fmt.Printf("data: %s\n", data)
		_, err = connection.WriteToUDP([]byte(data), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
