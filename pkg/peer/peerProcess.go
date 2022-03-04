// =============================================================
/*
	CPSC 559 - Iteration 2
	peerFunc.go

	Erick Yip
	Chris Chen
*/

package peer

import (
	"559Project/pkg/sock"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type peerStruct struct {
	address   string
	source    string
	lastHeard time.Time
}

type receivedEvent struct {
	received     string
	source       string
	timeReceived time.Time
}

type sentEvent struct {
	sentTo   string
	peer     string
	timeSent time.Time
}

type snip struct {
	content   string
	timeStamp string
	source    string
}

var peerList []peerStruct
var recievedPeers []receivedEvent
var peersSent []sentEvent
var snipList []snip

var peerProcessAddr string
var currTimeStamp int = 0
var mutex sync.Mutex

func InitPeerProcess(address string, ctx context.Context) {
	// Setup our peer process: add ourselves to our peerlist and configure our threads,
	peerProcessAddr = address
	fmt.Printf("Peer process started at %s\n", peerProcessAddr)
	peerList = append(peerList, peerStruct{peerProcessAddr, peerProcessAddr, time.Now()})
	var wg sync.WaitGroup
	peerCtx, cancel := context.WithCancel(ctx)

	conn := sock.InitializeUdpServer(address)

	wg.Add(4)
	go func() {
		defer wg.Done()
		handleMessage(conn, peerCtx, cancel)
		fmt.Printf("Exiting handleMessage\n")
	}()

	go func() {
		defer wg.Done()
		readSnip(conn, peerCtx)
		fmt.Printf("Exiting readSnip\n")
	}()

	go func() {
		defer wg.Done()
		sendPeerList(conn, peerCtx)
		fmt.Printf("Exiting sendpeerList\n")
	}()

	go func() {
		defer wg.Done()
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
		case <-ctx.Done():
			return
		default:
			//fmt.Printf("Waiting for message\n")
			msg, addr, err := sock.ReceiveUdpMessage(conn)
			//fmt.Println("Received ", msg, " from ", addr)
			if err != nil {
				fmt.Printf("Error detected: %v\n", err)
				continue
			}

			if len(msg) < 4 {
				fmt.Printf("Message invalid")
				continue
			}

			// Peers are updated in our own peerlist when we receive messages from them
			index := searchPeerList(addr)
			if index != -1 {
				peerList[index].lastHeard = time.Now()
			}

			switch string(msg[0:4]) {
			case "stop":
				fmt.Printf("Received stop command, exiting...\n")
				conn.Close()
				cancel() // Stop all our other running threads when we get a "stop" message
				return
			case "snip":
				//fmt.Println("Storing snippet...")
				source := strings.TrimSuffix(string(msg[4:]), "\n")
				go storeSnip(source, addr)
			case "peer":
				//fmt.Println("Storing peer address...")
				source := strings.TrimSpace(strings.TrimSuffix(string(msg[4:]), "\n"))
				go addPeer(source, addr)
			}
		}
	}
}

// =============================================================
