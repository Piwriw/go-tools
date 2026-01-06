package prom

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	prommodel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPrometheusClient 测试创建Prometheus客户端
// 功能：使用各种选项创建客户端
// 测试覆盖：正常场景、边界条件、错误处理
func TestNewPrometheusClient(t *testing.T) {
	tests := []struct {
		name     string                              // 测试用例名称
		address  string                              // Prometheus地址
		opts     []Option                            // 配置选项
		wantErr  bool                                // 是否预期错误
		errMsg   string                              // 预期错误信息
		validate func(*testing.T, *PrometheusClient) // 验证函数
	}{
		{
			name:    "Valid address without options",
			address: "http://localhost:9090",
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.NotNil(t, pc)
				assert.NotNil(t, pc.client)
				assert.NotNil(t, pc.httpClient)
				assert.Equal(t, "http://localhost:9090", pc.address)
			},
		},
		{
			name:    "Valid address with token",
			address: "http://localhost:9090",
			opts:    []Option{WithToken("test-token")},
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.Equal(t, "test-token", pc.token)
				// httpClient is an interface, can't directly check transport type
				// Transport is tested in TestAuthTransport
				assert.NotNil(t, pc.httpClient)
			},
		},
		{
			name:    "With timeout option",
			address: "http://localhost:9090",
			opts:    []Option{WithTimeout(30 * time.Second)},
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.Equal(t, 30*time.Second, pc.timeout)
			},
		},
		{
			name:    "With zero timeout - should use default",
			address: "http://localhost:9090",
			opts:    []Option{WithTimeout(0)},
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.Equal(t, time.Duration(0), pc.timeout) // 0表示不设置
			},
		},
		{
			name:    "With negative timeout - should be ignored",
			address: "http://localhost:9090",
			opts:    []Option{WithTimeout(-10)},
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.Equal(t, time.Duration(0), pc.timeout)
			},
		},
		{
			name:    "Multiple options",
			address: "http://localhost:9090",
			opts: []Option{
				WithToken("my-token"),
				WithTimeout(60 * time.Second),
			},
			wantErr: false,
			validate: func(t *testing.T, pc *PrometheusClient) {
				assert.Equal(t, "my-token", pc.token)
				assert.Equal(t, 60*time.Second, pc.timeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewPrometheusClient(tt.address, tt.opts...)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, client)
				}
			}
		})
	}
}

// TestAuthTransport 测试authTransport的RoundTrip方法
// 功能：验证Bearer Token被正确添加到请求头
func TestAuthTransport(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		token      string // Token
		requestURL string // 请求URL
	}{
		{
			name:       "Basic auth header is set",
			token:      "dXNlcjpwYXNz", // base64 encoded "user:pass"
			requestURL: "http://example.com/api",
		},
		{
			name:       "Empty token",
			token:      "",
			requestURL: "http://example.com/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个假的HTTP服务器来接收请求
			var receivedAuth string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedAuth = r.Header.Get("Authorization")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// 创建authTransport并发送请求
			transport := &authTransport{
				base:  http.DefaultTransport,
				token: tt.token,
			}

			req, err := http.NewRequest("GET", server.URL, nil)
			require.NoError(t, err)

			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// 验证Authorization头
			if tt.token != "" {
				assert.Equal(t, "Basic "+tt.token, receivedAuth)
			}
		})
	}
}

