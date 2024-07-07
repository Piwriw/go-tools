package main

import "syscall"

func NewUserNS() {
	NewLinuxNamespace(syscall.NEWUSER)
}
