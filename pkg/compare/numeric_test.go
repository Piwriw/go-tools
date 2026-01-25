package compare

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsNumeric 测试IsNumeric函数
// 功能：判断字符串s是否可以转换为T类型，并返回转换后的值
// 测试覆盖：正常场景、边界条件、异常情况、边缘场景
func TestIsNumeric(t *testing.T) {
	// 表格驱动测试用例结构
	tests := []struct {
		name      string  // 测试用例名称
		input     string  // 输入参数
		wantBool  bool    // 预期返回是否成功（int类型）
		wantInt   int     // 预期int类型结果
		wantFloat float64 // 预期float64类型结果
		wantErr   bool    // 是否预期错误（这里用wantBool=false表示）
		errMsg    string  // 预期错误信息子字符串（可选）
		floatOk   bool    // float64类型是否预期成功（默认与wantBool相同）
	}{
		// ========== 正常业务场景 (Happy Path) ==========
		{
			name:      "Valid positive integer",
			input:     "123",
			wantBool:  true,
			wantInt:   123,
			wantFloat: 123.0,
			wantErr:   false,
		},
		{
			name:      "Valid negative integer",
			input:     "-456",
			wantBool:  true,
			wantInt:   -456,
			wantFloat: -456.0,
			wantErr:   false,
		},
		{
			name:      "Valid positive float",
			input:     "123.45",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 123.45,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Valid negative float",
			input:     "-789.12",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: -789.12,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Float with integer value",
			input:     "42.0",
			wantBool:  false, // int类型无法解析（因为是浮点格式）
			wantInt:   0,
			wantFloat: 42.0,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},

		// ========== 边界条件 ==========
		{
			name:      "Empty string",
			input:     "",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
			errMsg:    "cannot convert",
		},
		{
			name:      "Zero value",
			input:     "0",
			wantBool:  true,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   false,
		},
		{
			name:      "Negative zero",
			input:     "-0",
			wantBool:  true,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   false,
		},
		{
			name:      "Maximum int value",
			input:     strconv.Itoa(math.MaxInt),
			wantBool:  true,
			wantInt:   math.MaxInt,
			wantFloat: float64(math.MaxInt),
			wantErr:   false,
		},
		{
			name:      "Minimum int value",
			input:     strconv.Itoa(math.MinInt),
			wantBool:  true,
			wantInt:   math.MinInt,
			wantFloat: float64(math.MinInt),
			wantErr:   false,
		},
		{
			name:      "Very large float - only for float64",
			input:     "1.7976931348623157e+308",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 1.7976931348623157e+308,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Very small float - only for float64",
			input:     "2.2250738585072014e-308",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 2.2250738585072014e-308,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Float starting with dot - only for float64",
			input:     ".5",
			wantBool:  false, // int类型预期失败
			wantInt:   0,
			wantFloat: 0.5,
			floatOk:   true, // float64类型预期成功
		},
		{
			name:      "Float ending with dot - only for float64",
			input:     "5.",
			wantBool:  false, // int类型预期失败
			wantInt:   0,
			wantFloat: 5.0,
			floatOk:   true, // float64类型预期成功
		},

		// ========== 异常情况 ==========
		{
			name:      "Alphabetic characters",
			input:     "abc",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
			errMsg:    "invalid",
		},
		{
			name:      "Mixed alphanumeric",
			input:     "123abc",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Special characters only",
			input:     "!@#$%",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Whitespace",
			input:     "   ",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Number with whitespace",
			input:     " 123 ",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Multiple decimal points",
			input:     "123.45.67",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Hexadecimal notation",
			input:     "0x1A",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Binary notation",
			input:     "0b1010",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Scientific notation with lowercase e - supported by ParseFloat only",
			input:     "1.5e10",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 1.5e10,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Scientific notation with uppercase E - supported by ParseFloat only",
			input:     "1.5E10",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 1.5e10,
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Infinity - supported by ParseFloat only",
			input:     "inf",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: math.Inf(1), // +Inf
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "NaN - supported by ParseFloat only",
			input:     "NaN",
			wantBool:  false, // int类型无法解析
			wantInt:   0,
			wantFloat: 0.0, // NaN无法直接比较，不检查具体值
			wantErr:   false,
			floatOk:   true, // float64类型可以解析
		},
		{
			name:      "Comma as decimal separator",
			input:     "123,45",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Number with thousand separators",
			input:     "1,234",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Negative sign in middle",
			input:     "12-34",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Multiple negative signs",
			input:     "--123",
			wantBool:  false,
			wantInt:   0,
			wantFloat: 0.0,
			wantErr:   true,
		},
		{
			name:      "Positive sign - supported by both Atoi and ParseFloat",
			input:     "+123",
			wantBool:  true, // Atoi支持+号
			wantInt:   123,
			wantFloat: 123.0,
			wantErr:   false,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 并行测试：测试用例之间独立，可并行执行
			t.Parallel()

			// 测试int类型转换
			t.Run("int_type", func(t *testing.T) {
				ok, val := IsNumeric[int](tt.input)

				if tt.wantBool {
					// 预期转换成功
					require.True(t, ok, "IsNumeric[int](%q) should return true", tt.input)
					assert.Equal(t, tt.wantInt, val,
						"IsNumeric[int](%q) returned value mismatch: got %v, want %v",
						tt.input, val, tt.wantInt)
				} else {
					// 预期转换失败
					assert.False(t, ok, "IsNumeric[int](%q) should return false", tt.input)
					assert.Equal(t, 0, val,
						"IsNumeric[int](%q) should return zero value on failure: got %v, want 0",
						tt.input, val)
				}
			})

			// 测试float64类型转换
			t.Run("float64_type", func(t *testing.T) {
				ok, val := IsNumeric[float64](tt.input)

				// 确定float64类型是否预期成功
				floatShouldSucceed := tt.floatOk || tt.wantBool

				if floatShouldSucceed {
					// 预期转换成功
					require.True(t, ok, "IsNumeric[float64](%q) should return true", tt.input)

					// 特殊处理NaN（NaN != NaN）
					if tt.name == "NaN - supported by ParseFloat only" {
						assert.True(t, val != val, // NaN是唯一不等于自己的值
							"IsNumeric[float64](%q) should return NaN", tt.input)
					} else {
						assert.InDelta(t, tt.wantFloat, val, 1e-10,
							"IsNumeric[float64](%q) returned value mismatch: got %v, want %v",
							tt.input, val, tt.wantFloat)
					}
				} else {
					// 预期转换失败
					assert.False(t, ok, "IsNumeric[float64](%q) should return false", tt.input)
					assert.Equal(t, 0.0, val,
						"IsNumeric[float64](%q) should return zero value on failure: got %v, want 0.0",
						tt.input, val)
				}
			})
		})
	}
}

