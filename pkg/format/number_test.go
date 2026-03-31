package format

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFloat64ToPercentFloat 测试将float64转换为百分比
// 测试场景：正常场景、边界值、负数、零值、精度测试
func TestFloat64ToPercentFloat(t *testing.T) {
	tests := []struct {
		name     string  // 测试用例名称
		input    float64 // 输入参数
		expected float64 // 预期返回结果
	}{
		{
			name:     "正常场景_50%",
			input:    0.5,
			expected: 50.0,
		},
		{
			name:     "正常场景_100%",
			input:    1.0,
			expected: 100.0,
		},
		{
			name:     "正常场景_小数百分比",
			input:    0.1234,
			expected: 12.34,
		},
		{
			name:     "正常场景_大于1的值",
			input:    1.5,
			expected: 150.0,
		},
		{
			name:     "边界条件_零值",
			input:    0,
			expected: 0,
		},
		{
			name:     "边界条件_负数",
			input:    -0.5,
			expected: -50.0,
		},
		{
			name:     "边界条件_非常小的值",
			input:    0.0001,
			expected: 0.01,
		},
		{
			name:     "精度测试_四舍五入",
			input:    0.333333,
			expected: 33.33,
		},
		{
			name:     "精度测试_进位",
			input:    0.335,
			expected: 33.5,
		},
		{
			name:     "精度测试_舍去",
			input:    0.334,
			expected: 33.4,
		},
		{
			name:     "边缘场景_大数值",
			input:    999.999,
			expected: 99999.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Float64ToPercentFloat(tt.input)
			assert.Equal(t, tt.expected, result, "转换后的百分比值应与预期一致")
		})
	}
}

