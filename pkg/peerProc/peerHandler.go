// =============================================================
/*
	CPSC 559 - Iteration 2
	peerHandler.go

	Erick Yip
	Chris Chen
*/
package peerProc

import (
	"559Project/pkg/sock"
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/**
* Helper function used to add a peer to our peerlist
*
* @param peer {string} The peer to add to our peerlist
* @param source {string} The peer who sent us information about the peer to add
 */
func AppendPeer(peer string, source string) {
	mutex.Lock()
	peerList = append(peerList, peerStruct{peer, source, time.Now()})
	mutex.Unlock()
}

//
func searchPeerList(peer string) int {
	for i := 0; i < len(peerList); i++ {
		if peerList[i].address == peer {
			return i
		}
	}
	return -1
}

func addPeer(receivedAddr string, source string) {
	peerIndex := searchPeerList(receivedAddr)
	sourceIndex := searchPeerList(source)
	if peerIndex == -1 && sock.CheckAddress(receivedAddr) {
		AppendPeer(receivedAddr, source)
	}

	// Add sender to list of received peers
	if sourceIndex == -1 && sock.CheckAddress(source) {
		AppendPeer(source, source)
	}

	recievedPeers = append(recievedPeers, receivedEvent{receivedAddr, source, time.Now()})
}

/**
* Helper function to remove peers from our peerlist when they are unresponsive.
* Peers are considered unresponsive when we do not receive snips or peer messages from them after a few seconds.
*
* @param ctx {context.Context} The context from the called instance during initialization. Used to gracefully exit our program.
 */
func checkInactivePeers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
		}

		// Prevent our other go functions from reading the peerlist while peers are being removed
		mutex.Lock()
		if len(peerList) > 0 {
			count := 0
			for i := 0; i < len(peerList); i++ {
				if time.Since(peerList[i].lastHeard) > 15*time.Second {
					count++
					peerList = append(peerList[:i], peerList[i+1:]...)
					fmt.Printf("Removed peer %s from peerlist\n", peerList[i].address)
				}
			}
			fmt.Printf("Removed %d inactive peers\n", count)
		} else {
			fmt.Printf("No peers to remove\n")
		}
		mutex.Unlock()
	}
}

// Periodically send a random peer to all peers in the peer list
func sendPeer(conn *net.UDPConn, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
			// Wait for 6 seconds
		case <-time.After(6 * time.Second):
		}
		mutex.Lock()
		if len(peerList) > 0 {
			count := 0
			index := 0
			// Find a random index to send to
			for {
				rand.Seed(time.Now().UnixNano())
				index = rand.Intn(len(peerList))
				if sock.CheckAddress(peerList[index].address) {
					break
				}
			}
			for i := 0; i < len(peerList); i++ {
				if sock.CheckAddress(peerList[i].address) {
					msg := "peer" + peerList[index].address
					sock.SendUdpMsg(peerList[i].address, msg, conn)

					// Append to the list of sent peers
					peersSent = append(peersSent, sentEvent{peerList[i].address, peerList[index].address, time.Now()})
					count++
					fmt.Printf("Sent peer %s to %s\n", peerList[index].address, peerList[i].address)
				}
			}

			fmt.Printf("Sent %s to %d peers at timeStamp %d\n", peerList[index].address, count, currTimeStamp)
		} else {
			fmt.Println("No peers to send to")
		}
		mutex.Unlock()
	}
}

// =============================================================
