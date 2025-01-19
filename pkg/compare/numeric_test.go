package compare

import (
	"testing"
)

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantBool bool
		wantVal  any
	}{
		// 测试 int 类型
		{
			name:     "Valid integer string for int",
			input:    "123",
			wantBool: true,
			wantVal:  123,
		},
		{
			name:     "Invalid integer string for int",
			input:    "abc",
			wantBool: false,
			wantVal:  0,
		},
		{
			name:     "Float string for int (should fail)",
			input:    "123.45",
			wantBool: false,
			wantVal:  0,
		},

		// 测试 float64 类型
		{
			name:     "Valid float string for float64",
			input:    "123.45",
			wantBool: true,
			wantVal:  123.45,
		},
		{
			name:     "Valid integer string for float64",
			input:    "123",
			wantBool: true,
			wantVal:  float64(123),
		},
		{
			name:     "Invalid float string for float64",
			input:    "abc",
			wantBool: false,
			wantVal:  float64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试 int 类型
			if tt.wantVal == 0 || tt.wantVal == 123 {
				ok, val := IsNumeric[int](tt.input)
				if ok != tt.wantBool || val != tt.wantVal {
					t.Errorf("IsNumeric[int](%q) = (%v, %v), want (%v, %v)",
						tt.input, ok, val, tt.wantBool, tt.wantVal)
				}
			}

			// 测试 float64 类型
			if tt.wantVal == float64(0) || tt.wantVal == 123.45 || tt.wantVal == float64(123) {
				ok, val := IsNumeric[float64](tt.input)
				if ok != tt.wantBool || val != tt.wantVal {
					t.Errorf("IsNumeric[float64](%q) = (%v, %v), want (%v, %v)",
						tt.input, ok, val, tt.wantBool, tt.wantVal)
				}
			}
		})
	}
}
