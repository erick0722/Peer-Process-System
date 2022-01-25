package tcp

import (
	"net"
	"os"
)

// Start the TCP connection to the server
func InitializeTCP(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	CheckError(err)
	return conn
}

// Send the given message to the server
func SendMessage(message string, conn net.Conn) {
	_, err := conn.Write([]byte(message))
	CheckError(err)
}

// Check for errors and exit if any
func CheckError(err error) {
	if err != nil {
		println("Error detected, exiting...", err.Error())
		os.Exit(1)
	}
}
