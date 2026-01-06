package prom

import (
	"fmt"
	"math"
	"testing"
	"time"

	prommodel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 辅助函数：创建测试用的SamplePair切片
func makeSamplePairs(values []float64) []prommodel.SamplePair {
	result := make([]prommodel.SamplePair, len(values))
	now := prommodel.Time(time.Now().Unix() * 1000)
	for i, v := range values {
		result[i] = prommodel.SamplePair{
			Timestamp: now + prommodel.Time(i*1000),
			Value:     prommodel.SampleValue(v),
		}
	}
	return result
}

// 辅助函数：比较两个SamplePair切片是否大致相等
func samplePairsEqual(a, b []prommodel.SamplePair) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Timestamp != b[i].Timestamp {
			return false
		}
		if math.Abs(float64(a[i].Value-b[i].Value)) > 1e-9 {
			return false
		}
	}
	return true
}

// TestResultPairDownSample 测试DownSample方法
// 功能：使用默认的均匀采样模式进行降采样
// 测试覆盖：正常场景、边界条件、各种数据大小
func TestResultPairDownSample(t *testing.T) {
	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
	}{
		{
			name:          "No downsampling needed - points <= maxPoints",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5}),
			maxPoints:     10,
			expectedCount: 5, // 不进行降采样
		},
		{
			name:          "Basic downsampling - 10 points to 3",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			maxPoints:     3,
			expectedCount: 3, // 降采样到3个点
		},
		{
			name:          "Large dataset downsampling",
			input:         makeSamplePairs(make([]float64, 1000)),
			maxPoints:     100,
			expectedCount: 100, // 降采样到100个点
		},
		{
			name:          "Empty data",
			input:         []prommodel.SamplePair{},
			maxPoints:     10,
			expectedCount: 0,
		},
		{
			name:          "Single point",
			input:         makeSamplePairs([]float64{42}),
			maxPoints:     10,
			expectedCount: 1,
		},
		{
			name:          "Exact maxPoints",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5}),
			maxPoints:     5,
			expectedCount: 5,
		},
		{
			name:          "Default MAXPoints",
			input:         makeSamplePairs(make([]float64, 2000)),
			maxPoints:     -1,   // 使用默认值MAXPoint
			expectedCount: 1024, // MAXPoint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResultPair{Values: tt.input}

			if tt.maxPoints == -1 {
				rp.DownSample()
			} else {
				rp.DownSample(tt.maxPoints)
			}

			assert.Equal(t, tt.expectedCount, len(rp.Values),
				"DownSample(%d) should produce %d points", tt.maxPoints, tt.expectedCount)

			// 验证最后一个点是原始数据的最后一个点
			if len(tt.input) > 0 && len(rp.Values) > 0 {
				lastOriginal := tt.input[len(tt.input)-1]
				lastSampled := rp.Values[len(rp.Values)-1]
				assert.Equal(t, lastOriginal.Timestamp, lastSampled.Timestamp,
					"Last sampled point should be the last original point")
			}
		})
	}
}

