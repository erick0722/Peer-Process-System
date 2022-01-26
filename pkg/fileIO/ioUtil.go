package fileIO

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// https://golang.cafe/blog/how-to-list-files-in-a-directory-in-go.html
func readDirectory(dirName string) string {
	files, err := ioutil.ReadDir(dirName)
	var sourceCode string
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		sourceCode += readFile(dirName + file.Name())
	}
	return sourceCode
}

// Format and return a string to match the code response format
func ParseCodeResponse() string {
	var language string = "golang"
	var endOfCode string = "..."

	sourceCode := readDirectory("pkg/fileIO/")
	sourceCode += readDirectory("pkg/tcp/")
	sourceCode += readDirectory("cmd/Iteration1/")
	codeResponse := fmt.Sprintf("%s\n%s\n%s\n", language, sourceCode, endOfCode)
	return codeResponse
}

// Read a file's content line-by-line and return it as string, separated by new-lines.
// Code was inspired from the following link: https://golangdocs.com/reading-files-in-golang
func readFile(fileName string) string {
	var sourceCode string = ""
	file, _ := os.Open(fileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sourceCode += fmt.Sprintf("%s\n", scanner.Text())
	}

	return sourceCode
}
