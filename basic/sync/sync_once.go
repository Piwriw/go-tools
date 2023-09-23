package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var once sync.Once
	body := func() {
		fmt.Println("Do once")
	}
	for i := 0; i < 10; i++ {
		once.Do(body)
	}
	time.Sleep(time.Second * 10)
}