// TestResultPairDownSampleWithOptions 测试DownSampleWithOptions方法
// 功能：使用配置选项进行降采样
// 测试覆盖：所有降采样模式、边界条件
func TestResultPairDownSampleWithOptions(t *testing.T) {
	tests := []struct {
		name          string                                                           // 测试用例名称
		input         []prommodel.SamplePair                                           // 输入数据
		opts          *DownSampleOptions                                               // 降采样选项
		expectedCount int                                                              // 预期输出数量
		validate      func(*testing.T, []prommodel.SamplePair, []prommodel.SamplePair) // 验证函数
	}{
		{
			name:  "Uniform mode - basic",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeUniform,
				MaxPoints: 3,
				Enabled:   true,
			},
			expectedCount: 3,
			validate: func(t *testing.T, original, sampled []prommodel.SamplePair) {
				// 验证等间隔采样
				assert.Equal(t, original[0].Timestamp, sampled[0].Timestamp)
				assert.Equal(t, original[len(original)-1].Timestamp, sampled[len(sampled)-1].Timestamp)
			},
		},
		{
			name:  "Max mode - preserve highest values",
			input: makeSamplePairs([]float64{1, 5, 2, 8, 3, 9, 4, 7, 6}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeMax,
				MaxPoints: 3,
				Enabled:   true,
			},
			expectedCount: 3,
		},
		{
			name:  "Min mode - preserve lowest values",
			input: makeSamplePairs([]float64{5, 1, 8, 2, 9, 3, 7, 4, 6}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeMin,
				MaxPoints: 3,
				Enabled:   true,
			},
			expectedCount: 3,
		},
		{
			name:  "Average mode - calculate averages",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeAverage,
				MaxPoints: 3,
				Enabled:   true,
			},
			expectedCount: 3,
		},
		{
			name:  "LTTB mode - visual downsampling",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeLTTB,
				MaxPoints: 5,
				Enabled:   true,
			},
			expectedCount: 5,
			validate: func(t *testing.T, original, sampled []prommodel.SamplePair) {
				// LTTB 应该保留第一个和最后一个点
				assert.Equal(t, original[0].Timestamp, sampled[0].Timestamp)
				assert.Equal(t, original[len(original)-1].Timestamp, sampled[len(sampled)-1].Timestamp)
			},
		},
		{
			name:  "Disabled downsampling",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeUniform,
				MaxPoints: 2,
				Enabled:   false, // 禁用
			},
			expectedCount: 10, // 不进行降采样
		},
		{
			name:          "Nil options",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5}),
			opts:          nil,
			expectedCount: 5, // 不进行降采样
		},
		{
			name:  "Zero MaxPoints - use default",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			opts: &DownSampleOptions{
				Mode:      DownSampleModeUniform,
				MaxPoints: 0,
				Enabled:   true,
			},
			expectedCount: 10, // 1024 is default, but input is smaller
		},
		{
			name:  "Empty mode - use uniform as default",
			input: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			opts: &DownSampleOptions{
				Mode:      "",
				MaxPoints: 3,
				Enabled:   true,
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResultPair{Values: tt.input}
			rp.DownSampleWithOptions(tt.opts)

			assert.Equal(t, tt.expectedCount, len(rp.Values),
				"DownSampleWithOptions produced unexpected count")

			if tt.validate != nil {
				tt.validate(t, tt.input, rp.Values)
			}
		})
	}
}

// TestUniformSampler 测试均匀采样器
// 功能：等间隔采样
func TestUniformSampler(t *testing.T) {
	sampler := &uniformSampler{}

	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
	}{
		{
			name:          "10 points to 3",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			maxPoints:     3,
			expectedCount: 3,
		},
		{
			name:          "100 points to 10",
			input:         makeSamplePairs(make([]float64, 100)),
			maxPoints:     10,
			expectedCount: 10,
		},
		{
			name:          "Fewer points than maxPoints",
			input:         makeSamplePairs([]float64{1, 2, 3}),
			maxPoints:     10,
			expectedCount: 3,
		},
		{
			name:          "Empty input",
			input:         []prommodel.SamplePair{},
			maxPoints:     10,
			expectedCount: 0,
		},
		{
			name:          "Single point",
			input:         makeSamplePairs([]float64{42}),
			maxPoints:     10,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sampler.Sample(tt.input, tt.maxPoints)
			assert.Equal(t, tt.expectedCount, len(result))
		})
	}
}

// TestMaxSampler 测试最大值采样器
// 功能：每个区间取最大值
func TestMaxSampler(t *testing.T) {
	sampler := &maxSampler{}

	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
		validate      func(*testing.T, []prommodel.SamplePair)
	}{
		{
			name:          "Basic max sampling",
			input:         makeSamplePairs([]float64{1, 5, 2, 8, 3, 9, 4, 7, 6}),
			maxPoints:     3,
			expectedCount: 3,
			validate: func(t *testing.T, result []prommodel.SamplePair) {
				// 验证结果中包含较大的值
				assert.True(t, result[0].Value >= 5)
			},
		},
		{
			name:          "All same values",
			input:         makeSamplePairs([]float64{5, 5, 5, 5, 5}),
			maxPoints:     2,
			expectedCount: 2,
		},
		{
			name:          "With negative values",
			input:         makeSamplePairs([]float64{-5, -1, -10, -3}),
			maxPoints:     2,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sampler.Sample(tt.input, tt.maxPoints)
			assert.Equal(t, tt.expectedCount, len(result))
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestMinSampler 测试最小值采样器
// 功能：每个区间取最小值
func TestMinSampler(t *testing.T) {
	sampler := &minSampler{}

	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
	}{
		{
			name:          "Basic min sampling",
			input:         makeSamplePairs([]float64{5, 1, 8, 2, 9, 3, 7, 4, 6}),
			maxPoints:     3,
			expectedCount: 3,
		},
		{
			name:          "All same values",
			input:         makeSamplePairs([]float64{5, 5, 5, 5, 5}),
			maxPoints:     2,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sampler.Sample(tt.input, tt.maxPoints)
			assert.Equal(t, tt.expectedCount, len(result))
		})
	}
}