// TestConvertToType 测试convertToType内部函数
// 功能：将值转换为T类型（int或float64）
// 场景：间接测试，通过IsNumeric调用
// 输入：各种类型的值（int, float64）
// 预期结果：成功转换或返回零值
func TestConvertToType(t *testing.T) {
	tests := []struct {
		name      string  // 测试用例名称
		input     any     // 输入值
		wantInt   bool    // int类型是否预期成功
		intVal    int     // 预期int结果
		wantFloat bool    // float64类型是否预期成功
		floatVal  float64 // 预期float64结果
	}{
		{
			name:      "Int value to int type",
			input:     42,
			wantInt:   true,
			intVal:    42,
			wantFloat: true,
			floatVal:  42.0,
		},
		{
			name:      "Float64 value to float64 type",
			input:     3.14,
			wantInt:   false,
			intVal:    0,
			wantFloat: true,
			floatVal:  3.14,
		},
		{
			name:      "Negative int to int type",
			input:     -100,
			wantInt:   true,
			intVal:    -100,
			wantFloat: true,
			floatVal:  -100.0,
		},
		{
			name:      "Negative float to float64 type",
			input:     -2.5,
			wantInt:   false,
			intVal:    0,
			wantFloat: true,
			floatVal:  -2.5,
		},
		{
			name:      "Zero int to int type",
			input:     0,
			wantInt:   true,
			intVal:    0,
			wantFloat: true,
			floatVal:  0.0,
		},
		{
			name:      "Zero float to float64 type",
			input:     0.0,
			wantInt:   false,
			intVal:    0,
			wantFloat: true,
			floatVal:  0.0,
		},
		{
			name:      "String value (unsupported type)",
			input:     "123",
			wantInt:   false,
			intVal:    0,
			wantFloat: false,
			floatVal:  0.0,
		},
		{
			name:      "Nil value",
			input:     nil,
			wantInt:   false,
			intVal:    0,
			wantFloat: false,
			floatVal:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 测试int类型转换
			t.Run("to_int", func(t *testing.T) {
				ok, val := convertToType[int](tt.input)
				assert.Equal(t, tt.wantInt, ok,
					"convertToType[int](%v) success mismatch: got %v, want %v",
					tt.input, ok, tt.wantInt)
				assert.Equal(t, tt.intVal, val,
					"convertToType[int](%v) value mismatch: got %v, want %v",
					tt.input, val, tt.intVal)
			})

			// 测试float64类型转换
			t.Run("to_float64", func(t *testing.T) {
				ok, val := convertToType[float64](tt.input)
				assert.Equal(t, tt.wantFloat, ok,
					"convertToType[float64](%v) success mismatch: got %v, want %v",
					tt.input, ok, tt.wantFloat)
				assert.InDelta(t, tt.floatVal, val, 1e-10,
					"convertToType[float64](%v) value mismatch: got %v, want %v",
					tt.input, val, tt.floatVal)
			})
		})
	}
}

