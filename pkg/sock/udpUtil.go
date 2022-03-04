package sock

import (
	"fmt"
	"net"
	"strings"
)

// func InitializeUdpClient(address string) net.Conn {
// 	udpAddr, err := net.ResolveUDPAddr("udp", address)
// 	checkError(err)
// 	conn, err := net.DialUDP("udp", nil, udpAddr)
// 	checkError(err)
// 	return conn
// }

func InitializeUdpServer(address string) *net.UDPConn {
	fmt.Printf("Initializing UDP server on %s\n", address)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	return conn

}

func ReceiveUdpMessage(conn *net.UDPConn) (string, string, error) {

	// Read from the connection
	data := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(data)
	//checkError(err)
	if err != nil {
		return "", "", err
	}
	msg := strings.TrimSpace(string(data[:len]))

	return msg, addr.String(), nil

}

func SendUdpMsg(addr string, msg string, conn *net.UDPConn) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	checkError(err)

	_, err = conn.WriteToUDP([]byte(msg), udpAddr)

	checkError(err)
}
