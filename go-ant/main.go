package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log/slog"
	"sync/atomic"
	"time"
)

var sum int32
var p *ants.Pool

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

//func ResizePool(pool *ants.Pool) {
//	fmt.Println(pool)
//	slog.Info("current cap is ", slog.Int("free", p.Free()))
//	slog.Info("current cap is ", slog.Int("cap", p.Cap()))
//	if p.Free() < 1024 {
//		p.Tune(pool.Cap() * 2)
//		slog.Info("current cap is ", slog.Int("free", p.Free()))
//		return
//	}
//	newCap := p.Cap() + p.Cap()/4
//	slog.Info("current cap is ", slog.Int("cap", newCap))
//	p.Tune(newCap)
//	return
//}

func demoFunc() {
	time.Sleep(30 * time.Second)
	fmt.Println("Hello World!")
	fmt.Println("current", p.Running())
}

func main() {

	var err error
	// ants will pre-malloc the whole capacity of pool when you invoke this method
	p, err = ants.NewPool(5, ants.WithPreAlloc(false), ants.WithNonblocking(true))
	if err != nil {
		panic(err)
	}
	slog.Info("pool start", p.Cap())
	p.Tune(499)
	slog.Info("pool start", p.Cap()/2)
	for i := 0; i < 200000; i++ {
		p.Submit(demoFunc)
	}

	time.Sleep(60 * time.Minute)
}

// WatchPool 监控Pool并且自动扩容，要求pool没有实现分配内存
func WatchPool(pool *ants.Pool, checkTimer time.Duration) {
	ticker := time.NewTicker(checkTimer)
	for {
		select {
		case <-ticker.C:
			slog.Info("Current Pool Detail", slog.Int("free count", pool.Free()), slog.Int("running count", pool.Running()), slog.Int("cap count", pool.Cap()))
			isResize := ResizePool(pool)
			if isResize {
				slog.Info("Pool is resize...,Current Pool Detail", slog.Int("free count", pool.Free()), slog.Int("running count", pool.Running()), slog.Int("cap count", pool.Cap()))
			}
			ticker.Reset(checkTimer)
		}
	}
}

// ResizePool Pool 扩容
/*
 当pool  空余还有 1/4 不扩容
         空余超过 1/2 并且 Cap大于1024 进行缩容 1/4 防止过多的内存占用
		 空余小于 1/4 并且cap小于1024  按照cap 俩倍扩容
		 剩下的情况 按照cap1.5倍进行扩容
*/
func ResizePool(pool *ants.Pool) bool {
	if pool.Free() > pool.Cap()/4 {
		return false
	} else if pool.Free() > pool.Cap()/2 && pool.Cap() > 1024 {
		newCap := pool.Cap() - pool.Cap()/4
		pool.Tune(newCap)
		return true
	}
	if pool.Cap() < 1024 && pool.Free() < pool.Cap()/4 {
		pool.Tune(pool.Cap() * 2)
	} else {
		newCap := pool.Cap() + pool.Cap()/2
		pool.Tune(newCap)
	}
	return true
}

func InitNonblockingPool(poolSize int) (*ants.Pool, error) {
	// 初始化 taskPool 用于异步插入数据库,预先为taskPool分配内存且使用非阻塞模式
	taskPool, err := ants.NewPool(poolSize, ants.WithPreAlloc(true), ants.WithNonblocking(true))
	if err != nil {
		slog.Error("taskPool init error", slog.Any("err", err))
		return nil, err
	}
	return taskPool, nil
}

// InitBlockingPool  阻塞模式连接池
func InitBlockingPool(poolSize int) (*ants.Pool, error) {
	pool, err := ants.NewPool(poolSize, ants.WithPreAlloc(true), ants.WithNonblocking(false))
	if err != nil {
		slog.Error("taskPool init error", slog.Any("err", err))
		return nil, err
	}
	return pool, nil
}

// InitNonblockingPoolNoPreAlloc 非阻塞 并且不提前分配空间
func InitNonblockingPoolNoPreAlloc(poolSize int) (*ants.Pool, error) {
	taskPool, err := ants.NewPool(poolSize, ants.WithPreAlloc(false), ants.WithNonblocking(true))
	if err != nil {
		slog.Error("taskPool init error", slog.Any("err", err))
		return nil, err
	}
	return taskPool, nil
}
