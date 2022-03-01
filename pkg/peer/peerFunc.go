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
	"context"
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
	active    bool
	lastHeard time.Time
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

func InitPeerProcess(address string, ctx context.Context, cancel context.CancelFunc) {

	//PeerList = append(PeerList, peerStruct{address, address, ""})
	var wg sync.WaitGroup
	var ch = make(chan bool)

	wg.Add(4)
	go func() {
		defer wg.Done()
		handleMessage(address, ctx, cancel)
		fmt.Printf("Exiting handleMessage\n")
	}()

	go func() {
		defer wg.Done()
		readSnip(ctx)
		fmt.Printf("Exiting readSnip\n")
	}()

	go func() {
		defer wg.Done()
		sendPeerList(ctx)
		fmt.Printf("Exiting sendPeerList\n")
	}()

	go func() {
		defer wg.Done()
		checkInactivePeers(ctx, ch)
		fmt.Printf("Exiting checkInactivePeers\n")
	}()

	wg.Wait()

}

func handleMessage(address string, ctx context.Context, cancel context.CancelFunc) {
	address, conn := sock.InitializeUdpServer(address)

	go func() {
		<-ctx.Done()
		sock.CloseUDP(conn)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Waiting for message\n")
			msg, addr, err := sock.ReceiveUdpMessage(address, conn)
			fmt.Println("Received ", msg, " from ", addr)
			if err != nil {
				fmt.Printf("Error detected: %v\n", err)
				continue
			}
			switch string(msg[0:4]) {
			case "stop":
				fmt.Printf("Received stop command, exiting...\n")
				sock.CloseUDP(conn)
				cancel()
				return
			case "snip":
				//fmt.Println("Storing snippet...")
				source := strings.TrimSuffix(msg[4:], "\n")
				go storeSnip(source, addr)
			case "peer":
				//fmt.Println("Storing peer address...")
				source := strings.Join(strings.Split(msg, "\n"), "")
				go addPeer(source[4:], addr)
			}
		}
	}
}

func AppendPeer(peer string, source string) {
	PeerList = append(PeerList, peerStruct{peer, source, true, time.Now()})
	//fmt.Printf("Appended %s, %s\n", PeerList[len(PeerList)-1].address, PeerList[len(PeerList)-1].source)
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

func readSnip(ctx context.Context) {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case input := <-ch:
			sendSnip(input)
		}
	}
}

func sendSnip(input string) {

	currTimeStampStr := strconv.Itoa(currTimeStamp)

	input = "snip" + currTimeStampStr + " " + input
	for i := 1; i < len(PeerList); i++ {
		if sock.CheckAddress(PeerList[i].address) && PeerList[i].active {
			conn := sock.InitializeUdpClient(PeerList[i].address)
			sock.SendMessage(input, conn)
			fmt.Printf("Sent [%s] to %s\n", input, PeerList[i].address)
		}
	}
	currTimeStamp++

}

func storeSnip(msg string, source string) {
	message := strings.Split(msg, " ")
	SnipList = append(SnipList, snip{message[1], message[0], source})
	_, index := searchPeerList(source)
	if index != -1 {
		PeerList[index].lastHeard = time.Now()
		PeerList[index].active = true
	}

	fmt.Printf("Received %s from %s at timeStamp %s\n", SnipList[len(SnipList)-1].content, SnipList[len(SnipList)-1].source, SnipList[len(SnipList)-1].timeStamp)
}

func checkInactivePeers(ctx context.Context, ch chan bool) {

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			// case <-ch:
			// 	continue
		}

		if len(PeerList) > 0 {
			for i := 0; i < len(PeerList); i++ {

				if time.Since(PeerList[i].lastHeard) > 10*time.Second && PeerList[i].active {

					fmt.Printf("Peer %s inactive, removing...\n", PeerList[i].address)
					//PeerList = append(PeerList[:i], PeerList[i+1:]...)
					PeerList[i].active = false
				}
			}
		}
	}
}

func sendPeerList(ctx context.Context) {

	for {
		//time.Sleep(8 * time.Second)
		select {
		case <-ctx.Done():
			return
		case <-time.After(8 * time.Second):
		}

		if len(PeerList) > 0 {
			for i := 0; i < len(PeerList); i++ {
				//send peerlist to everyone
				for j := 0; j < len(PeerList); j++ {
					if sock.CheckAddress(PeerList[j].address) && PeerList[j].active {
						conn := sock.InitializeUdpClient(PeerList[j].address)
						sock.SendMessage("peer"+PeerList[i].address, conn)
						//fmt.Printf("Sent %s to %s\n", PeerList[i].address, PeerList[j].address)
					}
				}
			}
			currTimeStamp++
			fmt.Printf("Sent peerlist at timeStamp %d\n", currTimeStamp)
		}
	}
}

// =============================================================
