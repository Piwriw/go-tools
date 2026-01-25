package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPClient 测试创建新HTTP客户端
// 功能：使用各种Options创建HTTP客户端
// 测试覆盖：正常场景、边界条件、各种配置选项
func TestNewHTTPClient(t *testing.T) {
	tests := []struct {
		name   string                    // 测试用例名称
		opts   []Options                 // 配置选项
		verify func(*testing.T, *Client) // 验证函数
	}{
		{
			name: "Default client",
			opts: nil,
			verify: func(t *testing.T, c *Client) {
				assert.NotNil(t, c)
				assert.NotNil(t, c.client)
				// contentType is empty by default, will be set to JSON when making requests
				assert.Equal(t, ContentType(""), c.contentType)
			},
		},
		{
			name: "With custom timeout",
			opts: []Options{WithTimeout(10 * time.Second)},
			verify: func(t *testing.T, c *Client) {
				assert.Equal(t, 10*time.Second, c.client.Timeout)
			},
		},
		{
			name: "With custom content type",
			opts: []Options{WithContentType("application/xml")},
			verify: func(t *testing.T, c *Client) {
				assert.Equal(t, ContentType("application/xml"), c.contentType)
			},
		},
		{
			name: "With idle conn timeout",
			opts: []Options{WithIdleConnTimeout(30 * time.Second)},
			verify: func(t *testing.T, c *Client) {
				transport, ok := c.client.Transport.(*http.Transport)
				require.True(t, ok, "Transport should be http.Transport")
				assert.Equal(t, 30*time.Second, transport.IdleConnTimeout)
			},
		},
		{
			name: "With max idle conns",
			opts: []Options{WithMaxIdleConns(100)},
			verify: func(t *testing.T, c *Client) {
				transport, ok := c.client.Transport.(*http.Transport)
				require.True(t, ok, "Transport should be http.Transport")
				assert.Equal(t, 100, transport.MaxIdleConns)
			},
		},
		{
			name: "With disable keep alives",
			opts: []Options{WithDisableKeepAlives(true)},
			verify: func(t *testing.T, c *Client) {
				transport, ok := c.client.Transport.(*http.Transport)
				require.True(t, ok, "Transport should be http.Transport")
				assert.True(t, transport.DisableKeepAlives)
			},
		},
		{
			name: "With check redirect",
			opts: []Options{WithCheckRedirect(func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			})},
			verify: func(t *testing.T, c *Client) {
				assert.NotNil(t, c.client.CheckRedirect)
			},
		},
		{
			name: "Multiple options",
			opts: []Options{
				WithTimeout(5 * time.Second),
				WithContentType("text/plain"),
				WithMaxIdleConns(50),
			},
			verify: func(t *testing.T, c *Client) {
				assert.Equal(t, 5*time.Second, c.client.Timeout)
				assert.Equal(t, ContentType("text/plain"), c.contentType)
				transport, ok := c.client.Transport.(*http.Transport)
				require.True(t, ok)
				assert.Equal(t, 50, transport.MaxIdleConns)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewHTTPClient(tt.opts...)
			tt.verify(t, client)
		})
	}
}

// TestClientGet 测试GET请求
// 功能：发送GET请求并获取响应
// 测试覆盖：正常场景、边界条件、错误场景
func TestClientGet(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		response   string // 模拟服务器响应
		statusCode int    // 模拟HTTP状态码
		wantErr    bool   // 是否预期错误
		errMsg     string // 预期错误信息子字符串
	}{
		{
			name:       "Successful GET request",
			response:   `{"message":"success"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Empty response",
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Large response",
			response:   strings.Repeat("x", 10000),
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "JSON response",
			response:   `{"status":"ok","data":{"key":"value"}}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "404 Not Found",
			response:   "Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    false, // GET不检查状态码
		},
		{
			name:       "500 Internal Server Error",
			response:   "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    false, // GET不检查状态码
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, string(JSON), r.Header.Get("Content-Type"))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// 创建客户端并发送请求
			client := NewHTTPClient()
			resp, err := client.Get(server.URL)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.response, string(resp))
			}
		})
	}
}

