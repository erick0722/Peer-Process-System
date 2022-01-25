package main

import (
	"testing"
)

func TestFileRead(t *testing.T) {
	sourceCode, err := readFile("main.go")
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(sourceCode)
	}

}

func TestParseCodeResponse(t *testing.T) {
	codeResponse := parseCodeResponse()
	t.Log(codeResponse)
}