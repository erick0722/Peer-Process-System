package peer

import (
	"fmt"
	"strconv"
)

/**
* Return a string representation of all the peers in our peerlist.
*
* @returns {string} The process' peerlist as a string.
 */
func ConcatPeerList() string {

	peerNum := strconv.Itoa(len(peerList))
	peerListStr := fmt.Sprintf("%s\n", peerNum)

	for i := 0; i < len(peerList); i++ {
		peerListStr += fmt.Sprintf("%s\n", peerList[i].address)
	}

	return peerListStr
}

/**
* Return a string representation of all the peers our process received from
* a UDP/IP message.
*
* @returns {string} The peers our process received from a UDP/IP message as a string.
 */
func ConcatRecvPeerList() string {

	peerNum := strconv.Itoa(len(recievedPeers))
	recvPeerListStr := fmt.Sprintf("%s\n", peerNum)

	for i := 0; i < len(recievedPeers); i++ {
		recvPeerListStr += fmt.Sprintf(
			"%s %s %s\n",
			recievedPeers[i].source,
			recievedPeers[i].received,
			recievedPeers[i].timeReceived.Format("2006-01-02 15:04:05"))
	}

	return recvPeerListStr
}

/**
* Return a string representation of all the peers our process sent to other processes as a
* UDP/IP message.
*
* @returns {string} The peers our process sent as a UDP/IP message in a string format.
 */
func ConcatPeersSent() string {

	peerNum := strconv.Itoa(len(peersSent))
	peersSentStr := fmt.Sprintf("%s\n", peerNum)

	for i := 0; i < len(peersSent); i++ {
		peersSentStr += fmt.Sprintf(
			"%s %s %s\n",
			peersSent[i].sentTo,
			peersSent[i].peer,
			peersSent[i].timeSent.Format("2006-01-02 15:04:05"))
	}

	return peersSentStr
}

/**
* Return a string representation of all the snips or messages our peer process
* received from other peers.
*
* @returns {string} All the messages our process received as a string.
 */
func ConcatSnipList() string {

	snipNum := strconv.Itoa(len(snipList))
	snipListStr := fmt.Sprintf("%s\n", snipNum)

	for i := 0; i < len(snipList); i++ {
		snipListStr += fmt.Sprintf(
			"%s %s %s\n",
			snipList[i].timeStamp,
			snipList[i].content,
			snipList[i].source)
	}

	return snipListStr
}

// =============================================================
