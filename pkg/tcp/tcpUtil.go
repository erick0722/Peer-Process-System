/*
	CPSC 559 - Iteration 1
	tcpUtil.go

	Erick Yip
	Chris Chen
*/

package tcp

import (
	"bytes"
	"net"
	"os"
)

// Start the TCP connection to the server
func InitializeTCP(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	return conn
}

// Receive message from the server
func ReceiveMessage(conn net.Conn) string {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	checkError(err)
	return string(bytes.Trim(reply, "\x00"))
}

// Send message to the server
func SendMessage(message string, conn net.Conn) {
	_, err := conn.Write([]byte(message))
	checkError(err)
}

// Check for errors
func checkError(err error) {
	if err != nil {
		println("Error detected, exiting...", err.Error())
		os.Exit(1)
	}
}