// TestAvgSampler 测试平均值采样器
// 功能：每个区间取平均值
func TestAvgSampler(t *testing.T) {
	sampler := &avgSampler{}

	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
		validate      func(*testing.T, []prommodel.SamplePair)
	}{
		{
			name:          "Basic average sampling",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5, 6}),
			maxPoints:     2,
			expectedCount: 2,
			validate: func(t *testing.T, result []prommodel.SamplePair) {
				// 第一个区间[1,2,3]的平均值应该是2
				assert.InDelta(t, 2.0, float64(result[0].Value), 0.1)
			},
		},
		{
			name:          "With zeros",
			input:         makeSamplePairs([]float64{0, 1, 0, 2, 0, 3}),
			maxPoints:     2,
			expectedCount: 2,
		},
		{
			name:          "With negative values",
			input:         makeSamplePairs([]float64{-1, 1, -2, 2, -3, 3}),
			maxPoints:     2,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sampler.Sample(tt.input, tt.maxPoints)
			assert.Equal(t, tt.expectedCount, len(result))
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestLTTBSampler 测试LTTB采样器
// 功能：Largest Triangle Three Buckets算法
func TestLTTBSampler(t *testing.T) {
	sampler := &lttbSampler{}

	tests := []struct {
		name          string                 // 测试用例名称
		input         []prommodel.SamplePair // 输入数据
		maxPoints     int                    // 最大点数
		expectedCount int                    // 预期输出数量
		validate      func(*testing.T, []prommodel.SamplePair, []prommodel.SamplePair)
	}{
		{
			name:          "Basic LTTB sampling",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			maxPoints:     5,
			expectedCount: 5,
			validate: func(t *testing.T, original, result []prommodel.SamplePair) {
				// LTTB应该保留第一个和最后一个点
				assert.Equal(t, original[0].Timestamp, result[0].Timestamp)
				assert.Equal(t, original[len(original)-1].Timestamp, result[len(result)-1].Timestamp)
			},
		},
		{
			name:          "Large dataset",
			input:         makeSamplePairs(make([]float64, 1000)),
			maxPoints:     50,
			expectedCount: 50,
		},
		{
			name:          "Fewer points than maxPoints",
			input:         makeSamplePairs([]float64{1, 2, 3}),
			maxPoints:     10,
			expectedCount: 3,
		},
		{
			name:          "Empty input",
			input:         []prommodel.SamplePair{},
			maxPoints:     10,
			expectedCount: 0,
		},
		{
			name:          "maxPoints <= 2 should fallback to uniform",
			input:         makeSamplePairs([]float64{1, 2, 3, 4, 5}),
			maxPoints:     2,
			expectedCount: 2,
		},
		{
			name:          "Single point",
			input:         makeSamplePairs([]float64{42}),
			maxPoints:     10,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sampler.Sample(tt.input, tt.maxPoints)
			assert.Equal(t, tt.expectedCount, len(result))
			if tt.validate != nil {
				tt.validate(t, tt.input, result)
			}
		})
	}
}

// TestDataDownSample 测试Data的DownSample方法
// 功能：对Data中所有ResultPair进行降采样
func TestDataDownSample(t *testing.T) {
	tests := []struct {
		name          string // 测试用例名称
		input         *Data  // 输入数据
		maxPoints     int    // 最大点数
		expectedCount []int  // 预期每个ResultPair的输出数量
	}{
		{
			name: "Multiple ResultPairs",
			input: &Data{
				Result: []ResultPair{
					{Values: makeSamplePairs([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})},
					{Values: makeSamplePairs([]float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1})},
				},
			},
			maxPoints:     3,
			expectedCount: []int{3, 3},
		},
		{
			name: "Empty Data",
			input: &Data{
				Result: []ResultPair{},
			},
			maxPoints:     10,
			expectedCount: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.DownSample(tt.maxPoints)

			assert.Equal(t, len(tt.expectedCount), len(tt.input.Result))
			for i, expected := range tt.expectedCount {
				assert.Equal(t, expected, len(tt.input.Result[i].Values),
					"ResultPair %d should have %d points", i, expected)
			}
		})
	}
}

