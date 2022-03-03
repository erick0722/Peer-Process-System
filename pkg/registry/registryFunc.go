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

func InitRegistryCommunicator(regAddress string, peerAddress string, ctx context.Context) {
	conn := sock.InitializeTcpClient(regAddress)
	fmt.Printf("Connected to server at %s\n", regAddress)

	var registry regServer = regServer{regAddress, []string{}, 0, time.Now()}
	scanner := bufio.NewScanner(conn)

	var teamName string = "It Takes Two\n" // Our team name

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
				registry.address = regAddress
				receivePeers(&registry, conn, scanner)
				clientMessage = "Storing peers...\n"

			case strings.Contains(serverReply, "get report"):
				//TODO: update this
				clientMessage = "Sending report...\n"
				report := generateReport(registry)
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
func receivePeers(server *regServer, conn net.Conn, scanner *bufio.Scanner) {

	reply := sock.ReceiveTcpMessage(conn, scanner)

	// Convert the number of peers to int
	numPeers, _ := strconv.Atoi(strings.Split(reply, " ")[0])

	fmt.Printf("%d peers received\n", numPeers)

	// Receive the peers
	for i := 0; i < numPeers; i++ {
		// check if the peer is already in the server's peerlist
		peerAddr := sock.ReceiveTcpMessage(conn, scanner)
		if !contains(server.peerList, peerAddr) {
			server.peerList = append(server.peerList, peerAddr)
			server.peerNum++
			peer.AppendPeer(peerAddr, server.address)
			fmt.Printf("%s added to peer list\n", peerAddr)
		} else {
			fmt.Printf("%s already in peer list\n", peerAddr)
		}

	}
	server.timeReceived = time.Now()
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
