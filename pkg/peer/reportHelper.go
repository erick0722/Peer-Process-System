package peer

import (
	"fmt"
	"strconv"
)

func ConcatPeerList() string {

	peerNum := strconv.Itoa(len(peerList))
	peerListStr := fmt.Sprintf("%s\n", peerNum)

	for i := 0; i < len(peerList); i++ {
		peerListStr += fmt.Sprintf("%s\n", peerList[i].address)
	}

	return peerListStr
}

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
