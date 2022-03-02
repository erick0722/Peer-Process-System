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

var PeerList []peerStruct
var RecievedPeers []receivedEvent
var PeersSent []sentEvent
var SnipList []snip

var peerProcessAddr string
var currTimeStamp = 0

func InitPeerProcess(address string, ctx context.Context, cancel context.CancelFunc) {

	peerProcessAddr = address
	fmt.Printf("Peer process started at %s\n", peerProcessAddr)
	PeerList = append(PeerList, peerStruct{peerProcessAddr, peerProcessAddr, true, time.Now()})
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
			//fmt.Printf("Waiting for message\n")
			msg, addr, err := sock.ReceiveUdpMessage(address, conn)
			//fmt.Println("Received ", msg, " from ", addr)
			if err != nil {
				fmt.Printf("Error detected: %v\n", err)
				continue
			}

			index := searchPeerList(addr)
			if index != -1 {
				PeerList[index].active = true
				PeerList[index].lastHeard = time.Now()
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

func searchPeerList(peer string) int {
	for i := 0; i < len(PeerList); i++ {
		if PeerList[i].address == peer {
			return i
		}
	}
	return -1
}

func addPeer(receivedAddr string, source string) {
	peerIndex := searchPeerList(receivedAddr)
	sourceIndex := searchPeerList(source)
	if peerIndex == -1 {
		AppendPeer(receivedAddr, source)
	}

	//add sender to list of received peers
	if sourceIndex == -1 {
		AppendPeer(source, source)
	}

	RecievedPeers = append(RecievedPeers, receivedEvent{receivedAddr, source, time.Now()})
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
	currTimeStamp++

	for i := 1; i < len(PeerList); i++ {
		if sock.CheckAddress(PeerList[i].address) && PeerList[i].active {
			if PeerList[i].address != peerProcessAddr {
				conn := sock.InitializeUdpClient(PeerList[i].address)
				sock.SendMessage(input, conn)
				conn.Close()
				fmt.Printf("Sent [%s] to %s\n", input, PeerList[i].address)
			}
		}
	}
}

func storeSnip(msg string, source string) {
	message := strings.Split(msg, " ")
	SnipList = append(SnipList, snip{message[1], message[0], source})
	index := searchPeerList(source)
	if index != -1 {
		PeerList[index].lastHeard = time.Now()
		PeerList[index].active = true
	}

	//convert message[0] to int
	timeStamp, _ := strconv.Atoi(message[0])

	currTimeStamp = findMax(currTimeStamp, timeStamp)

	fmt.Printf("Received %s from %s at timeStamp %s\n", SnipList[len(SnipList)-1].content, SnipList[len(SnipList)-1].source, SnipList[len(SnipList)-1].timeStamp)
}

func findMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
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
				if PeerList[i].address != peerProcessAddr {
					if time.Since(PeerList[i].lastHeard) > 10*time.Second && PeerList[i].active {
						fmt.Printf("Peer %s inactive, removing...\n", PeerList[i].address)
						//PeerList = append(PeerList[:i], PeerList[i+1:]...)
						PeerList[i].active = false
					}
				}
			}
			for i := 0; i < len(PeerList); i++ {
				if !PeerList[i].active {
					PeerList = append(PeerList[:i], PeerList[i+1:]...)
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
			currTimeStamp++
			for i := 0; i < len(PeerList); i++ {
				//send peerlist to everyone
				for j := 0; j < len(PeerList); j++ {
					if sock.CheckAddress(PeerList[j].address) && PeerList[j].active {
						if PeerList[j].address != peerProcessAddr {
							conn := sock.InitializeUdpClient(PeerList[j].address)
							sock.SendMessage("peer"+PeerList[i].address, conn)
							conn.Close()
							PeersSent = append(PeersSent, sentEvent{PeerList[i].address, PeerList[j].address, time.Now()})
						}
					}
				}
			}
			fmt.Printf("Sent peerlist at timeStamp %d\n", currTimeStamp)
		}
	}
}

// =============================================================
