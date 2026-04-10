package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnsureDeleted 测试确保删除功能
func TestEnsureDeleted(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (string, func())
		wantErr   bool
	}{
		{
			name: "正常场景_删除文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(testFile, []byte("content"), 0644)
				require.NoError(t, err)

				// 验证文件存在
				_, err = os.Stat(testFile)
				require.NoError(t, err)

				return testFile, func() {}
			},
			wantErr: false,
		},
		{
			name: "正常场景_删除目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "subdir")

				// 创建目录和文件
				err := os.Mkdir(subDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("content"), 0644)
				require.NoError(t, err)

				return subDir, func() {}
			},
			wantErr: false,
		},
		{
			name: "正常场景_删除嵌套目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				nestedDir := filepath.Join(tempDir, "a", "b", "c")

				err := os.MkdirAll(nestedDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(nestedDir, "file.txt"), []byte("content"), 0644)
				require.NoError(t, err)

				return filepath.Join(tempDir, "a"), func() {}
			},
			wantErr: false,
		},
		{
			name: "边界条件_路径不存在",
			setupFunc: func() (string, func()) {
				return "/nonexistent/path/that/does/not/exist", func() {}
			},
			wantErr: false, // 不存在应该返回 nil，不报错
		},
		{
			name: "边界条件_空目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				emptyDir := filepath.Join(tempDir, "empty")
				err := os.Mkdir(emptyDir, 0755)
				require.NoError(t, err)
				return emptyDir, func() {}
			},
			wantErr: false,
		},
		{
			name: "边界条件_删除只读文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				readOnlyFile := filepath.Join(tempDir, "readonly.txt")
				err := os.WriteFile(readOnlyFile, []byte("content"), 0444) // 只读
				require.NoError(t, err)
				return readOnlyFile, func() {}
			},
			wantErr: false, // RemoveAll 应该能删除只读文件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setupFunc()
			defer cleanup()

			err := EnsureDeleted(path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证路径已被删除
				_, err := os.Stat(path)
				assert.True(t, os.IsNotExist(err), "路径应该已被删除")
			}
		})
	}
}

// TestForceDeleteDir 测试强制删除目录
func TestForceDeleteDir(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (string, func())
		wantErr   bool
	}{
		{
			name: "正常场景_删除含只读文件的目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testDir := filepath.Join(tempDir, "testdir")

				// 创建目录
				err := os.Mkdir(testDir, 0755)
				require.NoError(t, err)

				// 创建只读文件
				readOnlyFile := filepath.Join(testDir, "readonly.txt")
				err = os.WriteFile(readOnlyFile, []byte("content"), 0444)
				require.NoError(t, err)

				// 创建普通文件
				normalFile := filepath.Join(testDir, "normal.txt")
				err = os.WriteFile(normalFile, []byte("content"), 0644)
				require.NoError(t, err)

				return testDir, func() {}
			},
			wantErr: false,
		},
		{
			name: "正常场景_删除嵌套只读文件目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "parent", "child")

				err := os.MkdirAll(subDir, 0755)
				require.NoError(t, err)

				// 在子目录创建只读文件
				readOnlyFile := filepath.Join(subDir, "readonly.txt")
				err = os.WriteFile(readOnlyFile, []byte("content"), 0400)
				require.NoError(t, err)

				return filepath.Join(tempDir, "parent"), func() {}
			},
			wantErr: false,
		},
		{
			name: "边界条件_路径是文件",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "file.txt")
				err := os.WriteFile(testFile, []byte("content"), 0644)
				require.NoError(t, err)
				return testFile, func() {}
			},
			wantErr: false,
		},
		{
			name: "边界条件_路径不存在",
			setupFunc: func() (string, func()) {
				return "/nonexistent/directory", func() {}
			},
			wantErr: false,
		},
		{
			name: "正常场景_空目录",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				emptyDir := filepath.Join(tempDir, "empty")
				err := os.Mkdir(emptyDir, 0755)
				require.NoError(t, err)
				return emptyDir, func() {}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setupFunc()
			defer cleanup()

			err := ForceDeleteDir(path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证路径已被删除
				_, err := os.Stat(path)
				assert.True(t, os.IsNotExist(err), "路径应该已被删除")
			}
		})
	}
}

