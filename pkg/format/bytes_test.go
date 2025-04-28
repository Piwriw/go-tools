package format

import "testing"

func TestKBToGB(t *testing.T) {
	tests := []struct {
		name     string
		kb       float64
		expected float64
	}{
		{"Zero KB", 0, 0},
		{"Exactly 1MB in KB", 1024, 0.00},
		{"Exactly 1GB in KB", 1024 * 1024, 1},
		{"1.5GB in KB", 1.5 * 1024 * 1024, 1.5},
		{"Rounding down", 1.234 * 1024 * 1024, 1.23},
		{"Rounding up", 1.235 * 1024 * 1024, 1.24},
		{"Large value", 5 * 1024 * 1024 * 1024, 5 * 1024},
		{"Negative value", -1024 * 1024, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KBToGB(tt.kb)
			if got != tt.expected {
				t.Errorf("KBToGB(%v) = %v, want %v", tt.kb, got, tt.expected)
			}
		})
	}
}

func TestKBToMB(t *testing.T) {
	tests := []struct {
		name     string
		kb       float64
		expected float64
	}{
		{"Zero KB", 0, 0},
		{"Exactly 1024KB", 1024, 1},
		{"1.5MB in KB", 1536, 1.5},
		{"Rounding down", 1234, 1.21},
		{"Rounding up", 1235, 1.21},
		{"Large value", 5 * 1024 * 1024, 5 * 1024},
		{"Negative value", -1024, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KBToMB(tt.kb)
			if got != tt.expected {
				t.Errorf("KBToMB(%v) = %v, want %v", tt.kb, got, tt.expected)
			}
		})
	}
}

func TestBytesToGB(t *testing.T) {
	tests := []struct {
		name     string
		bytes    float64
		expected float64
	}{
		{"Zero bytes", 0, 0},
		{"Exactly 1GB", 1024 * 1024 * 1024, 1},
		{"1.5GB", 1.5 * 1024 * 1024 * 1024, 1.5},
		{"Precise rounding down", 1.234 * 1024 * 1024 * 1024, 1.23},
		{"Precise rounding up", 1.235 * 1024 * 1024 * 1024, 1.24},
		{"Large value", 5 * 1024 * 1024 * 1024 * 1024, 5 * 1024},
		{"Negative value", -1024 * 1024 * 1024, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToGB(tt.bytes)
			if got != tt.expected {
				t.Errorf("BytesToGB(%v) = %v, want %v", tt.bytes, got, tt.expected)
			}
		})
	}
}

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		name     string
		bytes    float64
		expected float64
	}{
		{"Zero bytes", 0, 0},
		{"Exactly 1MB", 1024 * 1024, 1},
		{"1.5MB", 1.5 * 1024 * 1024, 1.5},
		{"Precise rounding down", 1.234 * 1024 * 1024, 1.23},
		{"Precise rounding up", 1.235 * 1024 * 1024, 1.24},
		{"Large value", 5 * 1024 * 1024 * 1024, 5 * 1024},
		{"Negative value", -1024 * 1024, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToMB(tt.bytes)
			if got != tt.expected {
				t.Errorf("BytesToMB(%v) = %v, want %v", tt.bytes, got, tt.expected)
			}
		})
	}
}