// TestDataDownSampleWithOptions 测试Data的DownSampleWithOptions方法
func TestDataDownSampleWithOptions(t *testing.T) {
	tests := []struct {
		name          string             // 测试用例名称
		input         *Data              // 输入数据
		opts          *DownSampleOptions // 降采样选项
		expectedCount []int              // 预期每个ResultPair的输出数量
	}{
		{
			name: "Multiple ResultPairs with LTTB",
			input: &Data{
				Result: []ResultPair{
					{Values: makeSamplePairs(make([]float64, 100))},
					{Values: makeSamplePairs(make([]float64, 100))},
				},
			},
			opts: &DownSampleOptions{
				Mode:      DownSampleModeLTTB,
				MaxPoints: 10,
				Enabled:   true,
			},
			expectedCount: []int{10, 10},
		},
		{
			name: "Disabled downsampling",
			input: &Data{
				Result: []ResultPair{
					{Values: makeSamplePairs([]float64{1, 2, 3, 4, 5})},
				},
			},
			opts: &DownSampleOptions{
				Mode:      DownSampleModeUniform,
				MaxPoints: 2,
				Enabled:   false, // 禁用
			},
			expectedCount: []int{5}, // 不进行降采样
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.DownSampleWithOptions(tt.opts)

			assert.Equal(t, len(tt.expectedCount), len(tt.input.Result))
			for i, expected := range tt.expectedCount {
				assert.Equal(t, expected, len(tt.input.Result[i].Values))
			}
		})
	}
}

