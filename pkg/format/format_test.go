package format

import (
	"encoding/json"
	"testing"
	"time"
)

// TestBytesToMBExact 测试 BytesToMBExact 函数
// 测试场景：各种整数类型转换为MB
func TestBytesToMBExact(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		want        float64
		description string
	}{
		{name: "int类型", input: int(1048576), want: 1.0, description: "1MB"},
		{name: "int64类型", input: int64(2097152), want: 2.0, description: "2MB"},
		{name: "uint类型", input: uint(5242880), want: 5.0, description: "5MB"},
		{name: "float32类型", input: float32(1572864), want: 1.5, description: "1.5MB"},
		{name: "float64类型", input: float64(2097152), want: 2.0, description: "2MB"},
		{name: "零值", input: 0, want: 0.0, description: "0字节"},
		{name: "负数int", input: int(-1048576), want: -1.0, description: "-1MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				if got := BytesToMBExact(v); got != tt.want {
					t.Errorf("BytesToMBExact() = %v, want %v", got, tt.want)
				}
			case int64:
				if got := BytesToMBExact(v); got != tt.want {
					t.Errorf("BytesToMBExact() = %v, want %v", got, tt.want)
				}
			case uint:
				if got := BytesToMBExact(v); got != tt.want {
					t.Errorf("BytesToMBExact() = %v, want %v", got, tt.want)
				}
			case float32:
				if got := BytesToMBExact(v); got != tt.want {
					t.Errorf("BytesToMBExact() = %v, want %v", got, tt.want)
				}
			case float64:
				if got := BytesToMBExact(v); got != tt.want {
					t.Errorf("BytesToMBExact() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestBytesToGBExact 测试 BytesToGBExact 函数
func TestBytesToGBExact(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  float64
	}{
		{name: "1GB", input: int64(1073741824), want: 1.0},
		{name: "2GB", input: int64(2147483648), want: 2.0},
		{name: "1.5GB", input: float64(1610612736), want: 1.5},
		{name: "零值", input: 0, want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int64:
				if got := BytesToGBExact(v); got != tt.want {
					t.Errorf("BytesToGBExact() = %v, want %v", got, tt.want)
				}
			case float64:
				if got := BytesToGBExact(v); got != tt.want {
					t.Errorf("BytesToGBExact() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestBytesToMB 测试 BytesToMB 函数（保留两位小数）
func TestBytesToMB(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{name: "int类型1.5MB", input: int(1572864), want: int(1)},
		{name: "float64类型1.56MB", input: float64(1635778), want: float64(1.56)},
		{name: "int类型精确MB", input: int64(2097152), want: int64(2)},
		{name: "float32类型", input: float32(1572864), want: float32(1.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				got := BytesToMB(v)
				want := tt.want.(int)
				if got != want {
					t.Errorf("BytesToMB() = %v, want %v", got, want)
				}
			case float64:
				got := BytesToMB(v)
				want := tt.want.(float64)
				if got != want {
					t.Errorf("BytesToMB() = %v, want %v", got, want)
				}
			}
		})
	}
}

// TestBytesToGB 测试 BytesToGB 函数（保留两位小数）
func TestBytesToGB(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{name: "int类型1GB", input: int64(1073741824), want: int64(1)},
		{name: "float64类型1.5GB", input: float64(1610612736), want: float64(1.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int64:
				got := BytesToGB(v)
				want := tt.want.(int64)
				if got != want {
					t.Errorf("BytesToGB() = %v, want %v", got, want)
				}
			case float64:
				got := BytesToGB(v)
				want := tt.want.(float64)
				if got != want {
					t.Errorf("BytesToGB() = %v, want %v", got, want)
				}
			}
		})
	}
}

// TestKBToMB 测试 KBToMB 函数
func TestKBToMB(t *testing.T) {
	tests := []struct {
		name string
		kb   float64
		want float64
	}{
		{name: "1024KB转1MB", kb: 1024, want: 1.0},
		{name: "2048KB转2MB", kb: 2048, want: 2.0},
		{name: "512KB转0.5MB", kb: 512, want: 0.5},
		{name: "1536KB保留两位小数", kb: 1536, want: 1.5},
		{name: "零值", kb: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KBToMB(tt.kb); got != tt.want {
				t.Errorf("KBToMB() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestKBToGB 测试 KBToGB 函数
func TestKBToGB(t *testing.T) {
	tests := []struct {
		name string
		kb   float64
		want float64
	}{
		{name: "1048576KB转1GB", kb: 1048576, want: 1.0},
		{name: "2097152KB转2GB", kb: 2097152, want: 2.0},
		{name: "524288KB转0.5GB", kb: 524288, want: 0.5},
		{name: "零值", kb: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KBToGB(tt.kb); got != tt.want {
				t.Errorf("KBToGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFloat64ToPercentFloat 测试 Float64ToPercentFloat 函数
func TestFloat64ToPercentFloat(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want float64
	}{
		{name: "1转100%", f: 1.0, want: 100.0},
		{name: "0.5转50%", f: 0.5, want: 50.0},
		{name: "0.01转1%", f: 0.01, want: 1.0},
		{name: "0.001转0.1%", f: 0.001, want: 0.1},
		{name: "零值", f: 0, want: 0},
		{name: "负数", f: -0.5, want: -50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64ToPercentFloat(tt.f); got != tt.want {
				t.Errorf("Float64ToPercentFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToFloat64 测试 ToFloat64 函数
func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    float64
		wantErr bool
		errMsg  string
	}{
		{name: "float64类型", value: float64(3.14), want: 3.14, wantErr: false},
		{name: "float32类型", value: float32(2.71), want: 2.71, wantErr: false},
		{name: "int类型", value: int(42), want: 42.0, wantErr: false},
		{name: "int64类型", value: int64(100), want: 100.0, wantErr: false},
		{name: "uint类型", value: uint(25), want: 25.0, wantErr: false},
		{name: "bool true", value: true, want: 1.0, wantErr: false},
		{name: "bool false", value: false, want: 0.0, wantErr: false},
		{name: "string数字", value: "123.45", want: 123.45, wantErr: false},
		{name: "json.Number", value: json.Number("99.9"), want: 99.9, wantErr: false},
		{name: "nil", value: nil, want: 0, wantErr: true, errMsg: "cannot convert nil"},
		{name: "无效string", value: "abc", want: 0, wantErr: true, errMsg: "cannot convert string"},
		{name: "超大uint64", value: uint64(1 << 54), want: 0, wantErr: true, errMsg: "value too large"},
		{name: "不支持类型", value: []int{1, 2}, want: 0, wantErr: true, errMsg: "unsupported type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Errorf("期望错误，但未返回错误")
				}
				if tt.errMsg != "" && err != nil {
					// 检查错误信息是否包含预期子串
					errStr := err.Error()
					found := false
					for i := 0; i <= len(errStr)-len(tt.errMsg); i++ {
						if errStr[i:i+len(tt.errMsg)] == tt.errMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("错误信息不包含预期子串，期望包含'%s'，实际'%s'", tt.errMsg, errStr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("未期望错误，但返回了错误: %v", err)
				}
				if got != tt.want {
					t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestConvertMapKeysToUpper 测试 ConvertMapKeysToUpper 函数
func TestConvertMapKeysToUpper(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "正常转换",
			input: map[string]interface{}{"name": "Alice", "age": 30},
			want:  map[string]interface{}{"NAME": "Alice", "AGE": 30},
		},
		{
			name:  "已有大写键",
			input: map[string]interface{}{"NAME": "Bob", "age": 25},
			want:  map[string]interface{}{"NAME": "Bob", "AGE": 25},
		},
		{
			name:  "空map",
			input: map[string]interface{}{},
			want:  map[string]interface{}{},
		},
		{
			name:  "混合大小写",
			input: map[string]interface{}{"NaMe": "Charlie", "AgE": 35},
			want:  map[string]interface{}{"NAME": "Charlie", "AGE": 35},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertMapKeysToUpper(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("长度不匹配: got %d, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("key %s: got %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

// TestConvertToUpperSingle 测试 ConvertToUpperSingle 函数
func TestConvertToUpperSingle(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{name: "小写转大写", str: "hello", want: "HELLO"},
		{name: "已有大写", str: "WORLD", want: "WORLD"},
		{name: "混合", str: "HeLLo", want: "HELLO"},
		{name: "空字符串", str: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToUpperSingle(tt.str); got != tt.want {
				t.Errorf("ConvertToUpperSingle() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConvertToUpperMultiple 测试 ConvertToUpperMultiple 函数
func TestConvertToUpperMultiple(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "多个字符串转大写",
			input: []string{"hello", "world", "test"},
			want:  []string{"HELLO", "WORLD", "TEST"},
		},
		{
			name:  "空切片",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "混合大小写",
			input: []string{"HeLLo", "WoRLd"},
			want:  []string{"HELLO", "WORLD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToUpperMultiple(tt.input...)
			if len(got) != len(tt.want) {
				t.Errorf("长度不匹配: got %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d]: got %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestToString 测试 ToString 函数
func TestToString(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{name: "float64", value: float64(3.14), want: "3.140000"},
		{name: "int", value: int(42), want: "42"},
		{name: "string", value: "hello", want: "hello"},
		{name: "bool true", value: true, want: "true"},
		{name: "bool false", value: false, want: "false"},
		{name: "nil", value: nil, want: "null"},
		{name: "json.Number", value: json.Number("123.45"), want: "123.45"},
		{name: "time.Time", value: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), want: "2024-01-01T00:00:00Z"},
		{name: "[]byte", value: []byte("test"), want: "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseCSTTime 测试 ParseCSTTime 函数
func TestParseCSTTime(t *testing.T) {
	tests := []struct {
		name        string
		layout      string
		cstTime     string
		wantErr     bool
		description string
	}{
		{
			name:        "正常解析CST时间",
			layout:      "2006-01-02 15:04:05",
			cstTime:     "2024-01-01 12:00:00",
			wantErr:     false,
			description: "解析标准CST时间格式",
		},
		{
			name:        "错误的时间格式",
			layout:      "2006-01-02 15:04:05",
			cstTime:     "invalid",
			wantErr:     true,
			description: "输入无效时间字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCSTTime(tt.layout, tt.cstTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSTTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTimeStampDateTime 测试 TimeStampDateTime 函数
func TestTimeStampDateTime(t *testing.T) {
	tests := []struct {
		name        string
		timestamp   int64
		wantErr     bool
		description string
	}{
		{name: "正常时间戳", timestamp: 1704067200, wantErr: false, description: "2024-01-01 00:00:00 UTC"},
		{name: "零时间戳", timestamp: 0, wantErr: false, description: "1970-01-01 00:00:00 UTC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeStampDateTime(tt.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStampDateTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTimeStampFormat 测试 TimeStampFormat 函数
func TestTimeStampFormat(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		layout    string
		wantErr   bool
	}{
		{name: "自定义格式", timestamp: 1704067200, layout: "2006-01-02", wantErr: false},
		{name: "空layout", timestamp: 1704067200, layout: "", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeStampFormat(tt.timestamp, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStampFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCSTTime 测试 CSTTime 函数
func TestCSTTime(t *testing.T) {
	t.Run("获取当前CST时间", func(t *testing.T) {
		_, err := CSTTime()
		if err != nil {
			t.Errorf("CSTTime() error = %v", err)
		}
	})
}

// TestTimeZoneTime 测试 TimeZoneTime 函数
func TestTimeZoneTime(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{name: "默认时区", timezone: "", wantErr: false},
		{name: "指定时区", timezone: "America/New_York", wantErr: false},
		{name: "无效时区", timezone: "Invalid/Timezone", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeZoneTime(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeZoneTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConvertToWeekdays 测试 ConvertToWeekdays 函数
func TestConvertToWeekdays(t *testing.T) {
	tests := []struct {
		name     string
		intArray []int
		wantLen  int
	}{
		{name: "正常转换", intArray: []int{0, 1, 2, 3, 4, 5, 6}, wantLen: 7},
		{name: "空切片", intArray: []int{}, wantLen: 0},
		{name: "部分星期", intArray: []int{0, 1, 2}, wantLen: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToWeekdays(tt.intArray)
			if len(got) != tt.wantLen {
				t.Errorf("ConvertToWeekdays() 长度 = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

// TestTimeParsedFormat 测试 TimeParsedFormat 函数
func TestTimeParsedFormat(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantErr bool
	}{
		{name: "正常时间", timeStr: "15:04:05", wantErr: false},
		{name: "午夜", timeStr: "00:00:00", wantErr: false},
		{name: "无效时间", timeStr: "invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeParsedFormat(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeParsedFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// BenchmarkBytesToMBExact 性能基准测试
func BenchmarkBytesToMBExact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToMBExact(1048576)
	}
}

// BenchmarkToFloat64 性能基准测试
func BenchmarkToFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToFloat64(123.45)
	}
}

// BenchmarkConvertMapKeysToUpper 性能基准测试
func BenchmarkConvertMapKeysToUpper(b *testing.B) {
	input := map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertMapKeysToUpper(input)
	}
}

// ExampleBytesToMB 示例函数
func ExampleBytesToMB() {
	// 将字节转换为MB
	mb := BytesToMB(1048576)
	println(mb)
}

// ExampleToFloat64 示例函数
func ExampleToFloat64() {
	// 将各种类型转换为float64
	f, _ := ToFloat64(123)
	println(f)
}

// ExampleConvertMapKeysToUpper 示例函数
func ExampleConvertMapKeysToUpper() {
	// 将map的所有key转为大写
	input := map[string]interface{}{"name": "Alice", "age": 30}
	output := ConvertMapKeysToUpper(input)
	println(output["NAME"])
}
