package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func NewLinuxNamespace(cloneflags uintptr) {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: cloneflags,
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
