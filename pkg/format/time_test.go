package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseCSTTime 测试解析中国标准时间
// 测试场景：正常场景、边界值、异常格式
func TestParseCSTTime(t *testing.T) {
	tests := []struct {
		name        string                      // 测试用例名称
		layout      string                      // 时间格式
		cstTime     string                      // CST时间字符串
		wantErr     bool                        // 是否预期发生错误
		errContains string                      // 预期错误信息包含的字符串
		validate    func(*testing.T, time.Time) // 验证函数
	}{
		{
			name:    "正常场景_标准格式",
			layout:  "2006-01-02 15:04:05",
			cstTime: "2024-01-01 12:00:00",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 2024, tm.Year())
				assert.Equal(t, time.January, tm.Month())
				assert.Equal(t, 1, tm.Day())
				assert.Equal(t, 12, tm.Hour())
			},
		},
		{
			name:    "正常场景_带日期时间格式",
			layout:  time.RFC3339,
			cstTime: "2024-01-01T12:00:00+08:00",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 2024, tm.Year())
				assert.Equal(t, int(time.January), int(tm.Month()))
				assert.Equal(t, 1, tm.Day())
			},
		},
		{
			name:    "边界条件_午夜时间",
			layout:  "2006-01-02 15:04:05",
			cstTime: "2024-01-01 00:00:00",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 0, tm.Hour())
				assert.Equal(t, 0, tm.Minute())
				assert.Equal(t, 0, tm.Second())
			},
		},
		{
			name:    "边界条件_一天结束时间",
			layout:  "2006-01-02 15:04:05",
			cstTime: "2024-01-01 23:59:59",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 23, tm.Hour())
				assert.Equal(t, 59, tm.Minute())
				assert.Equal(t, 59, tm.Second())
			},
		},
		{
			name:    "边界条件_最小年份",
			layout:  "2006-01-02 15:04:05",
			cstTime: "0001-01-01 00:00:00",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 1, tm.Year())
			},
		},
		{
			name:        "异常场景_格式不匹配",
			layout:      "2006-01-02 15:04:05",
			cstTime:     "2024/01/01 12:00:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_无效时间",
			layout:      "2006-01-02 15:04:05",
			cstTime:     "2024-13-01 12:00:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_空字符串",
			layout:      "2006-01-02 15:04:05",
			cstTime:     "",
			wantErr:     true,
			errContains: "",
		},
		{
			name:    "正常场景_只有时间",
			layout:  "15:04:05",
			cstTime: "12:30:45",
			wantErr: false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, 12, tm.Hour())
				assert.Equal(t, 30, tm.Minute())
				assert.Equal(t, 45, tm.Second())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCSTTime(tt.layout, tt.cstTime)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
			} else {
				require.NoError(t, err, "预期不发生错误")
				if tt.validate != nil {
					tt.validate(t, result)
				}
				// 验证时区确实是上海时区
				loc, _ := time.LoadLocation("Asia/Shanghai")
				assert.Equal(t, loc.String(), result.Location().String(), "时区应为Asia/Shanghai")
			}
		})
	}
}