// TestDeleteByPattern 测试按模式删除文件
func TestDeleteByPattern(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func() string
		pattern      string
		wantCount    int
		wantErr      bool
		shouldExist  []string // 删除后应该存在的文件
		shouldDelete []string // 删除后应该不存在的文件
	}{
		{
			name: "正常场景_删除所有txt文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"test1.txt", "test2.txt", "test.log"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			pattern:      "*.txt",
			wantCount:    2,
			wantErr:      false,
			shouldExist:  []string{"test.log"},
			shouldDelete: []string{"test1.txt", "test2.txt"},
		},
		{
			name: "正常场景_使用?通配符",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"file1.txt", "file2.txt", "file10.txt"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			pattern:      "file?.txt",
			wantCount:    2, // file1.txt, file2.txt（file10.txt 不匹配）
			wantErr:      false,
			shouldExist:  []string{"file10.txt"},
			shouldDelete: []string{"file1.txt", "file2.txt"},
		},
		{
			name: "边界条件_没有匹配文件",
			setupFunc: func() string {
				return t.TempDir()
			},
			pattern:     "*.txt",
			wantCount:   0,
			wantErr:     false,
			shouldExist: []string{},
		},
		{
			name: "边界条件_模式匹配目录但不删除",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("content"), 0644)
				_ = os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)
				return tempDir
			},
			pattern:      "test*",
			wantCount:    1, // 只删除文件，不删除目录
			wantErr:      false,
			shouldExist:  []string{"subdir"},
			shouldDelete: []string{"test.txt"},
		},
		{
			name: "边界条件_完整路径模式",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "exact.txt"), []byte("content"), 0644)
				return tempDir
			},
			pattern:      "exact.txt",
			wantCount:    1,
			wantErr:      false,
			shouldDelete: []string{"exact.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := tt.setupFunc()
			fullPattern := filepath.Join(tempDir, tt.pattern)

			count, err := DeleteByPattern(fullPattern)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)

				// 验证文件存在性
				for _, f := range tt.shouldExist {
					_, err := os.Stat(filepath.Join(tempDir, f))
					assert.NoError(t, err, "文件应该存在: %s", f)
				}
				for _, f := range tt.shouldDelete {
					_, err := os.Stat(filepath.Join(tempDir, f))
					assert.Error(t, err, "文件应该不存在: %s", f)
					assert.True(t, os.IsNotExist(err))
				}
			}
		})
	}
}

// TestDeleteByPatternInDir 测试在目录中按模式删除
func TestDeleteByPatternInDir(t *testing.T) {
	t.Run("正常场景_删除目录中的文件", func(t *testing.T) {
		tempDir := t.TempDir()

		// 创建文件
		files := []string{"test1.txt", "test2.txt", "test.log"}
		for _, f := range files {
			err := os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
			require.NoError(t, err)
		}

		count, err := DeleteByPatternInDir(tempDir, "*.txt")
		assert.NoError(t, err)
		assert.Equal(t, 2, count)

		// 验证
		_, err = os.Stat(filepath.Join(tempDir, "test1.txt"))
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(filepath.Join(tempDir, "test.log"))
		assert.NoError(t, err)
	})

	t.Run("边界条件_目录不存在", func(t *testing.T) {
		count, err := DeleteByPatternInDir("/nonexistent/dir", "*.txt")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// ExampleEnsureDeleted 示例代码
func ExampleEnsureDeleted() {
	err := EnsureDeleted("/tmp/test.txt")
	if err != nil {
		println(err.Error())
	}
}

// ExampleDeleteByPattern 示例代码
func ExampleDeleteByPattern() {
	count, err := DeleteByPattern("/tmp/*.log")
	if err != nil {
		println(err.Error())
	} else {
		println(fmt.Sprintf("Deleted %d files", count))
	}
}
