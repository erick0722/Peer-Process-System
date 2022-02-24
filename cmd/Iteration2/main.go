// =============================================================
/*
	CPSC 559 - Iteration 2
	main.go

	Erick Yip
	Chris Chen
*/

package main

import (
	"559Project/pkg/peer"
	"559Project/pkg/registry"
	"fmt"
	"sync"
)

func main() {

	var regAddr, peerAddr string
	var wg sync.WaitGroup

	//ask for the registry and peer process' address
	fmt.Println("Please enter the registry's address: ")
	fmt.Scanln(&regAddr)
	// regAddr = "localhost:55921"

	fmt.Println("Please enter the peer process's address: ")
	fmt.Scanln(&peerAddr)
	// peerAddr = "localhost:3000"

	// Connect to the server via TCP
	wg.Add(2)

	go func() {
		registry.InitRegistryCommunicator(regAddr, peerAddr)
		fmt.Println("Registry Communicator exited")
		wg.Done()
	}()

	go func() {
		peer.InitPeerProcess(peerAddr)
		fmt.Println("Peer Process exited")
		wg.Done()
	}()

	wg.Wait()
}

// =============================================================
