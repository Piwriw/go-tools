package main

import "syscall"

func NewUserNS() {
	NewLinuxNamespace(syscall.CLONE_NEWNET)
}