// TestRowGetVal 测试Row的GetVal方法
// 功能：获取Row中的字段值
func TestRowGetVal(t *testing.T) {
	tests := []struct {
		name      string      // 测试用例名称
		row       Row         // 输入Row
		fieldName string      // 要获取的字段名
		expected  interface{} // 预期值
		wantErr   bool        // 是否预期错误
		errMsg    string      // 预期错误信息
	}{
		{
			name:      "Get existing string field",
			row:       Row{"name": "test", "value": 42},
			fieldName: "name",
			expected:  "test",
			wantErr:   false,
		},
		{
			name:      "Get existing int field",
			row:       Row{"name": "test", "value": 42},
			fieldName: "value",
			expected:  42,
			wantErr:   false,
		},
		{
			name:      "Get non-existent field",
			row:       Row{"name": "test"},
			fieldName: "nonexistent",
			expected:  nil,
			wantErr:   true,
			errMsg:    "no field",
		},
		{
			name:      "Get field from empty row",
			row:       Row{},
			fieldName: "any",
			expected:  nil,
			wantErr:   true,
			errMsg:    "no field",
		},
		{
			name:      "Get nil value field",
			row:       Row{"field": nil},
			fieldName: "field",
			expected:  nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.row.GetVal(tt.fieldName)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

// TestRowGetValStr 测试Row的GetValStr方法
// 功能：获取Row中的字符串类型字段值
func TestRowGetValStr(t *testing.T) {
	tests := []struct {
		name      string // 测试用例名称
		row       Row    // 输入Row
		fieldName string // 要获取的字段名
		expected  string // 预期值
		wantErr   bool   // 是否预期错误
		errMsg    string // 预期错误信息
	}{
		{
			name:      "Get existing string field",
			row:       Row{"name": "test", "value": "42"},
			fieldName: "name",
			expected:  "test",
			wantErr:   false,
		},
		{
			name:      "Get non-string field",
			row:       Row{"value": 42},
			fieldName: "value",
			expected:  "",
			wantErr:   true,
			errMsg:    "not of type string",
		},
		{
			name:      "Get non-existent field",
			row:       Row{"name": "test"},
			fieldName: "other",
			expected:  "",
			wantErr:   true,
			errMsg:    "no field",
		},
		{
			name:      "Get empty string field",
			row:       Row{"name": ""},
			fieldName: "name",
			expected:  "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.row.GetValStr(tt.fieldName)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

// TestRowGetValue 测试Row的GetValue方法
// 功能：获取Row中的Prometheus值
func TestRowGetValue(t *testing.T) {
	tests := []struct {
		name     string  // 测试用例名称
		row      Row     // 输入Row
		expected float64 // 预期值
		wantErr  bool    // 是否预期错误
		errMsg   string  // 预期错误信息
	}{
		{
			name:     "Get valid value",
			row:      Row{PROMVALUEKET: 123.45},
			expected: 123.45,
			wantErr:  false,
		},
		{
			name:     "Get zero value",
			row:      Row{PROMVALUEKET: 0.0},
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "Get negative value",
			row:      Row{PROMVALUEKET: -42.5},
			expected: -42.5,
			wantErr:  false,
		},
		{
			name:     "Get int value - should error (only float64 accepted)",
			row:      Row{PROMVALUEKET: 42},
			expected: 0,
			wantErr:  true,
			errMsg:   "not of type float64",
		},
		{
			name:     "Value field not found",
			row:      Row{"name": "test"},
			expected: 0,
			wantErr:  true,
			errMsg:   "no field value found",
		},
		{
			name:     "Value is not a number",
			row:      Row{PROMVALUEKET: "not a number"},
			expected: 0,
			wantErr:  true,
			errMsg:   "not of type float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.row.GetValue()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

// TestResPromQLString 测试ResPromQL的String方法
func TestResPromQLString(t *testing.T) {
	tests := []struct {
		name     string          // 测试用例名称
		input    prommodel.Value // 输入Prometheus值
		expected string          // 预期字符串输出
	}{
		{
			name:     "Vector value",
			input:    prommodel.Vector{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}
			result := rp.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// BenchmarkDownSamplerUniform 基准测试 - 均匀采样
func BenchmarkDownSamplerUniform(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &uniformSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sampler.Sample(input, 100)
	}
}

// BenchmarkDownSamplerMax 基准测试 - 最大值采样
func BenchmarkDownSamplerMax(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &maxSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sampler.Sample(input, 100)
	}
}

// BenchmarkDownSamplerMin 基准测试 - 最小值采样
func BenchmarkDownSamplerMin(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &minSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sampler.Sample(input, 100)
	}
}

// BenchmarkDownSamplerAvg 基准测试 - 平均值采样
func BenchmarkDownSamplerAvg(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &avgSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sampler.Sample(input, 100)
	}
}

// BenchmarkDownSamplerLTTB 基准测试 - LTTB采样
func BenchmarkDownSamplerLTTB(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &lttbSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sampler.Sample(input, 100)
	}
}

// BenchmarkDownSamplerParallel 并行基准测试
func BenchmarkDownSamplerParallel(b *testing.B) {
	input := makeSamplePairs(make([]float64, 10000))
	sampler := &uniformSampler{}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = sampler.Sample(input, 100)
		}
	})
}

// ExampleResultPair_DownSample 示例函数 - 基本降采样
func ExampleResultPair_DownSample() {
	// 创建一个包含1000个数据点的结果对
	values := make([]prommodel.SamplePair, 1000)
	for i := range values {
		values[i] = prommodel.SamplePair{
			Timestamp: prommodel.Time(i * 1000),
			Value:     prommodel.SampleValue(float64(i)),
		}
	}

	rp := &ResultPair{Values: values}

	// 降采样到100个点（使用默认的均匀采样模式）
	rp.DownSample(100)

	fmt.Printf("Downsampled from 1000 to %d points\n", len(rp.Values))
}

// ExampleResultPair_DownSampleWithOptions 示例函数 - 使用选项降采样
func ExampleResultPair_DownSampleWithOptions() {
	values := make([]prommodel.SamplePair, 1000)
	for i := range values {
		values[i] = prommodel.SamplePair{
			Timestamp: prommodel.Time(i * 1000),
			Value:     prommodel.SampleValue(float64(i)),
		}
	}

	rp := &ResultPair{Values: values}

	// 使用LTTB算法进行降采样（保留视觉特征）
	opts := &DownSampleOptions{
		Mode:      DownSampleModeLTTB,
		MaxPoints: 50,
		Enabled:   true,
	}
	rp.DownSampleWithOptions(opts)

	fmt.Printf("Downsampled using LTTB to %d points\n", len(rp.Values))
}

// ExampleData_DownSample 示例函数 - 批量降采样
func ExampleData_DownSample() {
	data := &Data{
		Result: []ResultPair{
			{Values: makeSamplePairs(make([]float64, 1000))},
			{Values: makeSamplePairs(make([]float64, 1000))},
		},
	}

	// 对所有时间序列进行降采样
	data.DownSample(100)

	fmt.Printf("Downsampled %d series to %d points each\n",
		len(data.Result), len(data.Result[0].Values))
}

// ExampleRow_GetVal 示例函数 - 获取字段值
func ExampleRow_GetVal() {
	row := Row{
		"name":  "metric_name",
		"value": 123.45,
		"tags":  "tag1,tag2",
	}

	// 获取字段值
	if name, err := row.GetVal("name"); err == nil {
		fmt.Printf("Name: %v\n", name)
	}

	if value, err := row.GetValue(); err == nil {
		fmt.Printf("Value: %f\n", value)
	}
}
