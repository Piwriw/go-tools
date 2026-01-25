package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// ConsoleHandler 控制台输出处理器
type ConsoleHandler struct {
	EnableColor bool
	EnableJSON  bool
	logger      *log.Logger
}

// NewConsoleHandler 创建新的控制台处理器
func NewConsoleHandler() *ConsoleHandler {
	return &ConsoleHandler{
		EnableColor: true,
		EnableJSON:  false,
		logger:      log.Default(),
	}
}

// Handle 处理审计事件
func (h *ConsoleHandler) Handle(ctx context.Context, event *Event) error {
	if h.EnableJSON {
		return h.handleJSON(ctx, event)
	}
	return h.handleText(ctx, event)
}

// handleText 文本格式输出
func (h *ConsoleHandler) handleText(ctx context.Context, event *Event) error {
	var sb strings.Builder

	// 时间戳
	sb.WriteString(fmt.Sprintf("[%s] ", event.Timestamp))

	// 操作类型（带颜色）
	if h.EnableColor {
		sb.WriteString(h.colorizeOperation(event.Operation))
	} else {
		sb.WriteString(fmt.Sprintf("%-8s", event.Operation))
	}

	// 表名和主键
	sb.WriteString(fmt.Sprintf(" | Table: %-20s", event.Table))
	if event.PrimaryKey != "" {
		sb.WriteString(fmt.Sprintf(" | PK: %s", event.PrimaryKey))
	}

	// 用户信息
	if event.UserID != "" || event.Username != "" {
		sb.WriteString(" | User: ")
		if event.Username != "" {
			sb.WriteString(event.Username)
		}
		if event.UserID != "" {
			sb.WriteString(fmt.Sprintf("(%s)", event.UserID))
		}
	}

	// 请求信息
	if event.RequestID != "" {
		sb.WriteString(fmt.Sprintf(" | ReqID: %s", event.RequestID))
	}

	h.logger.Println(sb.String())

	// 旧值
	if len(event.OldValues) > 0 {
		h.logValues("OLD VALUES", event.OldValues)
	}

	// 新值
	if len(event.NewValues) > 0 {
		h.logValues("NEW VALUES", event.NewValues)
	}

	// SQL（如果存在）
	if event.SQL != "" && !h.EnableJSON {
		h.logger.Printf("  SQL: %s", event.SQL)
		if len(event.SQLArgs) > 0 {
			h.logger.Printf("  Args: %v", event.SQLArgs)
		}
	}

	return nil
}

// handleJSON JSON 格式输出
func (h *ConsoleHandler) handleJSON(ctx context.Context, event *Event) error {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	h.logger.Println(string(data))
	return nil
}

// logValues 记录值
func (h *ConsoleHandler) logValues(label string, values map[string]any) {
	h.logger.Printf("  %s:", label)
	for k, v := range values {
		h.logger.Printf("    %s: %v", k, v)
	}
}

// colorizeOperation 为操作类型添加颜色
func (h *ConsoleHandler) colorizeOperation(op Operation) string {
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
	)

	var color string
	switch op {
	case OperationCreate:
		color = colorGreen
	case OperationUpdate:
		color = colorBlue
	case OperationDelete:
		color = colorRed
	case OperationQuery:
		color = colorYellow
	default:
		color = colorReset
	}

	return fmt.Sprintf("%s%-8s%s", color, op, colorReset)
}

// SetLogger 设置自定义 logger
func (h *ConsoleHandler) SetLogger(logger *log.Logger) {
	h.logger = logger
}

// SetColor 启用或禁用颜色
func (h *ConsoleHandler) SetColor(enable bool) {
	h.EnableColor = enable
}

// SetJSON 启用或禁用 JSON 格式
func (h *ConsoleHandler) SetJSON(enable bool) {
	h.EnableJSON = enable
}

// ==================== Middleware Chain Support ====================

// ChainMiddleware 链式中间件
type ChainMiddleware struct {
	handler EventHandler
}

// NewChainMiddleware 创建链式中间件包装器
func NewChainMiddleware(handler EventHandler) *ChainMiddleware {
	return &ChainMiddleware{handler: handler}
}

// Handle 实现 EventHandler 接口
func (c *ChainMiddleware) Handle(ctx context.Context, event *Event) error {
	return c.handler.Handle(ctx, event)
}

// Then 添加下一个处理器
func (c *ChainMiddleware) Then(next EventHandler) EventHandler {
	return EventHandlerFunc(func(ctx context.Context, event *Event) error {
		if err := c.Handle(ctx, event); err != nil {
			return err
		}
		return next.Handle(ctx, event)
	})
}

// ==================== Filter Middleware ====================

// FilterMiddleware 过滤中间件
type FilterMiddleware struct {
	predicate func(*Event) bool
	handler   EventHandler
}

// NewFilterMiddleware 创建过滤中间件
func NewFilterMiddleware(predicate func(*Event) bool, handler EventHandler) *FilterMiddleware {
	return &FilterMiddleware{
		predicate: predicate,
		handler:   handler,
	}
}

// Handle 实现 EventHandler 接口
func (f *FilterMiddleware) Handle(ctx context.Context, event *Event) error {
	if f.predicate != nil && !f.predicate(event) {
		return nil // 跳过
	}
	return f.handler.Handle(ctx, event)
}

// ==================== Retry Middleware ====================

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	handler    EventHandler
	maxRetries int
	delay      time.Duration
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(handler EventHandler, maxRetries int, delay time.Duration) *RetryMiddleware {
	return &RetryMiddleware{
		handler:    handler,
		maxRetries: maxRetries,
		delay:      delay,
	}
}

// Handle 实现 EventHandler 接口
func (r *RetryMiddleware) Handle(ctx context.Context, event *Event) error {
	var err error
	for i := 0; i <= r.maxRetries; i++ {
		err = r.handler.Handle(ctx, event)
		if err == nil {
			return nil
		}

		if i < r.maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(r.delay):
				// 继续重试
			}
		}
	}
	return fmt.Errorf("after %d retries: %w", r.maxRetries, err)
}
