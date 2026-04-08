package format

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// CustomStringer 实现Stringer接口的自定义类型
type CustomStringer struct {
	Value string
}

func (c CustomStringer) String() string {
	return "custom: " + c.Value
}

// TestConvertMapKeysToUpper 测试将map的key转换为大写
// 测试场景：正常场景、空map、嵌套map、特殊字符
func TestConvertMapKeysToUpper(t *testing.T) {
	tests := []struct {
		name     string                 // 测试用例名称
		input    map[string]interface{} // 输入参数
		expected map[string]interface{} // 预期返回结果
	}{
		{
			name:     "正常场景_简单map",
			input:    map[string]interface{}{"name": "test", "age": 18},
			expected: map[string]interface{}{"NAME": "test", "AGE": 18},
		},
		{
			name:     "正常场景_混合大小写",
			input:    map[string]interface{}{"FirstName": "John", "lastName": "Doe", "AGE": 30},
			expected: map[string]interface{}{"FIRSTNAME": "John", "LASTNAME": "Doe", "AGE": 30},
		},
		{
			name:     "边界条件_空map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "边界条件_nil",
			input:    nil,
			expected: map[string]interface{}{}, // nil输入返回空map
		},
		{
			name:     "边界条件_特殊字符",
			input:    map[string]interface{}{"key-with-dash": 1, "key_with_underscore": 2, "key.with.dot": 3},
			expected: map[string]interface{}{"KEY-WITH-DASH": 1, "KEY_WITH_UNDERSCORE": 2, "KEY.WITH.DOT": 3},
		},
		{
			name:     "边界条件_中文key",
			input:    map[string]interface{}{"姓名": "张三", "年龄": 25},
			expected: map[string]interface{}{"姓名": "张三", "年龄": 25}, // 中文无大小写
		},
		{
			name:     "复杂场景_各种值类型",
			input:    map[string]interface{}{"string": "value", "int": 123, "float": 45.67, "bool": true, "nil": nil},
			expected: map[string]interface{}{"STRING": "value", "INT": 123, "FLOAT": 45.67, "BOOL": true, "NIL": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertMapKeysToUpper(tt.input)
			assert.Equal(t, tt.expected, result, "转换后的map应与预期一致")
			if tt.input != nil {
				assert.Equal(t, len(tt.input), len(result), "map长度应保持一致")
			}
		})
	}
}

// TestCovertToUpperSingle 测试单个字符串转大写
// 测试场景：正常场景、空字符串、特殊字符、中文、数字
func TestCovertToUpperSingle(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		input    string // 输入参数
		expected string // 预期返回结果
	}{
		{
			name:     "正常场景_小写转大写",
			input:    "hello",
			expected: "HELLO",
		},
		{
			name:     "正常场景_混合大小写",
			input:    "HeLLo WoRLd",
			expected: "HELLO WORLD",
		},
		{
			name:     "边界条件_空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "边界条件_单个字符",
			input:    "a",
			expected: "A",
		},
		{
			name:     "边界条件_已是大写",
			input:    "ALREADY UPPERCASE",
			expected: "ALREADY UPPERCASE",
		},
		{
			name:     "特殊字符_包含数字和符号",
			input:    "test123!@#",
			expected: "TEST123!@#",
		},
		{
			name:     "特殊字符_中文混合",
			input:    "test测试",
			expected: "TEST测试",
		},
		{
			name:     "特殊字符_仅中文",
			input:    "测试中文",
			expected: "测试中文",
		},
		{
			name:     "特殊字符_仅数字",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "特殊字符_特殊符号",
			input:    "!@#$%^&*()",
			expected: "!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToUpperSingle(tt.input)
			assert.Equal(t, tt.expected, result, "转换后的字符串应与预期一致")
		})
	}
}

