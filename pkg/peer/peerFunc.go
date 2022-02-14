// =============================================================
/*
	CPSC 559 - Iteration 2
	peerFunc.go

	Erick Yip
	Chris Chen
*/

package peer

import (
	"559Project/pkg/sock"
	"fmt"
)

func InitPeerProcess(address string) {
	conn := sock.InitializeUdpServer(address)

	for {
		msg, addr := sock.ReceiveUdpMessage(address, conn)
		fmt.Println("Received ", msg, " from ", addr)

		switch string(msg[0:4]) {
		case "stop":
			fmt.Println("Received stop command, exiting...")
		case "snip":
			fmt.Println("Received snip command, exiting...")
		case "peer":
			fmt.Println("Received peer command, exiting...")
		}
	}
}

// =============================================================
