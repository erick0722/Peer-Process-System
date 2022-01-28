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
	"bufio"
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

	if len(os.Args) != 2 {
		fmt.Println("Missing <server address>")
		os.Exit(1)
	}

	// Get the server's address from the command line
	address := os.Args[1]

	var server regServer = regServer{"", []string{}, 0, ""}

	var teamName string = "It Takes Two\n" // Our team name

	// Connect to the server via TCP
	conn := tcp.InitializeTCP(address)
	fmt.Printf("Connected to server at %s\n", address)

	scanner := bufio.NewScanner(conn)

forLoop:
	for {
		serverReply := tcp.ReceiveMessage(conn, scanner)

		var clientMessage string

		fmt.Printf("Server message = %s\n", serverReply)

		// Check the server's response and take corresponding actions
		switch {
		case strings.Contains(serverReply, "get team name"):
			clientMessage = "Sending team name...\n"
			tcp.SendMessage(teamName, conn)

		case strings.Contains(serverReply, "get code"):
			clientMessage = "Sending code...\n"
			codeResponse := fileIO.ParseCodeResponse()
			tcp.SendMessage(codeResponse, conn)

		case strings.Contains(serverReply, "receive peers"):
			server.address = address
			receivePeers(&server, conn, scanner)
			clientMessage = "Storing peers...\n"

		case strings.Contains(serverReply, "get report"):
			clientMessage = "Sending report...\n"
			report := generateReport(server)
			tcp.SendMessage(report, conn)

		case strings.Contains(serverReply, "close"):
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
func receivePeers(server *regServer, conn net.Conn, scanner *bufio.Scanner) {

	reply := tcp.ReceiveMessage(conn, scanner)

	// Convert the number of peers to int
	numPeers, _ := strconv.Atoi(strings.Split(reply, " ")[0])

	fmt.Printf("%d peers received\n", numPeers)

	// Receive the peers
	for i := 0; i < numPeers; i++ {
		// check if the peer is already in the server's peerlist
		peer := tcp.ReceiveMessage(conn, scanner)
		if !contains(server.peerList, peer) {
			server.peerList = append(server.peerList, peer)
			server.peerNum++
			fmt.Printf("%s added to peer list\n", peer)
		} else {
			fmt.Printf("%s already in peer list\n", peer)
		}
	}

	server.timeReceived = time.Now().Format("2006-01-02 15:04:05")
}

// Checks if a string is present in a slice
func contains(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// Generate a report for the list of peers and sources
func generateReport(server regServer) string {

	// Return nothing if no peers have been received (address empty)
	if server.address == "" {
		return "0\n0\n\n0\n"
	}

	// Convert the number of peers to string
	peerNumString := strconv.Itoa(server.peerNum)
	report := fmt.Sprintf("%s\n", peerNumString)

	// Concat the list of peers
	report += concatPeers(server)

	// Format the report
	report += fmt.Sprintf("1\n%s\n%s\n%s\n", server.address, server.timeReceived, peerNumString)

	report += concatPeers(server)

	fmt.Printf("%s", report)

	return report
}

func concatPeers(server regServer) string {
	var peerList string

	for i := 0; i < server.peerNum; i++ {
		peerList += fmt.Sprintf("%s\n", server.peerList[i])
	}

	return peerList
}
