package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

func TaskFunc() {
	time.Sleep(100 * time.Second)
	fmt.Println("Task Func")
}
func main() {
	// 开启
	p, _ := ants.NewPool(1000, ants.WithNonblocking(true))
	group := sync.WaitGroup{}
	group.Add(1000)
	for i := 0; i < 1000; i++ {
		p.Submit(TaskFunc)
		group.Done()
	}

	fmt.Println(p.Running())
	group.Wait()

}
