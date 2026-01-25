package utils

import (
	"maps"
	"slices"
	"testing"
	"time"
)

// TestCopy 测试 Copy 函数
// 测试场景：nil值、基本类型、切片、map、结构体、指针、嵌套结构
func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		description string
	}{
		{
			name:        "nil值",
			input:       nil,
			description: "复制nil应返回nil",
		},
		{
			name:        "整数",
			input:       42,
			description: "复制整数",
		},
		{
			name:        "浮点数",
			input:       3.14,
			description: "复制浮点数",
		},
		{
			name:        "字符串",
			input:       "hello",
			description: "复制字符串",
		},
		{
			name:        "布尔值",
			input:       true,
			description: "复制布尔值",
		},
		{
			name:        "空切片",
			input:       []int{},
			description: "复制空切片",
		},
		{
			name:        "整数切片",
			input:       []int{1, 2, 3, 4, 5},
			description: "复制整数切片",
		},
		{
			name:        "字符串切片",
			input:       []string{"a", "b", "c"},
			description: "复制字符串切片",
		},
		{
			name:        "空map",
			input:       map[string]int{},
			description: "复制空map",
		},
		{
			name:        "字符串到整数的map",
			input:       map[string]int{"a": 1, "b": 2},
			description: "复制map",
		},
		{
			name:        "Time类型",
			input:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			description: "复制time.Time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Copy(tt.input)

			// 验证结果不为nil（除非输入是nil）
			if tt.input == nil {
				if result != nil {
					t.Errorf("期望nil, 实际=%v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Copy返回nil")
			}

			// 验证类型相同
			if typeof(result) != typeof(tt.input) {
				t.Errorf("类型不匹配: 期望=%T, 实际=%T", tt.input, result)
			}
		})
	}
}

// TestCopy_Struct 测试结构体深拷贝
func TestCopy_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	original := Person{Name: "Alice", Age: 30}
	copied := Copy(original)

	result, ok := copied.(Person)
	if !ok {
		t.Fatal("类型断言失败")
	}

	if result.Name != original.Name {
		t.Errorf("Name: 期望=%s, 实际=%s", original.Name, result.Name)
	}
	if result.Age != original.Age {
		t.Errorf("Age: 期望=%d, 实际=%d", original.Age, result.Age)
	}
}

// TestCopy_NestedStruct 测试嵌套结构体深拷贝
func TestCopy_NestedStruct(t *testing.T) {
	type Address struct {
		City string
	}

	type Person struct {
		Name    string
		Address Address
	}

	original := Person{
		Name:    "Bob",
		Address: Address{City: "Beijing"},
	}
	copied := Copy(original)

	result, ok := copied.(Person)
	if !ok {
		t.Fatal("类型断言失败")
	}

	if result.Name != original.Name {
		t.Errorf("Name: 期望=%s, 实际=%s", original.Name, result.Name)
	}
	if result.Address.City != original.Address.City {
		t.Errorf("Address.City: 期望=%s, 实际=%s", original.Address.City, result.Address.City)
	}
}

