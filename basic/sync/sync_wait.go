package main

import (
	"fmt"
	"sync"
)

func main() {
	wp := sync.WaitGroup{}
	wp.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("Done")
			wp.Done()
		}()
	}
	wp.Wait()
	fmt.Println("wait end")
}
