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
	"bufio"
	"fmt"
	"os"
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

type snip struct {
	content   string
	timeStamp string
}

var PeerList []peerStruct
var RecievedPeers []receivedEvent
var SnipList []snip

func InitPeerProcess(address string) {
	conn := sock.InitializeUdpServer(address)

	go readSnip()
	go countdown()

	for {
		msg, addr := sock.ReceiveUdpMessage(address, conn)
		fmt.Println("Received ", msg, " from ", addr)

		switch string(msg[0:4]) {
		case "stop":
			fmt.Println("Received stop command, exiting...")
			return
		case "snip":
			fmt.Println("Storing snippet...")
			source := strings.TrimSuffix(msg[4:], "\n")
			go storeSnip(source, addr)
		case "peer":
			fmt.Println("Storing peer address...")
			source := strings.Join(strings.Split(msg, "\n"), "")
			go addPeer(source[4:], addr)
		}
	}
}

func AppendPeer(peer string, source string) {
	PeerList = append(PeerList, peerStruct{peer, source})
	fmt.Printf("Appended %s, %s\n", PeerList[len(PeerList)-1].address, PeerList[len(PeerList)-1].source)
}

func peerExists(peer string) bool {
	for i := 0; i < len(PeerList); i++ {
		if PeerList[i].address == peer {
			return true
		}
	}
	return false
}

func addPeer(receivedAddr string, source string) {
	if !peerExists(receivedAddr) {
		AppendPeer(receivedAddr, source)
	}

	//add sender to list of received peers
	if !peerExists(source) {
		AppendPeer(source, source)
	}

	addRecvEvent(receivedAddr, source, time.Now().Format("2006-01-02 15:04:05"))
}

func addRecvEvent(receivedAddr string, source string, timeReceived string) {
	RecievedPeers = append(RecievedPeers, receivedEvent{receivedAddr, source, timeReceived})
	fmt.Printf("Received Peer %d: %s, %s\n", len(RecievedPeers)-1, RecievedPeers[len(RecievedPeers)-1].received, RecievedPeers[len(RecievedPeers)-1].timeReceived)
}

func readSnip() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// read from user
		var input string
		// scan the next line
		scanner.Scan()
		// save the line in a variable
		input = scanner.Text()
		// send the snip to all peers
		sendSnip(input)
	}
}

func sendSnip(input string) {
	var wg sync.WaitGroup
	input = "snip" + "1" + input
	fmt.Printf("Sending snip: %s\n", input)
	for i := 0; i < len(PeerList); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn := sock.InitializeUdpClient(PeerList[i].address)
			sock.SendMessage(input, conn)
		}(i)
	}
	wg.Wait()
}

func storeSnip(msg string, source string) {
	//get everything after the first character of msg

	SnipList = append(SnipList, snip{msg[1:], string(msg[0])})
	fmt.Printf("Received Snip %d: %s, %s\n", len(SnipList)-1, SnipList[len(SnipList)-1].content, SnipList[len(SnipList)-1].timeStamp)
}

func countdown() {
	for {
		//wait 10 seconds
		time.Sleep(10 * time.Second)
		//send peerlist to everyone
		sendPeerList()
	}
}

func sendPeerList() {
	var wg sync.WaitGroup
	for i := 0; i < len(PeerList); i++ {
		//send peerlist to everyone
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < len(PeerList); j++ {
				conn := sock.InitializeUdpClient(PeerList[j].address)
				sock.SendMessage(PeerList[i].address, conn)
				fmt.Printf("Sent %s to %s\n", PeerList[i].address, PeerList[j].address)
			}
		}(i)
	}
	wg.Wait()
}

// =============================================================
