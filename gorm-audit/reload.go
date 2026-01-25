package audit

import (
	"os"
	"strings"
)

const (
	envAuditLevel = "GORM_AUDIT_LEVEL"
)

// getAuditLevelFromEnv 从环境变量读取审计级别
func getAuditLevelFromEnv() AuditLevel {
	value := os.Getenv(envAuditLevel)
	if value == "" {
		return AuditLevelChangesOnly // 默认值
	}

	switch strings.ToLower(value) {
	case "all":
		return AuditLevelAll
	case "changes_only":
		return AuditLevelChangesOnly
	case "none":
		return AuditLevelNone
	default:
		return AuditLevelChangesOnly // 无效值返回默认
	}
}
