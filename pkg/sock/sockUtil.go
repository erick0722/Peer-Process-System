// =============================================================
/*
	CPSC 559 - Iteration 2
	sockUtil.go

	Erick Yip
	Chris Chen
*/

package sock

import (
	"net"
	"strconv"
	"strings"
)

// Send message to the server
func SendMessage(message string, conn net.Conn) {
	_, err := conn.Write([]byte(message))
	checkError(err)
}

// Check if the ip address is valid
func CheckAddress(addr string) bool {
	// check if the port number is valid
	address := strings.Split(addr, ":")
	if len(address) == 2 {
		_, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return false
		}

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
		println("Error detected!", err.Error())
	}
}

// =============================================================
