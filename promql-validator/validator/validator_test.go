// Copyright 2025 The Prometheus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidate_ValidQueries 测试有效的 PromQL 查询
func TestValidate_ValidQueries(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectedValid   bool
		expectedType    string
		expectedMetrics []string
		expectedFuncs   []string
	}{
		{
			name:          "简单指标名",
			query:         "http_requests_total",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{}, // 空切片而非 nil
		},
		{
			name:          "带标签选择器",
			query:         `http_requests_total{job="api",method="GET"}`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{}, // 空切片而非 nil
		},
		{
			name:          "范围查询",
			query:         "http_requests_total[5m]",
			expectedValid: true,
			expectedType:  "matrix",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{}, // 空切片而非 nil
		},
		{
			name:          "函数调用 - rate",
			query:         "rate(http_requests_total[5m])",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{"rate"},
		},
		{
			name:          "函数调用 - increase",
			query:         "increase(http_requests_total[5m])",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{"increase"},
		},
		{
			name:          "聚合操作 - sum by",
			query:         `sum by (job) (http_requests_total)`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{}, // sum 是 AggregateExpr，不是 Call
		},
		{
			name:          "聚合操作 - count without",
			query:         `count without (instance) (http_requests_total)`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{}, // count 是 AggregateExpr，不是 Call
		},
		{
			name:          "二元运算 - 乘法",
			query:         "http_requests_total * 1000",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "二元运算 - 加法",
			query:         "http_requests_total + http_errors_total",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
				"http_errors_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:            "标量",
			query:           "1.23",
			expectedValid:   true,
			expectedType:    "scalar",
			expectedMetrics: []string{},
			expectedFuncs:   []string{},
		},
		{
			name:            "负标量",
			query:           "-42",
			expectedValid:   true,
			expectedType:    "scalar",
			expectedMetrics: []string{},
			expectedFuncs:   []string{},
		},
		{
			name:          "布尔运算",
			query:         "http_requests_total > 100",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "多函数组合",
			query:         `sum(rate(http_requests_total[5m])) by (job)`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{"rate"}, // 只有 rate 是 Call，sum 是 AggregateExpr
		},
		{
			name:          "标签匹配 - 正则",
			query:         `http_requests_total{job=~"api.*"}`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "标签匹配 - 不等于",
			query:         `http_requests_total{job!="api"}`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "空 by 子句（允许）",
			query:         "sum by () (http_requests_total)",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "偏移量修饰符",
			query:         "http_requests_total offset 1h",
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
		{
			name:          "时间戳比较",
			query:         `http_requests_total and http_requests_total offset 1h`,
			expectedValid: true,
			expectedType:  "vector",
			expectedMetrics: []string{
				"http_requests_total",
			},
			expectedFuncs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.query)

			assert.Equal(t, tt.expectedValid, result.Valid, "验证状态应该匹配")
			assert.Equal(t, tt.expectedType, string(result.ExprType), "表达式类型应该匹配")
			assert.Equal(t, tt.expectedMetrics, result.Metrics, "指标列表应该匹配")
			assert.Equal(t, tt.expectedFuncs, result.Functions, "函数列表应该匹配")
			assert.Empty(t, result.Errors, "有效查询不应有错误")
		})
	}
}

// TestValidate_InvalidQueries 测试无效的 PromQL 查询
func TestValidate_InvalidQueries(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectedValid   bool
		expectErrors    bool
		expectedErrMsgs []string // 期望包含的错误消息子串
	}{
		{
			name:          "空字符串",
			query:         "",
			expectedValid: false,
			expectErrors:  true,
			expectedErrMsgs: []string{
				"empty query",
			},
		},
		{
			name:          "只有空格",
			query:         "   ",
			expectedValid: false,
			expectErrors:  true,
			expectedErrMsgs: []string{
				"empty query",
			},
		},
		{
			name:          "只有制表符和换行",
			query:         "\t\n  \r",
			expectedValid: false,
			expectErrors:  true,
			expectedErrMsgs: []string{
				"empty query",
			},
		},
		{
			name:          "语法错误 - 缺少右括号",
			query:         "sum(http_requests_total",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 多个缺少括号",
			query:         "sum(rate(http_requests_total[5m])",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 无效字符",
			query:         "@invalid",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 未闭合引号",
			query:         `http_requests_total{job="api}`,
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 无效标签名",
			query:         `http_requests_total{123bad="value"}`,
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 空范围",
			query:         "http_requests_total[]",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 无效范围单位",
			query:         "http_requests_total[5x]",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 不匹配的括号类型",
			query:         "sum(http_requests_total]",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 孤立的右括号",
			query:         "http_requests_total)",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 无效运算符",
			query:         "http_requests_total %% 1000",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 未知函数",
			query:         "unknown_func(http_requests_total)",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 未知函数",
			query:         "unknown_func(http_requests_total)",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 函数参数数量错误",
			query:         "rate(http_requests_total)",
			expectedValid: false,
			expectErrors:  true,
		},
		{
			name:          "语法错误 - 标签选择器格式错误",
			query:         `http_requests_total{job}`,
			expectedValid: false,
			expectErrors:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.query)

			assert.Equal(t, tt.expectedValid, result.Valid, "验证状态应该为无效")

			if tt.expectErrors {
				assert.NotEmpty(t, result.Errors, "应该有错误信息")
			}

			// 检查是否包含期望的错误消息
			if len(tt.expectedErrMsgs) > 0 {
				errorMsgs := make([]string, len(result.Errors))
				for i, err := range result.Errors {
					errorMsgs[i] = err.Message
				}
				for _, expectedMsg := range tt.expectedErrMsgs {
					found := false
					for _, errMsg := range errorMsgs {
						if contains(errMsg, expectedMsg) {
							found = true
							break
						}
					}
					assert.True(t, found, "错误信息应包含 '%s', 实际错误: %v", expectedMsg, errorMsgs)
				}
			}
		})
	}
}