// TestClientGetWithParams 测试带参数的GET请求
// 功能：发送带查询参数的GET请求
// 测试覆盖：正常场景、边界条件、特殊字符
func TestClientGetWithParams(t *testing.T) {
	tests := []struct {
		name      string            // 测试用例名称
		baseURL   string            // 基础URL
		params    map[string]string // 查询参数
		wantQuery string            // 预期查询字符串
		response  string            // 模拟响应
		wantErr   bool              // 是否预期错误
	}{
		{
			name:      "Single parameter",
			baseURL:   "http://example.com/api",
			params:    map[string]string{"key": "value"},
			wantQuery: "key=value",
			response:  `{"result":"ok"}`,
			wantErr:   false,
		},
		{
			name:      "Multiple parameters",
			baseURL:   "http://example.com/api",
			params:    map[string]string{"key1": "value1", "key2": "value2"},
			wantQuery: "key1=value1&key2=value2",
			response:  `{"result":"ok"}`,
			wantErr:   false,
		},
		{
			name:      "Empty parameters",
			baseURL:   "http://example.com/api",
			params:    map[string]string{},
			wantQuery: "",
			response:  `{"result":"ok"}`,
			wantErr:   false,
		},
		{
			name:     "Special characters in value",
			baseURL:  "http://example.com/api",
			params:   map[string]string{"q": "hello world"},
			response: `{"result":"ok"}`,
			wantErr:  false,
		},
		{
			name:     "URL with existing query",
			baseURL:  "http://example.com/api?existing=param",
			params:   map[string]string{"new": "value"},
			response: `{"result":"ok"}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, string(JSON), r.Header.Get("Content-Type"))
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// 创建客户端并发送请求
			client := NewHTTPClient()
			resp, err := client.GetWithParams(server.URL, tt.params)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.response, string(resp))
			}
		})
	}
}

// TestClientPost 测试POST请求
// 功能：发送POST请求并获取响应
// 测试覆盖：正常场景、边界条件、各种数据类型
func TestClientPost(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		data       []byte // 请求数据
		response   string // 模拟响应
		statusCode int    // HTTP状态码
		wantErr    bool   // 是否预期错误
	}{
		{
			name:       "Successful POST with JSON",
			data:       []byte(`{"key":"value"}`),
			response:   `{"success":true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST with empty data",
			data:       []byte{},
			response:   `{"success":true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST with nil data",
			data:       nil,
			response:   `{"success":true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST with large data",
			data:       []byte(strings.Repeat("x", 100000)),
			response:   `{"success":true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST with binary data",
			data:       []byte{0x00, 0x01, 0x02, 0x03},
			response:   `{"success":true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, string(JSON), r.Header.Get("Content-Type"))

				// 读取请求体
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				// 对于 nil data，bytes.NewBuffer 返回空 buffer
				expectedData := tt.data
				if tt.data == nil {
					expectedData = []byte{}
				}
				assert.Equal(t, expectedData, body)

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// 创建客户端并发送请求
			client := NewHTTPClient()
			resp, err := client.Post(server.URL, tt.data)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.response, string(resp))
			}
		})
	}
}

// TestClientPut 测试PUT请求
// 功能：发送PUT请求并获取响应
// 测试覆盖：正常场景、边界条件
func TestClientPut(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		data     []byte // 请求数据
		response string // 模拟响应
		wantErr  bool   // 是否预期错误
	}{
		{
			name:     "Successful PUT request",
			data:     []byte(`{"updated":true}`),
			response: `{"success":true}`,
			wantErr:  false,
		},
		{
			name:     "PUT with empty data",
			data:     []byte{},
			response: `{"success":true}`,
			wantErr:  false,
		},
		{
			name:     "PUT with nil data",
			data:     nil,
			response: `{"success":true}`,
			wantErr:  false,
			// bytes.NewBuffer(nil) creates an empty buffer, not nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, string(JSON), r.Header.Get("Content-Type"))

				// 读取请求体
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				// 对于 nil data，bytes.NewBuffer 返回空 buffer
				expectedData := tt.data
				if tt.data == nil {
					expectedData = []byte{}
				}
				assert.Equal(t, expectedData, body)

				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// 创建客户端并发送请求
			client := NewHTTPClient()
			resp, err := client.Put(server.URL, tt.data)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.response, string(resp))
			}
		})
	}
}

