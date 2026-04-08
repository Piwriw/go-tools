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

func TestMustToInt(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result, err := MustToInt("56")
		assert.NoError(t, err)
		assert.Equal(t, 56, result)
	})

	t.Run("转换失败返回错误", func(t *testing.T) {
		result, err := MustToInt("12.3")
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})
}

func TestToFloat64WithDefault(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result := ToFloat64("12.5")
		assert.Equal(t, 12.5, result)
	})

	t.Run("转换失败使用默认值", func(t *testing.T) {
		result := ToFloat64("abc", 3.14)
		assert.Equal(t, 3.14, result)
	})
}

func TestMustToFloat64(t *testing.T) {
	// keep unique from existing number_test.go symbols
	t.Run("转换成功", func(t *testing.T) {
		result, err := MustToFloat64(json.Number("12.5"))
		assert.NoError(t, err)
		assert.Equal(t, 12.5, result)
	})

	t.Run("转换失败返回错误", func(t *testing.T) {
		result, err := MustToFloat64("abc")
		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
	})
}

func TestToStringWithDefault(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result := ToString(123)
		assert.Equal(t, "123", result)
	})

	t.Run("转换失败使用默认值", func(t *testing.T) {
		result := ToString(func() {}, "fallback")
		assert.Equal(t, "fallback", result)
	})
}

func TestMustToString(t *testing.T) {
	// keep unique from existing string_test.go symbols
	t.Run("转换成功", func(t *testing.T) {
		result, err := MustToString(true)
		assert.NoError(t, err)
		assert.Equal(t, "true", result)
	})

	t.Run("转换失败返回错误", func(t *testing.T) {
		result, err := MustToString(func() {})
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
}

func TestToBoolWithDefault(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result := ToBool("true")
		assert.Equal(t, true, result)
	})

	t.Run("转换失败使用默认值", func(t *testing.T) {
		result := ToBool("abc", true)
		assert.Equal(t, true, result)
	})
}

func TestMustToBool(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result, err := MustToBool(1)
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("转换失败返回错误", func(t *testing.T) {
		result, err := MustToBool("abc")
		assert.Error(t, err)
		assert.Equal(t, false, result)
	})
}

// ============================================================================
// ToMapStringAny
// ============================================================================

func TestMustToMapStringAny(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name    string
		input   any
		want    map[string]any
		wantErr bool
	}{
		{
			name:  "map[string]any 直接返回",
			input: map[string]any{"key": "value"},
			want:  map[string]any{"key": "value"},
		},
		{
			name:  "map[string]string",
			input: map[string]string{"name": "Tom"},
			want:  map[string]any{"name": "Tom"},
		},
		{
			name:  "map[string]int",
			input: map[string]int{"count": 10},
			want:  map[string]any{"count": 10},
		},
		{
			name:  "struct 通过 JSON 转换",
			input: user{Name: "Alice", Age: 30},
			want:  map[string]any{"name": "Alice", "age": float64(30)},
		},
		{
			name:  "指向 struct 的指针",
			input: &user{Name: "Bob", Age: 25},
			want:  map[string]any{"name": "Bob", "age": float64(25)},
		},
		{
			name:  "JSON 字符串",
			input: `{"hello":"world"}`,
			want:  map[string]any{"hello": "world"},
		},
		{
			name:  "JSON []byte",
			input: []byte(`{"a":"b"}`),
			want:  map[string]any{"a": "b"},
		},
		{
			name:    "nil",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "nil 指针",
			input:   (*user)(nil),
			wantErr: true,
		},
		{
			name:    "无效 JSON 字符串",
			input:   "not json",
			wantErr: true,
		},
		{
			name:    "不支持的类型",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MustToMapStringAny(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToMapStringAny(t *testing.T) {
	t.Run("转换成功", func(t *testing.T) {
		result := ToMapStringAny(map[string]any{"key": "val"})
		assert.Equal(t, map[string]any{"key": "val"}, result)
	})

	t.Run("转换失败返回 nil", func(t *testing.T) {
		result := ToMapStringAny(123)
		assert.Nil(t, result)
	})

	t.Run("转换失败使用默认值", func(t *testing.T) {
		fallback := map[string]any{"default": true}
		result := ToMapStringAny(123, fallback)
		assert.Equal(t, fallback, result)
	})
}

