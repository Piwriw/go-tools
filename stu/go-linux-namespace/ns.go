package main

import (
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func NewLinuxNamespace(cloneflags uintptr) {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &unix.SysProcAttr{
		cloneflags: cloneflags,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	NewUTSNS()
}
