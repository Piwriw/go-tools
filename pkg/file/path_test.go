package file

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetFileExt 测试获取文件扩展名
func TestGetFileExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_单个扩展名",
			path:     "/path/to/file.txt",
			expected: ".txt",
		},
		{
			name:     "正常场景_多个扩展名",
			path:     "/path/to/archive.tar.gz",
			expected: ".gz",
		},
		{
			name:     "边界条件_无扩展名",
			path:     "/path/to/file",
			expected: "",
		},
		{
			name:     "边界条件_隐藏文件",
			path:     "/path/to/.hidden",
			expected: ".hidden", // filepath.Ext 返回最后一个点后的内容
		},
		{
			name:     "边界条件_隐藏文件带扩展名",
			path:     "/path/to/.hidden.txt",
			expected: ".txt",
		},
		{
			name:     "边界条件_点文件",
			path:     ".gitignore",
			expected: ".gitignore", // filepath.Ext 对点文件返回整个名称
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExt(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetFileName 测试获取文件名（不含扩展名）
func TestGetFileName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_带扩展名",
			path:     "/path/to/file.txt",
			expected: "file",
		},
		{
			name:     "正常场景_多个扩展名",
			path:     "/path/to/archive.tar.gz",
			expected: "archive.tar",
		},
		{
			name:     "边界条件_无扩展名",
			path:     "/path/to/file",
			expected: "file",
		},
		{
			name:     "边界条件_只有文件名",
			path:     "file.txt",
			expected: "file",
		},
		{
			name:     "边界条件_目录路径",
			path:     "/path/to/dir/",
			expected: "dir", // filepath.Base("/path/to/dir/") 返回 "dir"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetFileNameWithExt 测试获取完整文件名
func TestGetFileNameWithExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_完整路径",
			path:     "/path/to/file.txt",
			expected: "file.txt",
		},
		{
			name:     "边界条件_只有文件名",
			path:     "file.txt",
			expected: "file.txt",
		},
		{
			name:     "边界条件_目录路径",
			path:     "/path/to/dir/",
			expected: "dir",
		},
		{
			name:     "边界条件_根目录",
			path:     "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileNameWithExt(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetDirName 测试获取目录名
func TestGetDirName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_多层路径",
			path:     "/path/to/file.txt",
			expected: "to",
		},
		{
			name:     "正常场景_单层路径",
			path:     "/file.txt",
			expected: "/",
		},
		{
			name:     "边界条件_只有文件名",
			path:     "file.txt",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDirName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetDirPath 测试获取目录路径
func TestGetDirPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_多层路径",
			path:     "/path/to/file.txt",
			expected: "/path/to",
		},
		{
			name:     "边界条件_根目录文件",
			path:     "/file.txt",
			expected: "/",
		},
		{
			name:     "边界条件_只有文件名",
			path:     "file.txt",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDirPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJoinPath 测试路径拼接
func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "正常场景_多个组件",
			parts:    []string{"path", "to", "file.txt"},
			expected: filepath.Join("path", "to", "file.txt"),
		},
		{
			name:     "正常场景_绝对路径",
			parts:    []string{"/abs", "path", "file.txt"},
			expected: filepath.Join("/abs", "path", "file.txt"),
		},
		{
			name:     "边界条件_包含空字符串",
			parts:    []string{"path", "", "file.txt"},
			expected: filepath.Join("path", "file.txt"),
		},
		{
			name:     "边界条件_全部为空",
			parts:    []string{"", "", ""},
			expected: "",
		},
		{
			name:     "边界条件_单个组件",
			parts:    []string{"single"},
			expected: "single",
		},
		{
			name:     "边界条件_无组件",
			parts:    []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinPath(tt.parts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAbsPath 测试获取绝对路径
func TestAbsPath(t *testing.T) {
	t.Run("正常场景_相对路径转绝对路径", func(t *testing.T) {
		// 创建临时目录并切换
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			_ = os.Chdir(originalWd)
		}()

		tempDir := t.TempDir()
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		abs, err := AbsPath("test.txt")
		assert.NoError(t, err)
		assert.True(t, filepath.IsAbs(abs))
		assert.Contains(t, abs, "test.txt")
	})

	t.Run("正常场景_已经是绝对路径", func(t *testing.T) {
		tempDir := t.TempDir()
		abs, err := AbsPath(tempDir)
		assert.NoError(t, err)
		assert.True(t, filepath.IsAbs(abs))
	})

	t.Run("异常场景_空路径", func(t *testing.T) {
		_, err := AbsPath("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "路径不能为空")
	})
}

// TestMustAbsPath 测试获取绝对路径（失败返回空）
func TestMustAbsPath(t *testing.T) {
	t.Run("正常场景_成功获取绝对路径", func(t *testing.T) {
		tempDir := t.TempDir()
		abs := MustAbsPath(tempDir)
		assert.NotEmpty(t, abs)
		assert.True(t, filepath.IsAbs(abs))
	})

	t.Run("边界条件_空路径返回空", func(t *testing.T) {
		abs := MustAbsPath("")
		assert.Empty(t, abs)
	})
}

// TestResolveHome 测试展开主目录路径
func TestResolveHome(t *testing.T) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" && runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	}

	// 跳过没有 HOME 环境变量的情况
	if homeDir == "" {
		t.Skip("HOME 环境变量未设置")
	}

	tests := []struct {
		name     string
		path     string
		contains string
	}{
		{
			name:     "正常场景_展开~",
			path:     "~/Documents",
			contains: "Documents",
		},
		{
			name:     "正常场景_~/斜杠",
			path:     "~/",
			contains: "",
		},
		{
			name:     "正常场景_只有~",
			path:     "~",
			contains: "",
		},
		{
			name:     "边界条件_不以~开头",
			path:     "/absolute/path",
			contains: "/absolute/path",
		},
		{
			name:     "边界条件_相对路径",
			path:     "relative/path",
			contains: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveHome(tt.path)
			if tt.path == "~" {
				assert.Contains(t, result, homeDir)
			} else if tt.path[0] != '~' {
				assert.Equal(t, tt.path, result)
			} else {
				assert.Contains(t, result, tt.contains)
				assert.NotContains(t, result, "~")
			}
		})
	}
}

// TestIsAbs 测试判断是否为绝对路径
func TestIsAbs(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "正常场景_绝对路径",
			path:     "/path/to/file",
			expected: true,
		},
		{
			name:     "正常场景_相对路径",
			path:     "path/to/file",
			expected: false,
		},
		{
			name:     "边界条件_当前目录",
			path:     ".",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAbs(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRelPath 测试获取相对路径
func TestRelPath(t *testing.T) {
	t.Run("正常场景_计算相对路径", func(t *testing.T) {
		base := "/path/to"
		target := "/path/to/file.txt"

		rel, err := RelPath(base, target)
		assert.NoError(t, err)
		assert.Equal(t, "file.txt", rel)
	})

	t.Run("正常场景_跨目录相对路径", func(t *testing.T) {
		base := "/path/to/dir"
		target := "/path/other/file.txt"

		rel, err := RelPath(base, target)
		assert.NoError(t, err)
		// 结果因系统而异，但不应为空
		assert.NotEmpty(t, rel)
	})
}

// TestCleanPath 测试规范化路径
func TestCleanPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "正常场景_清理冗余分隔符",
			path:     "path//to///file",
			expected: filepath.Join("path", "to", "file"),
		},
		{
			name:     "正常场景_清理.",
			path:     "path/./to/file",
			expected: filepath.Join("path", "to", "file"),
		},
		{
			name:     "正常场景_清理..",
			path:     "path/to/../file",
			expected: filepath.Join("path", "file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSplitPath 测试分割路径
func TestSplitPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedDir string
		expectedFile string
	}{
		{
			name:        "正常场景_文件路径",
			path:        "/path/to/file.txt",
			expectedDir: "/path/to/",
			expectedFile: "file.txt",
		},
		{
			name:        "边界条件_只有文件名",
			path:        "file.txt",
			expectedDir: "",
			expectedFile: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, file := SplitPath(tt.path)
			assert.Equal(t, tt.expectedDir, dir)
			assert.Equal(t, tt.expectedFile, file)
		})
	}
}

// TestGetCommonPath 测试获取公共路径
func TestGetCommonPath(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		contains string
	}{
		{
			name:     "正常场景_有公共路径",
			paths:    []string{"/path/to/file1.txt", "/path/to/file2.txt"},
			contains: filepath.Join("path", "to"),
		},
		{
			name:     "边界条件_单个路径",
			paths:    []string{"/path/to/file.txt"},
			contains: "path",
		},
		{
			name:     "边界条件_无公共路径",
			paths:    []string{"/path1/file", "/path2/file"},
			contains: "",
		},
		{
			name:     "边界条件_空列表",
			paths:    []string{},
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommonPath(tt.paths...)
			if tt.contains == "" {
				assert.Equal(t, "", result)
			} else {
				assert.Contains(t, result, tt.contains)
			}
		})
	}
}

