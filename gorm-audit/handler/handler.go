package handler

import "context"

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event *Event) error
}

// EventHandlerFunc 函数式事件处理器
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

// Event 事件定义（避免循环导入，重新定义）
type Event struct {
	Timestamp  string
	Operation  string
	Table      string
	PrimaryKey string
	Data       map[string]any
}