// TestToFloat64 测试将任意类型转换为float64
// 测试场景：各种数字类型、字符串转换、边界值、错误情况
func TestToFloat64(t *testing.T) {
	tests := []struct {
		name        string      // 测试用例名称
		input       interface{} // 输入参数
		expected    float64     // 预期返回结果
		wantErr     bool        // 是否预期发生错误
		errContains string      // 预期错误信息包含的字符串
	}{
		// float64类型
		{
			name:     "float64_正常值",
			input:    float64(123.456),
			expected: 123.456,
			wantErr:  false,
		},
		{
			name:     "float64_负数",
			input:    float64(-45.67),
			expected: -45.67,
			wantErr:  false,
		},
		{
			name:     "float64_零值",
			input:    float64(0),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "float64_最大值",
			input:    math.MaxFloat64,
			expected: math.MaxFloat64,
			wantErr:  false,
		},
		// float32类型
		{
			name:     "float32_正常值",
			input:    float32(78.9),
			expected: 78.9000015258789,
			wantErr:  false,
		},
		// int类型
		{
			name:     "int_正常值",
			input:    int(123),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "int_负数",
			input:    int(-456),
			expected: -456,
			wantErr:  false,
		},
		{
			name:     "int_零值",
			input:    int(0),
			expected: 0,
			wantErr:  false,
		},
		// int8类型
		{
			name:     "int8_最大值127",
			input:    int8(127),
			expected: 127,
			wantErr:  false,
		},
		{
			name:     "int8_最小值-128",
			input:    int8(-128),
			expected: -128,
			wantErr:  false,
		},
		// int16类型
		{
			name:     "int16_最大值32767",
			input:    int16(32767),
			expected: 32767,
			wantErr:  false,
		},
		// int32类型
		{
			name:     "int32_最大值2147483647",
			input:    int32(2147483647),
			expected: 2147483647,
			wantErr:  false,
		},
		// int64类型
		{
			name:     "int64_正常值",
			input:    int64(922337203685477580),
			expected: 922337203685477580,
			wantErr:  false,
		},
		// uint类型
		{
			name:     "uint_正常值",
			input:    uint(123),
			expected: 123,
			wantErr:  false,
		},
		// uint8类型
		{
			name:     "uint8_最大值255",
			input:    uint8(255),
			expected: 255,
			wantErr:  false,
		},
		// uint16类型
		{
			name:     "uint16_最大值65535",
			input:    uint16(65535),
			expected: 65535,
			wantErr:  false,
		},
		// uint32类型
		{
			name:     "uint32_最大值4294967295",
			input:    uint32(4294967295),
			expected: 4294967295,
			wantErr:  false,
		},
		// uint64类型
		{
			name:     "uint64_正常值",
			input:    uint64(9007199254740991),
			expected: 9007199254740991,
			wantErr:  false,
		},
		{
			name:        "uint64_超出精度范围",
			input:       uint64(9007199254740992),
			expected:    0,
			wantErr:     true,
			errContains: "value too large",
		},
		// json.Number类型
		{
			name:     "json.Number_整数",
			input:    json.Number("123"),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "json.Number_浮点数",
			input:    json.Number("123.456"),
			expected: 123.456,
			wantErr:  false,
		},
		{
			name:     "json.Number_负数",
			input:    json.Number("-78.9"),
			expected: -78.9,
			wantErr:  false,
		},
		{
			name:        "json.Number_无效字符串",
			input:       json.Number("abc"),
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert json.Number",
		},
		// string类型
		{
			name:     "string_整数",
			input:    "123",
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "string_浮点数",
			input:    "123.456",
			expected: 123.456,
			wantErr:  false,
		},
		{
			name:     "string_负数",
			input:    "-45.67",
			expected: -45.67,
			wantErr:  false,
		},
		{
			name:     "string_科学计数法",
			input:    "1.23e2",
			expected: 123,
			wantErr:  false,
		},
		{
			name:        "string_无效格式",
			input:       "abc",
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert string",
		},
		{
			name:        "string_空字符串",
			input:       "",
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert string",
		},
		// bool类型
		{
			name:     "bool_true",
			input:    true,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "bool_false",
			input:    false,
			expected: 0,
			wantErr:  false,
		},
		// nil类型
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert nil",
		},
		// 不支持的类型
		{
			name:        "不支持_切片",
			input:       []int{1, 2, 3},
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert []int to float64",
		},
		{
			name:        "不支持_map",
			input:       map[string]int{"key": 1},
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert map[string]int to float64",
		},
		{
			name:        "不支持_结构体",
			input:       struct{ Name string }{"test"},
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert struct { Name string } to float64",
		},
		{
			name:        "不支持_函数",
			input:       func() {},
			expected:    0,
			wantErr:     true,
			errContains: "cannot convert func() to float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MustToFloat64(tt.input)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "错误信息应包含指定字符串")
				}
			} else {
				require.NoError(t, err, "预期不发生错误")
				assert.Equal(t, tt.expected, result, "转换后的值应与预期一致")
			}
		})
	}
}

// TestToFloat64_Parallel 并发安全测试_ToFloat64
func TestToFloat64_Parallel(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{"int", int(123), 123},
		{"float64", float64(45.67), 45.67},
		{"string", "123.45", 123.45},
		{"bool", true, 1},
	}

	for _, tt := range tests {
		tt := tt // 创建局部变量副本
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 并发调用100次，确保没有数据竞争
			for i := 0; i < 100; i++ {
				result, err := MustToFloat64(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestToFloat64_EdgeCases 边缘场景测试
func TestToFloat64_EdgeCases(t *testing.T) {
	t.Run("uint64精度边界", func(t *testing.T) {
		// 测试float64能精确表示的最大整数边界
		maxExactFloat64 := uint64(1<<53 - 1) // 9007199254740991
		result, err := MustToFloat64(maxExactFloat64)
		require.NoError(t, err)
		assert.Equal(t, float64(maxExactFloat64), result)

		// 超出边界
		tooLarge := uint64(1 << 53) // 9007199254740992
		result, err = MustToFloat64(tooLarge)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "value too large")
	})

	t.Run("浮点数精度", func(t *testing.T) {
		// 测试浮点数精度保留
		val := float64(123.45678901234567)
		result, err := MustToFloat64(val)
		require.NoError(t, err)
		assert.InDelta(t, val, result, 1e-9)
	})

	t.Run("字符串科学计数法", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected float64
		}{
			{"1e10", 1e10},
			{"1.5e-5", 1.5e-5},
			{"-2.5E3", -2.5e3},
		}

		for _, tc := range testCases {
			result, err := MustToFloat64(tc.input)
			require.NoError(t, err)
			assert.InDelta(t, tc.expected, result, 1e-9)
		}
	})
}

