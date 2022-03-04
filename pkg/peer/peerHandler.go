package peer

import (
	"559Project/pkg/sock"
	"context"
	"fmt"
	"math/rand"
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
	if peerIndex == -1 {
		AppendPeer(receivedAddr, source)
	}

	// Add sender to list of received peers
	if sourceIndex == -1 {
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
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
		}

		// Prevent our other go functions from reading the peerlist while peers are being removed
		mutex.Lock()
		if len(peerList) > 0 {
			count = 0
			for i := 0; i < len(peerList); i++ {
				if peerList[i].address != peerProcessAddr {
					if time.Since(peerList[i].lastHeard) > 10*time.Second {
						count++
						peerList = append(peerList[:i], peerList[i+1:]...)
					}
				}
			}
			fmt.Printf("Removed %d inactive peers\n", count)
		}
		mutex.Unlock()
	}
}

func sendPeerList(ctx context.Context) {
	count := 0
	for {
		//time.Sleep(8 * time.Second)
		select {
		case <-ctx.Done():
			return
		case <-time.After(6 * time.Second):
		}
		mutex.Lock()
		if len(peerList) > 0 {
			count = 0
			currTimeStamp++
			//find a random index to send to
			index := rand.Intn(len(peerList))
			for i := 0; i < len(peerList); i++ {
				if peerList[i].address != peerProcessAddr {
					if sock.CheckAddress(peerList[index].address) {
						conn := sock.InitializeUdpClient(peerList[i].address)
						sock.SendMessage("peer"+peerList[index].address, conn)
						conn.Close()
						peersSent = append(peersSent, sentEvent{peerList[i].address, peerList[index].address, time.Now()})
						count++
					}
				}
			}
			// for i := 0; i < len(peerList); i++ {
			// 	//send peerlist to everyone
			// 	for j := 0; j < len(peerList); j++ {
			// 		if sock.CheckAddress(peerList[j].address) {
			// 			if peerList[j].address != peerProcessAddr {
			// 				conn := sock.InitializeUdpClient(peerList[j].address)
			// 				sock.SendMessage("peer"+peerList[i].address, conn)
			// 				conn.Close()
			// 				peersSent = append(peersSent, sentEvent{peerList[i].address, peerList[j].address, time.Now()})
			// 			}
			// 		}
			// 	}
			// 	count++
			// }
			fmt.Printf("Sent %s to %d peers at timeStamp %d\n", peerList[index].address, count, currTimeStamp)
		}
		mutex.Unlock()
	}
}