// TestCovertToUpperMultiple 测试多个字符串转大写
// 测试场景：正常场景、空参数、单参数、混合类型
func TestCovertToUpperMultiple(t *testing.T) {
	tests := []struct {
		name     string   // 测试用例名称
		input    []string // 输入参数
		expected []string // 预期返回结果
	}{
		{
			name:     "正常场景_多个字符串",
			input:    []string{"hello", "world", "test"},
			expected: []string{"HELLO", "WORLD", "TEST"},
		},
		{
			name:     "正常场景_混合大小写",
			input:    []string{"HeLLo", "WoRLd", "TeST"},
			expected: []string{"HELLO", "WORLD", "TEST"},
		},
		{
			name:     "边界条件_空参数",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "边界条件_单参数",
			input:    []string{"single"},
			expected: []string{"SINGLE"},
		},
		{
			name:     "边界条件_包含空字符串",
			input:    []string{"", "test", ""},
			expected: []string{"", "TEST", ""},
		},
		{
			name:     "特殊场景_包含特殊字符",
			input:    []string{"test123", "!@#$%", "中文test"},
			expected: []string{"TEST123", "!@#$%", "中文TEST"},
		},
		{
			name:     "边缘场景_大量参数",
			input:    []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			expected: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建输入副本用于函数调用
			inputCopy := make([]string, len(tt.input))
			copy(inputCopy, tt.input)

			result := ConvertToUpperMultiple(inputCopy...)

			assert.Equal(t, tt.expected, result, "转换后的字符串切片应与预期一致")
			assert.Equal(t, len(tt.input), len(result), "切片长度应保持一致")
		})
	}
}

// TestToString 测试将任意类型转换为字符串
// 测试场景：各种基础类型、复合类型、指针类型、自定义类型
func TestToString(t *testing.T) {
	// 定义一个自定义结构体用于测试
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name     string      // 测试用例名称
		input    interface{} // 输入参数
		expected string      // 预期返回结果
	}{
		// 浮点数类型
		{
			name:     "浮点数_float64",
			input:    float64(123.456),
			expected: "123.456000",
		},
		{
			name:     "浮点数_float32",
			input:    float32(78.9),
			expected: "78.900002", // float32精度限制导致的实际输出
		},
		{
			name:     "浮点数_负数",
			input:    float64(-45.67),
			expected: "-45.670000",
		},
		// 整数类型
		{
			name:     "整数_int",
			input:    int(123),
			expected: "123",
		},
		{
			name:     "整数_int8",
			input:    int8(127),
			expected: "127",
		},
		{
			name:     "整数_int16",
			input:    int16(32767),
			expected: "32767",
		},
		{
			name:     "整数_int32",
			input:    int32(2147483647),
			expected: "2147483647",
		},
		{
			name:     "整数_int64_最大值",
			input:    int64(9223372036854775807),
			expected: "9223372036854775807",
		},
		{
			name:     "整数_uint",
			input:    uint(123),
			expected: "123",
		},
		{
			name:     "整数_uint8",
			input:    uint8(255),
			expected: "255",
		},
		{
			name:     "整数_uint16",
			input:    uint16(65535),
			expected: "65535",
		},
		{
			name:     "整数_uint32",
			input:    uint32(4294967295),
			expected: "4294967295",
		},
		// json.Number类型
		{
			name:     "json.Number_整数",
			input:    json.Number("123"),
			expected: "123",
		},
		{
			name:     "json.Number_浮点数",
			input:    json.Number("123.456"),
			expected: "123.456",
		},
		// 字符串类型
		{
			name:     "字符串_普通文本",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "字符串_空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "字符串_中文",
			input:    "你好世界",
			expected: "你好世界",
		},
		// 布尔类型
		{
			name:     "布尔值_true",
			input:    true,
			expected: "true",
		},
		{
			name:     "布尔值_false",
			input:    false,
			expected: "false",
		},
		// nil类型
		{
			name:     "nil",
			input:    nil,
			expected: "null",
		},
		// time.Time类型
		{
			name:     "time.Time_标准时间",
			input:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "2024-01-01T12:00:00Z",
		},
		// []byte类型
		{
			name:     "字节切片_普通文本",
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			name:     "字节切片_空切片",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "字节_slice_nil",
			input:    []byte(nil),
			expected: "",
		},
		// 结构体类型
		{
			name:     "结构体_简单结构",
			input:    Person{Name: "John", Age: 30},
			expected: `{"name":"John","age":30}`,
		},
		// map类型
		{
			name:     "map_字符串map",
			input:    map[string]interface{}{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "map_空map",
			input:    map[string]interface{}{},
			expected: `{}`,
		},
		// slice类型
		{
			name:     "slice_整数切片",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "slice_空切片",
			input:    []int{},
			expected: `[]`,
		},
		// 指针类型 - 注意：基本类型指针返回"unknown type"
		{
			name:     "指针_指向int",
			input:    func() *int { i := 123; return &i }(),
			expected: "", // 当前 ToString 失败时返回空字符串
		},
		{
			name:     "指针_指向结构体",
			input:    func() *Person { p := Person{Name: "Jane"}; return &p }(),
			expected: `{"name":"Jane","age":0}`,
		},
		// Stringer接口类型 - 注意：struct会被JSON序列化
		{
			name:     "Stringer接口_自定义类型",
			input:    CustomStringer{Value: "test"},
			expected: `{"Value":"test"}`, // struct会被JSON序列化
		},
		// 未知类型
		{
			name:     "未知类型_函数",
			input:    func() {},
			expected: "", // 当前 ToString 失败时返回空字符串
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			assert.Equal(t, tt.expected, result, "转换后的字符串应与预期一致")
		})
	}
}

