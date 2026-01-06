package compare

import (
	"sort"
	"testing"
)

// TestIsNumeric_Int 测试 IsNumeric 函数，类型参数为 int
// 测试场景：有效整数、无效整数、浮点数字符串
func TestIsNumeric_Int(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   bool
		wantValue   int
		description string
	}{
		{
			name:        "有效正整数",
			input:       "123",
			wantValid:   true,
			wantValue:   123,
			description: "输入有效的正整数字符串",
		},
		{
			name:        "有效负整数",
			input:       "-456",
			wantValid:   true,
			wantValue:   -456,
			description: "输入有效的负整数字符串",
		},
		{
			name:        "零",
			input:       "0",
			wantValid:   true,
			wantValue:   0,
			description: "输入零",
		},
		{
			name:        "无效字符串",
			input:       "abc",
			wantValid:   false,
			wantValue:   0,
			description: "输入非数字字符串",
		},
		{
			name:        "空字符串",
			input:       "",
			wantValid:   false,
			wantValue:   0,
			description: "输入空字符串",
		},
		{
			name:        "浮点数字符串返回int",
			input:       "123.45",
			wantValid:   true,
			wantValue:   123,
			description: "输入浮点数字符串，转换为int",
		},
		{
			name:        "带空格的字符串",
			input:       " 123 ",
			wantValid:   false,
			wantValue:   0,
			description: "输入带空格的字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotValue := IsNumeric[int](tt.input)
			if gotValid != tt.wantValid {
				t.Errorf("IsNumeric() valid = %v, want %v", gotValid, tt.wantValid)
			}
			if gotValid && gotValue != tt.wantValue {
				t.Errorf("IsNumeric() value = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

// TestIsNumeric_Float64 测试 IsNumeric 函数，类型参数为 float64
// 测试场景：有效浮点数、无效浮点数、整数字符串
func TestIsNumeric_Float64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   bool
		wantValue   float64
		description string
	}{
		{
			name:        "有效正浮点数",
			input:       "123.45",
			wantValid:   true,
			wantValue:   123.45,
			description: "输入有效的正浮点数字符串",
		},
		{
			name:        "有效负浮点数",
			input:       "-789.12",
			wantValid:   true,
			wantValue:   -789.12,
			description: "输入有效的负浮点数字符串",
		},
		{
			name:        "整数字符串返回float64",
			input:       "42",
			wantValid:   true,
			wantValue:   42.0,
			description: "输入整数字符串，转换为float64",
		},
		{
			name:        "科学计数法",
			input:       "1.23e10",
			wantValid:   true,
			wantValue:   1.23e10,
			description: "输入科学计数法表示的数字",
		},
		{
			name:        "无效字符串",
			input:       "hello",
			wantValid:   false,
			wantValue:   0,
			description: "输入非数字字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotValue := IsNumeric[float64](tt.input)
			if gotValid != tt.wantValid {
				t.Errorf("IsNumeric() valid = %v, want %v", gotValid, tt.wantValid)
			}
			if gotValid && gotValue != tt.wantValue {
				t.Errorf("IsNumeric() value = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

// TestTopN 测试 TopN 函数
// 测试场景：空切片、切片长度小于n、切片长度等于n、切片长度大于n、有并列值
func TestTopN(t *testing.T) {
	tests := []struct {
		name        string
		items       []int
		n           int
		want        []int
		description string
	}{
		{
			name:        "空切片",
			items:       []int{},
			n:           5,
			want:        nil,
			description: "输入空切片应返回nil",
		},
		{
			name:        "切片长度小于n",
			items:       []int{1, 2, 3},
			n:           5,
			want:        []int{1, 2, 3},
			description: "切片长度小于n应返回整个切片",
		},
		{
			name:        "切片长度等于n",
			items:       []int{1, 2, 3},
			n:           3,
			want:        []int{1, 2, 3},
			description: "切片长度等于n应返回整个切片",
		},
		{
			name:        "切片长度大于n",
			items:       []int{1, 2, 3, 4, 5},
			n:           3,
			want:        []int{1, 2, 3},
			description: "切片长度大于n应返回前n个",
		},
		{
			name:        "有并列值-全部相同",
			items:       []int{5, 5, 5, 5, 5},
			n:           2,
			want:        []int{5, 5, 5, 5, 5},
			description: "所有值相同时应返回全部",
		},
		{
			name:        "有并列值-部分相同",
			items:       []int{10, 9, 9, 8, 7},
			n:           3,
			want:        []int{10, 9, 9, 9, 9},
			description: "有并列值时应包含所有并列项",
		},
		{
			name:        "n为0",
			items:       []int{1, 2, 3},
			n:           0,
			want:        nil,
			description: "n为0应返回nil",
		},
		{
			name:        "n为负数",
			items:       []int{1, 2, 3},
			n:           -1,
			want:        nil,
			description: "n为负数应返回nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopN(tt.items, tt.n)
			if len(got) != len(tt.want) {
				t.Errorf("TopN() 长度 = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if i < len(tt.want) && got[i] != tt.want[i] {
					t.Errorf("TopN()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestTopNSort 测试 TopNSort 函数
// 测试场景：需要排序的切片、有并列值
func TestTopNSort(t *testing.T) {
	// 定义一个实现 Sortable 接口的类型
	type IntSlice []int

	func (s IntSlice) Len() int           { return len(s) }
	func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
	func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
	func (s IntSlice) At(i int) int       { return s[i] }
	func (s IntSlice) Slice(start, end int) Sortable[int] {
		return IntSlice(s[start:end])
	}

	tests := []struct {
		name        string
		items       IntSlice
		n           int
		wantCount   int
		description string
	}{
		{
			name:        "空切片",
			items:       IntSlice{},
			n:           5,
			wantCount:   0,
			description: "输入空切片应返回空",
		},
		{
			name:        "需要排序的切片",
			items:       IntSlice{5, 2, 8, 1, 9},
			n:           3,
			wantCount:   3,
			description: "输入未排序切片应返回前3名",
		},
		{
			name:        "有并列值",
			items:       IntSlice{5, 5, 3, 2, 1},
			n:           2,
			wantCount:   2,
			description: "有并列值应包含所有并列项",
		},
		{
			name:        "n为0",
			items:       IntSlice{1, 2, 3},
			n:           0,
			wantCount:   0,
			description: "n为0应返回空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopNSort(tt.items, tt.n)
			if len(got) != tt.wantCount {
				t.Errorf("TopNSort() 长度 = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

// TestToSlice 测试 toSlice 辅助函数
func TestToSlice(t *testing.T) {
	type IntSlice []int

	func (s IntSlice) Len() int           { return len(s) }
	func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
	func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
	func (s IntSlice) At(i int) int       { return s[i] }
	func (s IntSlice) Slice(start, end int) Sortable[int] {
		return IntSlice(s[start:end])
	}

	t.Run("正常转换", func(t *testing.T) {
		slice := IntSlice{1, 2, 3, 4, 5}
		got := toSlice[int, IntSlice](slice)

		if len(got) != 5 {
			t.Errorf("期望长度=5, 实际=%d", len(got))
		}
		for i := range got {
			if got[i] != i+1 {
				t.Errorf("got[%d] = %v, want %d", i, got[i], i+1)
			}
		}
	})
}

// BenchmarkIsNumeric_Int 性能基准测试
func BenchmarkIsNumeric_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsNumeric[int]("12345")
	}
}

// BenchmarkIsNumeric_Float64 性能基准测试
func BenchmarkIsNumeric_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsNumeric[float64]("12345.67")
	}
}

// BenchmarkTopN 性能基准测试
func BenchmarkTopN(b *testing.B) {
	items := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TopN(items, 100)
	}
}

// BenchmarkTopNSort 性能基准测试
func BenchmarkTopNSort(b *testing.B) {
	type IntSlice []int
	func (s IntSlice) Len() int           { return len(s) }
	func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
	func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
	func (s IntSlice) At(i int) int       { return s[i] }
	func (s IntSlice) Slice(start, end int) Sortable[int] {
		return IntSlice(s[start:end])
	}

	items := make(IntSlice, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TopNSort(items, 100)
	}
}

// ExampleIsNumeric 示例函数
func ExampleIsNumeric() {
	// 判断字符串是否为整数
	valid, value := IsNumeric[int]("123")
	if valid {
		println(value)
	}

	// 判断字符串是否为浮点数
	validF, valueF := IsNumeric[float64]("123.45")
	if validF {
		println(valueF)
	}
}

// ExampleTopN 示例函数
func ExampleTopN() {
	// 从已排序的切片中取前3名（支持并列）
	scores := []int{100, 95, 95, 90, 85}
	top3 := TopN(scores, 3)
	// top3 会包含所有95分的成绩
	_ = top3
}

// ExampleTopNSort 示例函数
func ExampleTopNSort() {
	// 定义一个可排序的类型
	type Student struct {
		Name   string
		Score  int
	}

	type Students []Student
	func (s Students) Len() int           { return len(s) }
	func (s Students) Less(i, j int) bool { return s[i].Score < s[j].Score }
	func (s Students) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
	func (s Students) At(i int) int       { return s[i].Score }
	func (s Students) Slice(start, end int) Sortable[int] {
		return Students(s[start:end])
	}

	students := Students{
		{Name: "Alice", Score: 95},
		{Name: "Bob", Score: 87},
		{Name: "Charlie", Score: 92},
	}

	// 使用标准库排序
	sort.Sort(students)
	_ = students
}
