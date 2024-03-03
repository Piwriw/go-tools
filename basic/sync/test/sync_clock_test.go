package test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

type Counter struct {
	value int64
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

func (c *Counter) GetValue() int64 {
	return atomic.LoadInt64(&c.value)
}
func TestClock(t *testing.T) {
	var wg sync.WaitGroup
	counter := Counter{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			counter.Increment()
			fmt.Println(counter.GetValue())
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			counter.Decrement()
			fmt.Println(counter.GetValue())
		}
	}()
	wg.Wait()
	fmt.Println(counter.GetValue())
}