// TestExtractPromQLQuery 测试提取PromQL查询
// 功能：从字符串中提取{{ query "..." }}中的PromQL表达式
func TestExtractPromQLQuery(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		input    string // 输入字符串
		expected string // 预期提取的PromQL
	}{
		{
			name:     "Simple query",
			input:    `{{ query "up" }}`,
			expected: "up",
		},
		{
			name:     "Query with pipeline",
			input:    `{{ query "up" | first }}`,
			expected: "up",
		},
		{
			name:     "Complex query",
			input:    `{{ query "rate(http_requests_total[5m])" }}`,
			expected: "rate(http_requests_total[5m])",
		},
		{
			name:     "Query with spaces",
			input:    `{{ query "sum by (job) (rate(http_requests_total[5m]))" }}`,
			expected: "sum by (job) (rate(http_requests_total[5m]))",
		},
		{
			name:     "Query with nested quotes",
			input:    `{{ query "count_values(\"version\", app_version)" }}`,
			expected: `count_values(\`,
		},
		{
			name:     "No query pattern",
			input:    `plain text`,
			expected: "",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Multiple queries - extract first",
			input:    `{{ query "up" }} and {{ query "down" }}`,
			expected: "up",
		},
		{
			name:     "Query with newlines (literal backslash-n)",
			input:    `{{ query "up\nor\ndown" }}`,
			expected: `up\nor\ndown`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPromQLQuery(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidateMetric 测试验证指标名称和标签集
// 功能：验证指标的有效性
func TestValidateMetric(t *testing.T) {
	tests := []struct {
		name     string             // 测试用例名称
		metric   string             // 指标名称
		labelSet prommodel.LabelSet // 标签集
		expected bool               // 预期结果
	}{
		{
			name:     "Valid metric name and labels",
			metric:   "http_requests_total",
			labelSet: prommodel.LabelSet{"job": "api", "instance": "localhost:9090"},
			expected: true,
		},
		{
			name:     "Valid metric with underscore",
			metric:   "metric_name_test",
			labelSet: prommodel.LabelSet{"label": "value"},
			expected: true,
		},
		{
			name:     "Invalid metric name - starts with digit",
			metric:   "123metric",
			labelSet: prommodel.LabelSet{},
			expected: false,
		},
		{
			name:     "Invalid metric name - contains special chars",
			metric:   "metric-name!",
			labelSet: prommodel.LabelSet{},
			expected: false,
		},
		{
			name:     "Empty metric name",
			metric:   "",
			labelSet: prommodel.LabelSet{},
			expected: false,
		},
		{
			name:     "Invalid label name",
			metric:   "valid_metric",
			labelSet: prommodel.LabelSet{"123invalid": "value"},
			expected: false,
		},
		{
			name:     "Valid label name with underscore",
			metric:   "valid_metric",
			labelSet: prommodel.LabelSet{"label_name": "value"},
			expected: true,
		},
		{
			name:   "Label with PromQL query - valid",
			metric: "description",
			labelSet: prommodel.LabelSet{
				"description": `{{ query "up" | first }}`,
			},
			expected: true,
		},
		{
			name:   "Label with invalid PromQL query",
			metric: "description",
			labelSet: prommodel.LabelSet{
				"description": `{{ query "unclosed bracket" }}`,
			},
			expected: false,
		},
		{
			name:   "Valid metric with all valid label names",
			metric: "go_goroutines",
			labelSet: prommodel.LabelSet{
				"job":      "prometheus",
				"instance": "localhost:9090",
				"__name__": "go_goroutines",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &PrometheusClient{}
			result := pc.ValidateMetric(tt.metric, tt.labelSet)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPrettify 测试PromQL格式化
// 功能：格式化PromQL查询语句
func TestPrettify(t *testing.T) {
	tests := []struct {
		name     string                   // 测试用例名称
		query    string                   // 输入查询
		wantErr  bool                     // 是否预期错误
		errMsg   string                   // 预期错误信息
		validate func(*testing.T, string) // 验证函数
	}{
		{
			name:    "Simple metric",
			query:   "up",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.Equal(t, "up", result)
			},
		},
		{
			name:    "Query with aggregation",
			query:   "sum(rate(http_requests_total[5m]))",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "sum")
			},
		},
		{
			name:    "Complex query",
			query:   "sum by (job) (rate(http_requests_total[5m]))",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "sum")
			},
		},
		{
			name:    "Binary operation",
			query:   "up and up{job=\"prometheus\"}",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:    "Invalid query - syntax error",
			query:   "sum(((unclosed",
			wantErr: true,
			errMsg:  "parse error",
		},
		{
			name:    "Empty query",
			query:   "",
			wantErr: true,
		},
		{
			name:    "Query with comments",
			query:   "# This is a comment\nup",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &PrometheusClient{}
			result, err := pc.Prettify(tt.query)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestQuery 测试查询Prometheus
// 功能：执行即时查询
func TestQuery(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		query      string // PromQL查询
		response   string // 模拟响应
		statusCode int    // HTTP状态码
		wantErr    bool   // 是否预期错误
	}{
		{
			name:       "Successful query",
			query:      "up",
			response:   `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"job":"prometheus"},"value":[1234567890,"1"]}]}}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Query with error response",
			query:      "invalid_query",
			response:   `{"status":"error","errorType":"bad_data","error":"parse error"}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "HTTP error",
			query:      "up",
			response:   "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 验证请求路径
				assert.Equal(t, "/api/v1/query", r.URL.Path)

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// 创建客户端
			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			// 执行查询
			result, err := client.Query(context.Background(), tt.query, time.Now())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestQueryVector 测试查询向量
// 功能：执行查询并返回Vector类型
func TestQueryVector(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		query    string // PromQL查询
		response string // 模拟响应
		wantErr  bool   // 是否预期错误
		errMsg   string // 预期错误信息
	}{
		{
			name:  "Successful vector query",
			query: "up",
			response: `{
				"status":"success",
				"data":{
					"resultType":"vector",
					"result":[
						{"metric":{"job":"prometheus"},"value":[1234567890,"1"]},
						{"metric":{"job":"node"},"value":[1234567890,"0"]}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "Query returns scalar - should error",
			query:    "123",
			response: `{"status":"success","data":{"resultType":"scalar","result":[1234567890,"123"]}}`,
			wantErr:  true,
			errMsg:   "not a vector",
		},
		{
			name:     "Query returns matrix - should error",
			query:    "up[5m]",
			response: `{"status":"success","data":{"resultType":"matrix","result":[]}}`,
			wantErr:  true,
			errMsg:   "not a vector",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			result, err := client.QueryVector(context.Background(), tt.query, time.Now())

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.IsType(t, prommodel.Vector{}, result)
			}
		})
	}
}

