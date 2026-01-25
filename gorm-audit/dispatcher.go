package audit

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
)

// EventHandler 事件处理器接口（audit 包内部使用）
type EventHandler interface {
	Handle(ctx context.Context, event *AuditEvent) error
}

// Dispatcher 事件分发器
type Dispatcher struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

// NewDispatcher 创建新的分发器
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make([]EventHandler, 0),
	}
}

// Add 添加事件处理器
func (d *Dispatcher) Add(handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, handler)
}

// Dispatch 分发事件到所有处理器
func (d *Dispatcher) Dispatch(ctx context.Context, event *AuditEvent) {
	d.mu.RLock()
	handlers := make([]EventHandler, len(d.handlers))
	copy(handlers, d.handlers)
	d.mu.RUnlock()

	for _, handler := range handlers {
		if handler != nil {
			go d.safeHandle(ctx, handler, event)
		}
	}
}

// safeHandle 安全执行事件处理器，带 panic 恢复
func (d *Dispatcher) safeHandle(ctx context.Context, handler EventHandler, event *AuditEvent) {
	defer func() {
		if r := recover(); r != nil {
			d.handlePanic(r, handler, event)
		}
	}()

	_ = handler.Handle(ctx, event)
}

// handlePanic 处理 panic 情况
func (d *Dispatcher) handlePanic(r any, handler EventHandler, event *AuditEvent) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	msg := fmt.Sprintf(
		"[AUDIT] panic recovered: %v, handler: %T, table: %s, operation: %s\nstack: %s",
		r, handler, event.Table, event.Operation, stack,
	)

	log.Println(msg)
}
