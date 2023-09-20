package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

/*
TestReadAll test read file
*/
func TestReadAll(t *testing.T) {

	all, err := io.ReadAll(strings.NewReader("sss"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res:%s\n", all)
}

/*
TestReadDir Read dir files
*/
func TestReadDir(t *testing.T) {
	dirs, err := os.ReadDir("../io")
	if err != nil {
		t.Error(err)
		return
	}
	for _, file := range dirs {
		fmt.Printf("%v", file)
	}
}

/*
TestReadFile read from file
*/
func TestReadFile(t *testing.T) {
	file, err := os.ReadFile("./writeAt.txt")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res:%s\n", file)
}

/*
TestWriteFile Test write file
*/
func TestWriteFile(t *testing.T) {
	if err := os.WriteFile("writefile.txt", []byte("TestWriteFile"), 0666); err != nil {
		t.Error(err)
		return
	}
}