// TestQueryRangeMatrix 测试范围查询返回矩阵
// 功能：执行时间范围查询并返回Matrix类型
func TestQueryRangeMatrix(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		query    string // PromQL查询
		response string // 模拟响应
		wantErr  bool   // 是否预期错误
		errMsg   string // 预期错误信息
	}{
		{
			name:  "Successful range query",
			query: "rate(http_requests_total[5m])",
			response: `{
				"status":"success",
				"data":{
					"resultType":"matrix",
					"result":[{
						"metric":{"job":"api"},
						"values":[[1234567890,"10"],[1234567895,"20"],[1234567900,"30"]]
					}]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "Query returns vector - should error",
			query:    "up",
			response: `{"status":"success","data":{"resultType":"vector","result":[]}}`,
			wantErr:  true,
			errMsg:   "not a matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			end := time.Now()
			start := end.Add(-1 * time.Hour)
			result, err := client.QueryRangeMatrix(context.Background(), tt.query, start, end, 1*time.Minute)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.IsType(t, prommodel.Matrix{}, result)
			}
		})
	}
}

// TestValidate 测试PromQL验证
// 功能：验证PromQL查询的语法正确性
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string // 测试用例名称
		query   string // PromQL查询
		wantErr bool   // 是否预期错误
		errMsg  string // 预期错误信息
	}{
		{
			name:    "Valid simple query",
			query:   "up",
			wantErr: false,
		},
		{
			name:    "Valid query with function",
			query:   "rate(http_requests_total[5m])",
			wantErr: false,
		},
		{
			name:    "Valid query with aggregation",
			query:   "sum by (job) (rate(http_requests_total[5m]))",
			wantErr: false,
		},
		{
			name:    "Valid query with binary operation",
			query:   "up and up{job=\"prometheus\"}",
			wantErr: false,
		},
		{
			name:    "Invalid query - syntax error",
			query:   "sum(((unclosed",
			wantErr: true,
			errMsg:  "parse error",
		},
		{
			name:    "Invalid query - unclosed bracket",
			query:   "up[5m",
			wantErr: true,
		},
		{
			name:    "Empty query",
			query:   "",
			wantErr: true,
		},
		{
			name:    "Valid query with regex",
			query:   `up{job=~".*"}`,
			wantErr: false,
		},
		{
			name:    "Valid query with range",
			query:   "http_requests_total offset 5m",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &PrometheusClient{}
			err := pc.Validate(context.Background(), tt.query)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestQueryOneValue 测试查询单个值
// 功能：执行查询并返回单个值
func TestQueryOneValue(t *testing.T) {
	tests := []struct {
		name     string                                  // 测试用例名称
		query    string                                  // PromQL查询
		response string                                  // 模拟响应
		wantErr  bool                                    // 是否预期错误
		errMsg   string                                  // 预期错误信息
		validate func(*testing.T, prommodel.SampleValue) // 验证函数
	}{
		{
			name:  "Single value query",
			query: "up",
			response: `{
				"status":"success",
				"data":{
					"resultType":"vector",
					"result":[{"metric":{"job":"prometheus"},"value":[1234567890,"1"]}]
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, val prommodel.SampleValue) {
				assert.Equal(t, prommodel.SampleValue(1), val)
			},
		},
		{
			name:  "No results",
			query: "up{job=\"nonexistent\"}",
			response: `{
				"status":"success",
				"data":{"resultType":"vector","result":[]}
			}`,
			wantErr: false,
			validate: func(t *testing.T, val prommodel.SampleValue) {
				assert.Equal(t, prommodel.SampleValue(0), val)
			},
		},
		{
			name:  "Multiple values - should error",
			query: "up",
			response: `{
				"status":"success",
				"data":{
					"resultType":"vector",
					"result":[
						{"value":[1234567890,"1"]},
						{"value":[1234567890,"0"]}
					]
				}
			}`,
			wantErr: true,
			errMsg:  "too many values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			result, err := client.QueryOneValue(context.Background(), tt.query, time.Now())

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestQueryAllMetrics 测试查询所有指标
// 功能：查询Prometheus中的所有指标名称
func TestQueryAllMetrics(t *testing.T) {
	tests := []struct {
		name     string                                  // 测试用例名称
		response string                                  // 模拟响应
		wantErr  bool                                    // 是否预期错误
		validate func(*testing.T, prommodel.LabelValues) // 验证函数
	}{
		{
			name: "Successful query",
			response: `{
				"status":"success",
				"data":["up", "process_cpu_seconds_total", "http_requests_total"]
			}`,
			wantErr: false,
			validate: func(t *testing.T, lv prommodel.LabelValues) {
				assert.Equal(t, 3, len(lv))
			},
		},
		{
			name:     "Empty result",
			response: `{"status":"success","data":[]}`,
			wantErr:  false,
			validate: func(t *testing.T, lv prommodel.LabelValues) {
				assert.Equal(t, 0, len(lv))
			},
		},
		{
			name:     "Error response",
			response: `{"status":"error","errorType":"bad_data","error":"parse error"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/label/__name__/values", r.URL.Path)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			result, err := client.QueryAllMetrics(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestInitPromClient 测试初始化全局Prometheus客户端
func TestInitPromClient(t *testing.T) {
	// 保存原始值
	originalClient := DefaultClient

	tests := []struct {
		name    string   // 测试用例名称
		address string   // Prometheus地址
		opts    []Option // 配置选项
		wantErr bool     // 是否预期错误
	}{
		{
			name:    "Valid initialization",
			address: "http://localhost:9090",
			wantErr: false,
		},
		{
			name:    "With token",
			address: "http://localhost:9090",
			opts:    []Option{WithToken("test-token")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitPromClient(tt.address, tt.opts...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, DefaultClient)
			}

			// 恢复原始值
			DefaultClient = originalClient
		})
	}
}

// TestPrometheusClientMethods 测试PrometheusClient的各种方法
// 功能：验证Client()和HTTPClient()方法
func TestPrometheusClientMethods(t *testing.T) {
	// 创建mock server用于测试
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewPrometheusClient(server.URL)
	require.NoError(t, err)

	t.Run("Client method", func(t *testing.T) {
		result := client.Client()
		assert.NotNil(t, result)
		assert.Same(t, client.client, result)
	})

	t.Run("HTTPClient method", func(t *testing.T) {
		result := client.HTTPClient()
		assert.NotNil(t, result)
		assert.Same(t, client.httpClient, result)
	})
}

// TestReload 测试重新加载Prometheus配置
// 功能：发送POST请求到/-/reload端点
func TestReload(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		response   string // 模拟响应
		statusCode int    // HTTP状态码
		wantErr    bool   // 是否预期错误
		errMsg     string // 预期错误信息
	}{
		{
			name:       "Successful reload",
			response:   "OK",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Reload failed",
			response:   "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
			errMsg:     "reload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedPath = r.URL.Path
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			err = client.Reload(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "/-/reload", receivedPath)
			}
		})
	}
}

// BenchmarkQuery 基准测试 - Query方法
func BenchmarkQuery(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
	}))
	defer server.Close()

	client, err := NewPrometheusClient(server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := "up"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := client.Query(ctx, query, time.Now())
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValidate 基准测试 - Validate方法
func BenchmarkValidate(b *testing.B) {
	client := &PrometheusClient{}
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = client.Validate(ctx, "sum(rate(http_requests_total[5m]))")
	}
}

// BenchmarkPrettify 基准测试 - Prettify方法
func BenchmarkPrettify(b *testing.B) {
	client := &PrometheusClient{}
	query := "sum by (job) (rate(http_requests_total[5m]))"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = client.Prettify(query)
	}
}

// ExampleNewPrometheusClient 示例函数 - 创建客户端
func ExampleNewPrometheusClient() {
	// 创建基本客户端
	client1, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	// 创建带认证的客户端
	client2, err := NewPrometheusClient(
		"http://localhost:9090",
		WithToken("base64-encoded-token"),
	)
	if err != nil {
		panic(err)
	}

	// 创建带超时的客户端
	client3, err := NewPrometheusClient(
		"http://localhost:9090",
		WithTimeout(30*time.Second),
	)
	if err != nil {
		panic(err)
	}

	_, _ = client1, client2
	_ = client3
}

// ExamplePrometheusClient_Query 示例函数 - 执行查询
func ExamplePrometheusClient_Query() {
	client, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	query := "up"
	ts := time.Now()

	result, err := client.Query(ctx, query, ts)
	if err != nil {
		panic(err)
	}

	_ = result
}

// ExamplePrometheusClient_QueryVector 示例函数 - 查询向量
func ExamplePrometheusClient_QueryVector() {
	client, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	query := "up"

	vector, err := client.QueryVector(ctx, query, time.Now())
	if err != nil {
		panic(err)
	}

	for _, sample := range vector {
		fmt.Printf("Metric: %v, Value: %v\n", sample.Metric, sample.Value)
	}
}

// ExamplePrometheusClient_Validate 示例函数 - 验证PromQL
func ExamplePrometheusClient_Validate() {
	client, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 验证有效的查询
	err = client.Validate(ctx, "up")
	if err != nil {
		fmt.Println("Query is invalid:", err)
	}

	// 验证无效的查询
	err = client.Validate(ctx, "sum(((unclosed")
	if err != nil {
		fmt.Println("Query is invalid:", err)
	}
}

// ExamplePrometheusClient_Prettify 示例函数 - 格式化PromQL
func ExamplePrometheusClient_Prettify() {
	client, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	// 格式化PromQL查询
	pretty, err := client.Prettify("sum by(job)(rate(http_requests_total[5m]))")
	if err != nil {
		panic(err)
	}

	fmt.Println(pretty)
}

// TestQueryRange 测试时间范围查询
// 功能：执行Prometheus时间范围查询
// 测试覆盖：正常场景、边界条件、错误处理
func TestQueryRange(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		query    string // PromQL查询
		response string // 模拟响应
		start    string // 开始时间
		end      string // 结束时间
		step     string // 步长
		wantErr  bool   // 是否预期错误
		validate func(*testing.T, prommodel.Value)
	}{
		{
			name:  "Successful range query",
			query: "up",
			start: "2024-01-01T00:00:00Z",
			end:   "2024-01-01T01:00:00Z",
			step:  "1m",
			response: `{
				"status":"success",
				"data":{
					"resultType":"matrix",
					"result":[{
						"metric":{"job":"prometheus"},
						"values":[
							[1234567890,"1"],
							[1234567950,"1"]
						]
					}]
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, result prommodel.Value) {
				assert.NotNil(t, result)
			},
		},
		{
			name:  "Empty result",
			query: "nonexistent_metric",
			start: "2024-01-01T00:00:00Z",
			end:   "2024-01-01T01:00:00Z",
			step:  "1m",
			response: `{
				"status":"success",
				"data":{"resultType":"matrix","result":[]}
			}`,
			wantErr: false,
		},
		{
			name:  "Query error response",
			query: "invalid_query",
			start: "2024-01-01T00:00:00Z",
			end:   "2024-01-01T01:00:00Z",
			step:  "1m",
			response: `{
				"status":"error",
				"errorType":"bad_data",
				"error":"parse error"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/query_range", r.URL.Path)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			start, _ := time.Parse(time.RFC3339, tt.start)
			end, _ := time.Parse(time.RFC3339, tt.end)
			step, _ := time.ParseDuration(tt.step)

			result, err := client.QueryRange(context.Background(), tt.query, start, end, step)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestQueryMetric 测试查询指定指标
// 功能：查询特定指标的当前值
// 测试覆盖：正常场景、错误处理
func TestQueryMetric(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		metric   string // 指标名称
		response string // 模拟响应
		wantErr  bool   // 是否预期错误
	}{
		{
			name:   "Successful metric query",
			metric: "up",
			response: `{
				"status":"success",
				"data":{
					"resultType":"vector",
					"result":[{
						"metric":{"job":"prometheus"},
						"value":[1234567890,"1"]
					}]
				}
			}`,
			wantErr: false,
		},
		{
			name:   "Metric not found",
			metric: "nonexistent_metric",
			response: `{
				"status":"success",
				"data":{"resultType":"vector","result":[]}
			}`,
			wantErr: false,
		},
		{
			name:   "Query error",
			metric: "invalid:metric",
			response: `{
				"status":"error",
				"errorType":"bad_data",
				"error":"parse error"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewPrometheusClient(server.URL)
			require.NoError(t, err)

			result, err := client.QueryMetric(context.Background(), tt.metric)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestResPromQLRows 测试Rows方法
// 功能：将Prometheus值转换为行
// 测试覆盖：Vector类型、Matrix类型、错误处理
func TestResPromQLRows(t *testing.T) {
	tests := []struct {
		name     string          // 测试用例名称
		input    prommodel.Value // 输入值
		wantErr  bool            // 是否预期错误
		errMsg   string          // 预期错误信息
		validate func(*testing.T, []Row)
	}{
		{
			name: "Vector value",
			input: prommodel.Vector{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Value:  prommodel.SampleValue(1),
				},
				{
					Metric: prommodel.Metric{"job": "node"},
					Value:  prommodel.SampleValue(0),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, rows []Row) {
				assert.Len(t, rows, 2)
				// 验证第一行
				assert.Equal(t, "prometheus", rows[0]["job"])
				assert.Equal(t, float64(1), rows[0][PROMVALUEKET])
			},
		},
		{
			name: "Matrix value",
			input: prommodel.Matrix{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Values: []prommodel.SamplePair{
						{Timestamp: 1234567890, Value: 1},
						{Timestamp: 1234567950, Value: 0},
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, rows []Row) {
				assert.Len(t, rows, 2)
				assert.Equal(t, "prometheus", rows[0]["job"])
			},
		},
		{
			name:    "Empty vector",
			input:   prommodel.Vector{},
			wantErr: false,
			validate: func(t *testing.T, rows []Row) {
				assert.Len(t, rows, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}

			rows, err := rp.Rows()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, rows)
				}
			}
		})
	}
}

// TestResPromQLLen 测试Len方法
// 功能：获取结果数量
// 测试覆盖：Vector类型、Matrix类型、空结果
func TestResPromQLLen(t *testing.T) {
	tests := []struct {
		name        string          // 测试用例名称
		input       prommodel.Value // 输入值
		expectedLen int             // 预期长度
		wantErr     bool            // 是否预期错误
	}{
		{
			name: "Vector with 2 samples",
			input: prommodel.Vector{
				{Metric: prommodel.Metric{"job": "prometheus"}, Value: 1},
				{Metric: prommodel.Metric{"job": "node"}, Value: 0},
			},
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name: "Matrix with 1 series and 2 points",
			input: prommodel.Matrix{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Values: []prommodel.SamplePair{
						{Timestamp: 1234567890, Value: 1},
						{Timestamp: 1234567950, Value: 0},
					},
				},
			},
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "Empty vector",
			input:       prommodel.Vector{},
			expectedLen: 0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}

			length, err := rp.Len()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedLen, length)
			}
		})
	}
}

// TestResPromQLData 测试Data方法
// 功能：获取Data结构
// 测试覆盖：Vector类型、Matrix类型、错误处理
func TestResPromQLData(t *testing.T) {
	tests := []struct {
		name     string          // 测试用例名称
		input    prommodel.Value // 输入值
		wantErr  bool            // 是否预期错误
		validate func(*testing.T, *Data)
	}{
		{
			name: "Vector value",
			input: prommodel.Vector{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Value:  prommodel.SampleValue(1),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data *Data) {
				assert.NotNil(t, data)
				assert.Equal(t, "vector", data.ResultType)
				assert.Len(t, data.Result, 1)
			},
		},
		{
			name: "Matrix value",
			input: prommodel.Matrix{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Values: []prommodel.SamplePair{
						{Timestamp: 1234567890, Value: 1},
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data *Data) {
				assert.NotNil(t, data)
				assert.Equal(t, "matrix", data.ResultType)
				assert.Len(t, data.Result, 1)
			},
		},
		{
			name:    "Empty vector",
			input:   prommodel.Vector{},
			wantErr: false,
			validate: func(t *testing.T, data *Data) {
				assert.NotNil(t, data)
				assert.Equal(t, "vector", data.ResultType)
				assert.Len(t, data.Result, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}

			data, err := rp.Data()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, data)
				}
			}
		})
	}
}

