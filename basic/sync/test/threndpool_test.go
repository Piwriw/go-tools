package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
题目：编写一个并发程序，实现一个简单的线程池。

线程池的大小固定为3个线程。
线程池可以接收任务，并并发执行这些任务。
每个任务是一个函数，接收一个整数参数并返回一个整数结果。
线程池应该返回每个任务的执行结果。
主函数应该等待所有任务完成后打印结果。
*/

type Task func(int) int

type ThreadPool struct {
	size     int
	tasks    chan Task
	results  chan int
	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewThreadPool(size int) *ThreadPool {
	pool := &ThreadPool{
		size:     size,
		tasks:    make(chan Task, size),
		results:  make(chan int),
		shutdown: make(chan struct{}),
	}

	pool.startWorkers()

	return pool
}

func (pool *ThreadPool) startWorkers() {
	for i := 0; i < pool.size; i++ {
		pool.wg.Add(1)
		go func(idx int) {
			defer pool.wg.Done()
			for {
				select {
				case task := <-pool.tasks:
					result := task(idx)
					pool.results <- result
				case <-pool.shutdown:
					fmt.Printf("Go:%d is Dead", idx)
					return
				}
			}
		}(i + 1)
	}
}

func (pool *ThreadPool) Shutdown() {
	fmt.Println("Shutdown")
	pool.wg.Wait()    // 等待所有任务完成
	close(pool.tasks) // 关闭任务通道，表示没有更多任务
	close(pool.shutdown)
	close(pool.results)
}

func (pool *ThreadPool) Submit(task Task) {
	pool.tasks <- task
}

func TestPV(t *testing.T) {
	pool := NewThreadPool(3)

	go func() {
		for {
			// 等待任务完成并打印结果
			select {
			case result, ok := <-pool.results:
				fmt.Printf("Result:%d，ok:%v\n", result, ok)
				if !ok {
					pool.shutdown <- struct{}{}
					return
				}
			}
		}
	}()

	// 提交任务
	for i := 0; i < 10; i++ {
		pool.Submit(func(n int) int {
			return n * n
		})
	}
	go pool.Shutdown()
	time.Sleep(10 * time.Second)
}
