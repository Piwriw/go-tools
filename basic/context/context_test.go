package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go operation(ctx)

	time.Sleep(3 * time.Second) // 在取消前等待一段时间
	cancel()                    // 取消操作

	// 给足够的时间来查看 operation 如何响应取消信号
	time.Sleep(10 * time.Second)
}
func operation(ctx context.Context) {
	select {
	case <-time.After(500 * time.Millisecond): // 模拟耗时操作
		fmt.Println("operation completed")
	case <-ctx.Done():
		fmt.Println("operation canceled")
	}
}