// BenchmarkConvertMapKeysToUpper 性能基准测试_ConvertMapKeysToUpper
func BenchmarkConvertMapKeysToUpper(b *testing.B) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "SmallMap",
			input: map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"},
		},
		{
			name: "MediumMap",
			input: func() map[string]interface{} {
				m := make(map[string]interface{})
				for i := 0; i < 100; i++ {
					m[fmt.Sprintf("key%d", i)] = i
				}
				return m
			}(),
		},
		{
			name: "LargeMap",
			input: func() map[string]interface{} {
				m := make(map[string]interface{})
				for i := 0; i < 1000; i++ {
					m[fmt.Sprintf("key%d", i)] = i
				}
				return m
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ConvertMapKeysToUpper(tt.input)
			}
		})
	}
}

// BenchmarkCovertToUpperSingle 性能基准测试_CovertToUpperSingle
func BenchmarkCovertToUpperSingle(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "ShortString",
			input: "hello",
		},
		{
			name:  "MediumString",
			input: "this is a medium length string for testing",
		},
		{
			name: "LongString",
			input: func() string {
				s := ""
				for i := 0; i < 1000; i++ {
					s += "a"
				}
				return s
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ConvertToUpperSingle(tt.input)
			}
		})
	}
}

// BenchmarkToString 性能基准测试_ToString
func BenchmarkToString(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "Int",
			input: 123,
		},
		{
			name:  "Float64",
			input: 123.456,
		},
		{
			name:  "String",
			input: "hello world",
		},
		{
			name:  "Struct",
			input: Person{Name: "John", Age: 30},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ToString(tt.input)
			}
		})
	}
}

// ExampleConvertMapKeysToUpper 示例代码_ConvertMapKeysToUpper
func ExampleConvertMapKeysToUpper() {
	input := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}

	result := ConvertMapKeysToUpper(input)
	// Output: map[AGE:30 CITY:New York NAME:John]
	fmt.Println(result)
}

// ExampleConvertToUpperSingle 示例代码_ConvertToUpperSingle
func ExampleConvertToUpperSingle() {
	result := ConvertToUpperSingle("hello world")
	// Output: HELLO WORLD
	fmt.Println(result)
}

// ExampleConvertToUpperMultiple 示例代码_ConvertToUpperMultiple
func ExampleConvertToUpperMultiple() {
	result := ConvertToUpperMultiple("hello", "world", "test")
	// Output: [HELLO WORLD TEST]
	fmt.Println(result)
}

// ExampleToString 示例代码_ToString
func ExampleToString() {
	// 各种类型转换为字符串
	fmt.Println(ToString(123))     // int
	fmt.Println(ToString(45.67))   // float64
	fmt.Println(ToString("hello")) // string
	fmt.Println(ToString(true))    // bool
	fmt.Println(ToString(nil))     // nil

	// Output:
	// 123
	// 45.670000
	// hello
	// true
	// null
}

// ============================================================================
// IsBlank / DefaultIfEmpty
// ============================================================================

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", true},
		{"spaces", "   ", true},
		{"tabs", "\t\t", true},
		{"newlines", "\n\r", true},
		{"mixed whitespace", " \t\n ", true},
		{"non-blank", "hello", false},
		{"spaces with content", " a ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsBlank(tt.input))
		})
	}
}

func TestDefaultIfEmpty(t *testing.T) {
	assert.Equal(t, "default", DefaultIfEmpty("", "default"))
	assert.Equal(t, "hello", DefaultIfEmpty("hello", "default"))
	assert.Equal(t, "default", DefaultIfEmpty("", "default"))
	assert.Equal(t, "a", DefaultIfEmpty("a", "default"))
}

