package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"559Project/pkg/tcp"
)

// Struct to store the server's address, list of peers, number of peers, and the time the peers are received
type regServer struct {
	address      string
	peerList     []string
	peerNum      int
	timeReceived string
}

func main() {

	var registry regServer
	registry.address = "localhost:55921"

	var teamName string = "It Takes Two\n" // Our team name

	// Connect to the server via TCP
	conn := tcp.InitializeTCP(registry.address)

forLoop:
	for {
		reply := make([]byte, 1024)

		// Receive the server's response
		_, err := conn.Read(reply)

		tcp.CheckError(err)

		// Trim off the extra characters
		serverReply := string(bytes.Trim(reply, "\x00"))
		var clientMessage string

		fmt.Printf("Server message = %s", serverReply)

		// Check the server's response and take corresponding actions
		switch {
		case serverReply == "get team name\n":
			clientMessage = "Sending team name...\n"
			tcp.SendMessage(teamName, conn)
		case serverReply == "get code\n":
			clientMessage = "Sending code...\n"
			codeResponse := parseCodeResponse()
			tcp.SendMessage(codeResponse, conn)
		case serverReply == "receive peers\n":
			clientMessage = "Receiving peers...\n"
			registry.peerNum, registry.peerList = receivePeers(conn)
			registry.timeReceived = saveCurrTime()
			clientMessage = "Peers received\n"
		case serverReply == "get report\n":
			clientMessage = "Sending report...\n"
			report := generateReport(registry)
			tcp.SendMessage(report, conn)
		case serverReply == "close\n":
			fmt.Printf("%s", "Closing...\n")
			conn.Close()
			break forLoop
		default:
			clientMessage = "Unknown message\n"
		}

		fmt.Printf("%s", clientMessage)
	}
}

// Receive the peers from the server and store them into the peerList
func receivePeers(conn net.Conn) (int, []string) {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	tcp.CheckError(err)

	numPeers, _ := strconv.Atoi(strings.Split(string(reply), "\n")[0])

	fmt.Printf("numPeers=%d\n", numPeers)

	peerList := make([]string, numPeers)
	for i := 0; i < numPeers; i++ {
		peerList[i] = strings.Split(string(reply), "\n")[i+1]
		fmt.Printf("peerList[%d]=%s\n", i, peerList[i])
	}
	return numPeers, peerList
}

// Save the current time to a string
func saveCurrTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Format and return a string to match the code response format
func parseCodeResponse() string {
	var language string = "golang"
	var sourceFile string = "main.go"
	var endOfCode string = "..."

	sourceCode, _ := readFile(sourceFile)

	codeResponse := fmt.Sprintf("%s\n%s\n%s\n", language, sourceCode, endOfCode)
	return codeResponse
}

func readCode() (string, error) {
	// for recursive file in ./cmd, read file and path to string

	// for recursive file in ./pkg, read file and path to string

	return "", nil
}

// Read a file's content line-by-line and return it as string, separated by new-lines.
// Code was inspired from the following link: https://golangdocs.com/reading-files-in-golang
func readFile(srcName string) (string, error) {
	var sourceCode string = ""
	file, err := os.Open(srcName)
	tcp.CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sourceCode += scanner.Text() + "\n"
	}

	return sourceCode, nil
}

// Generate a report for the list of peers and sources
func generateReport(registry regServer) string {
	// Convert the number of peers to string
	peerNumString := strconv.Itoa(registry.peerNum)
	report := peerNumString + "\n"

	for i := 0; i < registry.peerNum; i++ {
		report += registry.peerList[i] + "\n"
	}

	// report += fmt.Sprintf("1\n%s\n%s\n%s\n", registry.address, registry.timeReceived, peerNumString)
	report += "1\n" + registry.address + "\n" + registry.timeReceived + "\n" + peerNumString + "\n"

	for i := 0; i < registry.peerNum; i++ {
		report += registry.peerList[i] + "\n"
	}
	fmt.Printf("Report:\n%s", report)

	return report
}
