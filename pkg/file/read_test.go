package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReadLines 测试按行读取文件
func TestReadLines(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func())
		wantErr     bool
		wantLines   int
		firstLine   string
		lastLine    string
		errContains string
	}{
		{
			name: "正常场景_多行文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "line1\nline2\nline3"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 3,
			firstLine: "line1",
			lastLine:  "line3",
		},
		{
			name: "正常场景_带CRLF的文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "line1\r\nline2\r\n"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 2,
			firstLine: "line1",
			lastLine:  "line2",
		},
		{
			name: "边界条件_空文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(testFile, []byte{}, 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 0,
		},
		{
			name: "边界条件_单行文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "single.txt")
				content := "single line"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 1,
			firstLine: "single line",
			lastLine:  "single line",
		},
		{
			name: "边界条件_只有换行符的文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "newlines.txt")
				content := "\n\n\n"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 3,
		},
		{
			name: "边界条件_较大行",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "large.txt")
				// 创建一个较大的行（在 scanner 默认缓冲区 64KB 内）
				longLine := strings.Repeat("a", 4*1024) // 4KB
				err := os.WriteFile(testFile, []byte(longLine), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr:   false,
			wantLines: 1,
			firstLine: strings.Repeat("a", 4*1024), // 期望读取完整的行
			lastLine:  strings.Repeat("a", 4*1024),
		},
		{
			name: "异常场景_文件不存在",
			setupFunc: func() (string, func()) {
				return "/nonexistent/file.txt", func() {}
			},
			wantErr:     true,
			errContains: "打开文件失败",
		},
		{
			name: "异常场景_路径是目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				return tempDir, func() {}
			},
			wantErr:     true,
			errContains: "读取文件失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFunc()
			defer cleanup()

			lines, err := ReadLines(filePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, lines, tt.wantLines)
				if len(lines) > 0 {
					assert.Equal(t, tt.firstLine, lines[0])
					assert.Equal(t, tt.lastLine, lines[len(lines)-1])
				}
				// 确保返回的不是 nil
				assert.NotNil(t, lines)
			}
		})
	}
}

// TestReadFileAsString 测试读取文件为字符串
func TestReadFileAsString(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func())
		wantContent string
		wantErr     bool
		errContains string
	}{
		{
			name: "正常场景_读取普通文本",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "Hello, World!"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantContent: "Hello, World!",
			wantErr:     false,
		},
		{
			name: "边界条件_空文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(testFile, []byte{}, 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantContent: "",
			wantErr:     false,
		},
		{
			name: "正常场景_包含特殊字符",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "special.txt")
				content := "你好\n世界\t\r\n"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantContent: "你好\n世界\t\r\n",
			wantErr:     false,
		},
		{
			name: "正常场景_二进制内容",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "binary.dat")
				content := []byte{0x00, 0x01, 0x02, 0xFF}
				err := os.WriteFile(testFile, content, 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantContent: string([]byte{0x00, 0x01, 0x02, 0xFF}),
			wantErr:     false,
		},
		{
			name: "异常场景_文件不存在",
			setupFunc: func() (string, func()) {
				return "/nonexistent/file.txt", func() {}
			},
			wantErr:     true,
			errContains: "读取文件失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFunc()
			defer cleanup()

			content, err := ReadFileAsString(filePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContent, content)
			}
		})
	}
}

// TestReadFileWithRetry 测试带重试的文件读取
func TestReadFileWithRetry(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func())
		config      *RetryConfig
		wantErr     bool
		errContains string
		minAttempts int // 最小尝试次数（用于测试重试逻辑）
	}{
		{
			name: "正常场景_一次成功",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "test content"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			config:  &RetryConfig{MaxAttempts: 3, WaitTime: 10 * time.Millisecond},
			wantErr: false,
		},
		{
			name: "正常场景_使用默认配置",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "test content"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			config:  nil,
			wantErr: false,
		},
		{
			name: "边界条件_MaxAttempts为1",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "test content"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			config:  &RetryConfig{MaxAttempts: 1, WaitTime: 10 * time.Millisecond},
			wantErr: false,
		},
		{
			name: "边界条件_MaxAttempts为0应修正为1",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := "test content"
				err := os.WriteFile(testFile, []byte(content), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			config:  &RetryConfig{MaxAttempts: 0, WaitTime: 10 * time.Millisecond},
			wantErr: false,
		},
		{
			name: "正常场景_读取空文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(testFile, []byte{}, 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			config:  &RetryConfig{MaxAttempts: 2, WaitTime: 10 * time.Millisecond},
			wantErr: false,
		},
		{
			name: "异常场景_文件不存在_重试多次",
			setupFunc: func() (string, func()) {
				return "/nonexistent/file.txt", func() {}
			},
			config:      &RetryConfig{MaxAttempts: 3, WaitTime: 5 * time.Millisecond},
			wantErr:     true,
			errContains: "读取文件失败（已重试 3 次）",
		},
		{
			name: "异常场景_文件不存在_一次尝试",
			setupFunc: func() (string, func()) {
				return "/nonexistent/file2.txt", func() {}
			},
			config:      &RetryConfig{MaxAttempts: 1, WaitTime: 10 * time.Millisecond},
			wantErr:     true,
			errContains: "读取文件失败（已重试 1 次）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFunc()
			defer cleanup()

			data, err := ReadFileWithRetry(filePath, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}

// TestReadFileAsStringWithRetry 测试带重试的字符串读取
func TestReadFileAsStringWithRetry(t *testing.T) {
	t.Run("正常场景_读取成功", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := "test content"
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		result, err := ReadFileAsStringWithRetry(testFile, nil)
		assert.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("异常场景_文件不存在", func(t *testing.T) {
		_, err := ReadFileAsStringWithRetry("/nonexistent/file.txt", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "读取文件失败（已重试")
	})
}

// BenchmarkReadLines 性能测试
func BenchmarkReadLines(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "bench.txt")

	// 创建一个包含 10000 行的文件
	var content string
	for i := 0; i < 10000; i++ {
		content += fmt.Sprintf("Line %d content here\n", i)
	}
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadLines(testFile)
	}
}

// ExampleReadLines 示例代码
func ExampleReadLines() {
	lines, err := ReadLines("example.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
