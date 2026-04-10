package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFileExists 测试文件存在性检查
func TestFileExists(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() string
		exists bool
	}{
		{
			name: "正常场景_文件存在",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				_ = os.WriteFile(testFile, []byte("content"), 0644)
				return testFile
			},
			exists: true,
		},
		{
			name: "正常场景_目录存在",
			setup: func() string {
				return t.TempDir()
			},
			exists: true,
		},
		{
			name: "边界条件_路径不存在",
			setup: func() string {
				return "/nonexistent/path"
			},
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := FileExists(path)
			assert.Equal(t, tt.exists, result)
		})
	}
}

// TestIsDir 测试目录判断
func TestIsDir(t *testing.T) {
	tests := []struct {
		name  string
		setup func() string
		isDir bool
	}{
		{
			name: "正常场景_是目录",
			setup: func() string {
				return t.TempDir()
			},
			isDir: true,
		},
		{
			name: "正常场景_不是目录",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "file.txt")
				_ = os.WriteFile(testFile, []byte("content"), 0644)
				return testFile
			},
			isDir: false,
		},
		{
			name: "边界条件_路径不存在",
			setup: func() string {
				return "/nonexistent/path"
			},
			isDir: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsDir(path)
			assert.Equal(t, tt.isDir, result)
		})
	}
}

// TestIsFile 测试文件判断
func TestIsFile(t *testing.T) {
	t.Run("正常场景_是文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "file.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		result := IsFile(testFile)
		assert.True(t, result)
	})

	t.Run("正常场景_不是文件（目录）", func(t *testing.T) {
		tempDir := t.TempDir()
		result := IsFile(tempDir)
		assert.False(t, result)
	})

	t.Run("边界条件_路径不存在", func(t *testing.T) {
		result := IsFile("/nonexistent/file")
		assert.False(t, result)
	})
}

// TestIsEmptyDir 测试空目录判断
func TestIsEmptyDir(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() string
		isEmpty  bool
		wantErr  bool
	}{
		{
			name: "正常场景_空目录",
			setup: func() string {
				return t.TempDir()
			},
			isEmpty: true,
			wantErr: false,
		},
		{
			name: "正常场景_非空目录",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "file.txt")
				_ = os.WriteFile(testFile, []byte("content"), 0644)
				return tempDir
			},
			isEmpty: false,
			wantErr: false,
		},
		{
			name: "正常场景_包含子目录的目录",
			setup: func() string {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "subdir")
				_ = os.Mkdir(subDir, 0755)
				return tempDir
			},
			isEmpty: false,
			wantErr: false,
		},
		{
			name: "边界条件_路径不存在",
			setup: func() string {
				return "/nonexistent/dir"
			},
			isEmpty: false,
			wantErr: false,
		},
		{
			name: "边界条件_路径是文件",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "file.txt")
				_ = os.WriteFile(testFile, []byte("content"), 0644)
				return testFile
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()

			isEmpty, err := IsEmptyDir(path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.isEmpty, isEmpty)
			}
		})
	}
}

// TestGetFileSize 测试获取文件大小
func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() string
		wantSize  int64
		wantErr   bool
		errContains string
	}{
		{
			name: "正常场景_获取文件大小",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				content := make([]byte, 1024)
				_ = os.WriteFile(testFile, content, 0644)
				return testFile
			},
			wantSize: 1024,
			wantErr:  false,
		},
		{
			name: "边界条件_空文件",
			setup: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "empty.txt")
				_ = os.WriteFile(testFile, []byte{}, 0644)
				return testFile
			},
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "异常场景_路径是目录",
			setup: func() string {
				return t.TempDir()
			},
			wantErr:      true,
			errContains: "是目录",
		},
		{
			name: "异常场景_路径不存在",
			setup: func() string {
				return "/nonexistent/file.txt"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()

			size, err := GetFileSize(path)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSize, size)
			}
		})
	}
}

