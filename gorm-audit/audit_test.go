package audit

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewAudit(t *testing.T) {
	// 测试 nil 配置
	p := New(nil)
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}
	if p.config.Level != AuditLevelChangesOnly {
		t.Errorf("expected default level ChangesOnly, got %v", p.config.Level)
	}

	// 测试自定义配置
	customConfig := &Config{
		Level:        AuditLevelAll,
		IncludeQuery: true,
	}
	p = New(customConfig)
	if p.config.Level != AuditLevelAll {
		t.Errorf("expected level All, got %v", p.config.Level)
	}
	if !p.config.IncludeQuery {
		t.Error("expected IncludeQuery to be true")
	}
}

func TestAuditName(t *testing.T) {
	p := New(nil)
	if p.Name() != "audit" {
		t.Errorf("expected name 'audit', got %v", p.Name())
	}
}

func TestAuditInitialize(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	p := New(nil)
	err = p.Initialize(db)
	if err != nil {
		t.Errorf("failed to initialize plugin: %v", err)
	}
}

func TestAuditInitializeWithQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	p := New(&Config{
		IncludeQuery: true,
	})
	err = p.Initialize(db)
	if err != nil {
		t.Errorf("failed to initialize plugin with query: %v", err)
	}
}

func TestAuditReload(t *testing.T) {
	// 设置环境变量
	t.Setenv("GORM_AUDIT_LEVEL", "all")

	audit := New(&Config{
		Level: AuditLevelNone,
	})

	// 初始级别应该是 None
	if audit.GetLevel() != AuditLevelNone {
		t.Errorf("expected initial level None, got %v", audit.GetLevel())
	}

	// 重新加载
	if err := audit.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// 级别应该变为 All
	if audit.GetLevel() != AuditLevelAll {
		t.Errorf("expected level All after reload, got %v", audit.GetLevel())
	}
}

func TestAuditReloadInvalidEnv(t *testing.T) {
	// 设置无效的环境变量
	t.Setenv("GORM_AUDIT_LEVEL", "invalid_value")

	audit := New(&Config{
		Level: AuditLevelAll,
	})

	// 重新加载应该回退到默认值
	if err := audit.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// 应该回退到默认值 changes_only
	if audit.GetLevel() != AuditLevelChangesOnly {
		t.Errorf("expected default level, got %v", audit.GetLevel())
	}
}

func TestAuditReloadEmptyEnv(t *testing.T) {
	// 清空环境变量
	t.Setenv("GORM_AUDIT_LEVEL", "")

	audit := New(&Config{
		Level: AuditLevelAll,
	})

	// 重新加载应该回退到默认值
	if err := audit.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// 应该回退到默认值 changes_only
	if audit.GetLevel() != AuditLevelChangesOnly {
		t.Errorf("expected default level, got %v", audit.GetLevel())
	}
}
