package audit

import (
	"testing"
)

func TestAuditLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    AuditLevel
		expected string
	}{
		{"All", AuditLevelAll, "all"},
		{"ChangesOnly", AuditLevelChangesOnly, "changes_only"},
		{"None", AuditLevelNone, "none"},
		{"Unknown", AuditLevel(100), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultContextKeys(t *testing.T) {
	keys := DefaultContextKeys()

	if keys.UserID == nil {
		t.Error("UserID should not be nil")
	}
	if keys.Username == nil {
		t.Error("Username should not be nil")
	}
}

func TestDefaultWorkerPoolConfig(t *testing.T) {
	config := DefaultWorkerPoolConfig()

	if config.WorkerCount != 10 {
		t.Errorf("expected WorkerCount 10, got %d", config.WorkerCount)
	}
	if config.QueueSize != 1000 {
		t.Errorf("expected QueueSize 1000, got %d", config.QueueSize)
	}
	if config.Timeout != 5000 {
		t.Errorf("expected Timeout 5000, got %d", config.Timeout)
	}
}