// TestTimeStampDateTime 测试时间戳转换为日期时间格式
// 测试场景：正常场景、边界值、特殊时间戳
func TestTimeStampDateTime(t *testing.T) {
	tests := []struct {
		name        string                   // 测试用例名称
		timestamp   int64                    // Unix时间戳
		wantErr     bool                     // 是否预期发生错误
		errContains string                   // 预期错误信息包含的字符串
		validate    func(*testing.T, string) // 验证函数
	}{
		{
			name:      "正常场景_标准时间戳",
			timestamp: 1704067200, // 2024-01-01 00:00:00 UTC
			wantErr:   false,
			validate: func(t *testing.T, result string) {
				// 验证返回格式为 YYYY-MM-DD HH:MM:SS
				assert.Len(t, result, 19)
				assert.Contains(t, result, "-")
				assert.Contains(t, result, ":")
			},
		},
		{
			name:      "边界条件_时间戳0",
			timestamp: 0,
			wantErr:   false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:      "边界条件_负数时间戳",
			timestamp: -86400, // 1970-01-01前一天
			wantErr:   false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:      "正常场景_大数值时间戳",
			timestamp: 1735689600, // 2025-01-01 00:00:00 UTC
			wantErr:   false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeStampDateTime(tt.timestamp)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
			} else {
				require.NoError(t, err, "预期不发生错误")
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestTimeStampFormat 测试时间戳转换为自定义格式
// 测试场景：各种格式、边界值
func TestTimeStampFormat(t *testing.T) {
	tests := []struct {
		name      string // 测试用例名称
		timestamp int64  // Unix时间戳
		layout    string // 自定义格式
		wantErr   bool   // 是否预期发生错误
		expected  string // 预期返回结果（可选）
	}{
		{
			name:      "正常场景_RFC3339格式",
			timestamp: 1704067200,
			layout:    time.RFC3339,
			wantErr:   false,
		},
		{
			name:      "正常场景_自定义格式",
			timestamp: 1704067200,
			layout:    "2006-01-02",
			wantErr:   false,
		},
		{
			name:      "正常场景_只有时间",
			timestamp: 1704067200,
			layout:    "15:04:05",
			wantErr:   false,
		},
		{
			name:      "边界条件_时间戳0",
			timestamp: 0,
			layout:    time.RFC3339,
			wantErr:   false,
		},
		{
			name:      "正常场景_年月日格式",
			timestamp: 1704067200,
			layout:    "2006年01月02日",
			wantErr:   false,
		},
		{
			name:      "正常场景_完整日期时间",
			timestamp: 1704067200,
			layout:    "2006-01-02 15:04:05",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeStampFormat(tt.timestamp, tt.layout)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
			} else {
				require.NoError(t, err, "预期不发生错误")
				assert.NotEmpty(t, result, "结果不应为空")
			}
		})
	}
}

// TestCSTTime 测试获取当前北京时间
// 测试场景：正常获取、时区验证
func TestCSTTime(t *testing.T) {
	t.Run("正常场景_获取北京时间", func(t *testing.T) {
		result, err := CSTTime()
		require.NoError(t, err, "预期不发生错误")

		// 验证返回的是合理的时间
		assert.False(t, result.IsZero(), "时间不应为零值")

		// 验证时区确实是上海时区
		loc, _ := time.LoadLocation("Asia/Shanghai")
		assert.Equal(t, loc.String(), result.Location().String(), "时区应为Asia/Shanghai")

		// 验证时间在合理范围内（不能是太久之前或之后）
		now := time.Now()
		diff := now.Sub(result)
		assert.Less(t, diff.Abs(), time.Second, "时间差应小于1秒")
	})
}

// TestTimeZoneTime 测试获取指定时区的当前时间
// 测试场景：各种时区、空字符串、无效时区
func TestTimeZoneTime(t *testing.T) {
	tests := []struct {
		name        string                      // 测试用例名称
		timezone    string                      // 时区字符串
		wantErr     bool                        // 是否预期发生错误
		errContains string                      // 预期错误信息包含的字符串
		validate    func(*testing.T, time.Time) // 验证函数
	}{
		{
			name:     "正常场景_纽约时区",
			timezone: "America/New_York",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				loc, _ := time.LoadLocation("America/New_York")
				assert.Equal(t, loc.String(), tm.Location().String())
			},
		},
		{
			name:     "正常场景_伦敦时区",
			timezone: "Europe/London",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				loc, _ := time.LoadLocation("Europe/London")
				assert.Equal(t, loc.String(), tm.Location().String())
			},
		},
		{
			name:     "正常场景_东京时区",
			timezone: "Asia/Tokyo",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				loc, _ := time.LoadLocation("Asia/Tokyo")
				assert.Equal(t, loc.String(), tm.Location().String())
			},
		},
		{
			name:     "边界条件_空字符串_默认上海时区",
			timezone: "",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				loc, _ := time.LoadLocation("Asia/Shanghai")
				assert.Equal(t, loc.String(), tm.Location().String())
			},
		},
		{
			name:     "边界条件_UTC时区",
			timezone: "UTC",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				assert.Equal(t, "UTC", tm.Location().String())
			},
		},
		{
			name:        "异常场景_无效时区",
			timezone:    "Invalid/Timezone",
			wantErr:     true,
			errContains: "failed to load location",
		},
		{
			name:     "正常场景_洛杉矶时区",
			timezone: "America/Los_Angeles",
			wantErr:  false,
			validate: func(t *testing.T, tm time.Time) {
				loc, _ := time.LoadLocation("America/Los_Angeles")
				assert.Equal(t, loc.String(), tm.Location().String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeZoneTime(tt.timezone)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "错误信息应包含指定字符串")
				}
			} else {
				require.NoError(t, err, "预期不发生错误")
				assert.False(t, result.IsZero(), "时间不应为零值")
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestConvertToWeekdays 测试将int数组转换为Weekday数组
// 测试场景：正常场景、边界值、空数组
func TestConvertToWeekdays(t *testing.T) {
	tests := []struct {
		name     string         // 测试用例名称
		intArray []int          // 输入int数组
		expected []time.Weekday // 预期返回结果
	}{
		{
			name:     "正常场景_标准一周",
			intArray: []int{0, 1, 2, 3, 4, 5, 6},
			expected: []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
		},
		{
			name:     "边界条件_空数组",
			intArray: []int{},
			expected: []time.Weekday{},
		},
		{
			name:     "正常场景_工作日",
			intArray: []int{1, 2, 3, 4, 5},
			expected: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		{
			name:     "边界条件_单个元素",
			intArray: []int{0},
			expected: []time.Weekday{time.Sunday},
		},
		{
			name:     "正常场景_周末",
			intArray: []int{0, 6},
			expected: []time.Weekday{time.Sunday, time.Saturday},
		},
		{
			name:     "边界条件_只有6",
			intArray: []int{6},
			expected: []time.Weekday{time.Saturday},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToWeekdays(tt.intArray)
			assert.Equal(t, tt.expected, result, "转换结果应与预期一致")
			assert.Equal(t, len(tt.intArray), len(result), "数组长度应一致")
		})
	}
}

// TestTimeParsedFormat 测试解析时间字符串
// 测试场景：正常场景、边界值、异常格式
func TestTimeParsedFormat(t *testing.T) {
	tests := []struct {
		name        string                       // 测试用例名称
		timeStr     string                       // 时间字符串
		wantErr     bool                         // 是否预期发生错误
		errContains string                       // 预期错误信息包含的字符串
		validate    func(*testing.T, *time.Time) // 验证函数
	}{
		{
			name:    "正常场景_标准时间",
			timeStr: "12:00:00",
			wantErr: false,
			validate: func(t *testing.T, tm *time.Time) {
				assert.Equal(t, 12, tm.Hour())
				assert.Equal(t, 0, tm.Minute())
				assert.Equal(t, 0, tm.Second())
			},
		},
		{
			name:    "正常场景_午夜",
			timeStr: "00:00:00",
			wantErr: false,
			validate: func(t *testing.T, tm *time.Time) {
				assert.Equal(t, 0, tm.Hour())
				assert.Equal(t, 0, tm.Minute())
				assert.Equal(t, 0, tm.Second())
			},
		},
		{
			name:    "边界条件_一天结束",
			timeStr: "23:59:59",
			wantErr: false,
			validate: func(t *testing.T, tm *time.Time) {
				assert.Equal(t, 23, tm.Hour())
				assert.Equal(t, 59, tm.Minute())
				assert.Equal(t, 59, tm.Second())
			},
		},
		{
			name:    "正常场景_中午",
			timeStr: "12:30:45",
			wantErr: false,
			validate: func(t *testing.T, tm *time.Time) {
				assert.Equal(t, 12, tm.Hour())
				assert.Equal(t, 30, tm.Minute())
				assert.Equal(t, 45, tm.Second())
			},
		},
		{
			name:        "异常场景_格式错误_缺少秒",
			timeStr:     "12:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_格式错误_多余字符",
			timeStr:     "12:00:00:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_无效小时",
			timeStr:     "24:00:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_无效分钟",
			timeStr:     "12:60:00",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_无效秒",
			timeStr:     "12:00:60",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_空字符串",
			timeStr:     "",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "异常场景_非数字",
			timeStr:     "ab:cd:ef",
			wantErr:     true,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeParsedFormat(tt.timeStr)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
			} else {
				require.NoError(t, err, "预期不发生错误")
				require.NotNil(t, result, "结果不应为nil")
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestTimeFunctions_Parallel 并发安全测试_时间函数
func TestTimeFunctions_Parallel(t *testing.T) {
	t.Run("ParseCSTTime并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			_, err := ParseCSTTime("2006-01-02 15:04:05", "2024-01-01 12:00:00")
			require.NoError(t, err)
		}
	})

	t.Run("CSTTime并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			_, err := CSTTime()
			require.NoError(t, err)
		}
	})

	t.Run("TimeZoneTime并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			_, err := TimeZoneTime("America/New_York")
			require.NoError(t, err)
		}
	})

	t.Run("ConvertToWeekdays并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := ConvertToWeekdays([]int{0, 1, 2})
			assert.Len(t, result, 3)
		}
	})
}

