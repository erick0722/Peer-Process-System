/*
	CPSC 559 - Iteration 1
	main.go

	Erick Yip
	Chris Chen
*/

package main

import (
	"559Project/pkg/fileIO"
	"559Project/pkg/tcp"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Struct to store the server's address, list of peers, number of peers, and the time the peers are received
type regServer struct {
	address      string
	peerList     []string
	peerNum      int
	timeReceived string
}

func main() {

	// Get the server's address from the command line
	address := os.Args[1]

	var registry regServer = regServer{"", []string{}, 0, ""}

	var teamName string = "It Takes Two\n" // Our team name

	// Connect to the server via TCP
	conn := tcp.InitializeTCP(address)

forLoop:
	for {
		serverReply := tcp.ReceiveMessage(conn)

		var clientMessage string

		fmt.Printf("Server message = %s", serverReply)

		// Check the server's response and take corresponding actions
		switch {
		case serverReply == "get team name\n":
			clientMessage = "Sending team name...\n"
			tcp.SendMessage(teamName, conn)

		case serverReply == "get code\n":
			clientMessage = "Sending code...\n"
			codeResponse := fileIO.ParseCodeResponse()
			tcp.SendMessage(codeResponse, conn)

		case serverReply == "receive peers\n":
			registry.address = address
			registry.peerNum, registry.peerList, registry.timeReceived = receivePeers(conn)
			clientMessage = "Peers received\n"

		case serverReply == "get report\n":
			clientMessage = "Sending report...\n"
			report := generateReport(registry)
			tcp.SendMessage(report, conn)

		case serverReply == "close\n":
			fmt.Printf("%s", "Closing...\n")
			conn.Close()
			break forLoop

		default:
			clientMessage = "Unknown message\n"
		}

		fmt.Printf("%s", clientMessage)
	}
}

// Receive the peers from the server and store them into the peerList
func receivePeers(conn net.Conn) (int, []string, string) {

	reply := tcp.ReceiveMessage(conn)

	numPeers, _ := strconv.Atoi(strings.Split(string(reply), "\n")[0])

	fmt.Printf("numPeers=%d\n", numPeers)

	peerList := make([]string, numPeers)
	for i := 0; i < numPeers; i++ {
		peerList[i] = strings.Split(string(reply), "\n")[i+1]
		fmt.Printf("peerList[%d]=%s\n", i, peerList[i])
	}
	return numPeers, peerList, time.Now().Format("2006-01-02 15:04:05")
}

// Generate a report for the list of peers and sources
func generateReport(registry regServer) string {

	// Return nothing if no peers have been received (address empty)
	if registry.address == "" {
		return "0\n0\n\n0\n"
	}

	// Convert the number of peers to string
	peerNumString := strconv.Itoa(registry.peerNum)
	report := fmt.Sprintf("%s\n", peerNumString)

	// Concat the list of peers
	for i := 0; i < registry.peerNum; i++ {
		report += fmt.Sprintf("%s\n", registry.peerList[i])
	}

	report += fmt.Sprintf("1\n%s\n%s\n%s\n", registry.address, registry.timeReceived, peerNumString)

	for i := 0; i < registry.peerNum; i++ {
		report += fmt.Sprintf("%s\n", registry.peerList[i])
	}

	fmt.Printf("Report:\n%s", report)

	return report
}
