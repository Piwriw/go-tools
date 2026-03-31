package format

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		defaults []int
		expected int
	}{
		{
			name:     "int_正常值",
			input:    int(123),
			expected: 123,
		},
		{
			name:     "int64_负数",
			input:    int64(-456),
			expected: -456,
		},
		{
			name:     "uint_正常值",
			input:    uint(789),
			expected: 789,
		},
		{
			name:     "float64_整数值",
			input:    12.0,
			expected: 12,
		},
		{
			name:     "json.Number_整数值",
			input:    json.Number("34"),
			expected: 34,
		},
		{
			name:     "string_整数",
			input:    "56",
			expected: 56,
		},
		{
			name:     "string_浮点整数值",
			input:    "78.0",
			expected: 78,
		},
		{
			name:     "bool_true",
			input:    true,
			expected: 1,
		},
		{
			name:     "bool_false",
			input:    false,
			expected: 0,
		},
		{
			name:     "float64_非整数值_未传默认值",
			input:    12.3,
			expected: 0,
		},
		{
			name:     "float64_非整数值_使用默认值",
			input:    12.3,
			defaults: []int{9},
			expected: 9,
		},
		{
			name:     "string_非整数值_使用默认值",
			input:    "12.3",
			defaults: []int{8},
			expected: 8,
		},
		{
			name:     "nil_使用默认值",
			input:    nil,
			defaults: []int{7},
			expected: 7,
		},
		{
			name:     "不支持类型_多个默认值仅取第一个",
			input:    []int{1, 2, 3},
			defaults: []int{5, 6},
			expected: 5,
		},
		{
			name:     "string_超出int范围_使用默认值",
			input:    "999999999999999999999999",
			defaults: []int{4},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt(tt.input, tt.defaults...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
