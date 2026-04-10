package format

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBytesToMBExact 测试字节精确转换为MB
// 测试场景：各种整数类型、浮点数类型、边界值
func TestBytesToMBExact(t *testing.T) {
	tests := []struct {
		name     string      // 测试用例名称
		input    interface{} // 输入参数
		expected float64     // 预期返回结果
	}{
		// int类型
		{
			name:     "int_1MB",
			input:    int(1024 * 1024),
			expected: 1.0,
		},
		{
			name:     "int_0",
			input:    int(0),
			expected: 0,
		},
		{
			name:     "int_负数",
			input:    int(-1024 * 1024),
			expected: -1.0,
		},
		// int64类型
		{
			name:     "int64_大数值",
			input:    int64(1024 * 1024 * 1024),
			expected: 1024.0,
		},
		// uint类型
		{
			name:     "uint_1MB",
			input:    uint(1024 * 1024),
			expected: 1.0,
		},
		// uint64类型
		{
			name:     "uint64_大数值",
			input:    uint64(1024 * 1024 * 1024 * 10),
			expected: 10240.0,
		},
		// float32类型
		{
			name:     "float32_1.5MB",
			input:    float32(1.5 * 1024 * 1024),
			expected: 1.5,
		},
		// float64类型
		{
			name:     "float64_2.5MB",
			input:    float64(2.5 * 1024 * 1024),
			expected: 2.5,
		},
		{
			name:     "float64_小数值",
			input:    float64(512 * 1024),
			expected: 0.5,
		},
		// 各种整数类型
		{
			name:     "int8_最大值",
			input:    int8(127),
			expected: 127.0 / (1024 * 1024),
		},
		{
			name:     "int16_1KB",
			input:    int16(1024),
			expected: 1.0 / 1024,
		},
		{
			name:     "int32_10MB",
			input:    int32(10 * 1024 * 1024),
			expected: 10.0,
		},
		{
			name:     "uint8_255",
			input:    uint8(255),
			expected: 255.0 / (1024 * 1024),
		},
		{
			name:     "uint16_最大值",
			input:    uint16(65535),
			expected: 65535.0 / (1024 * 1024),
		},
		{
			name:     "uint32_100MB",
			input:    uint32(100 * 1024 * 1024),
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case int8:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case int16:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case int32:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case int64:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint8:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint16:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint32:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint64:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case float32:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case float64:
				result := BytesToMBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			}
		})
	}
}

// TestBytesToGBExact 测试字节精确转换为GB
// 测试场景：各种整数类型、浮点数类型、边界值
func TestBytesToGBExact(t *testing.T) {
	tests := []struct {
		name     string      // 测试用例名称
		input    interface{} // 输入参数
		expected float64     // 预期返回结果
	}{
		{
			name:     "int_1GB",
			input:    int(1024 * 1024 * 1024),
			expected: 1.0,
		},
		{
			name:     "int_0",
			input:    int(0),
			expected: 0,
		},
		{
			name:     "int64_10GB",
			input:    int64(10 * 1024 * 1024 * 1024),
			expected: 10.0,
		},
		{
			name:     "uint64_100GB",
			input:    uint64(100 * 1024 * 1024 * 1024),
			expected: 100.0,
		},
		{
			name:     "float32_1.5GB",
			input:    float32(1.5 * 1024 * 1024 * 1024),
			expected: 1.5,
		},
		{
			name:     "float64_2.5GB",
			input:    float64(2.5 * 1024 * 1024 * 1024),
			expected: 2.5,
		},
		{
			name:     "float64_512MB",
			input:    float64(512 * 1024 * 1024),
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := BytesToGBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case int64:
				result := BytesToGBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case uint64:
				result := BytesToGBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case float32:
				result := BytesToGBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			case float64:
				result := BytesToGBExact(v)
				assert.InDelta(t, tt.expected, result, 0.0001, "转换结果应与预期一致")
			}
		})
	}
}

