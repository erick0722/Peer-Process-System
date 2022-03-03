package sock

import (
	"bufio"
	"net"
)

// Start the TCP connection to the server
func InitializeTcpClient(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	return conn
}

//Receive message from the server
func ReceiveTcpMessage(conn net.Conn, scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}
