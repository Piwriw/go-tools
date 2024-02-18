package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func main() {
	//getName()
}
func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Please in your name"
	fmt.Fprintf(w, msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You did'nt enter your name")
	}
	return name, nil
}
