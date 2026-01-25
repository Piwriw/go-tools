package audit

import "time"

// Operation 定义审计操作类型的枚举
type Operation string

const (
	OperationCreate Operation = "create"
	OperationUpdate Operation = "update"
	OperationDelete Operation = "delete"
	OperationQuery  Operation = "query"
)

// String 实现 Stringer 接口
func (o Operation) String() string {
	return string(o)
}

// IsValid 验证操作类型是否有效
func (o Operation) IsValid() bool {
	switch o {
	case OperationCreate, OperationUpdate, OperationDelete, OperationQuery:
		return true
	}
	return false
}

// AuditEvent 审计事件
type AuditEvent struct {
	Timestamp  time.Time
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
