package audit

import (
	"context"
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
