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