// TestParseError_Error 测试 ParseError.Error() 方法
func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ParseError
		expected string
	}{
		{
			name:     "带位置信息的错误",
			err:      ParseError{Pos: 10, Message: "unexpected token"},
			expected: "position 10: unexpected token",
		},
		{
			name:     "不带位置信息的错误",
			err:      ParseError{Pos: 0, Message: "empty query"},
			expected: "empty query",
		},
		{
			name:     "负数位置（视为无位置）",
			err:      ParseError{Pos: -1, Message: "invalid syntax"},
			expected: "invalid syntax",
		},
		{
			name:     "空消息（自动填充 unknown error）",
			err:      ParseError{Pos: 5, Message: ""},
			expected: "position 5: unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Equal(t, tt.expected, result, "错误消息格式应该匹配")
		})
	}
}

// TestWarning_String 测试 Warning.String() 方法
func TestWarning_String(t *testing.T) {
	tests := []struct {
		name     string
		warning  Warning
		expected string
	}{
		{
			name:     "普通警告",
			warning:  Warning{Message: "deprecated function usage"},
			expected: "deprecated function usage",
		},
		{
			name:     "空警告消息",
			warning:  Warning{Message: ""},
			expected: "",
		},
		{
			name:     "带特殊字符的警告",
			warning:  Warning{Message: "warning: label `job` is high cardinality"},
			expected: "warning: label `job` is high cardinality",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.warning.String()
			assert.Equal(t, tt.expected, result, "警告消息应该匹配")
		})
	}
}

// TestValidationResult_String 测试 ValidationResult.String() 方法
func TestValidationResult_String(t *testing.T) {
	tests := []struct {
		name        string
		result      *ValidationResult
		contains    []string // 期望包含的字符串
		notContains []string // 期望不包含的字符串
	}{
		{
			name: "有效查询 - 简单指标",
			result: &ValidationResult{
				Valid:     true,
				ExprType:  "vector",
				Metrics:   []string{"http_requests_total"},
				Functions: []string{},
				Errors:    []ParseError{},
				Warnings:  []Warning{},
			},
			contains: []string{
				"✅ Query is VALID",
				"Expression type: vector",
				"Metrics: [http_requests_total]",
			},
			notContains: []string{
				"❌ Query is INVALID",
				"Errors:",
				"Warnings:",
			},
		},
		{
			name: "有效查询 - 带函数",
			result: &ValidationResult{
				Valid:     true,
				ExprType:  "vector",
				Metrics:   []string{"http_requests_total"},
				Functions: []string{"rate", "sum"},
				Errors:    []ParseError{},
				Warnings:  []Warning{},
			},
			contains: []string{
				"✅ Query is VALID",
				"Functions: [rate sum]",
			},
			notContains: []string{
				"❌ Query is INVALID",
			},
		},
		{
			name: "无效查询 - 有错误",
			result: &ValidationResult{
				Valid:     false,
				ExprType:  "",
				Metrics:   []string{},
				Functions: []string{},
				Errors: []ParseError{
					{Pos: 5, Message: "unexpected token"},
				},
				Warnings: []Warning{},
			},
			contains: []string{
				"❌ Query is INVALID",
				"Errors:",
				"position 5: unexpected token",
			},
			notContains: []string{
				"✅ Query is VALID",
			},
		},
		{
			name: "有警告的查询",
			result: &ValidationResult{
				Valid:     true,
				ExprType:  "vector",
				Metrics:   []string{"http_requests_total"},
				Functions: []string{},
				Errors:    []ParseError{},
				Warnings: []Warning{
					{Message: "high cardinality label"},
					{Message: "expensive query"},
				},
			},
			contains: []string{
				"✅ Query is VALID",
				"Warnings:",
				"- high cardinality label",
				"- expensive query",
			},
		},
		{
			name: "空指标和函数列表",
			result: &ValidationResult{
				Valid:     true,
				ExprType:  "scalar",
				Metrics:   []string{},
				Functions: []string{},
				Errors:    []ParseError{},
				Warnings:  []Warning{},
			},
			contains: []string{
				"✅ Query is VALID",
				"Expression type: scalar",
			},
			notContains: []string{
				"Metrics:",
				"Functions:",
			},
		},
		{
			name: "多个错误",
			result: &ValidationResult{
				Valid:     false,
				ExprType:  "",
				Metrics:   []string{},
				Functions: []string{},
				Errors: []ParseError{
					{Pos: 0, Message: "empty query"},
					{Pos: 10, Message: "syntax error"},
				},
				Warnings: []Warning{},
			},
			contains: []string{
				"❌ Query is INVALID",
				"- empty query",
				"- position 10: syntax error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result.String()

			for _, expected := range tt.contains {
				assert.Contains(t, result, expected, "结果应包含: %s", expected)
			}

			for _, notExpected := range tt.notContains {
				assert.NotContains(t, result, notExpected, "结果不应包含: %s", notExpected)
			}
		})
	}
}

