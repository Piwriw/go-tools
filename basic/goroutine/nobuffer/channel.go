package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
基于无buffer 的channel分离机制
*/
func doGoroutine(ch chan int) {
	for task := range ch {
		fmt.Println("go task", task, "goroutine count", runtime.NumGoroutine())
		// 结束了一个任务
		wg.Done()
	}

}
func sendTask(task int, ch chan int) {
	wg.Add(1)
	// 任务发给channel
	ch <- task
}

var wg sync.WaitGroup

func main() {
	//无buffer channel
	ch := make(chan int)
	// Goroutine数量
	goCnt := 3
	for i := 0; i < goCnt; i++ {
		go doGoroutine(ch)
	}
	// 业务数量
	taskCnt := 10
	for i := 0; i < taskCnt; i++ {
		sendTask(i, ch)
	}
	wg.Wait()
}
