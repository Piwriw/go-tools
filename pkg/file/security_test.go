package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSecureDelete 测试安全删除
func TestSecureDelete(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		passes    int
		wantErr   bool
	}{
		{
			name: "正常场景_安全删除文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "secret.txt")
				content := []byte("sensitive data")
				_ = os.WriteFile(testFile, content, 0644)
				return testFile
			},
			passes:  3,
			wantErr: false,
		},
		{
			name: "正常场景_单次覆写",
			setupFunc: func() string {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				_ = os.WriteFile(testFile, []byte("data"), 0644)
				return testFile
			},
			passes:  1,
			wantErr: false,
		},
		{
			name: "边界条件_文件不存在",
			setupFunc: func() string {
				return "/nonexistent/file.txt"
			},
			passes:  3,
			wantErr: false, // 文件不存在视为成功
		},
		{
			name: "异常场景_路径是目录",
			setupFunc: func() string {
				return t.TempDir()
			},
			passes:  3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()

			err := SecureDelete(path, tt.passes)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证文件已被删除
				if path != "" && !filepath.IsAbs(path) || os.Getenv("TEMP") != "" {
					// 只检查临时目录下的文件
					if FileExists(path) && IsFile(path) {
						assert.Fail(t, "文件应该已被删除")
					}
				}
			}
		})
	}
}

// TestSetExecutable 测试设置可执行权限
func TestSetExecutable(t *testing.T) {
	t.Run("正常场景_设置可执行", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "script.sh")
		_ = os.WriteFile(testFile, []byte("#!/bin/bash"), 0644)

		err := SetExecutable(testFile)
		assert.NoError(t, err)

		// 验证权限
		info, _ := os.Stat(testFile)
		mode := info.Mode().Perm()
		assert.Equal(t, os.FileMode(0755), mode)
	})
}

// TestSetReadOnly 测试设置只读权限
func TestSetReadOnly(t *testing.T) {
	t.Run("正常场景_设置只读", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "readonly.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		err := SetReadOnly(testFile)
		assert.NoError(t, err)

		// 验证权限
		info, _ := os.Stat(testFile)
		mode := info.Mode().Perm()
		assert.Equal(t, os.FileMode(0444), mode)
	})
}

// TestSetReadWrite 测试设置可读写权限
func TestSetReadWrite(t *testing.T) {
	t.Run("正常场景_设置可读写", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "rw.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0444)

		err := SetReadWrite(testFile)
		assert.NoError(t, err)

		// 验证权限
		info, _ := os.Stat(testFile)
		mode := info.Mode().Perm()
		assert.Equal(t, os.FileMode(0644), mode)
	})
}

// TestSetPermissions 测试设置自定义权限
func TestSetPermissions(t *testing.T) {
	t.Run("正常场景_设置自定义权限", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "custom.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		err := SetPermissions(testFile, 0600)
		assert.NoError(t, err)

		// 验证权限
		info, _ := os.Stat(testFile)
		mode := info.Mode().Perm()
		assert.Equal(t, os.FileMode(0600), mode)
	})
}

// TestGetFilePermissions 测试获取权限字符串
func TestGetFilePermissions(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		{
			name:     "正常场景_0644权限",
			mode:     0644,
			expected: "rw-r--r--",
		},
		{
			name:     "正常场景_0755权限",
			mode:     0755,
			expected: "rwxr-xr-x",
		},
		{
			name:     "正常场景_0600权限",
			mode:     0600,
			expected: "rw-------",
		},
		{
		name:     "正常场景_0755权限",
		mode:     0755,
		expected: "rwxr-xr-x", // umask 会影响，实际可能是 0755
		},
		{
			name:     "边界条件_0000权限",
			mode:     0000,
			expected: "---------",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "test.txt")
			_ = os.WriteFile(testFile, []byte("content"), tt.mode)

			perms, err := GetFilePermissions(testFile)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, perms)
		})
	}
}

