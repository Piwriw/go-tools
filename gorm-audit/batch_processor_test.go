package audit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

// MockBatchHandler 模拟处理器
type MockBatchHandler struct {
	CallCount int
}

func (m *MockBatchHandler) Handle(ctx context.Context, event *handler.Event) error {
	m.CallCount++
	return nil
}

func TestNewBatchProcessor(t *testing.T) {
	handler := &MockBatchHandler{}
	config := &WorkerPoolConfig{
		QueueSize:     100,
		BatchSize:     10,
		FlushInterval: 1 * time.Second,
		EnableBatch:   true,
	}

	bp := NewBatchProcessor(handler, config)

	if bp == nil {
		t.Fatal("expected non-nil BatchProcessor")
	}

	if bp.batchSize != 10 {
		t.Errorf("expected batch size 10, got %d", bp.batchSize)
	}

	if cap(bp.buffer) != 100 {
		t.Errorf("expected buffer capacity 100, got %d", cap(bp.buffer))
	}
}

func TestBatchProcessorSizeTrigger(t *testing.T) {
	mockHandler := &MockBatchHandler{}
	config := &WorkerPoolConfig{
		QueueSize:     100,
		BatchSize:     5, // 小批量便于测试
		FlushInterval: 10 * time.Second,
		EnableBatch:   true,
	}

	bp := NewBatchProcessor(mockHandler, config)
	bp.Start()
	defer bp.Close()

	// 发送 5 个事件
	for i := 0; i < 5; i++ {
		event := &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: handler.OperationCreate,
			Table:     "users",
		}
		bp.Dispatch(event)
	}

	// 等待异步处理
	time.Sleep(100 * time.Millisecond)

	// 验证所有事件都被处理
	if mockHandler.CallCount != 5 {
		t.Errorf("expected 5 handler calls, got %d", mockHandler.CallCount)
	}

	// 验证统计信息
	stats := bp.Stats()
	if stats.TotalEvents != 5 {
		t.Errorf("expected 5 total events, got %d", stats.TotalEvents)
	}
	if stats.TotalBatches != 1 {
		t.Errorf("expected 1 batch, got %d", stats.TotalBatches)
	}
}

func TestBatchProcessorTimeTrigger(t *testing.T) {
	mockHandler := &MockBatchHandler{}
	config := &WorkerPoolConfig{
		QueueSize:     100,
		BatchSize:     1000, // 大批量不会触发
		FlushInterval: 200 * time.Millisecond, // 短时间便于测试
		EnableBatch:   true,
	}

	bp := NewBatchProcessor(mockHandler, config)
	bp.Start()
	defer bp.Close()

	// 只发送 2 个事件（不足以触发批量）
	for i := 0; i < 2; i++ {
		event := &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: handler.OperationCreate,
			Table:     "users",
		}
		bp.Dispatch(event)
	}

	// 等待时间触发
	time.Sleep(300 * time.Millisecond)

	// 验证事件被处理
	if mockHandler.CallCount != 2 {
		t.Errorf("expected 2 handler calls, got %d", mockHandler.CallCount)
	}

	stats := bp.Stats()
	if stats.TotalBatches != 1 {
		t.Errorf("expected 1 batch, got %d", stats.TotalBatches)
	}
}

func TestBatchProcessorClose(t *testing.T) {
	mockHandler := &MockBatchHandler{}
	config := &WorkerPoolConfig{
		QueueSize:     100,
		BatchSize:     100,
		FlushInterval: 1 * time.Second,
		EnableBatch:   true,
	}

	bp := NewBatchProcessor(mockHandler, config)
	bp.Start()

	// 发送 3 个事件
	for i := 0; i < 3; i++ {
		event := &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: handler.OperationCreate,
			Table:     "users",
		}
		if !bp.Dispatch(event) {
			t.Errorf("failed to dispatch event %d", i)
		}
	}

	// 给一点时间让事件从 channel 被消费到 batch 中
	time.Sleep(50 * time.Millisecond)

	// 关闭应该刷新剩余事件
	bp.Close()

	// Close() 内部已经等待了 wg.Wait()，所以不需要额外等待
	// 但为了确保 goroutine 完全退出，给一小段时间
	time.Sleep(10 * time.Millisecond)

	// 验证所有事件都被处理
	if mockHandler.CallCount != 3 {
		t.Errorf("expected 3 handler calls after close, got %d", mockHandler.CallCount)
	}

	// 注意：Close() 不是幂等的，多次调用会 panic
	// 实际使用中应该保证只调用一次
}

func TestBatchProcessorErrorHandling(t *testing.T) {
	// 创建会返回错误的 handler
	errorHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
		return fmt.Errorf("handler error")
	})

	config := &WorkerPoolConfig{
		QueueSize:     100,
		BatchSize:     3,
		FlushInterval: 1 * time.Second,
		EnableBatch:   true,
	}

	bp := NewBatchProcessor(errorHandler, config)
	bp.Start()
	defer bp.Close()

	// 发送 3 个事件
	for i := 0; i < 3; i++ {
		event := &handler.Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Operation: handler.OperationCreate,
			Table:     "users",
		}
		bp.Dispatch(event)
	}

	// 等待处理
	time.Sleep(100 * time.Millisecond)

	// 验证错误被记录
	stats := bp.Stats()
	if stats.TotalErrors != 3 {
		t.Errorf("expected 3 errors, got %d", stats.TotalErrors)
	}

	// 验证所有事件都被尝试处理
	if stats.TotalEvents != 3 {
		t.Errorf("expected 3 total events, got %d", stats.TotalEvents)
	}
}
