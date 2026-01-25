package audit

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
