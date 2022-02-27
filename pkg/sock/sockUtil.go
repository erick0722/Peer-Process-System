// =============================================================
/*
	CPSC 559 - Iteration 2
	sockUtil.go

	Erick Yip
	Chris Chen
*/

package sock

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Start the TCP connection to the server
func InitializeTcpClient(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	return conn
}

func InitializeUdpClient(address string) net.Conn {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	return conn
}

// https://varshneyabhi.wordpress.com/2014/12/23/simple-udp-clientserver-in-golang/
func InitializeUdpServer(address string) *net.UDPConn {
	fmt.Print("Initializing UDP server...\n")
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	return conn
}

func ReceiveUdpMessage(address string, conn *net.UDPConn) (string, string) {

	// Read from the connection
	data := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(data)
	checkError(err)
	return string(data), addr.String()

}

//Receive message from the server
func ReceiveTcpMessage(conn net.Conn, scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

// Send message to the server
func SendMessage(message string, conn net.Conn) {
	_, err := conn.Write([]byte(message))
	checkError(err)
}

func CheckAddress(addr string) bool {
	// check if the port number is valid
	address := strings.Split(addr, ":")
	if len(address) == 2 {
		port, _ := strconv.Atoi(address[1])
		if port > 0 && port < 65536 {
			return true
		}
	}
	return false
}

// Check for errors
func checkError(err error) {
	if err != nil {
		println("Error detected, exiting...", err.Error())
		os.Exit(1)
	}
}

// =============================================================