// TestBytesToMB 测试字节转换为MB（保留两位小数）
// 测试场景：四舍五入、边界值、负数
func TestBytesToMB(t *testing.T) {
	tests := []struct {
		name     string      // 测试用例名称
		input    interface{} // 输入参数
		expected interface{} // 预期返回结果
	}{
		{
			name:     "int_1MB",
			input:    int(1024 * 1024),
			expected: int(1),
		},
		{
			name:     "int_0",
			input:    int(0),
			expected: int(0),
		},
		{
			name:     "int_1.5MB",
			input:    int(1536 * 1024),
			expected: int(1), // 1.5
		},
		{
			name:     "float64_1.234MB",
			input:    float64(1.234 * 1024 * 1024),
			expected: float64(1.23),
		},
		{
			name:     "float64_1.235MB",
			input:    float64(1.235 * 1024 * 1024),
			expected: float64(1.24),
		},
		{
			name:     "float64_负数",
			input:    float64(-1024 * 1024),
			expected: float64(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := BytesToMB(v)
				assert.Equal(t, tt.expected, result, "转换结果应与预期一致")
			case float64:
				result := BytesToMB(v)
				assert.InDelta(t, tt.expected, result, 0.01, "转换结果应与预期一致")
			}
		})
	}
}

// TestBytesToGB 测试字节转换为GB（保留两位小数）
// 测试场景：四舍五入、边界值、大数值
func TestBytesToGB(t *testing.T) {
	tests := []struct {
		name     string      // 测试用例名称
		input    interface{} // 输入参数
		expected interface{} // 预期返回结果
	}{
		{
			name:     "int_1GB",
			input:    int(1024 * 1024 * 1024),
			expected: int(1),
		},
		{
			name:     "int_0",
			input:    int(0),
			expected: int(0),
		},
		{
			name:     "float64_1.234GB",
			input:    float64(1.234 * 1024 * 1024 * 1024),
			expected: float64(1.23),
		},
		{
			name:     "float64_1.235GB",
			input:    float64(1.235 * 1024 * 1024 * 1024),
			expected: float64(1.24),
		},
		{
			name:     "float64_大数值",
			input:    float64(5 * 1024 * 1024 * 1024 * 1024),
			expected: float64(5 * 1024),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := BytesToGB(v)
				assert.Equal(t, tt.expected, result, "转换结果应与预期一致")
			case float64:
				result := BytesToGB(v)
				assert.InDelta(t, tt.expected, result, 0.01, "转换结果应与预期一致")
			}
		})
	}
}

// TestKBToMB 测试KB转换为MB
// 测试场景：正常转换、四舍五入、边界值
func TestKBToMB(t *testing.T) {
	tests := []struct {
		name     string  // 测试用例名称
		kb       float64 // 输入参数
		expected float64 // 预期返回结果
	}{
		{
			name:     "正常转换_1024KB",
			kb:       1024,
			expected: 1,
		},
		{
			name:     "正常转换_1.5MB",
			kb:       1536,
			expected: 1.5,
		},
		{
			name:     "边界条件_0KB",
			kb:       0,
			expected: 0,
		},
		{
			name:     "边界条件_1KB",
			kb:       1,
			expected: 0,
		},
		{
			name:     "精度测试_四舍五入",
			kb:       1234,
			expected: 1.21,
		},
		{
			name:     "精度测试_进位",
			kb:       1235,
			expected: 1.21,
		},
		{
			name:     "大数值_5GB",
			kb:       5 * 1024 * 1024,
			expected: 5 * 1024,
		},
		{
			name:     "负数_-1MB",
			kb:       -1024,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KBToMB(tt.kb)
			assert.Equal(t, tt.expected, result, "转换结果应与预期一致")
		})
	}
}