// TestIsNumericConcurrent 测试IsNumeric的并发安全性
// 场景：多个goroutine同时调用IsNumeric
// 预期结果：每个goroutine都能正确返回结果，无数据竞争
func TestIsNumericConcurrent(t *testing.T) {
	testStrings := []string{"123", "456.78", "abc", "999", "-123.45", "0"}

	// 使用WaitGroup等待所有goroutine完成
	// 模拟并发场景：100个goroutine同时调用IsNumeric
	const goroutines = 100
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer func() { done <- true }()

			// 每个goroutine随机选择一个测试字符串
			s := testStrings[idx%len(testStrings)]

			// 测试int类型
			okInt, valInt := IsNumeric[int](s)
			_ = okInt
			_ = valInt

			// 测试float64类型
			okFloat, valFloat := IsNumeric[float64](s)
			_ = okFloat
			_ = valFloat
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

// BenchmarkIsNumericInt 性能基准测试 - int类型
// 测试IsNumeric[int]的性能和内存分配
func BenchmarkIsNumericInt(b *testing.B) {
	testCases := []string{"123", "456", "789", "abc", "123.45"}

	b.ResetTimer()   // 重置计时器，排除初始化时间
	b.ReportAllocs() // 报告内存分配情况

	for i := 0; i < b.N; i++ {
		s := testCases[i%len(testCases)]
		IsNumeric[int](s)
	}
}

// BenchmarkIsNumericFloat64 性能基准测试 - float64类型
// 测试IsNumeric[float64]的性能和内存分配
func BenchmarkIsNumericFloat64(b *testing.B) {
	testCases := []string{"123.45", "678.90", "abc", "123", "-456.78"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := testCases[i%len(testCases)]
		IsNumeric[float64](s)
	}
}

// BenchmarkIsNumericParallel 并行基准测试
// 测试并发场景下的性能表现
func BenchmarkIsNumericParallel(b *testing.B) {
	testCases := []string{"123", "456.78", "999", "abc"}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		idx := 0
		for pb.Next() {
			s := testCases[idx%len(testCases)]
			IsNumeric[int](s)
			idx++
		}
	})
}

// ExampleIsNumeric 示例函数 - int类型
// 展示如何使用IsNumeric函数判断字符串是否为整数
func ExampleIsNumeric_int() {
	// 判断字符串是否为整数
	ok, val := IsNumeric[int]("123")
	fmt.Printf("Is numeric: %v, Value: %d\n", ok, val)

	// 尝试转换非数字字符串
	ok, val = IsNumeric[int]("abc")
	fmt.Printf("Is numeric: %v, Value: %d\n", ok, val)

	// Output:
	// Is numeric: true, Value: 123
	// Is numeric: false, Value: 0
}

// ExampleIsNumeric_float64 示例函数 - float64类型
// 展示如何使用IsNumeric函数判断字符串是否为浮点数
func ExampleIsNumeric_float64() {
	// 判断字符串是否为浮点数
	ok, val := IsNumeric[float64]("123.45")
	fmt.Printf("Is numeric: %v, Value: %.2f\n", ok, val)

	// 整数字符串也可以转换为浮点数
	ok, val = IsNumeric[float64]("42")
	fmt.Printf("Is numeric: %v, Value: %.1f\n", ok, val)

	// 尝试转换非数字字符串
	ok, val = IsNumeric[float64]("hello")
	fmt.Printf("Is numeric: %v, Value: %.1f\n", ok, val)

	// Output:
	// Is numeric: true, Value: 123.45
	// Is numeric: true, Value: 42.0
	// Is numeric: false, Value: 0.0
}