// TestMetadataCollection 测试元数据收集
func TestMetadataCollection(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectedMetrics []string
		expectedFuncs   []string
	}{
		{
			name:            "单个指标",
			query:           "http_requests_total",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{},
		},
		{
			name:            "多个相同指标 - 应去重",
			query:           "http_requests_total + http_requests_total",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{},
		},
		{
			name:            "多个不同指标",
			query:           "http_requests_total + http_errors_total",
			expectedMetrics: []string{"http_requests_total", "http_errors_total"},
			expectedFuncs:   []string{},
		},
		{
			name:            "单个函数",
			query:           "rate(http_requests_total[5m])",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{"rate"},
		},
		{
			name:            "多个相同函数 - 应去重",
			query:           "sum(rate(http_requests_total[5m])) + sum(increase(http_requests_total[5m]))",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{"rate", "increase"}, // sum 是 AggregateExpr，不是 Call
		},
		{
			name:            "VectorSelector 收集",
			query:           "http_requests_total{job=\"api\"}",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{},
		},
		{
			name:            "MatrixSelector 收集",
			query:           "http_requests_total[5m]",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{},
		},
		{
			name:            "嵌套表达式中的指标",
			query:           "sum(rate(http_requests_total[5m])) + sum(rate(http_errors_total[5m]))",
			expectedMetrics: []string{"http_requests_total", "http_errors_total"},
			expectedFuncs:   []string{"rate"}, // sum 是 AggregateExpr，不是 Call
		},
		{
			name:            "复杂聚合",
			query:           "sum by (job) (rate(http_requests_total{code=~\"5..\"}[5m]))",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{"rate"}, // sum 是 AggregateExpr，不是 Call
		},
		{
			name:            "多种聚合组合",
			query:           "sum(http_requests_total) + count(http_requests_total) + avg(http_requests_total)",
			expectedMetrics: []string{"http_requests_total"},
			expectedFuncs:   []string{}, // sum/count/avg 都是 AggregateExpr，不是 Call
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.query)

			assert.True(t, result.Valid, "查询应该有效: %s", tt.query)
			assert.ElementsMatch(t, tt.expectedMetrics, result.Metrics, "指标列表应该匹配")
			assert.ElementsMatch(t, tt.expectedFuncs, result.Functions, "函数列表应该匹配")
		})
	}
}

