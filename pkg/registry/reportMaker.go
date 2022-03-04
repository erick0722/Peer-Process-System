package registry

import (
	"559Project/pkg/peer"
	"fmt"
	"strconv"
)

// Generate a report for the list of peers and sources
func generateReport() string {

	// Return nothing if no peers have been received (address empty)
	if reg.address == "" {
		return "0\n0\n\n0\n"
	}

	// Add the current list of peers to the report
	report := peer.ConcatPeerList()

	// Add the lists we have received to the report
	report += fmt.Sprintf("1\n%s\n%s\n%s\n", reg.address, reg.timeReceived.Format("2006-01-02 15:04:05"), strconv.Itoa(reg.peerNum))
	report += concatRegPeers()

	// Add the peers received via UDP/IP to the report
	report += peer.ConcatRecvPeerList()

	// Add all the ppers we sent via UDP/IP to the report
	report += peer.ConcatPeersSent()

	// Add all snippets we received to the report
	report += peer.ConcatSnipList()

	fmt.Printf("%s", report)
	return report
}

/**
* Combine all the peers the registry gave to our peerlist as a string
*
* @param server {regServer} The data for the registry server
* @returns {string} A string type of all the registry sent peers
 */
func concatRegPeers() string {
	var peerList string

	for i := 0; i < reg.peerNum; i++ {
		peerList += fmt.Sprintf("%s\n", reg.peerList[i])
	}

	return peerList
}

// =============================================================