// TestGetFileModTime 测试获取文件修改时间
func TestGetFileModTime(t *testing.T) {
	t.Run("正常场景_获取修改时间", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		modTime, err := GetFileModTime(testFile)
		assert.NoError(t, err)
		assert.Greater(t, modTime, int64(0))
	})

	t.Run("异常场景_路径不存在", func(t *testing.T) {
		_, err := GetFileModTime("/nonexistent/file")
		assert.Error(t, err)
	})
}

// TestGetFileHash 测试获取文件哈希
func TestGetFileHash(t *testing.T) {
	t.Run("正常场景_计算MD5", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("hello world")
		_ = os.WriteFile(testFile, content, 0644)

		// MD5 of "hello world" is "5eb63bbbe01eeed093cb22bb8f5acdc3"
		hash, err := GetFileHash(testFile, HashMD5)
		assert.NoError(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hash)
	})

	t.Run("正常场景_计算SHA256", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("hello world")
		_ = os.WriteFile(testFile, content, 0644)

		// SHA256 of "hello world" is "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
		hash, err := GetFileHash(testFile, HashSHA256)
		assert.NoError(t, err)
		assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", hash)
	})

	t.Run("异常场景_不支持的哈希类型", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		_, err := GetFileHash(testFile, "unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不支持的哈希类型")
	})

	t.Run("异常场景_文件不存在", func(t *testing.T) {
		_, err := GetFileHash("/nonexistent/file", HashMD5)
		assert.Error(t, err)
	})
}

// TestGetFileMD5 测试获取文件MD5
func TestGetFileMD5(t *testing.T) {
	t.Run("正常场景_计算MD5", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("test")
		_ = os.WriteFile(testFile, content, 0644)

		hash, err := GetFileMD5(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", hash)
	})
}

// TestGetFileSHA256 测试获取文件SHA256
func TestGetFileSHA256(t *testing.T) {
	t.Run("正常场景_计算SHA256", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("test")
		_ = os.WriteFile(testFile, content, 0644)

		hash, err := GetFileSHA256(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", hash)
	})
}

// TestFilesEqual 测试文件内容比较
func TestFilesEqual(t *testing.T) {
	t.Run("正常场景_相同文件", func(t *testing.T) {
		tempDir := t.TempDir()
		content := []byte("same content")

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")

		_ = os.WriteFile(file1, content, 0644)
		_ = os.WriteFile(file2, content, 0644)

		equal, err := FilesEqual(file1, file2)
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("正常场景_不同文件", func(t *testing.T) {
		tempDir := t.TempDir()

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")

		_ = os.WriteFile(file1, []byte("content1"), 0644)
		_ = os.WriteFile(file2, []byte("content2"), 0644)

		equal, err := FilesEqual(file1, file2)
		assert.NoError(t, err)
		assert.False(t, equal)
	})

	t.Run("异常场景_文件不存在", func(t *testing.T) {
		_, err := FilesEqual("/nonexistent/file1", "/nonexistent/file2")
		assert.Error(t, err)
	})
}

// TestGetFileInfo 测试获取文件详细信息
func TestGetFileInfo(t *testing.T) {
	t.Run("正常场景_获取文件信息", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := make([]byte, 1024)
		_ = os.WriteFile(testFile, content, 0644)

		info, err := GetFileInfo(testFile)
		assert.NoError(t, err)
		assert.Equal(t, testFile, info.Path)
		assert.Equal(t, int64(1024), info.Size)
		assert.False(t, info.IsDir)
		assert.NotEmpty(t, info.Mode)
		assert.Greater(t, info.ModTime, int64(0))
	})

	t.Run("正常场景_获取目录信息", func(t *testing.T) {
		tempDir := t.TempDir()

		info, err := GetFileInfo(tempDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir)
	})

	t.Run("异常场景_路径不存在", func(t *testing.T) {
		_, err := GetFileInfo("/nonexistent/file")
		assert.Error(t, err)
	})
}

// ExampleFileExists 示例代码
func ExampleFileExists() {
	if FileExists("/tmp/test.txt") {
		println("File exists")
	}
}

// ExampleGetFileMD5 示例代码
func ExampleGetFileMD5() {
	hash, err := GetFileMD5("/tmp/file.txt")
	if err != nil {
		println(err.Error())
		return
	}
	println(hash)
}