// BenchmarkFloat64ToPercentFloat 性能基准测试_Float64ToPercentFloat
func BenchmarkFloat64ToPercentFloat(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Float64ToPercentFloat(0.1234)
	}
}

// BenchmarkToFloat64 性能基准测试_ToFloat64
func BenchmarkToFloat64(b *testing.B) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "Int",
			input: int(123),
		},
		{
			name:  "Float64",
			input: float64(123.456),
		},
		{
			name:  "String",
			input: "123.456",
		},
		{
			name:  "JSONNumber",
			input: json.Number("123.456"),
		},
		{
			name:  "Bool",
			input: true,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = MustToFloat64(tt.input)
			}
		})
	}
}

// BenchmarkToFloat64_Parallel 并发性能基准测试
func BenchmarkToFloat64_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = MustToFloat64(float64(123.456))
		}
	})
}

// ExampleFloat64ToPercentFloat 示例代码_Float64ToPercentFloat
func ExampleFloat64ToPercentFloat() {
	result := Float64ToPercentFloat(0.5)
	fmt.Println(result)

	result = Float64ToPercentFloat(0.1234)
	fmt.Println(result)

	// Output:
	// 50
	// 12.34
}

// ExampleToFloat64 示例代码_ToFloat64
func ExampleToFloat64() {
	// 各种类型转换为float64
	val1 := ToFloat64(int(123))
	fmt.Printf("%.2f\n", val1)

	val2 := ToFloat64(float32(45.67))
	fmt.Printf("%.2f\n", val2)

	val3 := ToFloat64("123.456")
	fmt.Printf("%.3f\n", val3)

	val4 := ToFloat64(true)
	fmt.Printf("%.0f\n", val4)

	// Output:
	// 123.00
	// 45.67
	// 123.456
	// 1
}

// TestIsNumber 测试判断字符串是否为数字
// 测试场景：整数、浮点数、科学计数法、边界值、非法输入
func TestIsNumber(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		input    string // 输入参数
		expected bool   // 预期返回结果
	}{
		// 整数
		{
			name:     "正整数",
			input:    "123",
			expected: true,
		},
		{
			name:     "负整数",
			input:    "-456",
			expected: true,
		},
		{
			name:     "零",
			input:    "0",
			expected: true,
		},
		{
			name:     "带正号的整数",
			input:    "+789",
			expected: true,
		},
		// 浮点数
		{
			name:     "正浮点数",
			input:    "123.456",
			expected: true,
		},
		{
			name:     "负浮点数",
			input:    "-78.9",
			expected: true,
		},
		{
			name:     "小数点前无数字",
			input:    ".5",
			expected: true,
		},
		{
			name:     "小数点后无数字",
			input:    "5.",
			expected: true,
		},
		{
			name:     "带正号的浮点数",
			input:    "+12.34",
			expected: true,
		},
		// 科学计数法
		{
			name:     "科学计数法_小写e",
			input:    "1e10",
			expected: true,
		},
		{
			name:     "科学计数法_大写E",
			input:    "1.5E3",
			expected: true,
		},
		{
			name:     "科学计数法_负指数",
			input:    "1.5e-5",
			expected: true,
		},
		{
			name:     "科学计数法_正指数",
			input:    "2.5e+5",
			expected: true,
		},
		// 边界值
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
		{
			name:     "纯空格",
			input:    "   ",
			expected: false,
		},
		{
			name:     "前后有空格",
			input:    " 123 ",
			expected: false,
		},
		// 非法输入
		{
			name:     "包含字母",
			input:    "abc",
			expected: false,
		},
		{
			name:     "字母数字混合",
			input:    "12a34",
			expected: false,
		},
		{
			name:     "多个小数点",
			input:    "1.2.3",
			expected: false,
		},
		{
			name:     "只有符号_加号",
			input:    "+",
			expected: false,
		},
		{
			name:     "只有符号_减号",
			input:    "-",
			expected: false,
		},
		{
			name:     "数字分隔符_下划线",
			input:    "1_000",
			expected: true,
		},
		{
			name:     "特殊字符_逗号",
			input:    "1,000",
			expected: false,
		},
		{
			name:     "十六进制",
			input:    "0x1F",
			expected: false,
		},
		{
			name:     "二进制",
			input:    "0b101",
			expected: false,
		},
		// 特殊值
		{
			name:     "NaN",
			input:    "NaN",
			expected: false,
		},
		{
			name:     "nan",
			input:    "nan",
			expected: false,
		},
		{
			name:     "Inf",
			input:    "Inf",
			expected: false,
		},
		{
			name:     "inf",
			input:    "inf",
			expected: false,
		},
		{
			name:     "+Inf",
			input:    "+Inf",
			expected: false,
		},
		{
			name:     "-Inf",
			input:    "-Inf",
			expected: false,
		},
		{
			name:     "Infinity",
			input:    "Infinity",
			expected: false,
		},
		{
			name:     "-Infinity",
			input:    "-Infinity",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumber(tt.input)
			assert.Equal(t, tt.expected, result, "IsNumber(%q) 结果应与预期一致", tt.input)
		})
	}
}