// ============================================================================
// Truncate / Substring
// ============================================================================

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		suffix string
		want   string
	}{
		{"short enough", "hello", 10, "...", "hello"},
		{"exact length", "hello", 5, "...", "hello"},
		{"truncate", "hello world", 8, "...", "hello..."},
		{"very small maxLen", "hello", 2, "...", ".."},
		{"empty input", "", 5, "...", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Truncate(tt.input, tt.maxLen, tt.suffix))
		})
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		name  string
		input string
		start int
		end   int
		want  string
	}{
		{"basic", "hello", 1, 4, "ell"},
		{"chinese", "你好世界", 1, 3, "好世"},
		{"start out of range", "hello", -1, 3, "hel"},
		{"end out of range", "hello", 2, 100, "llo"},
		{"start >= end", "hello", 3, 2, ""},
		{"empty string", "", 0, 5, ""},
		{"full range", "abc", 0, 3, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Substring(tt.input, tt.start, tt.end))
		})
	}
}

// ============================================================================
// Case Conversion
// ============================================================================

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"snake_case", "user_name", "userName"},
		{"kebab-case", "user-name", "userName"},
		{"multiple underscores", "first_name_last_name", "firstNameLastName"},
		{"already camel", "userName", "userName"},
		{"single word", "name", "name"},
		{"empty", "", ""},
		{"PascalCase input", "UserName", "UserName"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToCamelCase(tt.input))
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"camelCase", "userName", "user_name"},
		{"PascalCase", "UserName", "user_name"},
		{"multiple caps", "firstNameLastName", "first_name_last_name"},
		{"already snake", "user_name", "user_name"},
		{"single word", "name", "name"},
		{"empty", "", ""},
		{"acronym", "HTMLParser", "h_t_m_l_parser"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToSnakeCase(tt.input))
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"camelCase", "userName", "user-name"},
		{"PascalCase", "UserName", "user-name"},
		{"already kebab", "user-name", "user-name"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToKebabCase(tt.input))
		})
	}
}

// ============================================================================
// Masking
// ============================================================================

func TestMaskString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		keepLen int
		mask    rune
		want    string
	}{
		{"normal", "13812345678", 3, '*', "138*****678"},
		{"short string", "abc", 2, '*', "***"},
		{"exact keepLen*2", "abcde", 2, '*', "ab*de"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MaskString(tt.input, tt.keepLen, tt.mask))
		})
	}
}

func TestMaskPhone(t *testing.T) {
	assert.Equal(t, "138*****678", MaskPhone("13812345678"))
	assert.Equal(t, "*****", MaskPhone("13456"))
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal", "test@example.com", "t***@example.com"},
		{"short local", "a@b.com", "a@b.com"},
		{"no at", "invalid-email", "invalid-email"},
		{"long local", "admin@company.org", "a****@company.org"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MaskEmail(tt.input))
		})
	}
}

// ============================================================================
// Padding
// ============================================================================

func TestPadLeft(t *testing.T) {
	assert.Equal(t, "00123", PadLeft("123", '0', 5))
	assert.Equal(t, "abc", PadLeft("abc", '0', 2)) // already long enough
	assert.Equal(t, "  abc", PadLeft("abc", ' ', 5))
}

func TestPadRight(t *testing.T) {
	assert.Equal(t, "12300", PadRight("123", '0', 5))
	assert.Equal(t, "abc", PadRight("abc", '0', 2))
	assert.Equal(t, "abc  ", PadRight("abc", ' ', 5))
}

// ============================================================================
// Reverse / RemoveSpaces
// ============================================================================

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "hello", "olleh"},
		{"chinese", "你好世界", "界世好你"},
		{"empty", "", ""},
		{"single", "a", "a"},
		{"palindrome", "aba", "aba"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Reverse(tt.input))
		})
	}
}

func TestRemoveSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"spaces", "a b c", "abc"},
		{"tabs and newlines", "a\tb\nc", "abc"},
		{"no spaces", "abc", "abc"},
		{"all spaces", "  \t\n ", ""},
		{"empty", "", ""},
		{"mixed", "h e l l o", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RemoveSpaces(tt.input))
		})
	}
}
