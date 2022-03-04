// =============================================================
/*
	CPSC 559 - Iteration 2
	snipHandler.go

	Erick Yip
	Chris Chen
*/

package peer

import (
	"559Project/pkg/sock"
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
* Reads a message from standard input to be sent to other peers.
* Users can send a message while our peer process reads messages, sends messages, etc.
*
* @param ctx {context.Context} 	The context initiated from the snipHandler.
* 							   	When the context's cancel function is called, will signal the
*								function to gracefully exit
 */
func readSnip(conn *net.UDPConn, ctx context.Context) {
	// Read input from the user
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
			sendSnip(input, conn)
		}
	}
}

/**
* Send a snip from our peer process to another peer process(es)
* @param input {string} The snip message read from standard input: contains correct formatting
 */
func sendSnip(input string, conn *net.UDPConn) {
	count := 0
	currTimeStamp++
	currTimeStampStr := strconv.Itoa(currTimeStamp)
	input = "snip" + currTimeStampStr + " " + input

	// Prevent other running threads from reading the peerlist while we send a snip
	mutex.Lock()
	for i := 0; i < len(peerList); i++ {
		if sock.CheckAddress(peerList[i].address) {
			sock.SendUdpMsg(peerList[i].address, input, conn)
			count++
		}
	}
	mutex.Unlock()
	fmt.Printf("Sent [%s] to %d peers\n", input, count)
}

/**
* Store the sent snip in our own process.
*
* @param msg {string} The snip that was sent to our peer
* @param source {string} The peer that sent the snip
 */
func storeSnip(msg string, source string) {
	message := strings.Split(msg, " ")
	if len(message) < 2 {
		fmt.Printf("Invalid message received: %s", msg)
		return
	}
	timeStamp, err := strconv.Atoi(message[0])
	if err != nil {
		fmt.Printf("Invalid message received: %s", msg)
		return
	}
	mutex.Lock()

	//join everything after the first space
	messageStr := strings.Join(message[1:], " ")

	// Update the sender of the message in our peerlist
	index := searchPeerList(source)
	if index != -1 {
		peerList[index].lastHeard = time.Now()
	}

	// Compare the timestamp in the message with our internal clock per Lamport's Timestamp Algorithm
	if source != peerProcessAddr {
		currTimeStamp = findMax(currTimeStamp, timeStamp)
	}
	currTimeStampStr := strconv.Itoa(currTimeStamp)

	snipList = append(snipList, snip{messageStr, currTimeStampStr, source})

	fmt.Printf("Received %s from %s at timeStamp %s\n", snipList[len(snipList)-1].content, snipList[len(snipList)-1].source, snipList[len(snipList)-1].timeStamp)
	mutex.Unlock()
}

/*
	Used to find max when implementing Lamport's Timestamp Algorithm for

	@param a {int} Integer 1
	@param b {int} Integer 2
	@returns {int} The largest between a or b. If integers are equal return b.
*/
func findMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// =============================================================
