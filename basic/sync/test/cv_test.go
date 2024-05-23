package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 题目：编写一个程序，使用 Goroutine 和通道实现以下场景：
// 有两个 Goroutine，分别为 producer 和 consumer。
// producer Goroutine 会生成一系列数字，并将其发送到一个通道 ch 中。
// consumer Goroutine 会从通道 ch 中接收数字，并计算它们的平方，并将平方结果发送到另一个通道 result 中。
// 主函数会从 result 通道中接收结果，并打印出来。
func producer(ch chan int) {
	for i := 0; i < 6; i++ {
		// 设计随机种子
		rand.NewSource(time.Now().UnixNano())
		randomNumber := rand.Intn(10) + 1 // 生成 1 到 10 之间的随机数
		ch <- randomNumber
	}
	close(ch)
}
func consumer(ch chan int, result chan int, wg *sync.WaitGroup) {
	for num := range ch {
		result <- num * num
	}
	wg.Done()
}
func Test(t *testing.T) {
	ch := make(chan int)
	result := make(chan int)
	var wg sync.WaitGroup
	go producer(ch)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go consumer(ch, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for res := range result {
		fmt.Println(res)
	}
}