// TestKBToGB 测试KB转换为GB
// 测试场景：正常转换、四舍五入、边界值
func TestKBToGB(t *testing.T) {
	tests := []struct {
		name     string  // 测试用例名称
		kb       float64 // 输入参数
		expected float64 // 预期返回结果
	}{
		{
			name:     "正常转换_1MB_in_KB",
			kb:       1024,
			expected: 0,
		},
		{
			name:     "正常转换_1GB",
			kb:       1024 * 1024,
			expected: 1,
		},
		{
			name:     "正常转换_1.5GB",
			kb:       1.5 * 1024 * 1024,
			expected: 1.5,
		},
		{
			name:     "边界条件_0KB",
			kb:       0,
			expected: 0,
		},
		{
			name:     "精度测试_1.234GB",
			kb:       1.234 * 1024 * 1024,
			expected: 1.23,
		},
		{
			name:     "精度测试_1.235GB",
			kb:       1.235 * 1024 * 1024,
			expected: 1.24,
		},
		{
			name:     "大数值_5TB",
			kb:       5 * 1024 * 1024 * 1024,
			expected: 5 * 1024,
		},
		{
			name:     "负数_-1GB",
			kb:       -1024 * 1024,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KBToGB(tt.kb)
			assert.Equal(t, tt.expected, result, "转换结果应与预期一致")
		})
	}
}

// TestBytesConversion_Parallel 并发安全测试_字节转换函数
func TestBytesConversion_Parallel(t *testing.T) {
	t.Run("BytesToMBExact并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := BytesToMBExact(1024 * 1024)
			assert.Equal(t, 1.0, result)
		}
	})

	t.Run("BytesToGBExact并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := BytesToGBExact(1024 * 1024 * 1024)
			assert.Equal(t, 1.0, result)
		}
	})

	t.Run("KBToMB并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := KBToMB(1024)
			assert.Equal(t, float64(1), result)
		}
	})

	t.Run("KBToGB并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := KBToGB(1024 * 1024)
			assert.Equal(t, float64(1), result)
		}
	})
}

// BenchmarkBytesToMBExact 性能基准测试_BytesToMBExact
func BenchmarkBytesToMBExact(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToMBExact(1024 * 1024)
	}
}

// BenchmarkBytesToGBExact 性能基准测试_BytesToGBExact
func BenchmarkBytesToGBExact(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToGBExact(1024 * 1024 * 1024)
	}
}

// BenchmarkBytesToMB 性能基准测试_BytesToMB
func BenchmarkBytesToMB(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToMB(1024 * 1024)
	}
}

// BenchmarkBytesToGB 性能基准测试_BytesToGB
func BenchmarkBytesToGB(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToGB(1024 * 1024 * 1024)
	}
}

// BenchmarkKBToMB 性能基准测试_KBToMB
func BenchmarkKBToMB(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = KBToMB(1024)
	}
}

// BenchmarkKBToGB 性能基准测试_KBToGB
func BenchmarkKBToGB(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = KBToGB(1024 * 1024)
	}
}

// ExampleBytesToMBExact 示例代码_BytesToMBExact
func ExampleBytesToMBExact() {
	// 将字节转换为MB（精确值）
	fmt.Printf("%.2f MB\n", BytesToMBExact(1024*1024)) // 1.00 MB
	fmt.Printf("%.2f MB\n", BytesToMBExact(1536*1024)) // 1.50 MB
	fmt.Printf("%.4f MB\n", BytesToMBExact(512*1024))  // 0.5000 MB

	// Output:
	// 1.00 MB
	// 1.50 MB
	// 0.5000 MB
}

// ExampleBytesToGBExact 示例代码_BytesToGBExact
func ExampleBytesToGBExact() {
	// 将字节转换为GB（精确值）
	fmt.Printf("%.2f GB\n", BytesToGBExact(1024*1024*1024)) // 1.00 GB
	fmt.Printf("%.2f GB\n", BytesToGBExact(256*1024*1024))  // 0.25 GB

	// Output:
	// 1.00 GB
	// 0.25 GB
}

// ExampleBytesToMB 示例代码_BytesToMB
func ExampleBytesToMB() {
	// 将字节转换为MB（保留两位小数）
	fmt.Println(BytesToMB(1024 * 1024))                  // 1
	fmt.Println(BytesToMB(1536 * 1024))                  // 1
	fmt.Println(BytesToMB(float64(1.234 * 1024 * 1024))) // 1.23

	// Output:
	// 1
	// 1
	// 1.23
}