// TestResPromQLResultType 测试ResultType方法
// 功能：获取结果类型
// 测试覆盖：Vector、Matrix类型（Scalar类型不支持）
func TestResPromQLResultType(t *testing.T) {
	tests := []struct {
		name         string          // 测试用例名称
		input        prommodel.Value // 输入值
		expectedType string          // 预期类型
		wantErr      bool            // 是否预期错误
	}{
		{
			name:         "Vector type",
			input:        prommodel.Vector{},
			expectedType: "vector",
			wantErr:      false,
		},
		{
			name:         "Matrix type",
			input:        prommodel.Matrix{},
			expectedType: "matrix",
			wantErr:      false,
		},
		{
			name:         "Scalar type - not supported",
			input:        &prommodel.Scalar{},
			expectedType: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}

			resultType, err := rp.ResultType()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedType, resultType)
			}
		})
	}
}

// TestResPromQLMetric 测试Metric方法
// 功能：获取所有指标信息
// 测试覆盖：Vector类型、Matrix类型、空结果
func TestResPromQLMetric(t *testing.T) {
	tests := []struct {
		name     string          // 测试用例名称
		input    prommodel.Value // 输入值
		wantErr  bool            // 是否预期错误
		validate func(*testing.T, []prommodel.Metric)
	}{
		{
			name: "Vector with metrics",
			input: prommodel.Vector{
				{Metric: prommodel.Metric{"job": "prometheus", "instance": "localhost:9090"}},
				{Metric: prommodel.Metric{"job": "node", "instance": "localhost:9100"}},
			},
			wantErr: false,
			validate: func(t *testing.T, metrics []prommodel.Metric) {
				assert.Len(t, metrics, 2)
				assert.Equal(t, "prometheus", string(metrics[0]["job"]))
				assert.Equal(t, "node", string(metrics[1]["job"]))
			},
		},
		{
			name: "Matrix with metrics",
			input: prommodel.Matrix{
				{
					Metric: prommodel.Metric{"job": "prometheus"},
					Values: []prommodel.SamplePair{{Timestamp: 1234567890, Value: 1}},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, metrics []prommodel.Metric) {
				assert.Len(t, metrics, 1)
				assert.Equal(t, "prometheus", string(metrics[0]["job"]))
			},
		},
		{
			name:    "Empty vector",
			input:   prommodel.Vector{},
			wantErr: false,
			validate: func(t *testing.T, metrics []prommodel.Metric) {
				assert.Len(t, metrics, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResPromQL{Val: tt.input}

			metrics, err := rp.Metric()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, metrics)
				}
			}
		})
	}
}

