/*
	CPSC 559 - Iteration 2
	main.go

	Erick Yip
	Chris Chen
*/

package main

import (
	"559Project/pkg/registry"
	"559Project/pkg/sock"
	"fmt"
	"os"
	"sync"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Missing <server address>")
		os.Exit(1)
	}

	// Get the server's address from the command line
	address := os.Args[1]

	// Connect to the server via TCP
	tcpConn := sock.InitializeTcpClient(address)
	fmt.Printf("Connected to server at %s\n", address)

	//udpConn := sock.InitializeUdpServer("136.159.5.22:8722")
	//fmt.Printf("Listening for UDP packets on %s\n", address)
	var peerGroup sync.WaitGroup
	peerGroup.Add(2)
	go func() {
		registry.RegistryCommunicator(address, tcpConn)
		fmt.Println("Registry Communicator exited")
		peerGroup.Done()
	}()

	go func() {

		peerGroup.Done()
	}()
	peerGroup.Wait()
}
