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
	sentTo      string
	peer	string
	timeSent    time.Time
}

type snip struct {
	content   string
	timeStamp string
	source    string
}



var peerList []peerStruct
var recievedPeers []receivedEvent
var snipList []snip
var peersSent []sentEvent

var peerProcessAddr string
var currTimeStamp = 0

func InitPeerProcess(address string, ctx context.Context, cancel context.CancelFunc) {

	peerProcessAddr = address
	fmt.Printf("Peer process started at %s\n", peerProcessAddr)
	peerList = append(peerList, peerStruct{peerProcessAddr, peerProcessAddr, true, time.Now()})
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
		sendpeerList(ctx)
		fmt.Printf("Exiting sendpeerList\n")
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

			index := searchpeerList(addr)
			if index != -1 {
				peerList[index].active = true
				peerList[index].lastHeard = time.Now()
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

func searchpeerList(peer string) (int) {
	for i := 0; i < len(peerList); i++ {
		if peerList[i].address == peer {
			return i
		}
	}
	return -1
}

func addPeer(receivedAddr string, source string) {
	peerIndex := searchpeerList(receivedAddr)
	sourceIndex := searchpeerList(source)
	if peerIndex == -1 {
		peerList = append(peerList, peerStruct{receivedAddr, source, true, time.Now()})
	}

	//add sender to list of received peers
	if sourceIndex == -1 {
		peerList = append(peerList, peerStruct{receivedAddr, source, true, time.Now()})
	}

	addRecvEvent(receivedAddr, source, time.Now())
}

func addRecvEvent(receivedAddr string, source string, timeReceived time.Time) {
	recievedPeers = append(recievedPeers, receivedEvent{receivedAddr, source, timeReceived})
	//fmt.Printf("Received Peer %d: %s, %s\n", len(recievedPeers)-1, recievedPeers[len(recievedPeers)-1].received, recievedPeers[len(recievedPeers)-1].timeReceived)
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

	for i := 1; i < len(peerList); i++ {
		if sock.CheckAddress(peerList[i].address) && peerList[i].active {
			if peerList[i].address != peerProcessAddr {
			conn := sock.InitializeUdpClient(peerList[i].address)
			sock.SendMessage(input, conn)
			conn.Close()
			fmt.Printf("Sent [%s] to %s\n", input, peerList[i].address)
		}
	}
	}
}

func storeSnip(msg string, source string) {
	message := strings.Split(msg, " ")
	snipList = append(snipList, snip{message[1], message[0], source})
	index := searchpeerList(source)
	if index != -1 {
		peerList[index].lastHeard = time.Now()
		peerList[index].active = true
	}

	//convert message[0] to int
	timeStamp, _ := strconv.Atoi(message[0])

	currTimeStamp = findMax(currTimeStamp, timeStamp)
	
	fmt.Printf("Received %s from %s at timeStamp %s\n", snipList[len(snipList)-1].content, snipList[len(snipList)-1].source, snipList[len(snipList)-1].timeStamp)
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

		if len(peerList) > 0 {
			for i := 0; i < len(peerList); i++ {
				if peerList[i].address != peerProcessAddr {
				if time.Since(peerList[i].lastHeard) > 10*time.Second && peerList[i].active {
					fmt.Printf("Peer %s inactive, removing...\n", peerList[i].address)
					//peerList = append(peerList[:i], peerList[i+1:]...)
					peerList[i].active = false
				}
			}
			}
			for i := 0; i < len(peerList); i++ {
				if !peerList[i].active {
					peerList = append(peerList[:i], peerList[i+1:]...)
				}
			}
		}
	}
}

func sendpeerList(ctx context.Context) {

	for {
		//time.Sleep(8 * time.Second)
		select {
		case <-ctx.Done():
			return
		case <-time.After(8 * time.Second):
		}

		if len(peerList) > 0 {
			currTimeStamp++
			for i := 0; i < len(peerList); i++ {
				//send peerlist to everyone
				for j := 0; j < len(peerList); j++ {
					if sock.CheckAddress(peerList[j].address) && peerList[j].active {
						if peerList[j].address != peerProcessAddr {
						conn := sock.InitializeUdpClient(peerList[j].address)
						sock.SendMessage("peer"+peerList[i].address, conn)
						peersSent = append(peersSent, sentEvent{peerList[i].address, peerList[j].address, time.Now()})
						conn.Close()

						//fmt.Printf("Sent %s to %s\n", peerList[i].address, peerList[j].address)
					}
				}
				}
			}
			fmt.Printf("Sent peerlist at timeStamp %d\n", currTimeStamp)
		}
	}
}

// =============================================================
