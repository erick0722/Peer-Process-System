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
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	var regAddr, peerAddr string
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

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
		defer wg.Done()
		registry.InitRegistryCommunicator(regAddr, peerAddr, ctx)
		fmt.Println("Registry Communicator exited")
	}()

	go func() {
		defer wg.Done()
		peer.InitPeerProcess(peerAddr, ctx, cancel)
		fmt.Println("Peer Process exited")
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(s, os.Interrupt, os.Kill)
	select {
	case <-s:
		fmt.Println("\nReceived SIGINT/SIGTERM. Exiting gracefully...")
		cancel()
	case <-ctx.Done():
		fmt.Println("\nContext cancelled. Exiting gracefully...")
	}
	wg.Wait()
}

// =============================================================
