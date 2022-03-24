// =============================================================
/*
	CPSC 559 - Iteration 3
	udpUtil.go

	Erick Yip
	Chris Chen
*/

package sock

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Initialize UDP server connection at the given address
func InitializeUdpServer(address string) *net.UDPConn {
	fmt.Printf("Initializing UDP server on %s\n", address)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	return conn
}

// Listen and read UDP message coming from other peers
func ReceiveUdpMessage(conn *net.UDPConn) (string, string, error) {
	// Set the read deadline to be 10 seconds max
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))

	// Read from the connection
	data := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return "", "", err
	}
	msg := strings.TrimSpace(string(data[:len]))
	return msg, addr.String(), err
}

// Send a message to the address
func SendUdpMsg(addr string, msg string, conn *net.UDPConn) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	checkError(err)

	_, err = conn.WriteToUDP([]byte(msg), udpAddr)

	checkError(err)
}

// =============================================================
