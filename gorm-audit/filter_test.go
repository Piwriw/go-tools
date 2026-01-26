package audit

import (
	"testing"

	"github.com/piwriw/gorm/gorm-audit/types"
)

func TestTableFilterWhitelist(t *testing.T) {
	filter := NewTableFilter(FilterModeWhitelist, []string{"users", "orders"})

	tests := []struct {
		name     string
		table    string
		expected bool
	}{
		{"users table", "users", true},
		{"orders table", "orders", true},
		{"products table", "products", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{Table: tt.table}
			result := filter.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTableFilterBlacklist(t *testing.T) {
	filter := NewTableFilter(FilterModeBlacklist, []string{"logs", "temp_*"})

	tests := []struct {
		name     string
		table    string
		expected bool
	}{
		{"users table", "users", true},
		{"logs table", "logs", false},
		{"temp_users table", "temp_users", false},
		{"temp_data table", "temp_data", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{Table: tt.table}
			result := filter.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestOperationFilter(t *testing.T) {
	filter := NewOperationFilter([]types.Operation{
		types.OperationCreate,
		types.OperationUpdate,
	})

	tests := []struct {
		name      string
		operation Operation
		expected  bool
	}{
		{"create operation", OperationCreate, true},
		{"update operation", OperationUpdate, true},
		{"delete operation", OperationDelete, false},
		{"query operation", OperationQuery, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{Operation: tt.operation}
			result := filter.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUserFilterWhitelist(t *testing.T) {
	filter := NewUserFilter(FilterModeWhitelist, []string{"user1", "user2"})

	tests := []struct {
		name     string
		userID   string
		expected bool
	}{
		{"user1", "user1", true},
		{"user2", "user2", true},
		{"user3", "user3", false},
		{"empty user", "", true}, // 默认审计
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{UserID: tt.userID}
			result := filter.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldFilter(t *testing.T) {
	filter := NewFieldFilter([]string{"email", "password"})

	tests := []struct {
		name      string
		oldValues map[string]any
		newValues map[string]any
		expected  bool
	}{
		{"email changed", map[string]any{"email": "old@example.com"}, map[string]any{"email": "new@example.com"}, true},
		{"password changed", map[string]any{"password": "old"}, map[string]any{"password": "new"}, true},
		{"name changed", map[string]any{"name": "old"}, map[string]any{"name": "new"}, false},
		{"no fields", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{
				OldValues: tt.oldValues,
				NewValues: tt.newValues,
			}
			result := filter.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompositeFilterAnd(t *testing.T) {
	tableFilter := NewTableFilter(FilterModeWhitelist, []string{"users"})
	opFilter := NewOperationFilter([]types.Operation{types.OperationCreate})
	composite := NewCompositeFilter(FilterLogicAnd, tableFilter, opFilter)

	tests := []struct {
		name      string
		table     string
		operation Operation
		expected  bool
	}{
		{"users + create", "users", OperationCreate, true},
		{"users + update", "users", OperationUpdate, false},
		{"orders + create", "orders", OperationCreate, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &AuditEvent{
				Table:     tt.table,
				Operation: tt.operation,
			}
			result := composite.ShouldAudit(event)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