// TestCopy_SliceIndependence 测试切片拷贝的独立性
func TestCopy_SliceIndependence(t *testing.T) {
	original := []int{1, 2, 3}
	copied := Copy(original)

	result, ok := copied.([]int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 修改原始切片
	original[0] = 999
	original = append(original, 4)

	// 验证拷贝未被影响
	if result[0] == 999 {
		t.Error("修改原始切片影响了拷贝")
	}
	if len(result) != 3 {
		t.Errorf("拷贝长度应仍为3, 实际=%d", len(result))
	}
}

// TestCopy_MapIndependence 测试map拷贝的独立性
func TestCopy_MapIndependence(t *testing.T) {
	original := map[string]int{"a": 1, "b": 2}
	copied := Copy(original)

	result, ok := copied.(map[string]int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 修改原始map
	original["a"] = 999
	delete(original, "b")
	original["c"] = 3

	// 验证拷贝未被影响
	if result["a"] == 999 {
		t.Error("修改原始map影响了拷贝")
	}
	if result["b"] != 2 {
		t.Error("删除原始map中的key影响了拷贝")
	}
	if _, exists := result["c"]; exists {
		t.Error("向原始map添加key影响了拷贝")
	}
}

// TestCopy_Pointer 测试指针深拷贝
func TestCopy_Pointer(t *testing.T) {
	value := 42
	original := &value
	copied := Copy(original)

	result, ok := copied.(*int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 验证值相同
	if *result != *original {
		t.Errorf("值: 期望=%d, 实际=%d", *original, *result)
	}

	// 修改原始指针的值
	*original = 999

	// 验证拷贝未被影响
	if *result == 999 {
		t.Error("修改原始指针的值影响了拷贝")
	}

	// 验证指针地址不同
	if result == original {
		t.Error("拷贝应指向不同的地址")
	}
}

// TestCopy_NilPointer 测试nil指针深拷贝
func TestCopy_NilPointer(t *testing.T) {
	var original *int
	copied := Copy(original)

	result, ok := copied.(*int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// nil指针应保持nil
	if result != nil {
		t.Error("nil指针拷贝后应为nil")
	}
}

// TestCopy_Interface 测试interface深拷贝
func TestCopy_Interface(t *testing.T) {
	var original interface{} = []int{1, 2, 3}
	copied := Copy(original)

	result, ok := copied.([]int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	if !slices.Equal(result, original.([]int)) {
		t.Error("interface切片拷贝值不匹配")
	}

	// 修改原始interface
	original = []int{4, 5, 6}

	// 验证拷贝未被影响
	if len(result) != 3 {
		t.Error("修改原始interface影响了拷贝")
	}
}

// TestCopy_SliceOfPointers 测试指针切片深拷贝
func TestCopy_SliceOfPointers(t *testing.T) {
	a, b := 1, 2
	original := []*int{&a, &b}
	copied := Copy(original)

	result, ok := copied.([]*int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 验证值相同
	if *result[0] != *original[0] {
		t.Error("拷贝值不匹配")
	}

	// 修改原始指针
	*original[0] = 999

	// 验证拷贝未被影响
	if *result[0] == 999 {
		t.Error("修改原始指针影响了拷贝")
	}
}

// TestCopy_MapOfSlices 测试map值为切片的深拷贝
func TestCopy_MapOfSlices(t *testing.T) {
	original := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5, 6},
	}
	copied := Copy(original)

	result, ok := copied.(map[string][]int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 修改原始map中的切片
	original["a"][0] = 999

	// 验证拷贝未被影响
	if result["a"][0] == 999 {
		t.Error("修改原始map中的切片影响了拷贝")
	}
}

// TestIface 测试 Iface 函数
// 测试场景：Iface是Copy的别名，行为应相同
func TestIface(t *testing.T) {
	original := []int{1, 2, 3}
	copied := Iface(original)

	result, ok := copied.([]int)
	if !ok {
		t.Fatal("类型断言失败")
	}

	if !slices.Equal(result, original) {
		t.Error("Iface拷贝值不匹配")
	}
}

// TestCopy_CustomDeepCopy 测试实现Interface接口的自定义深拷贝
func TestCopy_CustomDeepCopy(t *testing.T) {
	type CustomStruct struct {
		Value int
	}

	func (c CustomStruct) DeepCopy() interface{} {
		return CustomStruct{Value: c.Value * 2}
	}

	original := CustomStruct{Value: 10}
	copied := Copy(original)

	result, ok := copied.(CustomStruct)
	if !ok {
		t.Fatal("类型断言失败")
	}

	// 自定义DeepCopy将值乘以2
	if result.Value != 20 {
		t.Errorf("自定义DeepCopy: 期望=20, 实际=%d", result.Value)
	}
}

// BenchmarkCopy_Slice 性能基准测试
func BenchmarkCopy_Slice(b *testing.B) {
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(slice)
	}
}

// BenchmarkCopy_Map 性能基准测试
func BenchmarkCopy_Map(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		m[string(rune(i))] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(m)
	}
}

// BenchmarkCopy_Struct 性能基准测试
func BenchmarkCopy_Struct(b *testing.B) {
	type LargeStruct struct {
		A, B, C, D, E int
		F, G, H, I, J string
	}
	s := LargeStruct{
		A: 1, B: 2, C: 3, D: 4, E: 5,
		F: "a", G: "b", H: "c", I: "d", J: "e",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(s)
	}
}

// ExampleCopy 示例函数
func ExampleCopy() {
	// 复制基本类型
	n := Copy(42)
	println(n.(int))

	// 复制切片
	slice := Copy([]int{1, 2, 3})
	println(len(slice.([]int)))

	// 复制map
	m := Copy(map[string]int{"a": 1, "b": 2})
	println(len(m.(map[string]int)))
}

// typeof 返回值的类型（辅助函数）
func typeof(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case float64:
		return "float64"
	case string:
		return "string"
	case bool:
		return "bool"
	case []int:
		return "[]int"
	case map[string]int:
		return "map[string]int"
	default:
		return "unknown"
	}
}
