// =============================================================
/*
	CPSC 559 - Iteration 2
	ioUtil.go

	Erick Yip
	Chris Chen
*/

package fileIO

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// Read all files in a directory and return their content as string
// Code was inspired from: https://golang.cafe/blog/how-to-list-files-in-a-directory-in-go.html
func readDirectory(dirName string) string {

	var sourceCode string
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			sourceCode += readFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return sourceCode
}

// Format and return a string to match the code response format
func ParseCodeResponse() string {
	var language string = "golang"
	var endOfCode string = "..."

	sourceCode := readDirectory("pkg/")
	sourceCode += readDirectory("cmd/")

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

// =============================================================