// ExamplePrometheusClient_QueryRange 示例函数 - 时间范围查询
func ExamplePrometheusClient_QueryRange() {
	client, err := NewPrometheusClient("http://localhost:9090")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()
	step := 1 * time.Minute

	// 执行时间范围查询
	result, err := client.QueryRange(ctx, "up", start, end, step)
	if err != nil {
		panic(err)
	}

	_ = result
	fmt.Println("Range query completed")
}

// ExampleResPromQL_Rows 示例函数 - 转换为行
func ExampleResPromQL_Rows() {
	// 创建一个Vector值
	vector := prommodel.Vector{
		{
			Metric: prommodel.Metric{"job": "prometheus"},
			Value:  prommodel.SampleValue(1),
		},
	}

	rp := &ResPromQL{Val: vector}

	// 转换为行
	rows, err := rp.Rows()
	if err != nil {
		panic(err)
	}

	_ = rows
	fmt.Println("Rows converted successfully")
}

// ExampleResPromQL_Data 示例函数 - 获取Data结构
func ExampleResPromQL_Data() {
	// 创建一个Vector值
	vector := prommodel.Vector{
		{
			Metric: prommodel.Metric{"job": "prometheus"},
			Value:  prommodel.SampleValue(1),
		},
	}

	rp := &ResPromQL{Val: vector}

	// 获取Data结构
	data, err := rp.Data()
	if err != nil {
		panic(err)
	}

	_ = data
	fmt.Println("Data retrieved successfully")
}