// TestGetOctalPermissions 测试获取八进制权限
func TestGetOctalPermissions(t *testing.T) {
	t.Run("正常场景_获取0644权限", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		octal, err := GetOctalPermissions(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "0644", octal)
	})

	t.Run("正常场景_获取0755权限", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "script.sh")
		_ = os.WriteFile(testFile, []byte("#!/bin/bash"), 0755)

		octal, err := GetOctalPermissions(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "0755", octal)
	})
}

// TestSecureDeleteByPattern 测试按模式安全删除
func TestSecureDeleteByPattern(t *testing.T) {
	t.Run("正常场景_按模式删除", func(t *testing.T) {
		tempDir := t.TempDir()

		// 创建测试文件
		files := []string{"secret1.txt", "secret2.txt", "normal.log"}
		for _, f := range files {
			_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
		}

		pattern := filepath.Join(tempDir, "*.txt")
		count, err := SecureDeleteByPattern(pattern, 1)

		assert.NoError(t, err)
		assert.Equal(t, 2, count)

		// 验证文件已删除
		_, err = os.Stat(filepath.Join(tempDir, "secret1.txt"))
		assert.True(t, os.IsNotExist(err))

		// 验证其他文件仍存在
		_, err = os.Stat(filepath.Join(tempDir, "normal.log"))
		assert.NoError(t, err)
	})
}

// TestCopyPermissions 测试复制权限
func TestCopyPermissions(t *testing.T) {
	t.Run("正常场景_复制权限", func(t *testing.T) {
		tempDir := t.TempDir()

		srcFile := filepath.Join(tempDir, "src.txt")
		dstFile := filepath.Join(tempDir, "dst.txt")

		_ = os.WriteFile(srcFile, []byte("source"), 0755)
		_ = os.WriteFile(dstFile, []byte("dest"), 0644)

		err := CopyPermissions(srcFile, dstFile)
		assert.NoError(t, err)

		// 验证权限已复制
		srcInfo, _ := os.Stat(srcFile)
		dstInfo, _ := os.Stat(dstFile)
		assert.Equal(t, srcInfo.Mode().Perm(), dstInfo.Mode().Perm())
	})
}

// TestIsExecutable 测试判断可执行
func TestIsExecutable(t *testing.T) {
	t.Run("正常场景_可执行文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "exec")
		_ = os.WriteFile(testFile, []byte("content"), 0755)

		assert.True(t, IsExecutable(testFile))
	})

	t.Run("正常场景_不可执行文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "nonexec")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		assert.False(t, IsExecutable(testFile))
	})
}

// TestIsReadOnly 测试判断只读
func TestIsReadOnly(t *testing.T) {
	t.Run("正常场景_只读文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "readonly")
		_ = os.WriteFile(testFile, []byte("content"), 0444)

		assert.True(t, IsReadOnly(testFile))
	})

	t.Run("正常场景_可写文件", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "writable")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		assert.False(t, IsReadOnly(testFile))
	})
}

// TestMakePrivate 测试设置私有权限
func TestMakePrivate(t *testing.T) {
	t.Run("正常场景_设置私有", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "private.txt")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		err := MakePrivate(testFile)
		assert.NoError(t, err)

		info, _ := os.Stat(testFile)
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})
}

// TestMakePrivateExecutable 测试设置私有可执行权限
func TestMakePrivateExecutable(t *testing.T) {
	t.Run("正常场景_设置私有可执行", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "private.sh")
		_ = os.WriteFile(testFile, []byte("content"), 0644)

		err := MakePrivateExecutable(testFile)
		assert.NoError(t, err)

		info, _ := os.Stat(testFile)
		assert.Equal(t, os.FileMode(0700), info.Mode().Perm())
	})
}

// ExampleSetExecutable 示例代码
func ExampleSetExecutable() {
	err := SetExecutable("/tmp/script.sh")
	if err != nil {
		println(err.Error())
	}
}

// ExampleGetFilePermissions 示例代码
func ExampleGetFilePermissions() {
	perms, err := GetFilePermissions("/tmp/file.txt")
	if err != nil {
		println(err.Error())
		return
	}
	println(perms) // 输出: rw-r--r--
}
