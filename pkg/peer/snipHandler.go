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

/** 
* Reads a message from standard input to be sent to other peers. 
* Users can send a message while our peer process reads messages, sends messages, etc. 
* 
* @param ctx {context.Context} 	The context initiated from the snipHandler. 
* 							   	When the context's cancel function is called, will signal the
*								function to gracefully exit
*/
func readSnip(ctx context.Context) {
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
			sendSnip(input)
		}
	}
}

/** 
* Send a snip from our peer process to another peer process(es)
* @param input {string} The snip message read from standard input: contains correct formatting
*/
func sendSnip(input string) {
	count := 0
	currTimeStampStr := strconv.Itoa(currTimeStamp)
	input = "snip" + currTimeStampStr + " " + input
	currTimeStamp++

	// Prevent other running threads from reading the peerlist while we send a snip
	mutex.Lock()
	for i := 1; i < len(peerList); i++ {
		if sock.CheckAddress(peerList[i].address) {
			if peerList[i].address != peerProcessAddr {
				conn := sock.InitializeUdpClient(peerList[i].address)
				sock.SendMessage(input, conn)
				conn.Close()
				count++
			}
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
	snipList = append(snipList, snip{message[1], message[0], source})
	
	// Update the sender of the message in our peerlist 
	index := searchPeerList(source)
	if index != -1 {
		peerList[index].lastHeard = time.Now()
	}

	// Compare the timestamp in the message with our internal clock per Lamport's Timestamp Algorithm
	timeStamp, _ := strconv.Atoi(message[0])

	currTimeStamp = findMax(currTimeStamp, timeStamp) + 1

	fmt.Printf("Received %s from %s at timeStamp %s\n", snipList[len(snipList)-1].content, snipList[len(snipList)-1].source, snipList[len(snipList)-1].timeStamp)
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
