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

func ReceiveMessage(conn net.Conn) string {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	checkError(err)
	return string(bytes.Trim(reply, "\x00"))
}

// Send the given message to the server
func SendMessage(message string, conn net.Conn) {
	_, err := conn.Write([]byte(message))
	checkError(err)
}

// Check for errors and exit if any
func checkError(err error) {
	if err != nil {
		println("Error detected, exiting...", err.Error())
		os.Exit(1)
	}
}
