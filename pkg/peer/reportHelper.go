package peer

import "fmt"

func ConcatPeerList(listPeers []peerStruct) string {
	var peerList string

	for i := 0; i < len(listPeers); i++ {
		peerList += fmt.Sprintf("%s\n", listPeers[i].address)
	}

	return peerList
}

func ConcatRecvPeerList(recvPeerList []receivedEvent) string {
	var peerList string

	for i := 0; i < len(recvPeerList); i++ {
		peerList += fmt.Sprintf(
			"%s %s %s\n",
			recvPeerList[i].source,
			recvPeerList[i].received,
			recvPeerList[i].timeReceived.Format("2006-01-02 15:04:05"))
	}

	return peerList
}

func ConcatPeersSent(peersSent []sentEvent) string {
	var peerList string

	for i := 0; i < len(peersSent); i++ {
		peerList += fmt.Sprintf(
			"%s %s %s\n",
			peersSent[i].sentTo,
			peersSent[i].peer,
			peersSent[i].timeSent.Format("2006-01-02 15:04:05"))
	}

	return peerList
}

func ConcatSnipList(snipList []snip) string {
	var peerList string

	for i := 0; i < len(snipList); i++ {
		peerList += fmt.Sprintf(
			"%s %s %s\n",
			snipList[i].timeStamp,
			snipList[i].content,
			snipList[i].source)
	}

	return peerList
}
