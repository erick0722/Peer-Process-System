package peer

import (
	"559Project/pkg/sock"
	"context"
	"fmt"
	"time"
)

func AppendPeer(peer string, source string) {
	PeerList = append(PeerList, peerStruct{peer, source, time.Now()})
}

func searchPeerList(peer string) int {
	for i := 0; i < len(PeerList); i++ {
		if PeerList[i].address == peer {
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

	//add sender to list of received peers
	if sourceIndex == -1 {
		AppendPeer(source, source)
	}

	RecievedPeers = append(RecievedPeers, receivedEvent{receivedAddr, source, time.Now()})
}

func checkInactivePeers(ctx context.Context) {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
		}
		mutex.Lock()
		if len(PeerList) > 0 {
			count = 0
			for i := 0; i < len(PeerList); i++ {
				if PeerList[i].address != peerProcessAddr {
					if time.Since(PeerList[i].lastHeard) > 10*time.Second {
						count++
						PeerList = append(PeerList[:i], PeerList[i+1:]...)
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
		if len(PeerList) > 0 {
			count = 0
			currTimeStamp++
			for i := 0; i < len(PeerList); i++ {
				//send peerlist to everyone
				for j := 0; j < len(PeerList); j++ {
					if sock.CheckAddress(PeerList[j].address) {
						if PeerList[j].address != peerProcessAddr {
							conn := sock.InitializeUdpClient(PeerList[j].address)
							sock.SendMessage("peer"+PeerList[i].address, conn)
							conn.Close()
							PeersSent = append(PeersSent, sentEvent{PeerList[i].address, PeerList[j].address, time.Now()})
						}
					}
				}
				count++
			}
			fmt.Printf("Sent peerlist to %d peers at timeStamp %d\n", count, currTimeStamp)
		}
		mutex.Unlock()
	}
}
