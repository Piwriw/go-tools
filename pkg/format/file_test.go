package format

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDirSizeWithDU 测试获取目录大小
// 测试场景：正常目录、空目录、不存在目录、深层嵌套目录
func TestDirSizeWithDU(t *testing.T) {
	// 跳过Windows系统的测试，因为du命令在Windows上不可用
	if runtime.GOOS == "windows" {
		t.Skip("跳过Windows系统，du命令不可用")
	}

	// 检查系统是否安装了du命令
	if !isDUCommandAvailable() {
		t.Skip("du命令不可用，跳过测试")
	}

	tests := []struct {
		name        string                  // 测试用例名称
		setupFunc   func() (string, func()) // 目录设置函数和清理函数
		wantErr     bool                    // 是否预期发生错误
		errContains string                  // 预期错误信息包含的字符串
		minSize     int64                   // 最小预期大小（用于验证）
	}{
		{
			name: "正常场景_包含文件的目录",
			setupFunc: func() (string, func()) {
				// 创建临时目录
				tempDir := t.TempDir()

				// 创建测试文件（约1KB大小）
				testFile := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(testFile, make([]byte, 1024), 0644)
				require.NoError(t, err)

				return tempDir, func() {} // TempDir会自动清理
			},
			wantErr: false,
			minSize: 1024,
		},
		{
			name: "正常场景_包含多个文件的目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()

				// 创建多个测试文件
				sizes := []int{512, 1024, 2048, 4096}
				for i, size := range sizes {
					testFile := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
					err := os.WriteFile(testFile, make([]byte, size), 0644)
					require.NoError(t, err)
				}

				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 512 + 1024 + 2048 + 4096,
		},
		{
			name: "正常场景_嵌套目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()

				// 创建子目录
				subDir := filepath.Join(tempDir, "subdir")
				err := os.Mkdir(subDir, 0755)
				require.NoError(t, err)

				// 在根目录创建文件
				rootFile := filepath.Join(tempDir, "root.txt")
				err = os.WriteFile(rootFile, make([]byte, 1024), 0644)
				require.NoError(t, err)

				// 在子目录创建文件
				subFile := filepath.Join(subDir, "sub.txt")
				err = os.WriteFile(subFile, make([]byte, 2048), 0644)
				require.NoError(t, err)

				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 1024 + 2048,
		},
		{
			name: "边界条件_空目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 0,
		},
		{
			name: "边界条件_只有子目录的空目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()

				// 创建空子目录
				subDir := filepath.Join(tempDir, "empty_subdir")
				err := os.Mkdir(subDir, 0755)
				require.NoError(t, err)

				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 0,
		},
		{
			name: "异常场景_不存在的目录",
			setupFunc: func() (string, func()) {
				return "/nonexistent/directory/that/does/not/exist", func() {}
			},
			wantErr:     true,
			errContains: "",
		},
		{
			name: "边界条件_大文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()

				// 创建较大的文件（1MB）
				testFile := filepath.Join(tempDir, "large.txt")
				err := os.WriteFile(testFile, make([]byte, 1024*1024), 0644)
				require.NoError(t, err)

				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 1024 * 1024,
		},
		{
			name: "边界条件_符号链接",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()

				// 创建原始文件
				originalFile := filepath.Join(tempDir, "original.txt")
				err := os.WriteFile(originalFile, make([]byte, 512), 0644)
				require.NoError(t, err)

				// 创建符号链接
				linkFile := filepath.Join(tempDir, "link.txt")
				err = os.Symlink(originalFile, linkFile)
				if err != nil {
					// 某些系统可能不支持符号链接，跳过此测试
					return "", func() {}
				}

				return tempDir, func() {}
			},
			wantErr: false,
			minSize: 0, // 符号链接可能被计入也可能不计入，取决于du命令的实现
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setupFunc()
			defer cleanup()

			// 跳过空路径的情况（符号链接测试在某些系统上不支持）
			if path == "" {
				return
			}

			size, err := DirSizeWithDU(path)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
			} else {
				require.NoError(t, err, "预期不发生错误")
				assert.GreaterOrEqual(t, size, tt.minSize, "目录大小应大于等于最小预期值")
			}
		})
	}
}

// TestDirSizeWithDU_Parallel 并发安全测试_DirSizeWithDU
func TestDirSizeWithDU_Parallel(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("跳过Windows系统")
	}

	if !isDUCommandAvailable() {
		t.Skip("du命令不可用")
	}

	t.Run("并发获取不同目录大小", func(t *testing.T) {
		// 创建多个临时目录
		dirs := make([]string, 5)
		for i := range dirs {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
			err := os.WriteFile(testFile, make([]byte, 1024*(i+1)), 0644)
			require.NoError(t, err)
			dirs[i] = tempDir
		}

		// 使用t.Parallel()进行并发测试
		for i, dir := range dirs {
			t.Run(fmt.Sprintf("Dir%d", i), func(t *testing.T) {
				t.Parallel()
				size, err := DirSizeWithDU(dir)
				assert.NoError(t, err)
				assert.Greater(t, size, int64(0))
			})
		}
	})
}

// TestDirSizeWithDU_RealDirectory 测试真实场景
func TestDirSizeWithDU_RealDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("跳过Windows系统")
	}

	if !isDUCommandAvailable() {
		t.Skip("du命令不可用")
	}

	t.Run("当前目录", func(t *testing.T) {
		// 获取当前工作目录
		cwd, err := os.Getwd()
		require.NoError(t, err)

		size, err := DirSizeWithDU(cwd)
		require.NoError(t, err)
		assert.Greater(t, size, int64(0), "当前目录大小应大于0")
	})
}

// BenchmarkDirSizeWithDU 性能基准测试_DirSizeWithDU
func BenchmarkDirSizeWithDU(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("跳过Windows系统")
	}

	if !isDUCommandAvailable() {
		b.Skip("du命令不可用")
	}

	// 创建测试目录
	setupBenchmark := func(b *testing.B) string {
		tempDir := b.TempDir()

		// 创建多个文件和子目录
		for i := 0; i < 10; i++ {
			subDir := filepath.Join(tempDir, fmt.Sprintf("dir%d", i))
			err := os.Mkdir(subDir, 0755)
			require.NoError(b, err)

			for j := 0; j < 10; j++ {
				testFile := filepath.Join(subDir, fmt.Sprintf("file%d.txt", j))
				err := os.WriteFile(testFile, make([]byte, 1024), 0644)
				require.NoError(b, err)
			}
		}

		return tempDir
	}

	b.Run("SmallDirectory", func(b *testing.B) {
		tempDir := setupBenchmark(b)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = DirSizeWithDU(tempDir)
		}
	})
}

// isDUCommandAvailable 检查du命令是否可用
func isDUCommandAvailable() bool {
	_, err := exec.LookPath("du")
	return err == nil
}

// ExampleDirSizeWithDU 示例代码_DirSizeWithDU
func ExampleDirSizeWithDU() {
	// 注意：此示例需要du命令可用
	size, err := DirSizeWithDU("/tmp")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Directory size: %d bytes\n", size)
}