// TestClientPostFile 测试文件上传
// 功能：使用multipart/form-data上传文件
// 测试覆盖：正常场景、边界条件、各种文件类型
func TestClientPostFile(t *testing.T) {
	tests := []struct {
		name      string // 测试用例名称
		fieldName string // 表单字段名
		content   string // 文件内容
		wantErr   bool   // 是否预期错误
		errMsg    string // 预期错误信息
	}{
		{
			name:      "Successful file upload",
			fieldName: "file",
			content:   "Hello, World!",
			wantErr:   false,
		},
		{
			name:      "Upload empty file",
			fieldName: "file",
			content:   "",
			wantErr:   false,
		},
		{
			name:      "Upload large file",
			fieldName: "file",
			content:   strings.Repeat("x", 10000),
			wantErr:   false,
		},
		{
			name:      "Upload binary data",
			fieldName: "file",
			content:   string([]byte{0x00, 0x01, 0x02, 0x03}),
			wantErr:   false,
		},
		{
			name:      "Upload with special field name",
			fieldName: "data[file]",
			content:   "test content",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)

				// 解析multipart表单
				err := r.ParseMultipartForm(10 << 20) // 10MB
				require.NoError(t, err)

				// 验证Content-Type包含multipart/form-data
				assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

				// 获取文件
				file, handler, err := r.FormFile(tt.fieldName)
				require.NoError(t, err)
				defer file.Close()

				// 读取文件内容
				content, err := io.ReadAll(file)
				require.NoError(t, err)

				// 验证文件内容
				assert.Equal(t, tt.content, string(content))
				assert.NotEmpty(t, handler.Filename)

				w.Write([]byte(`{"success":true}`))
			}))
			defer server.Close()

			// 创建临时文件
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "test.txt")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// 创建客户端并上传文件
			client := NewHTTPClient()
			resp, err := client.PostFile(server.URL, tt.fieldName, filePath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, `{"success":true}`, string(resp))
			}
		})
	}
}

