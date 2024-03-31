package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
基于sync机制 和buffer channel实现最大协程数量控制
*/
var wg sync.WaitGroup

func doGoroutine(i int, ch chan bool) {
	fmt.Println("go func", i, "goroutine count", runtime.NumGoroutine())
	// 结束了一个任务
	wg.Done()
	<-ch
}
func main() {
	task_cnt := 10
	ch := make(chan bool, 3)
	for i := 0; i < task_cnt; i++ {
		wg.Add(1)
		ch <- true
		go doGoroutine(i, ch)
	}
	wg.Wait()
}
