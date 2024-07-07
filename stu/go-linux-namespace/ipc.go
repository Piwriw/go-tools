package main

import "syscall"

func NewIPCNS() {
	NewLinuxNamespace(syscall.CLONE_NEWIPC)
}