// TestClientPostFileError 测试文件上传错误场景
// 功能：测试文件不存在等错误情况
func TestClientPostFileError(t *testing.T) {
	tests := []struct {
		name      string // 测试用例名称
		fieldName string // 表单字段名
		filePath  string // 文件路径
		wantErr   bool   // 是否预期错误
		errMsg    string // 预期错误信息
	}{
		{
			name:      "File does not exist",
			fieldName: "file",
			filePath:  "/nonexistent/path/to/file.txt",
			wantErr:   true,
			errMsg:    "no such file or directory",
		},
		{
			name:      "Empty file path",
			fieldName: "file",
			filePath:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// 创建客户端并尝试上传不存在的文件
			client := NewHTTPClient()
			resp, err := client.PostFile(server.URL, tt.fieldName, tt.filePath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

// TestJSON 测试JSON方法
// 功能：设置客户端的Content-Type为JSON
func TestJSON(t *testing.T) {
	client := NewHTTPClient()
	result := client.JSON()

	assert.Same(t, client, result, "JSON should return the same client instance")
	assert.Equal(t, JSON, client.contentType, "contentType should be set to JSON")
}

// TestGlobalSetters 测试全局设置函数
// 功能：测试SetTimeout、SetContentType等全局设置
func TestGlobalSetters(t *testing.T) {
	// 保存原始值
	originalTimeout := DefaultClient.client.Timeout
	originalContentType := DefaultClient.contentType

	t.Run("SetTimeout", func(t *testing.T) {
		SetTimeout(30 * time.Second)
		assert.Equal(t, 30*time.Second, DefaultClient.client.Timeout)
	})

	t.Run("SetContentType", func(t *testing.T) {
		SetContentType("text/plain")
		assert.Equal(t, ContentType("text/plain"), DefaultClient.contentType)
	})

	t.Run("SetIdleConnTimeout", func(t *testing.T) {
		SetIdleConnTimeout(60 * time.Second)
		transport, ok := DefaultClient.client.Transport.(*http.Transport)
		require.True(t, ok)
		assert.Equal(t, 60*time.Second, transport.IdleConnTimeout)
	})

	t.Run("SetMaxIdleConns", func(t *testing.T) {
		SetMaxIdleConns(200)
		transport, ok := DefaultClient.client.Transport.(*http.Transport)
		require.True(t, ok)
		assert.Equal(t, 200, transport.MaxIdleConns)
	})

	// 恢复原始值
	DefaultClient.client.Timeout = originalTimeout
	DefaultClient.contentType = originalContentType
}

// TestClientPostWithInvalidURL 测试无效URL的POST请求
func TestClientPostWithInvalidURL(t *testing.T) {
	tests := []struct {
		name    string // 测试用例名称
		url     string // 无效URL
		data    []byte // 请求数据
		wantErr bool   // 是否预期错误
	}{
		{
			name:    "Invalid URL format",
			url:     "://invalid-url",
			data:    []byte(`{"test":"data"}`),
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			data:    []byte(`{"test":"data"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewHTTPClient()
			resp, err := client.Post(tt.url, tt.data)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Nil(t, resp)
		})
	}
}

// TestClientConcurrent 并发安全测试
// 功能：测试多个goroutine同时使用客户端
func TestClientConcurrent(t *testing.T) {
	t.Parallel()

	var requestCount int64
	var mu sync.Mutex

	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := NewHTTPClient()

	// 并发发送请求
	const goroutines = 100
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer func() { done <- true }()
			_, err := client.Get(server.URL)
			assert.NoError(t, err)
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// 验证所有请求都成功处理
	assert.Equal(t, int64(goroutines), requestCount)
}

// BenchmarkGet 性能基准测试 - GET请求
func BenchmarkGet(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	client := NewHTTPClient()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := client.Get(server.URL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPost 性能基准测试 - POST请求
func BenchmarkPost(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	client := NewHTTPClient()
	data := []byte(`{"test":"data"}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := client.Post(server.URL, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPostFile 性能基准测试 - 文件上传
func BenchmarkPostFile(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	// 创建临时文件
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte(strings.Repeat("x", 1000))
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		b.Fatal(err)
	}

	client := NewHTTPClient()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := client.PostFile(server.URL, "file", filePath)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetParallel 并行基准测试 - GET请求
func BenchmarkGetParallel(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	client := NewHTTPClient()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Get(server.URL)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// ExampleNewHTTPClient 示例函数 - 创建HTTP客户端
func ExampleNewHTTPClient() {
	// 使用默认配置创建客户端
	client1 := NewHTTPClient()

	// 使用自定义超时创建客户端
	client2 := NewHTTPClient(WithTimeout(10 * time.Second))

	// 使用多个选项创建客户端
	client3 := NewHTTPClient(
		WithTimeout(30*time.Second),
		WithContentType("application/json"),
		WithMaxIdleConns(100),
		WithIdleConnTimeout(60*time.Second),
	)

	_, _ = client1, client2
	_ = client3
}

// ExampleClientGet 示例函数 - GET请求
func ExampleClient_Get() {
	client := NewHTTPClient()

	// 发送简单的GET请求
	resp, err := client.Get("http://example.com/api/data")
	if err != nil {
		panic(err)
	}

	println(string(resp))
}

// ExampleClientPost 示例函数 - POST请求
func ExampleClient_Post() {
	client := NewHTTPClient()

	// 发送JSON数据
	data := []byte(`{"name":"John","age":30}`)
	resp, err := client.Post("http://example.com/api/users", data)
	if err != nil {
		panic(err)
	}

	println(string(resp))
}

// ExampleClientGetWithParams 示例函数 - 带参数的GET请求
func ExampleClient_GetWithParams() {
	client := NewHTTPClient()

	// 发送带查询参数的GET请求
	params := map[string]string{
		"page":     "1",
		"limit":    "20",
		"category": "books",
	}

	resp, err := client.GetWithParams("http://example.com/api/items", params)
	if err != nil {
		panic(err)
	}

	println(string(resp))
}

// ExampleClientPostFile 示例函数 - 文件上传
func ExampleClient_PostFile() {
	client := NewHTTPClient()

	// 上传文件
	resp, err := client.PostFile(
		"http://example.com/api/upload",
		"file",
		"/path/to/file.txt",
	)
	if err != nil {
		panic(err)
	}

	println(string(resp))
}

// ExampleWithOptions 示例函数 - 使用各种选项配置客户端
func ExampleWithOptions() {
	// 创建一个带有各种自定义选项的客户端
	client := NewHTTPClient(
		WithTimeout(10*time.Second),
		WithContentType("application/json"),
		WithIdleConnTimeout(30*time.Second),
		WithMaxIdleConns(100),
		WithDisableKeepAlives(false),
	)

	resp, err := client.Get("http://example.com/api")
	if err != nil {
		panic(err)
	}

	println(string(resp))
}