// BenchmarkParseCSTTime 性能基准测试_ParseCSTTime
func BenchmarkParseCSTTime(b *testing.B) {
	layout := "2006-01-02 15:04:05"
	cstTime := "2024-01-01 12:00:00"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseCSTTime(layout, cstTime)
	}
}

// BenchmarkTimeStampDateTime 性能基准测试_TimeStampDateTime
func BenchmarkTimeStampDateTime(b *testing.B) {
	timestamp := int64(1704067200)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TimeStampDateTime(timestamp)
	}
}

// BenchmarkCSTTime 性能基准测试_CSTTime
func BenchmarkCSTTime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CSTTime()
	}
}

// BenchmarkTimeZoneTime 性能基准测试_TimeZoneTime
func BenchmarkTimeZoneTime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TimeZoneTime("America/New_York")
	}
}

// BenchmarkConvertToWeekdays 性能基准测试_ConvertToWeekdays
func BenchmarkConvertToWeekdays(b *testing.B) {
	intArray := []int{0, 1, 2, 3, 4, 5, 6}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConvertToWeekdays(intArray)
	}
}

// ExampleParseCSTTime 示例代码_ParseCSTTime
func ExampleParseCSTTime() {
	// 解析中国标准时间字符串
	tm, err := ParseCSTTime("2006-01-02 15:04:05", "2024-01-01 12:00:00")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(tm)
	// Output: 2024-01-01 12:00:00 +0800 CST
}

