package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFindFilesByExt 测试按扩展名查找文件
func TestFindFilesByExt(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string
		ext         string
		wantCount   int
		wantErr     bool
		shouldExist []string
	}{
		{
			name: "正常场景_查找txt文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"test1.txt", "test2.txt", "test.log", "data.json"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			ext:         ".txt",
			wantCount:   2,
			wantErr:     false,
			shouldExist: []string{"test1.txt", "test2.txt"},
		},
		{
			name: "正常场景_查找json文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"data.json", "config.json", "test.txt"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			ext:         ".json",
			wantCount:   2,
			wantErr:     false,
			shouldExist: []string{"data.json", "config.json"},
		},
		{
			name: "正常场景_查找嵌套目录中的文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "subdir")
				_ = os.Mkdir(subDir, 0755)
				_ = os.WriteFile(filepath.Join(tempDir, "root.txt"), []byte("content"), 0644)
				_ = os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("content"), 0644)
				_ = os.WriteFile(filepath.Join(subDir, "other.log"), []byte("content"), 0644)
				return tempDir
			},
			ext:         ".txt",
			wantCount:   2,
			wantErr:     false,
			shouldExist: []string{"root.txt", "nested.txt"},
		},
		{
			name: "边界条件_无匹配文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "test.log"), []byte("content"), 0644)
				return tempDir
			},
			ext:       ".txt",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "边界条件_空目录",
			setupFunc: func() string {
				return t.TempDir()
			},
			ext:       ".txt",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc()

			matches, err := FindFilesByExt(dir, tt.ext)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, matches, tt.wantCount)

				if len(tt.shouldExist) > 0 {
					for _, expected := range tt.shouldExist {
						found := false
						for _, match := range matches {
							if filepath.Base(match) == expected {
								found = true
								break
							}
						}
						assert.True(t, found, "应该找到文件: %s", expected)
					}
				}
			}
		})
	}
}

// TestFindFilesByExts 测试按多个扩展名查找
func TestFindFilesByExts(t *testing.T) {
	t.Run("正常场景_查找txt和log文件", func(t *testing.T) {
		tempDir := t.TempDir()
		files := []string{"test1.txt", "test2.txt", "app.log", "error.log", "data.json"}
		for _, f := range files {
			_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
		}

		matches, err := FindFilesByExts(tempDir, []string{".txt", ".log"})
		assert.NoError(t, err)
		assert.Len(t, matches, 4)
	})

	t.Run("边界条件_空扩展名列表", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("content"), 0644)

		matches, err := FindFilesByExts(tempDir, []string{})
		assert.NoError(t, err)
		assert.Len(t, matches, 0)
	})
}

// TestFindFilesByName 测试按名称模式查找
func TestFindFilesByName(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		pattern   string
		wantCount int
	}{
		{
			name: "正常场景_使用*通配符",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"test1.txt", "test2.txt", "other.log"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			pattern:   "test*.txt",
			wantCount: 2,
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
			pattern:   "file?.txt",
			wantCount: 2,
		},
		{
			name: "正常场景_使用[]通配符",
			setupFunc: func() string {
				tempDir := t.TempDir()
				files := []string{"file1.txt", "file2.txt", "file3.txt", "filea.txt"}
				for _, f := range files {
					_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
				}
				return tempDir
			},
			pattern:   "file[1-2].txt",
			wantCount: 2,
		},
		{
			name: "边界条件_精确匹配",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "exact.txt"), []byte("content"), 0644)
				return tempDir
			},
			pattern:   "exact.txt",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc()

			matches, err := FindFilesByName(dir, tt.pattern)
			assert.NoError(t, err)
			assert.Len(t, matches, tt.wantCount)
		})
	}
}

// TestFindRecentFiles 测试查找最近修改的文件
func TestFindRecentFiles(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		within    time.Duration
		wantCount int
	}{
		{
			name: "正常场景_查找最近1小时修改的文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "recent.txt"), []byte("content"), 0644)
				return tempDir
			},
			within:    time.Hour,
			wantCount: 1,
		},
		{
			name: "边界条件_查找最近1秒修改的文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "now.txt"), []byte("content"), 0644)
				return tempDir
			},
			within:    time.Second,
			wantCount: 1,
		},
		{
			name: "边界条件_查找很久以前修改的文件",
			setupFunc: func() string {
				tempDir := t.TempDir()
				_ = os.WriteFile(filepath.Join(tempDir, "old.txt"), []byte("content"), 0644)
				return tempDir
			},
			within:    time.Nanosecond, // 极短时间，不太可能匹配
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc()

			matches, err := FindRecentFiles(dir, tt.within)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(matches), tt.wantCount)
		})
	}
}

