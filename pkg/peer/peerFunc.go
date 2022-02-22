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
	"fmt"
	"strings"
	"sync"
	"time"
)

type peerStruct struct {
	address string
	source  string
}

type receivedEvent struct {
	received     string
	source       string
	timeReceived string
}

var PeerList []peerStruct
var RecvOrder []receivedEvent

func InitPeerProcess(address string) {
	conn := sock.InitializeUdpServer(address)
	var wg sync.WaitGroup
	const TOTAL_THREADS int = 4

	for {
		msg, addr := sock.ReceiveUdpMessage(address, conn)
		fmt.Println("Received ", msg, " from ", addr)

		wg.Add(TOTAL_THREADS)

		switch string(msg[0:4]) {
		case "stop":
			fmt.Println("Received stop command, exiting...")
		case "snip":
			fmt.Println("Received snip command, exiting...")
		case "peer":
			fmt.Println("Storing peer address...")
			//trim off all the white spaces in msg
			//TODO: trim the msg, get rid of all white spaces and new lines
			// length of msg is 1024 rn
			source := strings.TrimSpace(msg)
			go addPeer(source[4:], addr)
		}
	}
}

func SetInitialPeerList(peerList []string, peerNum int) {
	PeerList = make([]peerStruct, peerNum)
	for i := 0; i < peerNum; i++ {
		PeerList[i].address = peerList[i]
		PeerList[i].source = ""
		fmt.Printf("Peer %d: %s\n", i, PeerList[i].address)
	}
}

func addPeer(receivedAddr string, source string) {
	fmt.Printf("Received peer %s from %s\n", receivedAddr, source)
	PeerList = append(PeerList, peerStruct{receivedAddr, ""})
	fmt.Printf("Peer %d: %s\n", len(PeerList)-1, PeerList[len(PeerList)-1].address)
	addRecvEvent(receivedAddr, source, time.Now().Format("2006-01-02 15:04:05"))
}

func addRecvEvent(receivedAddr string, source string, timeReceived string) {
	RecvOrder = append(RecvOrder, receivedEvent{receivedAddr, source, timeReceived})
}

// =============================================================
