// =============================================================
/*
	CPSC 559 - Iteration 2
	registryFunc.go

	Erick Yip
	Chris Chen
*/

package registry

import (
	"559Project/pkg/fileIO"
	"559Project/pkg/peer"
	"559Project/pkg/sock"
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Struct to store the server's address, list of peers, number of peers, and the time the peers are received
type regServer struct {
	address      string
	peerList     []string
	peerNum      int
	timeReceived time.Time
}

var reg regServer

/**
* Start a TCP connection to the registry server and handle the requests it sends to our process
*
* @param regAddress {string} IP address of the registry server
* @param peerAddress {string} IP address of our peer process
* @param ctx {context.Context} context type used to gracefully exit the connection to the registry
 */
func InitRegistryCommunicator(regAddress string, peerAddress string, ctx context.Context) {
	conn := sock.InitializeTcpClient(regAddress)
	fmt.Printf("Connected to server at %s\n", regAddress)

	scanner := bufio.NewScanner(conn)
	var teamName string = "It Takes Two\n"

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Closing connection to registry at %s\n", regAddress)
			conn.Close()
			return
		default:
			serverReply := sock.ReceiveTcpMessage(conn, scanner)

			var clientMessage string

			fmt.Printf("Server message = %s\n", serverReply)

			// Check the server's response and take corresponding actions
			switch {
			case strings.Contains(serverReply, "get team name"):
				clientMessage = "Sending team name...\n"
				sock.SendMessage(teamName, conn)

			case strings.Contains(serverReply, "get code"):
				clientMessage = "Sending code...\n"
				codeResponse := fileIO.ParseCodeResponse()
				sock.SendMessage(codeResponse, conn)

			case strings.Contains(serverReply, "get location"):
				clientMessage = "Sending udp server location...\n"
				sock.SendMessage(peerAddress+"\n", conn)

			case strings.Contains(serverReply, "receive peers"):
				reg.address = regAddress
				receivePeers(conn, scanner)
				clientMessage = "Storing peers...\n"

			case strings.Contains(serverReply, "get report"):
				clientMessage = "Sending report...\n"
				report := generateReport()
				sock.SendMessage(report, conn)

			case strings.Contains(serverReply, "close"):
				fmt.Printf("%s", "Closing...\n")
				conn.Close()
				return

			default:
				clientMessage = "Unknown message\n"
			}

			fmt.Printf("%s", clientMessage)
		}
	}
}

// Receive the peers from the server and store them into the peerList
func receivePeers(conn net.Conn, scanner *bufio.Scanner) {

	reply := sock.ReceiveTcpMessage(conn, scanner)

	// Convert the number of peers to int
	numPeers, _ := strconv.Atoi(strings.Split(reply, " ")[0])

	fmt.Printf("%d peers received\n", numPeers)

	// Receive the peers
	for i := 0; i < numPeers; i++ {
		// check if the peer is already in the server's peerlist
		peerAddr := sock.ReceiveTcpMessage(conn, scanner)
		if !contains(reg.peerList, peerAddr) && sock.CheckAddress(peerAddr) {
			reg.peerList = append(reg.peerList, peerAddr)
			reg.peerNum++
			peer.AppendPeer(peerAddr, reg.address)
			fmt.Printf("%s added to peer list\n", peerAddr)
		} else {
			fmt.Printf("%s already in peer list\n", peerAddr)
		}

	}
	reg.timeReceived = time.Now()
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

// =============================================================