// TestFindLargeFiles 测试查找大文件
func TestFindLargeFiles(t *testing.T) {
	t.Run("正常场景_查找大于1KB的文件", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "small.txt"), make([]byte, 512), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "large.txt"), make([]byte, 2048), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "huge.txt"), make([]byte, 10*1024), 0644)

		matches, err := FindLargeFiles(tempDir, 1024)
		assert.NoError(t, err)
		assert.Len(t, matches, 2)

		// 验证返回的文件信息
		for _, info := range matches {
			assert.GreaterOrEqual(t, info.Size, int64(1024))
			assert.NotEmpty(t, info.Path)
		}
	})

	t.Run("边界条件_无匹配文件", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "small.txt"), make([]byte, 100), 0644)

		matches, err := FindLargeFiles(tempDir, 1024)
		assert.NoError(t, err)
		assert.Len(t, matches, 0)
	})
}

// TestFindFilesBySizeRange 测试按大小范围查找
func TestFindFilesBySizeRange(t *testing.T) {
	t.Run("正常场景_查找1KB-10KB范围的文件", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "small.txt"), make([]byte, 512), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "medium.txt"), make([]byte, 2048), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "large.txt"), make([]byte, 20*1024), 0644)

		matches, err := FindFilesBySizeRange(tempDir, 1024, 10*1024)
		assert.NoError(t, err)
		assert.Len(t, matches, 1)
		assert.Equal(t, "medium.txt", filepath.Base(matches[0].Path))
	})

	t.Run("正常场景_无上限查找", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "file1.txt"), make([]byte, 100), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "file2.txt"), make([]byte, 5*1024), 0644)

		matches, err := FindFilesBySizeRange(tempDir, 50, 0) // maxSize=0 表示无上限
		assert.NoError(t, err)
		assert.Len(t, matches, 2)
	})
}

// TestFindFilesByPredicate 测试使用谓词函数查找
func TestFindFilesByPredicate(t *testing.T) {
	t.Run("正常场景_自定义谓词", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "a.txt"), make([]byte, 100), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "b.log"), make([]byte, 200), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "c.txt"), make([]byte, 300), 0644)

		// 查找大小大于150且扩展名为.txt的文件
		matches, err := FindFilesByPredicate(tempDir, func(path string, info os.FileInfo) bool {
			return info.Size() > 150 && filepath.Ext(path) == ".txt"
		})

		assert.NoError(t, err)
		assert.Len(t, matches, 1)
		assert.Equal(t, "c.txt", filepath.Base(matches[0]))
	})
}

// TestCountFilesByExt 测试统计文件扩展名
func TestCountFilesByExt(t *testing.T) {
	t.Run("正常场景_统计各类型文件", func(t *testing.T) {
		tempDir := t.TempDir()
		files := []string{"a.txt", "b.txt", "c.txt", "d.log", "e.json"}
		for _, f := range files {
			_ = os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644)
		}

		counts, err := CountFilesByExt(tempDir)
		assert.NoError(t, err)
		assert.Equal(t, 3, counts[".txt"])
		assert.Equal(t, 1, counts[".log"])
		assert.Equal(t, 1, counts[".json"])
	})

	t.Run("边界条件_包含无扩展名文件", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(tempDir, "noext"), []byte("content"), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "with.txt"), []byte("content"), 0644)

		counts, err := CountFilesByExt(tempDir)
		assert.NoError(t, err)
		assert.Equal(t, 1, counts[""])
		assert.Equal(t, 1, counts[".txt"])
	})
}

// BenchmarkFindFilesByExt 性能测试
func BenchmarkFindFilesByExt(b *testing.B) {
	tempDir := b.TempDir()

	// 创建100个文件
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, "file.txt")
		if i%2 == 0 {
			filename = filepath.Join(tempDir, "file.log")
		}
		_ = os.WriteFile(filename, make([]byte, 1024), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindFilesByExt(tempDir, ".txt")
	}
}
