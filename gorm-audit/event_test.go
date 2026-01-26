package audit

import (
	"testing"
	"time"
)

func TestOperationString(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected string
	}{
		{"Create", OperationCreate, "create"},
		{"Update", OperationUpdate, "update"},
		{"Delete", OperationDelete, "delete"},
		{"Query", OperationQuery, "query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOperationIsValid(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected bool
	}{
		{"Valid Create", OperationCreate, true},
		{"Valid Update", OperationUpdate, true},
		{"Valid Delete", OperationDelete, true},
		{"Valid Query", OperationQuery, true},
		{"Invalid", Operation("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAuditEventCreation(t *testing.T) {
	event := &AuditEvent{
		Timestamp:  time.Now(),
		Operation:  OperationCreate,
		Table:      "users",
		PrimaryKey: "1",
		OldValues:  make(map[string]any),
		NewValues:  map[string]any{"name": "test"},
	}

	if event.Operation != OperationCreate {
		t.Errorf("expected OperationCreate, got %v", event.Operation)
	}
	if event.Table != "users" {
		t.Errorf("expected table 'users', got %v", event.Table)
	}
}
