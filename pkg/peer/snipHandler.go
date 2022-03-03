package peer

import (
	"559Project/pkg/sock"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	count := 0
	currTimeStampStr := strconv.Itoa(currTimeStamp)
	input = "snip" + currTimeStampStr + " " + input
	currTimeStamp++
	mutex.Lock()
	for i := 1; i < len(PeerList); i++ {
		if sock.CheckAddress(PeerList[i].address) {
			if PeerList[i].address != peerProcessAddr {
				conn := sock.InitializeUdpClient(PeerList[i].address)
				sock.SendMessage(input, conn)
				conn.Close()
				count++
			}
		}
	}
	mutex.Unlock()
	fmt.Printf("Sent [%s] to %d peers\n", input, count)

}

func storeSnip(msg string, source string) {
	message := strings.Split(msg, " ")
	SnipList = append(SnipList, snip{message[1], message[0], source})
	index := searchPeerList(source)
	if index != -1 {
		PeerList[index].lastHeard = time.Now()
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