// TestHasParent 测试判断父子关系
func TestHasParent(t *testing.T) {
	tests := []struct {
		name     string
		parent   string
		child    string
		expected bool
	}{
		{
			name:     "正常场景_子目录在父目录下",
			parent:   "/path/to",
			child:    "/path/to/file.txt",
			expected: true,
		},
		{
			name:     "正常场景_子目录不在父目录下",
			parent:   "/path/to",
			child:    "/other/file.txt",
			expected: false,
		},
		{
			name:     "边界条件_相同路径",
			parent:   "/path/to",
			child:    "/path/to",
			expected: true, // filepath.Rel 返回 "."，不以 ".." 开头
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasParent(tt.parent, tt.child)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestChangeExt 测试更改扩展名
func TestChangeExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		newExt   string
		expected string
	}{
		{
			name:     "正常场景_更改扩展名",
			path:     "/path/to/file.txt",
			newExt:   ".log",
			expected: "/path/to/file.log",
		},
		{
			name:     "边界条件_无原扩展名",
			path:     "/path/to/file",
			newExt:   ".log",
			expected: "/path/to/file.log",
		},
		{
			name:     "正常场景_多个扩展名",
			path:     "/path/to/archive.tar.gz",
			newExt:   ".zip",
			expected: "/path/to/archive.tar.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ChangeExt(tt.path, tt.newExt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEnsureExt 测试确保扩展名
func TestEnsureExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		ext      string
		expected string
	}{
		{
			name:     "正常场景_添加扩展名",
			path:     "/path/to/file",
			ext:      "txt",
			expected: "/path/to/file.txt",
		},
		{
			name:     "正常场景_已有扩展名不替换",
			path:     "/path/to/file.txt",
			ext:      "txt",
			expected: "/path/to/file.txt",
		},
		{
			name:     "正常场景_已有扩展名替换",
			path:     "/path/to/file.log",
			ext:      "txt",
			expected: "/path/to/file.txt",
		},
		{
			name:     "正常场景_带点的扩展名",
			path:     "/path/to/file",
			ext:      ".txt",
			expected: "/path/to/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnsureExt(tt.path, tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}
