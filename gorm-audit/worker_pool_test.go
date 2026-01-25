package audit

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/piwriw/gorm/gorm-audit/types"
)

// MockHandler 用于测试的模拟处理器
type MockHandler struct {
	CallCount int32
	Delay     time.Duration
	Error     error
}

func (m *MockHandler) Handle(ctx context.Context, event *handler.Event) error {
	atomic.AddInt32(&m.CallCount, 1)
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	return m.Error
}

func TestNewWorkerPool(t *testing.T) {
	mockHandler := &MockHandler{}
	config := &WorkerPoolConfig{
		WorkerCount: 2,
		QueueSize:   100,
		Timeout:     100,
	}

	wp := NewWorkerPool(mockHandler, config)

	if wp == nil {
		t.Fatal("expected non-nil worker pool")
	}

	// 等待 workers 启动
	time.Sleep(50 * time.Millisecond)

	stats := wp.Stats()
	if stats.WorkerCount != 2 {
		t.Errorf("expected 2 workers, got %d", stats.WorkerCount)
	}
	if stats.QueueCapacity != 100 {
		t.Errorf("expected queue capacity 100, got %d", stats.QueueCapacity)
	}

	// 关闭工作池
	wp.Close()
}

func TestWorkerPoolDispatch(t *testing.T) {
	mockHandler := &MockHandler{}
	config := &WorkerPoolConfig{
		WorkerCount: 2,
		QueueSize:   10,
		Timeout:     1000,
	}

	wp := NewWorkerPool(mockHandler, config)
	defer wp.Close()

	// 分发事件
	event := &handler.Event{
		Timestamp: time.Now().Format(time.RFC3339),
		Operation: types.OperationCreate,
		Table:     "users",
	}

	// 分发应该成功
	if !wp.Dispatch(event) {
		t.Error("expected dispatch to succeed")
	}

	// 等待处理
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&mockHandler.CallCount) != 1 {
		t.Errorf("expected handler to be called once, got %d", mockHandler.CallCount)
	}
}

func TestWorkerPoolQueueFull(t *testing.T) {
	mockHandler := &MockHandler{
		Delay: 100 * time.Millisecond, // 模拟慢处理
	}
	config := &WorkerPoolConfig{
		WorkerCount: 1,
		QueueSize:   2, // 小队列
		Timeout:     1000,
	}

	wp := NewWorkerPool(mockHandler, config)
	defer wp.Close()

	events := make([]*handler.Event, 5)
	for i := 0; i < 5; i++ {
		events[i] = &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: types.OperationCreate,
			Table:     "users",
		}
	}

	// 分发所有事件
	successCount := 0
	for _, event := range events {
		if wp.Dispatch(event) {
			successCount++
		}
	}

	// 队列只能容纳 2 个，所以最多 2 个成功
	if successCount > 2 {
		t.Errorf("expected max 2 successful dispatches, got %d", successCount)
	}

	// 等待队列清空
	time.Sleep(300 * time.Millisecond)

	stats := wp.Stats()
	if stats.QueueLength > 0 {
		t.Logf("Queue still has %d items after wait", stats.QueueLength)
	}
}

func TestWorkerPoolPanicRecovery(t *testing.T) {
	// 包装为会 panic 的处理器
	wrappedHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
		panic("test panic in worker")
	})

	config := &WorkerPoolConfig{
		WorkerCount: 2,
		QueueSize:   10,
		Timeout:     1000,
	}

	wp := NewWorkerPool(wrappedHandler, config)
	defer wp.Close()

	event := &handler.Event{
		Timestamp: time.Now().Format(time.RFC3339),
		Operation: types.OperationCreate,
		Table:     "users",
	}

	// 分发事件
	wp.Dispatch(event)

	// 等待处理和 panic 恢复
	time.Sleep(200 * time.Millisecond)

	// 工作池应该仍然运行
	stats := wp.Stats()
	if stats.WorkerCount != 2 {
		t.Errorf("expected 2 workers after panic, got %d", stats.WorkerCount)
	}
}

func TestWorkerPoolTimeout(t *testing.T) {
	slowMockHandler := &MockHandler{
		Delay: 200 * time.Millisecond, // 超过超时时间
		Error: errors.New("timeout error"),
	}

	config := &WorkerPoolConfig{
		WorkerCount: 1,
		QueueSize:   10,
		Timeout:     100, // 100ms 超时
	}

	wp := NewWorkerPool(slowMockHandler, config)
	defer wp.Close()

	event := &handler.Event{
		Timestamp: time.Now().Format(time.RFC3339),
		Operation: types.OperationCreate,
		Table:     "users",
	}

	// 分发事件
	dispatched := wp.Dispatch(event)
	if !dispatched {
		t.Error("expected dispatch to succeed")
	}

	// 等待处理（应该超时）
	time.Sleep(300 * time.Millisecond)

	// 即使超时，worker 应该继续运行
	stats := wp.Stats()
	if stats.WorkerCount != 1 {
		t.Errorf("expected 1 worker after timeout, got %d", stats.WorkerCount)
	}
}

func TestWorkerPoolClose(t *testing.T) {
	mockHandler := &MockHandler{}
	config := &WorkerPoolConfig{
		WorkerCount: 2,
		QueueSize:   10,
		Timeout:     1000,
	}

	wp := NewWorkerPool(mockHandler, config)

	// 分发一些事件
	for i := 0; i < 3; i++ {
		event := &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: types.OperationCreate,
			Table:     "users",
		}
		wp.Dispatch(event)
	}

	// 关闭工作池
	wp.Close()

	// 再次关闭应该安全
	wp.Close()

	stats := wp.Stats()
	if stats.WorkerCount != 0 {
		t.Error("expected 0 workers after close")
	}
}
