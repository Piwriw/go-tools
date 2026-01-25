package audit

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// MockEventHandler 用于测试的模拟处理器
type MockEventHandler struct {
	CallCount int32
	Delay     time.Duration
	Panic     bool
	Error     error
}

func (m *MockEventHandler) Handle(ctx context.Context, event *AuditEvent) error {
	atomic.AddInt32(&m.CallCount, 1)
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	if m.Panic {
		panic("test panic")
	}
	return m.Error
}

func TestDispatcherAdd(t *testing.T) {
	d := NewDispatcher()
	handler := &MockEventHandler{}

	d.Add(handler)

	if len(d.handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(d.handlers))
	}
}

func TestDispatcherDispatch(t *testing.T) {
	d := NewDispatcher()
	handler := &MockEventHandler{}
	d.Add(handler)

	event := &AuditEvent{
		Timestamp:  time.Now(),
		Operation:  OperationCreate,
		Table:      "users",
	}

	d.Dispatch(context.Background(), event)

	// 等待异步处理
	time.Sleep(100 * time.Millisecond)

	// 注意：CallCount 检查会在 callback.go 实现后更新
}

func TestDispatcherDispatchMultiple(t *testing.T) {
	d := NewDispatcher()
	handler1 := &MockEventHandler{}
	handler2 := &MockEventHandler{}
	d.Add(handler1)
	d.Add(handler2)

	event := &AuditEvent{
		Timestamp:  time.Now(),
		Operation:  OperationUpdate,
		Table:      "users",
	}

	d.Dispatch(context.Background(), event)

	time.Sleep(100 * time.Millisecond)
}

func TestDispatcherPanicRecovery(t *testing.T) {
	d := NewDispatcher()
	handler := &MockEventHandler{Panic: true}
	d.Add(handler)

	event := &AuditEvent{
		Timestamp:  time.Now(),
		Operation:  OperationDelete,
		Table:      "users",
	}

	// 不应该 panic
	d.Dispatch(context.Background(), event)

	time.Sleep(100 * time.Millisecond)
}
