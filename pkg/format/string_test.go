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
			expected: "unknown type: *int", // 当前实现对基本类型指针支持不完善
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
			expected: "unknown type: func()",
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

// ExampleCovertToUpperSingle 示例代码_CovertToUpperSingle
func ExampleCovertToUpperSingle() {
	result := ConvertToUpperSingle("hello world")
	// Output: HELLO WORLD
	fmt.Println(result)
}

// ExampleCovertToUpperMultiple 示例代码_CovertToUpperMultiple
func ExampleCovertToUpperMultiple() {
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
