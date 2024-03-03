package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSync(t *testing.T) {
	//var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			count += i
		}
	}()
	go func() {
		defer wg.Done()
		for i := 101; i <= 200; i++ {
			count += i
		}
	}()
	wg.Wait()
	t.Log(count)
}

//func TestSync2(t *testing.T) {
//	var sumChann chan int
//	var wg sync.WaitGroup
//	sumChann = make(chan int, 2) // 使用带缓冲的通道，避免死锁
//	wg.Add(2)
//
//	go doCount(1, 100, sumChann)
//	go doCount(101, 200, sumChann)
//
//	wg.Wait()
//	close(sumChann)
//
//	sum := 0
//	for i := range sumChann {
//		sum += i
//	}
//
//	t.Log("sum", sum)
//}
//
//func doCount(start, end int, ch chan<- int) { // 修改函数签名，接收带缓冲的通道作为参数
//	count := 0
//	for i := start; i <= end; i++ {
//		count += i
//	}
//	ch <- count // 将结果写入通道
//	wg.Done()
//}

func TestSync3(t *testing.T) {
	var wg sync.WaitGroup
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i += 2 {
			<-ch1
			fmt.Println(i)
			ch2 <- struct{}{}
		}
		<-ch1
	}()

	go func() {
		for i := 2; i <= 100; i += 2 {
			<-ch2
			fmt.Println(i)
			ch1 <- struct{}{}
		}
		defer wg.Done()

	}()
	ch1 <- struct{}{}

	wg.Wait()

}

func TestSync4(t *testing.T) {
	var wg sync.WaitGroup
	send := make(chan int)
	revice := make(chan int)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case num := <-send:
				fmt.Println("Consumer received:", num)
				// 处理消费逻辑
			case <-revice:
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			num := rand.Intn(100)
			send <- num
			fmt.Println("Producer sent:", num)
			time.Sleep(time.Second)
		}
		close(revice)
	}()

	wg.Wait()

}

// 编程题：编写一个并发程序，实现生产者-消费者模型，并限制缓冲区大小。
// 要求：
//
// 实现一个有限大小的缓冲区，用于生产者和消费者之间的数据传输。
// 当缓冲区满时，生产者应该等待，直到有空间可用。
// 当缓冲区为空时，消费者应该等待，直到有数据可用
func TestSync5(t *testing.T) {
	taskCh := make(chan struct{})
	doneCh := make(chan struct{})

	go worker(taskCh, doneCh)

	for i := 0; i < 5; i++ {
		taskCh <- struct{}{} // 发送任务给worker
		<-doneCh             // 阻塞等待任务完成
	}
}

func worker(taskCh <-chan struct{}, doneCh chan<- struct{}) {
	for {
		<-taskCh // 等待接收任务
		fmt.Println("Worker is working...")
		time.Sleep(2 * time.Second)
		fmt.Println("Work done!")
		doneCh <- struct{}{} // 任务完成，发送信号
	}
}

//要求：
//
//线程池的大小固定为3个线程。
//线程池可以接收任务，并并发执行这些任务。
//每个任务是一个函数，接收一个整数参数并返回一个整数结果。
//线程池应该返回每个任务的执行结果。
//主函数应该等待所有任务完成后打印结果。
