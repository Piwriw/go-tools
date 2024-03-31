package main

import (
	"fmt"
	"runtime"
	"time"
)

func doGoroutine(i int, ch chan bool) {
	fmt.Println("go func", i, "goroutine count", runtime.NumGoroutine())
	// 结束了一个任务
	<-ch
}

/*
基于buffer的channel控制最大协程数量
*/
func main() {
	task_cnt := 10
	// 容量控制了 Goroutine 的数量
	ch := make(chan bool, 3)
	// for的数据决定了Goroutine的创建速度
	for i := 0; i < task_cnt; i++ {
		ch <- true
		go doGoroutine(i, ch)
	}
	// task_cnt 数量太小，主线程会优先退出，需要阻塞等待
	time.Sleep(100 * time.Second)
}
