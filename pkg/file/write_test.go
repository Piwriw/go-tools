package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStringToDir(t *testing.T) {
	if err := WriteStringToFile("/Users/joohwan/GolandProjects/go-tools/pkg/file", "hello world", nil); err != nil {
		t.Errorf("WriteStringToFile failed: %v", err)
	}
}

// TestWriteToFile 测试WriteToFile函数的各种场景
func TestWriteToFile(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 测试用例
	tests := []struct {
		name        string
		filePath    string
		data        []byte
		opts        *WriteOptions
		expectError bool
		setup       func(string) error         // 测试前的准备函数
		verify      func(string, []byte) error // 验证函数
	}{
		{
			name:        "正常写入文件",
			filePath:    filepath.Join(tempDir, "test1.txt"),
			data:        []byte("hello world"),
			opts:        nil,
			expectError: false,
			setup:       nil,
			verify:      nil,
		},
		{
			name:        "覆盖已存在文件",
			filePath:    filepath.Join(tempDir, "test2.txt"),
			data:        []byte("new content"),
			opts:        nil,
			expectError: false,
			setup: func(filePath string) error {
				initialContent := []byte("initial content")
				return WriteToFile(filePath, initialContent, nil)
			},
			verify: func(filePath string, expectedData []byte) error {
				content, err := os.ReadFile(filePath)
				if err != nil {
					return err
				}
				if string(content) != string(expectedData) {
					return fmt.Errorf("内容不匹配，期望 %s，实际 %s", string(expectedData), string(content))
				}
				return nil
			},
		},
		{
			name:     "追加模式写入到新文件",
			filePath: filepath.Join(tempDir, "test3.txt"),
			data:     []byte("appended content"),
			opts:     &WriteOptions{Append: true, Create: true, Perm: 0644}, // 确保明确指定可读权限

			expectError: false,
			setup:       nil,
			verify:      nil,
		},
		{
			name:        "追加模式写入到已存在文件",
			filePath:    filepath.Join(tempDir, "test4.txt"),
			data:        []byte(" appended"),
			opts:        &WriteOptions{Append: true},
			expectError: false,
			setup: func(filePath string) error {
				initialContent := []byte("initial content")
				return WriteToFile(filePath, initialContent, nil)
			},
			verify: func(filePath string, expectedData []byte) error {
				content, err := os.ReadFile(filePath)
				if err != nil {
					return err
				}
				expected := "initial content appended"
				if string(content) != expected {
					return fmt.Errorf("追加内容不正确，期望 %s，实际 %s", expected, string(content))
				}
				return nil
			},
		},
		{
			name:        "创建嵌套目录文件",
			filePath:    filepath.Join(tempDir, "nested", "dir", "test5.txt"),
			data:        []byte("nested content"),
			opts:        nil,
			expectError: false,
			setup:       nil,
			verify:      nil,
		},
		{
			name:        "不创建文件标志错误",
			filePath:    filepath.Join(tempDir, "nonexist.txt"),
			data:        []byte("should fail"),
			opts:        &WriteOptions{Create: false},
			expectError: true,
			setup:       nil,
			verify:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行测试前准备
			if tt.setup != nil {
				err := tt.setup(tt.filePath)
				if err != nil {
					t.Fatalf("测试准备失败: %v", err)
				}
			}

			// 执行测试
			err := WriteToFile(tt.filePath, tt.data, tt.opts)

			// 验证错误结果
			if tt.expectError && err == nil {
				t.Errorf("期望错误但未发生错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("未期望错误但发生了错误: %v", err)
			}

			// 如果没有错误且有验证函数，则执行验证
			if !tt.expectError && err == nil {
				if tt.verify != nil {
					err = tt.verify(tt.filePath, tt.data)
					if err != nil {
						t.Errorf("验证失败: %v", err)
					}
				} else {
					// 默认验证：检查文件内容是否与写入数据一致
					content, readErr := os.ReadFile(tt.filePath)
					if readErr != nil {
						t.Fatalf("读取文件失败: %v", readErr)
					}
					if string(content) != string(tt.data) {
						t.Errorf("内容不匹配，期望 %s，实际 %s", string(tt.data), string(content))
					}
				}
			}
		})
	}
}

// TestWriteStringToFile 测试WriteStringToFile函数
func TestWriteStringToFile(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "string_test.txt")
	content := "test string content"

	err := WriteStringToFile(filePath, content, nil)
	if err != nil {
		t.Errorf("WriteStringToFile执行失败: %v", err)
	}

	// 验证文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("内容不匹配，期望 %s，实际 %s", content, string(fileContent))
	}
}

// TestWriteOptions 测试默认写入选项
func TestWriteOptions(t *testing.T) {
	expected := WriteOptions{
		Perm:    0644,
		DirPerm: 0755,
		Append:  false,
		Create:  true,
	}

	if DefaultWriteOptions != expected {
		t.Errorf("默认写入选项不匹配，期望 %+v，实际 %+v", expected, DefaultWriteOptions)
	}
}
