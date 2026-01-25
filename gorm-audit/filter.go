package audit

import "strings"

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
