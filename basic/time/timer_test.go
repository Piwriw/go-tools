package main

import (
	"fmt"
	"testing"
	"time"
)

func dowork() {
	fmt.Println("Do Work")
}
func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-timer.C:
			dowork()
			// 重置定时器
			timer.Reset(time.Second * 5)
		}
	}
	timer.Stop()
}
func TestTimerAfter(t *testing.T) {
	timeChannel := time.After(10 * time.Second)
	select {
	case <-timeChannel:
		dowork()
	}
}
