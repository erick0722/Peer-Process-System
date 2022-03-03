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

	fmt.Println("Please enter the registry's address: ")
	fmt.Scanln(&regAddr)

	fmt.Println("Please enter the peer process's address: ")
	fmt.Scanln(&peerAddr)

	wg.Add(2)

	// Start a thread to communicate with the Registry 
	go func() {
		defer wg.Done()
		registry.InitRegistryCommunicator(regAddr, peerAddr, ctx)
		fmt.Println("Registry Communicator exited")
	}()

	// Start a thread to create our peer process
	go func() {
		defer wg.Done()
		peer.InitPeerProcess(peerAddr, ctx)
		fmt.Println("Peer Process exited, connecting to the registry again")
		registry.InitRegistryCommunicator(regAddr, peerAddr, ctx)
		cancel()
	}()

	// Shut our program down gracefully when CTRL+C is pressed or is interrupted 
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
