package main

import "syscall"

func NewUTSNS() {
	NewLinuxNamespace(syscall.CLONE_NEWUTS)
}
