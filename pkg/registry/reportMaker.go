package registry

import (
	"559Project/pkg/peer"
	"fmt"
	"strconv"
)

// Generate a report for the list of peers and sources
func generateReport(server regServer) string {

	// Return nothing if no peers have been received (address empty)
	if server.address == "" {
		return "0\n0\n\n0\n"
	}

	// Add the current list of peers to the report
	peerNum := strconv.Itoa(len(peer.PeerList))
	report := fmt.Sprintf("%s\n", peerNum)
	report += peer.ConcatPeerList(peer.PeerList)

	// Add the lists we have received to the report
	report += fmt.Sprintf("1\n%s\n%s\n%s\n", server.address, server.timeReceived, strconv.Itoa(server.peerNum))
	report += concatRegPeers(server)

	// Add the peers received via UDP/IP to the report
	receivedPeerNum := strconv.Itoa(len(peer.RecievedPeers))
	report += fmt.Sprintf("%s\n", receivedPeerNum)
	report += peer.ConcatRecvPeerList(peer.RecievedPeers)

	// Add all the ppers we sent via UDP/IP to the report
	peersSent := strconv.Itoa(len(peer.PeersSent))
	report += fmt.Sprintf("%s\n", peersSent)
	report += peer.ConcatPeersSent(peer.PeersSent)

	// Add all snippets we received to the report
	numSnippets := strconv.Itoa(len(peer.SnipList))
	report += fmt.Sprintf("%s\n", numSnippets)
	report += peer.ConcatSnipList(peer.SnipList)

	fmt.Printf("%s", report)
	return report
}

func concatRegPeers(server regServer) string {
	var peerList string

	for i := 0; i < server.peerNum; i++ {
		peerList += fmt.Sprintf("%s\n", server.peerList[i])
	}

	return peerList
}
