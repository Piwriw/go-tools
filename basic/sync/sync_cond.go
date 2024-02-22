package main

import (
	"log"
	"sync"
	"time"
)

/*
sync.Cond 一组Goroutine 需要等待的情况
*/
func main() {
	testSyncCond()
}

var done = false

func readCond(name string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		c.Wait()
	}
	log.Println(name, "starts reading")
	c.L.Unlock()
}

func writeCond(name string, c *sync.Cond) {
	log.Println(name, "starts writing")
	time.Sleep(time.Second)
	c.L.Lock()
	done = true
	c.L.Unlock()
	log.Println(name, "wakes all")
	c.Broadcast()
}

func testSyncCond() {
	cond := sync.NewCond(&sync.Mutex{})

	go readCond("reader1", cond)
	go readCond("reader2", cond)
	go readCond("reader3", cond)
	writeCond("writer", cond)

	time.Sleep(time.Second * 3)
}