// ExampleBytesToGB 示例代码_BytesToGB
func ExampleBytesToGB() {
	// 将字节转换为GB（保留两位小数）
	fmt.Println(BytesToGB(1024 * 1024 * 1024))                // 1
	fmt.Println(BytesToGB(float64(1.5 * 1024 * 1024 * 1024))) // 1.5

	// Output:
	// 1
	// 1.5
}

// ExampleKBToMB 示例代码_KBToMB
func ExampleKBToMB() {
	// 将KB转换为MB
	fmt.Println(KBToMB(1024)) // 1
	fmt.Println(KBToMB(1536)) // 1.5
	fmt.Println(KBToMB(500))  // 0.49

	// Output:
	// 1
	// 1.5
	// 0.49
}

// ExampleKBToGB 示例代码_KBToGB
func ExampleKBToGB() {
	// 将KB转换为GB
	fmt.Println(KBToGB(1024 * 1024)) // 1
	fmt.Println(KBToGB(1536 * 1024)) // 1.5

	// Output:
	// 1
	// 1.5
}

// --- 新增函数的测试 ---

func TestBytesToKB(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{name: "0字节", input: int(0), expected: 0},
		{name: "1KB", input: int(1024), expected: 1},
		{name: "1MB", input: int(1024 * 1024), expected: 1024},
		{name: "512字节", input: int(512), expected: 0.5},
		{name: "uint64", input: uint64(2048), expected: 2},
		{name: "float64_1.5KB", input: float64(1536), expected: 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				assert.InDelta(t, tt.expected, BytesToKB(v), 0.01)
			case uint64:
				assert.InDelta(t, tt.expected, BytesToKB(v), 0.01)
			case float64:
				assert.InDelta(t, tt.expected, BytesToKB(v), 0.01)
			}
		})
	}
}

func TestBytesToTB(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{name: "0字节", input: int(0), expected: 0},
		{name: "1TB", input: uint64(1024 * 1024 * 1024 * 1024), expected: 1},
		{name: "1.5TB", input: float64(1.5 * 1024 * 1024 * 1024 * 1024), expected: 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				assert.InDelta(t, tt.expected, BytesToTB(v), 0.01)
			case uint64:
				assert.InDelta(t, tt.expected, BytesToTB(v), 0.01)
			case float64:
				assert.InDelta(t, tt.expected, BytesToTB(v), 0.01)
			}
		})
	}
}

func TestMBToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected uint64
	}{
		{name: "0MB", input: 0, expected: 0},
		{name: "1MB", input: 1, expected: 1024 * 1024},
		{name: "1.5MB", input: 1.5, expected: uint64(1.5 * 1024 * 1024)},
		{name: "100MB", input: 100, expected: 100 * 1024 * 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, MBToBytes(tt.input))
		})
	}
}

func TestGBToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected uint64
	}{
		{name: "0GB", input: 0, expected: 0},
		{name: "1GB", input: 1, expected: 1024 * 1024 * 1024},
		{name: "2.5GB", input: 2.5, expected: uint64(2.5 * 1024 * 1024 * 1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GBToBytes(tt.input))
		})
	}
}

func TestTBToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected uint64
	}{
		{name: "0TB", input: 0, expected: 0},
		{name: "1TB", input: 1, expected: 1024 * 1024 * 1024 * 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TBToBytes(tt.input))
		})
	}
}

func TestHumanReadableBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{name: "0字节", input: int(0), expected: "0 B"},
		{name: "500字节", input: int(500), expected: "500 B"},
		{name: "1KB", input: int(1024), expected: "1.00 KB"},
		{name: "1MB", input: int(1024 * 1024), expected: "1.00 MB"},
		{name: "1GB", input: uint64(1024 * 1024 * 1024), expected: "1.00 GB"},
		{name: "1TB", input: uint64(1024 * 1024 * 1024 * 1024), expected: "1.00 TB"},
		{name: "1PB", input: uint64(1024 * 1024 * 1024 * 1024 * 1024), expected: "1.00 PB"},
		{name: "1.5MB", input: float64(1.5 * 1024 * 1024), expected: "1.50 MB"},
		{name: "512KB", input: int(512 * 1024), expected: "512.00 KB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				assert.Equal(t, tt.expected, HumanReadableBytes(v))
			case uint64:
				assert.Equal(t, tt.expected, HumanReadableBytes(v))
			case float64:
				assert.Equal(t, tt.expected, HumanReadableBytes(v))
			}
		})
	}
}

func TestParseHumanReadableBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint64
		expectError bool
	}{
		{name: "纯字节", input: "1024", expected: 1024},
		{name: "带B单位", input: "1024B", expected: 1024},
		{name: "带B单位和空格", input: "1024 B", expected: 1024},
		{name: "1KB", input: "1KB", expected: 1024},
		{name: "1 KB带空格", input: "1 KB", expected: 1024},
		{name: "1.5MB", input: "1.5MB", expected: uint64(1.5 * 1024 * 1024)},
		{name: "2 GB", input: "2 GB", expected: 2 * 1024 * 1024 * 1024},
		{name: "1TB小写", input: "1tb", expected: 1024 * 1024 * 1024 * 1024},
		{name: "1PB", input: "1PB", expected: 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "简写K", input: "1K", expected: 1024},
		{name: "简写M", input: "1M", expected: 1024 * 1024},
		{name: "简写G", input: "1G", expected: 1024 * 1024 * 1024},
		{name: "空字符串", input: "", expectError: true},
		{name: "无效单位", input: "1XX", expectError: true},
		{name: "无数字", input: "abc", expectError: true},
		{name: "负数", input: "-1MB", expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHumanReadableBytes(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// 新增函数的基准测试
func BenchmarkBytesToKB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BytesToKB(1024 * 1024)
	}
}

func BenchmarkBytesToTB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BytesToTB(uint64(1024 * 1024 * 1024 * 1024))
	}
}

func BenchmarkMBToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MBToBytes(1.5)
	}
}

func BenchmarkGBToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GBToBytes(2.5)
	}
}

func BenchmarkHumanReadableBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = HumanReadableBytes(uint64(1024 * 1024 * 1024))
	}
}

func BenchmarkParseHumanReadableBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseHumanReadableBytes("1.5 GB")
	}
}

// 新增函数的示例代码
func ExampleBytesToKB() {
	fmt.Println(BytesToKB(1024))        // 1
	fmt.Println(BytesToKB(512))         // 0.5
	fmt.Println(BytesToKB(1024 * 1024)) // 1024
	// Output:
	// 1
	// 0.5
	// 1024
}

func ExampleBytesToTB() {
	fmt.Println(BytesToTB(uint64(1024 * 1024 * 1024 * 1024))) // 1
	// Output:
	// 1
}

func ExampleMBToBytes() {
	fmt.Println(MBToBytes(1))   // 1048576
	fmt.Println(MBToBytes(1.5)) // 1572864
	// Output:
	// 1048576
	// 1572864
}

func ExampleGBToBytes() {
	fmt.Println(GBToBytes(1)) // 1073741824
	// Output:
	// 1073741824
}

func ExampleTBToBytes() {
	fmt.Println(TBToBytes(1)) // 1099511627776
	// Output:
	// 1099511627776
}

func ExampleHumanReadableBytes() {
	fmt.Println(HumanReadableBytes(int(0)))                     // 0 B
	fmt.Println(HumanReadableBytes(int(1024)))                  // 1.00 KB
	fmt.Println(HumanReadableBytes(int(1024 * 1024)))           // 1.00 MB
	fmt.Println(HumanReadableBytes(uint64(1024 * 1024 * 1024))) // 1.00 GB
	// Output:
	// 0 B
	// 1.00 KB
	// 1.00 MB
	// 1.00 GB
}

func ExampleParseHumanReadableBytes() {
	v, _ := ParseHumanReadableBytes("1.5GB")
	fmt.Println(v)
	v, _ = ParseHumanReadableBytes("1024 KB")
	fmt.Println(v)
	// Output:
	// 1610612736
	// 1048576
}
