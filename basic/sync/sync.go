package main

import (
	"sync"
	"time"
)

var m *sync.RWMutex
var value int

func main() {
	m = new(sync.RWMutex)
	go read(1)
	go write(2)
	go read(3)
	time.Sleep(5 * time.Second)
}

func read(i int) {
	m.Lock()
	time.Sleep(1 * time.Second)
	println("val: ", value)
	m.RUnlock()
}
func write(i int) {
	m.Lock()
	value = 10
	println("val: ", value)
	time.Sleep(1 * time.Second)
	m.RUnlock()
}
