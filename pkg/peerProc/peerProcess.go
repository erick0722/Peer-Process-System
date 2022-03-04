// =============================================================
/*
	CPSC 559 - Iteration 2
	peerProcess.go

	Erick Yip
	Chris Chen
*/

package peerProc

import (
	"559Project/pkg/sock"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Struct for storing the peer information
type peerStruct struct {
	address   string
	source    string
	lastHeard time.Time
}

// Struct for storing the peer receiving events
type receivedEvent struct {
	received     string
	source       string
	timeReceived time.Time
}

// Struct for storing the peers sent
type sentEvent struct {
	sentTo   string
	peer     string
	timeSent time.Time
}

// Struct for storing the snips being received
type snip struct {
	content   string
	timeStamp string
	source    string
}

var peerList []peerStruct
var recievedPeers []receivedEvent
var peersSent []sentEvent
var snipList []snip

var peerProcessAddr string // Address of the peer process (UDP)
var currTimeStamp int = 0  // Local logical time stamp
var mutex sync.Mutex

// Start the peer process and the threads
func InitPeerProcess(address string, ctx context.Context) {
	peerProcessAddr = address

	var wg sync.WaitGroup
	peerCtx, cancel := context.WithCancel(ctx)

	conn := sock.InitializeUdpServer(address)
	fmt.Printf("Peer process started at %s\n", peerProcessAddr)

	// Start 4 threads
	wg.Add(4)
	go func() {
		defer wg.Done()
		// Handle incoming messages from other peer processes
		handleMessage(conn, peerCtx, cancel)
		fmt.Printf("Exiting handleMessage\n")
	}()

	go func() {
		defer wg.Done()
		// Monitor and read message typed from the console
		readSnip(conn, peerCtx)
		fmt.Printf("Exiting readSnip\n")
	}()

	go func() {
		defer wg.Done()
		// Periodically send the peer list to other peers
		sendPeer(conn, peerCtx)
		fmt.Printf("Exiting sendpeerList\n")
	}()

	go func() {
		defer wg.Done()
		// Periodically check for inactive peers
		checkInactivePeers(peerCtx)
		fmt.Printf("Exiting checkInactivePeers\n")
	}()

	wg.Wait()
}

//
/**
*	Process the different messages our peer process receives from other peers in the system
*
*	@param address {string} The address of the peer who sent the message
*	@param ctx {context.Context} The context of our app, used to stop the other threads / gracefully exit the program
*	@param cancel {context.CancelFunc} The function used to initiate the cancel process for our context
 */
func handleMessage(conn *net.UDPConn, ctx context.Context, cancel context.CancelFunc) {

	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	for {
		select {
		// If the context is cancelled, exit the thread
		case <-ctx.Done():
			return
		default:
			// Receive message from other peer process
			msg, addr, err := sock.ReceiveUdpMessage(conn)
			if err != nil {
				fmt.Printf("Error detected: %v\n", err)
				continue
			}

			// Check if message is at least 4 characters
			if len(msg) < 4 {
				fmt.Printf("Message invalid")
				continue
			}

			// Update the sender's last heard time if they are in the list
			index := searchPeerList(addr)
			if index != -1 {
				peerList[index].lastHeard = time.Now()
			}

			switch string(msg[0:4]) {
			case "snip":
				// Handle snip message
				source := strings.TrimSuffix(string(msg[4:]), "\n")
				storeSnip(source, addr)
			case "peer":
				// Handle peer message
				source := strings.TrimSpace(strings.TrimSuffix(string(msg[4:]), "\n"))
				addPeer(source, addr)
			case "stop":
				// Handle stop message
				fmt.Printf("Received stop command, exiting...\n")
				conn.Close()
				cancel() // Stop all our other running threads when we get a "stop" message
				return
			}
		}
	}
}

// =============================================================
