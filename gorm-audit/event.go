package audit

import (
	"time"

	"github.com/piwriw/gorm/gorm-audit/types"
)

// Operation 导出共享的操作类型
type Operation = types.Operation

// 导出常量方便使用
const (
	OperationCreate = types.OperationCreate
	OperationUpdate = types.OperationUpdate
	OperationDelete = types.OperationDelete
	OperationQuery  = types.OperationQuery
)

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
