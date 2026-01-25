package audit

import (
	"strings"

	"github.com/piwriw/gorm/gorm-audit/types"
)

// FilterMode 过滤模式
type FilterMode int

const (
	FilterModeWhitelist FilterMode = iota // 白名单：只有匹配的才审计
	FilterModeBlacklist                  // 黑名单：匹配的不审计
)

// FilterLogic 组合过滤器逻辑
type FilterLogic int

const (
	FilterLogicAnd FilterLogic = iota // 所有过滤器都返回 true
	FilterLogicOr                    // 任一过滤器返回 true
)

// Filter 事件过滤器接口
type Filter interface {
	ShouldAudit(event *AuditEvent) bool
}

// String 返回过滤模式的字符串表示
func (m FilterMode) String() string {
	switch m {
	case FilterModeWhitelist:
		return "whitelist"
	case FilterModeBlacklist:
		return "blacklist"
	default:
		return "unknown"
	}
}

// TableFilter 表名过滤器
type TableFilter struct {
	mode   FilterMode
	tables map[string]bool // 用于快速查找
}

// NewTableFilter 创建表名过滤器
func NewTableFilter(mode FilterMode, tables []string) *TableFilter {
	tableMap := make(map[string]bool, len(tables))
	for _, table := range tables {
		// 支持通配符 * 后缀
		if strings.HasSuffix(table, "*") {
			prefix := strings.TrimSuffix(table, "*")
			tableMap[prefix+"*"] = true
		} else {
			tableMap[table] = true
		}
	}
	return &TableFilter{
		mode:   mode,
		tables: tableMap,
	}
}

// ShouldAudit 实现过滤逻辑
func (f *TableFilter) ShouldAudit(event *AuditEvent) bool {
	matched := f.matchTable(event.Table)

	switch f.mode {
	case FilterModeWhitelist:
		return matched // 白名单：匹配的才审计
	case FilterModeBlacklist:
		return !matched // 黑名单：匹配的不审计
	default:
		return true
	}
}

// matchTable 检查表名是否匹配（支持通配符）
func (f *TableFilter) matchTable(table string) bool {
	// 精确匹配
	if f.tables[table] {
		return true
	}

	// 通配符匹配
	for pattern := range f.tables {
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(table, prefix) {
				return true
			}
		}
	}

	return false
}

// OperationFilter 操作类型过滤器
type OperationFilter struct {
	operations map[types.Operation]bool
}

// NewOperationFilter 创建操作类型过滤器（白名单模式）
func NewOperationFilter(operations []types.Operation) *OperationFilter {
	opMap := make(map[types.Operation]bool, len(operations))
	for _, op := range operations {
		opMap[op] = true
	}
	return &OperationFilter{
		operations: opMap,
	}
}

// ShouldAudit 实现过滤逻辑
func (f *OperationFilter) ShouldAudit(event *AuditEvent) bool {
	// 如果没有配置操作列表，默认审计所有
	if len(f.operations) == 0 {
		return true
	}

	// 只审计配置的操作类型
	return f.operations[event.Operation]
}

// UserFilter 用户 ID 过滤器
type UserFilter struct {
	mode    FilterMode
	userIDs map[string]bool
}

// NewUserFilter 创建用户 ID 过滤器
func NewUserFilter(mode FilterMode, userIDs []string) *UserFilter {
	idMap := make(map[string]bool, len(userIDs))
	for _, id := range userIDs {
		idMap[id] = true
	}
	return &UserFilter{
		mode:    mode,
		userIDs: idMap,
	}
}

// ShouldAudit 实现过滤逻辑
func (f *UserFilter) ShouldAudit(event *AuditEvent) bool {
	// 如果没有用户信息，默认审计
	if event.UserID == "" {
		return true
	}

	matched := f.userIDs[event.UserID]

	switch f.mode {
	case FilterModeWhitelist:
		return matched // 白名单：只有匹配的用户才审计
	case FilterModeBlacklist:
		return !matched // 黑名单：匹配的用户不审计
	default:
		return true
	}
}
