package audit

import "time"

// AuditLevel 审计级别
type AuditLevel int

const (
	AuditLevelAll         AuditLevel = iota // 记录所有操作
	AuditLevelChangesOnly                   // 仅记录变更操作
	AuditLevelNone                          // 不记录
)

// String 实现 Stringer 接口
func (a AuditLevel) String() string {
	switch a {
	case AuditLevelAll:
		return "all"
	case AuditLevelChangesOnly:
		return "changes_only"
	case AuditLevelNone:
		return "none"
	default:
		return "unknown"
	}
}

// ContextKeyConfig Context 键配置
type ContextKeyConfig struct {
	UserID    any
	Username  any
	IP        any
	UserAgent any
	RequestID any
}

// DefaultContextKeys 返回默认的 context key 配置
func DefaultContextKeys() ContextKeyConfig {
	type contextKey string
	return ContextKeyConfig{
		UserID:    contextKey("user_id"),
		Username:  contextKey("username"),
		IP:        contextKey("ip"),
		UserAgent: contextKey("user_agent"),
		RequestID: contextKey("request_id"),
	}
}

// WorkerPoolConfig worker pool 配置
type WorkerPoolConfig struct {
	WorkerCount int // worker 数量
	QueueSize   int // 队列大小
	Timeout     int // 处理超时（毫秒）

	// 批量处理配置
	EnableBatch   bool          // 是否启用批量处理
	BatchSize     int           // 批量大小，默认 1000
	FlushInterval time.Duration // 刷新间隔，默认 10秒
}

// DefaultWorkerPoolConfig 返回默认的 worker pool 配置
func DefaultWorkerPoolConfig() *WorkerPoolConfig {
	return &WorkerPoolConfig{
		WorkerCount:   10,
		QueueSize:     1000,
		Timeout:       5000, // 5秒
		EnableBatch:   false,
		BatchSize:     1000,
		FlushInterval: 10 * time.Second,
	}
}
