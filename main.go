package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {

	// servAddr := "localhost:6666"
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
		clientMessage = "get code received\n"
		// TODO
	case serverReply == "receive peers\n":
		clientMessage = "receive peers received\n"
		// TODO
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
	conn.Close()
}

func sendTeamName(teamName string, conn net.Conn) () {
	fmt.Printf("client message=%s", teamName)

	_, err := conn.Write([]byte(teamName))
	if err != nil {
		print("Write to server failed:", err.Error())
		os.Exit(1)
	}

}