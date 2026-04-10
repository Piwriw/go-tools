package format

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatFileSize 测试文件大小格式化
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "正常场景_字节",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "正常场景_千字节",
			bytes:    1536,
			expected: "1.50 KB",
		},
		{
			name:     "正常场景_兆字节",
			bytes:    1048576 * 3,
			expected: "3.00 MB",
		},
		{
			name:     "正常场景_吉字节",
			bytes:    1073741824 * 2,
			expected: "2.00 GB",
		},
		{
			name:     "正常场景_太字节",
			bytes:    1099511627776,
			expected: "1.00 TB",
		},
		{
			name:     "边界条件_零字节",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "边界条件_1023字节",
			bytes:    1023,
			expected: "1023 B",
		},
		{
			name:     "边界条件_1024字节",
			bytes:    1024,
			expected: "1.00 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatFileSizeExact 测试精确小数位格式化
func TestFormatFileSizeExact(t *testing.T) {
	tests := []struct {
		name      string
		bytes     int64
		decimals  int
		expected  string
	}{
		{
			name:     "零小数位",
			bytes:    1536,
			decimals: 0,
			expected: "2 KB", // 1.5 KB 四舍五入
		},
		{
			name:     "一位小数",
			bytes:    1536,
			decimals: 1,
			expected: "1.5 KB",
		},
		{
			name:     "两位小数",
			bytes:    1536,
			decimals: 2,
			expected: "1.50 KB",
		},
		{
			name:     "三位小数",
			bytes:    1536,
			decimals: 3,
			expected: "1.500 KB",
		},
		{
			name:     "边界_负小数位",
			bytes:    1536,
			decimals: -1,
			expected: "2 KB", // 负值修正为0，四舍五入
		},
		{
			name:     "边界_超大小数位",
			bytes:    1536,
			decimals: 10,
			expected: "1.500 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileSizeExact(tt.bytes, tt.decimals)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatPath 测试路径格式化
func TestFormatPath(t *testing.T) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" && runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	}

	if homeDir == "" {
		t.Skip("HOME 环境变量未设置")
	}

	tests := []struct {
		name     string
		path     string
		options  []PathFormatOption
		contains string
	}{
		{
			name:     "正常场景_替换主目录",
			path:     filepath.Join(homeDir, "Documents", "file.txt"),
			options:  []PathFormatOption{WithHomeSymbol("~")},
			contains: "~/Documents/file.txt",
		},
		{
			name:     "正常场景_缩短路径",
			path:     "/very/long/path/to/file.txt",
			options:  []PathFormatOption{WithPathShortening(2)},
			contains: "file.txt",
		},
		{
			name:     "边界条件_不缩短",
			path:     "/path/to/file.txt",
			options:  []PathFormatOption{WithPathShortening(0)},
			contains: "/path/to/file.txt",
		},
		{
			name:     "边界条件_绝对路径",
			path:     "/absolute/path/file.txt",
			options:  []PathFormatOption{},
			contains: "/absolute/path/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPath(tt.path, tt.options...)
			assert.Contains(t, result, tt.contains)
		})
	}
}

// TestFormatPermissions 测试权限格式化
func TestFormatPermissions(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		{
			name:     "正常场景_0644",
			mode:     0644,
			expected: "rw-r--r--",
		},
		{
			name:     "正常场景_0755",
			mode:     0755,
			expected: "rwxr-xr-x",
		},
		{
			name:     "正常场景_0600",
			mode:     0600,
			expected: "rw-------",
		},
		{
			name:     "正常场景_0700",
			mode:     0700,
			expected: "rwx------",
		},
		{
			name:     "正常场景_0777",
			mode:     0777,
			expected: "rwxrwxrwx",
		},
		{
			name:     "边界条件_0000",
			mode:     0000,
			expected: "---------",
		},
		{
			name:     "边界条件_0444",
			mode:     0444,
			expected: "r--r--r--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPermissions(tt.mode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatPermissionsOctal 测试八进制权限格式化
func TestFormatPermissionsOctal(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		{
			name:     "正常场景_0644",
			mode:     0644,
			expected: "0644",
		},
		{
			name:     "正常场景_0755",
			mode:     0755,
			expected: "0755",
		},
		{
			name:     "正常场景_0600",
			mode:     0600,
			expected: "0600",
		},
		{
			name:     "边界条件_0000",
			mode:     0000,
			expected: "0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPermissionsOctal(tt.mode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatPathPermissions 测试获取路径权限并格式化
func TestFormatPathPermissions(t *testing.T) {
	t.Run("正常场景_获取文件权限", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		perms, err := FormatPathPermissions(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "rw-r--r--", perms)
	})

	t.Run("异常场景_文件不存在", func(t *testing.T) {
		_, err := FormatPathPermissions("/nonexistent/file")
		assert.Error(t, err)
	})
}

// TestGetFileType 测试获取文件类型
func TestGetFileType(t *testing.T) {
	t.Run("正常场景_普通文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "file.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

	 fileType := GetFileType(testFile)
	 assert.Equal(t, "file", fileType)
	})

	t.Run("正常场景_目录", func(t *testing.T) {
		tempDir := t.TempDir()
		fileType := GetFileType(tempDir)
		assert.Equal(t, "directory", fileType)
	})

	t.Run("正常场景_符号链接", func(t *testing.T) {
		tempDir := t.TempDir()
		originalFile := filepath.Join(tempDir, "original.txt")
		_ = os.WriteFile(originalFile, []byte("content"), 0644)

		linkFile := filepath.Join(tempDir, "link.txt")
		err := os.Symlink(originalFile, linkFile)
		if err != nil {
			t.Skip("无法创建符号链接")
		}

		fileType := GetFileType(linkFile)
		assert.Equal(t, "symlink", fileType)
	})

	t.Run("边界条件_不存在", func(t *testing.T) {
		fileType := GetFileType("/nonexistent/path")
		assert.Equal(t, "unknown", fileType)
	})
}

// TestFormatFileType 测试格式化文件类型字符
func TestFormatFileType(t *testing.T) {
	t.Run("正常场景_普通文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "file.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		typeChar := FormatFileType(testFile)
		assert.Equal(t, "-", typeChar)
	})

	t.Run("正常场景_目录", func(t *testing.T) {
		tempDir := t.TempDir()
		typeChar := FormatFileType(tempDir)
		assert.Equal(t, "d", typeChar)
	})

	t.Run("正常场景_符号链接", func(t *testing.T) {
		tempDir := t.TempDir()
		originalFile := filepath.Join(tempDir, "original.txt")
		_ = os.WriteFile(originalFile, []byte("content"), 0644)

		linkFile := filepath.Join(tempDir, "link.txt")
		err := os.Symlink(originalFile, linkFile)
		if err != nil {
			t.Skip("无法创建符号链接")
		}

		typeChar := FormatFileType(linkFile)
		assert.Equal(t, "l", typeChar)
	})

	t.Run("边界条件_不存在", func(t *testing.T) {
		typeChar := FormatFileType("/nonexistent/path")
		assert.Equal(t, "?", typeChar)
	})
}

// TestFormatFileInfo 测试格式化文件信息
func TestFormatFileInfo(t *testing.T) {
	t.Run("正常场景_格式化文件信息", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := make([]byte, 2048)
		_ = os.WriteFile(testFile, content, 0644)

		info, err := FormatFileInfo(testFile)
		assert.NoError(t, err)
		assert.Contains(t, info, "-rw-r--r--")
		assert.Contains(t, info, "KB")
		assert.Contains(t, info, "test.txt")
	})

	t.Run("正常场景_格式化目录信息", func(t *testing.T) {
		tempDir := t.TempDir()

		info, err := FormatFileInfo(tempDir)
		assert.NoError(t, err)
		assert.Contains(t, info, "d") // 目录以 d 开头
	})

	t.Run("异常场景_文件不存在", func(t *testing.T) {
		_, err := FormatFileInfo("/nonexistent/file")
		assert.Error(t, err)
	})
}

// BenchmarkFormatFileSize 性能测试
func BenchmarkFormatFileSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatFileSize(int64(1024 * 1024 * 100))
	}
}
