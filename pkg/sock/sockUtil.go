/*
	CPSC 559 - Iteration 1
	tcpUtil.go

	Erick Yip
	Chris Chen
*/

package sock

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	return conn
}

func ReceiveUdpMessage(address string) bool {
	var conn *net.UDPConn = InitializeUdpServer(address)

	for {
		// Read from the connection
		data := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(data)
		checkError(err)

		// Print the received data
		fmt.Println("Received ", string(data), " from ", addr)

		// compare the first 4 letters of the message with the string "exit" (testing purposes)
		if string(data[0:4]) == "exit" {
			fmt.Println("Exiting...")
			return true
		}
	}
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

// Check for errors
func checkError(err error) {
	if err != nil {
		println("Error detected, exiting...", err.Error())
		os.Exit(1)
	}
}
