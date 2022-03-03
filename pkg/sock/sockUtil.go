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
	"os"
	"strconv"
	"strings"
)

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
