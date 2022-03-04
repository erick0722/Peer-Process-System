// =============================================================
/*
	CPSC 559 - Iteration 2
	tcpUtil.go

	Erick Yip
	Chris Chen
*/

package sock

import (
	"bufio"
	"net"
)

// Start the TCP connection to the registry
func InitializeTcpClient(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	return conn
}

//Receive message from the registry
func ReceiveTcpMessage(conn net.Conn, scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

// =============================================================
