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
	"strconv"
	"strings"
	"sync"
	"time"
)

type peerStruct struct {
	address   string
	source    string
	lastHeard string
}

type receivedEvent struct {
	received     string
	source       string
	timeReceived string
}

type snip struct {
	content   string
	timeStamp string
	source    string
}

var PeerList []peerStruct
var RecievedPeers []receivedEvent
var SnipList []snip

var currTimeStamp int = 0

func InitPeerProcess(address string) {
	address, conn := sock.InitializeUdpServer(address)

	go readSnip()
	go sendPeerList()
	//go checkInactivePeers()

	PeerList = append(PeerList, peerStruct{address, address, ""})

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
	PeerList = append(PeerList, peerStruct{peer, source, time.Now().Format("2006-01-02 15:04:05")})
	fmt.Printf("Appended %s, %s\n", PeerList[len(PeerList)-1].address, PeerList[len(PeerList)-1].source)
}

func searchPeerList(peer string) (bool, int) {
	for i := 0; i < len(PeerList); i++ {
		if PeerList[i].address == peer {
			return true, i
		}
	}
	return false, -1
}

func addPeer(receivedAddr string, source string) {
	peerExist, _ := searchPeerList(receivedAddr)
	sourceExist, _ := searchPeerList(source)
	if !peerExist {
		AppendPeer(receivedAddr, source)
	}

	//add sender to list of received peers
	if !sourceExist {
		AppendPeer(source, source)
	}

	addRecvEvent(receivedAddr, source, time.Now().Format("2006-01-02 15:04:05"))
}

func addRecvEvent(receivedAddr string, source string, timeReceived string) {
	RecievedPeers = append(RecievedPeers, receivedEvent{receivedAddr, source, timeReceived})
	//fmt.Printf("Received Peer %d: %s, %s\n", len(RecievedPeers)-1, RecievedPeers[len(RecievedPeers)-1].received, RecievedPeers[len(RecievedPeers)-1].timeReceived)
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

	currTimeStampStr := strconv.Itoa(currTimeStamp)

	input = "snip" + currTimeStampStr + " " + input
	for i := 1; i < len(PeerList); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if sock.CheckAddress(PeerList[i].address) {
				conn := sock.InitializeUdpClient(PeerList[i].address)
				sock.SendMessage(input, conn)
				fmt.Printf("Sent [%s] to %s\n", input, PeerList[i].address)
			}
		}(i)
	}
	wg.Wait()
	currTimeStamp++

}

func storeSnip(msg string, source string) {
	SnipList = append(SnipList, snip{msg[1:], string(msg[0]), source})
	_, index := searchPeerList(source)
	if index != -1 {
		PeerList[index].lastHeard = time.Now().Format("2006-01-02 15:04:05")
	}

	fmt.Printf("Received %s from %s at timeStamp %s\n", SnipList[len(SnipList)-1].content, SnipList[len(SnipList)-1].source, SnipList[len(SnipList)-1].timeStamp)
}

func checkInactivePeers() {
	var wg sync.WaitGroup
	for {
		time.Sleep(10 * time.Second)
		if len(PeerList) > 0 {
			for i := 0; i < len(PeerList); i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					if PeerList[i].lastHeard != "" {
						diff, _ := time.Parse("2006-01-02 15:04:05", PeerList[i].lastHeard)
						if time.Since(diff) >= 10*time.Second {
							fmt.Printf("Removing peer %s\n", PeerList[i].address)
							PeerList = append(PeerList[:i], PeerList[i+1:]...)
						}
					}
				}(i)
			}
			wg.Wait()
		}
	}
}

func sendPeerList() {
	var wg sync.WaitGroup
	for {
		time.Sleep(8 * time.Second)
		if len(PeerList) > 0 {
			for i := 0; i < len(PeerList); i++ {
				//send peerlist to everyone
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					for j := 0; j < len(PeerList); j++ {
						if sock.CheckAddress(PeerList[j].address) {
							conn := sock.InitializeUdpClient(PeerList[j].address)
							sock.SendMessage(PeerList[i].address, conn)
							//fmt.Printf("Sent %s to %s\n", PeerList[i].address, PeerList[j].address)
						}
					}
				}(i)
			}
			wg.Wait()
			currTimeStamp++

		}
	}
}

// =============================================================
