package handler

import (
	"context"
	"time"
)

// Operation 操作类型
type Operation string

const (
	OperationCreate Operation = "CREATE"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"
	OperationQuery  Operation = "QUERY"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event *Event) error
}

// EventHandlerFunc 函数式事件处理器
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

// Event 审计事件
type Event struct {
	Timestamp  string
	Operation  Operation
	Table      string
	PrimaryKey string
	OldValues  map[string]any
	NewValues  map[string]any
	SQL        string
	SQLArgs    []any
	UserID     string
	Username   string
	IP         string
	UserAgent  string
	RequestID  string
}

// GetTimestamp 获取格式化的时间戳
func (e *Event) GetTimestamp() time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000", e.Timestamp)
	if err != nil {
		return time.Time{}
	}
	return t
}
