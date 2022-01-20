package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	servAddr := "localhost:55921"

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	for {

		reply := make([]byte, 1024)

		_, err = conn.Read(reply)
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		serverReply := string(bytes.Trim(reply, "\x00"))
		var clientMessage string
		println("reply from server=", serverReply)

		switch {
		case serverReply == "get team name\n":
			clientMessage = "it takes two\n"
			sendTeamName(clientMessage, conn)
		case serverReply == "get code\n":
			//TODO: send code over
			clientMessage = "go\ndickass\n...\n"
			sendTeamName(clientMessage, conn)
		case serverReply == "receive peers\n":
			clientMessage = "receive peers received\n"
			//peerNum, peerList := receivePeers(conn)
			_, _ = receivePeers(conn)
		case serverReply == "get report\n":
			clientMessage = "get report received\n"
			// TODO
		case serverReply == "close\n":
			clientMessage = "bye have a great day\n"
			conn.Close()
		default:
			clientMessage = "unknown message\n"
		}

		fmt.Printf("client message=%s", clientMessage)
	}
}

func sendTeamName(teamName string, conn net.Conn) {
	_, err := conn.Write([]byte(teamName))
	if err != nil {
		print("Write to server failed:", err.Error())
		os.Exit(1)
	}
}

func receivePeers(conn net.Conn) (int, []string) {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("reply from server=%s", reply)
	numPeers, _ := strconv.Atoi(strings.Split(string(reply), "\n")[0])
	fmt.Printf("numPeers=%d\n", numPeers)
	peerIPs := make([]string, numPeers)
	for i := 0; i < numPeers; i++ {
		peerIPs[i] = strings.Split(string(reply), "\n")[i+1]
		fmt.Printf("peerIPs[%d]=%s\n", i, peerIPs[i])
	}

	return numPeers, peerIPs
}