// TestIsNumber_Parallel 并发安全测试_IsNumber
func TestIsNumber_Parallel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"整数", "123", true},
		{"浮点数", "123.456", true},
		{"科学计数法", "1e10", true},
		{"非法输入", "abc", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 100; i++ {
				result := IsNumber(tt.input)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// BenchmarkIsNumber 性能基准测试_IsNumber
func BenchmarkIsNumber(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{"整数", "123456"},
		{"浮点数", "123.456789"},
		{"科学计数法", "1.23456789e-10"},
		{"非法输入", "not_a_number"},
		{"空字符串", ""},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = IsNumber(tt.input)
			}
		})
	}
}

// BenchmarkIsNumber_Parallel 并发性能基准测试
func BenchmarkIsNumber_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = IsNumber("123.456")
		}
	})
}

// ExampleIsNumber 示例代码_IsNumber
func ExampleIsNumber() {
	fmt.Println(IsNumber("123"))
	fmt.Println(IsNumber("-45.67"))
	fmt.Println(IsNumber("1.5e-10"))
	fmt.Println(IsNumber("abc"))
	fmt.Println(IsNumber(""))

	// Output:
	// true
	// true
	// true
	// false
	// false
}

// TestIsDigitsOnly 测试 IsDigitsOnly 函数
// 测试场景：纯数字、包含字母、空字符串、特殊字符、混合内容、边界条件
func TestIsDigitsOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"纯数字", "12345", true},
		{"纯数字零", "0", true},
		{"单数字", "7", true},
		{"包含字母", "123abc", false},
		{"包含空格", "123 456", false},
		{"空字符串", "", false},
		{"特殊字符", "123!456", false},
		{"负号", "-123", false},
		{"小数点", "12.34", false},
		{"字母开头", "a12345", false},
		{"字母结尾", "12345a", false},
		{"中文数字", "123四56", false},
		{"Unicode数字", "١٢٣", false}, // Arabic numerals
		{"大数字", "99999999999999999999", true},
		{"前导零", "000123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDigitsOnly(tt.input)
			assert.Equal(t, tt.expected, result, "IsDigitsOnly(%q) = %v; want %v", tt.input, result, tt.expected)
		})
	}
}

// BenchmarkIsDigitsOnly 性能基准测试
func BenchmarkIsDigitsOnly(b *testing.B) {
	inputs := []string{
		"12345",
		"12345abc",
		"",
		"99999999999999999999",
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			IsDigitsOnly(input)
		}
	}
}