// TestExpressionTypes 测试各种表达式类型
func TestExpressionTypes(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedType  string
		expectedValid bool
	}{
		{
			name:          "Instant Vector - 简单",
			query:         "http_requests_total",
			expectedType:  "vector",
			expectedValid: true,
		},
		{
			name:          "Instant Vector - 带标签",
			query:         `http_requests_total{job="api"}`,
			expectedType:  "vector",
			expectedValid: true,
		},
		{
			name:          "Range Vector",
			query:         "http_requests_total[5m]",
			expectedType:  "matrix",
			expectedValid: true,
		},
		{
			name:          "Scalar",
			query:         "1.23",
			expectedType:  "scalar",
			expectedValid: true,
		},
		{
			name:          "函数返回 Vector",
			query:         "rate(http_requests_total[5m])",
			expectedType:  "vector",
			expectedValid: true,
		},
		{
			name:          "聚合返回 Vector",
			query:         "sum(http_requests_total)",
			expectedType:  "vector",
			expectedValid: true,
		},
		{
			name:          "二元运算返回 Vector",
			query:         "http_requests_total * 1000",
			expectedType:  "vector",
			expectedValid: true,
		},
		{
			name:          "比较返回 Vector",
			query:         "http_requests_total > 100",
			expectedType:  "vector",
			expectedValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.query)

			assert.Equal(t, tt.expectedValid, result.Valid, "验证状态应该匹配")
			assert.Equal(t, tt.expectedType, string(result.ExprType), "表达式类型应该匹配")
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedValid bool
		description   string
	}{
		{
			name:          "超长有效指标名",
			query:         strings.Repeat("a", 10000),
			expectedValid: true, // 长字符串实际上是有效的指标名
			description:   "长指标名应该有效",
		},
		{
			name:          "只有运算符",
			query:         "+",
			expectedValid: false,
			description:   "只有运算符应该无效",
		},
		{
			name:          "嵌套括号",
			query:         "sum(sum(sum(http_requests_total)))",
			expectedValid: true,
			description:   "多层嵌套应该有效",
		},
		{
			name:          "多个连续运算符",
			query:         "http_requests_total * * 1000",
			expectedValid: false,
			description:   "连续运算符应该无效",
		},
		{
			name:          "空标签值",
			query:         `http_requests_total{job=""}`,
			expectedValid: true,
			description:   "空标签值应该有效",
		},
		{
			name:          "多个标签选择器",
			query:         `http_requests_total{job="api",method="GET",code="200"}`,
			expectedValid: true,
			description:   "多个标签应该有效",
		},
		{
			name:          "正则标签匹配",
			query:         `http_requests_total{job=~"api.*",method=~"GET|POST"}`,
			expectedValid: true,
			description:   "正则匹配应该有效",
		},
		{
			name:          "负时间范围",
			query:         "http_requests_total[5m] offset -5m",
			expectedValid: true,
			description:   "负偏移量应该有效",
		},
		{
			name:          "非常小的时间范围",
			query:         "http_requests_total[1s]",
			expectedValid: true,
			description:   "1秒范围应该有效",
		},
		{
			name:          "大时间范围",
			query:         "http_requests_total[365d]",
			expectedValid: true,
			description:   "大范围应该有效",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.query)
			assert.Equal(t, tt.expectedValid, result.Valid, tt.description)
		})
	}
}

// TestResultConsistency 测试结果一致性
func TestResultConsistency(t *testing.T) {
	t.Run("相同查询应该产生相同结果", func(t *testing.T) {
		query := "sum(rate(http_requests_total[5m]))"

		result1 := Validate(query)
		result2 := Validate(query)

		assert.Equal(t, result1.Valid, result2.Valid, "验证状态应该一致")
		assert.Equal(t, result1.ExprType, result2.ExprType, "表达式类型应该一致")
		assert.Equal(t, result1.Metrics, result2.Metrics, "指标列表应该一致")
		assert.Equal(t, result1.Functions, result2.Functions, "函数列表应该一致")
	})

	t.Run("不同查询应该产生不同结果", func(t *testing.T) {
		query1 := "http_requests_total"
		query2 := "http_errors_total"

		result1 := Validate(query1)
		result2 := Validate(query2)

		assert.NotEqual(t, result1.Metrics, result2.Metrics, "不同指标的列表应该不同")
	})
}

// TestEmptyResultFields 测试空结果字段
func TestEmptyResultFields(t *testing.T) {
	t.Run("无效查询应该初始化所有字段", func(t *testing.T) {
		result := Validate("")

		require.NotNil(t, result, "结果不应为 nil")
		assert.False(t, result.Valid, "应该无效")
		assert.NotNil(t, result.Errors, "错误列表应初始化")
		assert.NotNil(t, result.Warnings, "警告列表应初始化")
		assert.NotNil(t, result.Metrics, "指标列表应初始化")
		assert.NotNil(t, result.Functions, "函数列表应初始化")
	})

	t.Run("有效查询不应有错误", func(t *testing.T) {
		result := Validate("http_requests_total")

		assert.True(t, result.Valid, "应该有效")
		assert.Empty(t, result.Errors, "不应有错误")
		assert.Empty(t, result.Warnings, "不应有警告")
	})
}

// 辅助函数
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