// ExampleTimeStampDateTime 示例代码_TimeStampDateTime
func ExampleTimeStampDateTime() {
	// 将时间戳转换为日期时间格式
	result, err := TimeStampDateTime(1704067200)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}

// ExampleTimeStampFormat 示例代码_TimeStampFormat
func ExampleTimeStampFormat() {
	// 将时间戳转换为自定义格式
	result1, _ := TimeStampFormat(1704067200, time.RFC3339)
	fmt.Println(result1)

	result2, _ := TimeStampFormat(1704067200, "2006-01-02")
	fmt.Println(result2)

	// Output:
	// 2024-01-01T08:00:00+08:00
	// 2024-01-01
}

// ExampleCSTTime 示例代码_CSTTime
func ExampleCSTTime() {
	// 获取当前北京时间
	tm, err := CSTTime()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(tm.Format("2006-01-02 15:04:05 MST"))
}

// ExampleTimeZoneTime 示例代码_TimeZoneTime
func ExampleTimeZoneTime() {
	// 获取指定时区的当前时间
	nyTime, _ := TimeZoneTime("America/New_York")
	fmt.Println("New York:", nyTime.Format("15:04:05"))

	tokyoTime, _ := TimeZoneTime("Asia/Tokyo")
	fmt.Println("Tokyo:", tokyoTime.Format("15:04:05"))

	// Output varies by current time
}

// ExampleConvertToWeekdays 示例代码_ConvertToWeekdays
func ExampleConvertToWeekdays() {
	// 将int数组转换为Weekday数组
	weekdays := ConvertToWeekdays([]int{0, 1, 2, 3, 4, 5, 6})
	fmt.Println(weekdays)

	// Output:
	// [Sunday Monday Tuesday Wednesday Thursday Friday Saturday]
}

// ExampleTimeParsedFormat 示例代码_TimeParsedFormat
func ExampleTimeParsedFormat() {
	// 解析时间字符串
	tm, err := TimeParsedFormat("12:30:45")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(tm.Format("15:04:05"))

	// Output:
	// 12:30:45
}
