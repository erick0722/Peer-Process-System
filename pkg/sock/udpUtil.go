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
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	return conn

}

// Listen and read UDP message coming from other peers
func ReceiveUdpMessage(conn *net.UDPConn) (string, string, error) {

	// Read from the connection
	data := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return "", "", err
	}
	msg := strings.TrimSpace(string(data[:len]))
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	return msg, addr.String(), err

}

// Send a message to the address
func SendUdpMsg(addr string, msg string, conn *net.UDPConn) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	checkError(err)

	_, err = conn.WriteToUDP([]byte(msg), udpAddr)

	checkError(err)
}

// func WaitForStop(conn *net.UDPConn) (string, string, error) {
// 	ch := time.After(time.Second * 15)

// 	for {
// 		select {
// 		case <-ch:
// 			fmt.Printf("Finished waiting.\n")
// 			return "", "", nil
// 		default:
// 			fmt.Printf("Waiting to receive a stop...\n")
// 			msg, addr, err := ReceiveUdpMessage(conn)
// 			return msg, addr, err
// 		}
// 	}
// }

// =============================================================
