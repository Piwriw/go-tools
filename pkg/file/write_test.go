package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriteStringToDir 测试写入目录应该失败
func TestWriteStringToDir(t *testing.T) {
	tempDir := t.TempDir()
	// 尝试写入目录应该失败
	err := WriteStringToFile(tempDir, "hello world", nil)
	if err == nil {
		t.Errorf("期望写入目录应该失败，但没有错误")
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

// --- Tests for new functions ---

func TestAppendToFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "append.txt")

	// Append to non-existent file (should create)
	err := AppendToFile(filePath, []byte("hello"))
	if err != nil {
		t.Fatalf("AppendToFile failed: %v", err)
	}

	// Append to existing file
	err = AppendToFile(filePath, []byte(" world"))
	if err != nil {
		t.Fatalf("AppendToFile failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(content))
	}
}

func TestAppendStringToFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "append_str.txt")

	err := AppendStringToFile(filePath, "line1\n")
	if err != nil {
		t.Fatalf("AppendStringToFile failed: %v", err)
	}

	err = AppendStringToFile(filePath, "line2\n")
	if err != nil {
		t.Fatalf("AppendStringToFile failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "line1\nline2\n"
	if string(content) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(content))
	}
}

func TestWriteLinesToFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "lines.txt")

	lines := []string{"alpha", "beta", "gamma"}
	err := WriteLinesToFile(filePath, lines, nil)
	if err != nil {
		t.Fatalf("WriteLinesToFile failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "alpha\nbeta\ngamma\n"
	if string(content) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(content))
	}
}

func TestWriteLinesToFileEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty_lines.txt")

	err := WriteLinesToFile(filePath, []string{}, nil)
	if err != nil {
		t.Fatalf("WriteLinesToFile with empty slice failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("expected empty file, got '%s'", string(content))
	}
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "src.txt")
	dstPath := filepath.Join(tempDir, "dst.txt")

	// Create source file
	content := []byte("copy me")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	err := CopyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read dst file: %v", err)
	}
	if string(dstContent) != string(content) {
		t.Errorf("expected '%s', got '%s'", string(content), string(dstContent))
	}

	// Verify permissions
	srcInfo, _ := os.Stat(srcPath)
	dstInfo, _ := os.Stat(dstPath)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("permissions mismatch: src=%v, dst=%v", srcInfo.Mode(), dstInfo.Mode())
	}
}

func TestCopyFileNestedDir(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "src.txt")
	dstPath := filepath.Join(tempDir, "sub", "dir", "dst.txt")

	if err := os.WriteFile(srcPath, []byte("nested"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	err := CopyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile to nested dir failed: %v", err)
	}

	content, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read dst file: %v", err)
	}
	if string(content) != "nested" {
		t.Errorf("expected 'nested', got '%s'", string(content))
	}
}

func TestCopyFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	err := CopyFile(filepath.Join(tempDir, "noexist.txt"), filepath.Join(tempDir, "dst.txt"))
	if err == nil {
		t.Error("expected error for non-existent source file")
	}
}

func TestMoveFile(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "src.txt")
	dstPath := filepath.Join(tempDir, "dst.txt")

	if err := os.WriteFile(srcPath, []byte("move me"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	err := MoveFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("MoveFile failed: %v", err)
	}

	// Source should no longer exist
	if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
		t.Error("source file should have been removed")
	}

	// Destination should have the content
	content, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read dst file: %v", err)
	}
	if string(content) != "move me" {
		t.Errorf("expected 'move me', got '%s'", string(content))
	}
}

func TestMoveFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	err := MoveFile(filepath.Join(tempDir, "noexist.txt"), filepath.Join(tempDir, "dst.txt"))
	if err == nil {
		t.Error("expected error for non-existent source file")
	}
}

func TestTouch(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "touch.txt")

	// Touch non-existent file
	err := Touch(filePath)
	if err != nil {
		t.Fatalf("Touch failed: %v", err)
	}

	// File should exist and be empty
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}

	// Touch existing file should update mod time
	modBefore := info.ModTime()
	err = Touch(filePath)
	if err != nil {
		t.Fatalf("Touch (update) failed: %v", err)
	}
	info2, _ := os.Stat(filePath)
	if !info2.ModTime().After(modBefore) {
		t.Errorf("mod time should be updated: before=%v, after=%v", modBefore, info2.ModTime())
	}
}

func TestTouchNestedDir(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "a", "b", "touch.txt")

	err := Touch(filePath)
	if err != nil {
		t.Fatalf("Touch with nested dir failed: %v", err)
	}

	if !FileExists(filePath) {
		t.Error("file should exist after Touch")
	}
}

func TestTouchPreservesContent(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "existing.txt")

	if err := os.WriteFile(filePath, []byte("important data"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := Touch(filePath)
	if err != nil {
		t.Fatalf("Touch failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "important data" {
		t.Errorf("Touch should preserve existing content, got '%s'", string(content))
	}
}

func TestAtomicWriteToFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "atomic.txt")

	err := AtomicWriteToFile(filePath, []byte("atomic content"), 0644)
	if err != nil {
		t.Fatalf("AtomicWriteToFile failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "atomic content" {
		t.Errorf("expected 'atomic content', got '%s'", string(content))
	}
}

func TestAtomicWriteToFileOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "atomic_overwrite.txt")

	// Initial write
	if err := AtomicWriteToFile(filePath, []byte("old"), 0644); err != nil {
		t.Fatalf("initial write failed: %v", err)
	}

	// Overwrite atomically
	if err := AtomicWriteToFile(filePath, []byte("new content"), 0644); err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "new content" {
		t.Errorf("expected 'new content', got '%s'", string(content))
	}
}

func TestAtomicWriteToFileNestedDir(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sub", "atomic.txt")

	err := AtomicWriteToFile(filePath, []byte("nested atomic"), 0600)
	if err != nil {
		t.Fatalf("AtomicWriteToFile with nested dir failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "nested atomic" {
		t.Errorf("expected 'nested atomic', got '%s'", string(content))
	}

	// Verify permissions
	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %04o", info.Mode().Perm())
	}
}

func TestAtomicWriteToFileNoTempLeak(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "no_leak.txt")

	_ = AtomicWriteToFile(filePath, []byte("ok"), 0644)

	// No .tmp-* files should remain
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Errorf("temp file leaked: %s", e.Name())
		}
	}
}

// --- Tests for permission convenience writers ---

func TestWritePrivateFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "secret.key")

	err := WritePrivateFile(filePath, []byte("private data"))
	if err != nil {
		t.Fatalf("WritePrivateFile failed: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "private data" {
		t.Errorf("expected 'private data', got '%s'", string(content))
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %04o", info.Mode().Perm())
	}
}

func TestWritePrivateStringFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "token")

	err := WritePrivateStringFile(filePath, "abc123")
	if err != nil {
		t.Fatalf("WritePrivateStringFile failed: %v", err)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "abc123" {
		t.Errorf("expected 'abc123', got '%s'", string(content))
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %04o", info.Mode().Perm())
	}
}

func TestWritePrivateFileCreatesPrivateDir(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "secrets", "nested", "key.pem")

	err := WritePrivateFile(filePath, []byte("key"))
	if err != nil {
		t.Fatalf("WritePrivateFile with nested dir failed: %v", err)
	}

	// Parent dir should be 0700
	dirInfo, _ := os.Stat(filepath.Join(tempDir, "secrets"))
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("expected dir perm 0700, got %04o", dirInfo.Mode().Perm())
	}
}

func TestWriteExecutableFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "run.sh")

	err := WriteExecutableFile(filePath, []byte("#!/bin/sh\necho hi"))
	if err != nil {
		t.Fatalf("WriteExecutableFile failed: %v", err)
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0755 {
		t.Errorf("expected perm 0755, got %04o", info.Mode().Perm())
	}
}

func TestWriteSharedFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "shared.log")

	err := WriteSharedFile(filePath, []byte("shared data"))
	if err != nil {
		t.Fatalf("WriteSharedFile failed: %v", err)
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0666 {
		t.Errorf("expected perm 0666, got %04o", info.Mode().Perm())
	}
}

func TestWriteSharedStringFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "notes.txt")

	err := WriteSharedStringFile(filePath, "collaborative content")
	if err != nil {
		t.Fatalf("WriteSharedStringFile failed: %v", err)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "collaborative content" {
		t.Errorf("expected 'collaborative content', got '%s'", string(content))
	}
}

func TestWriteConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")

	err := WriteConfigFile(filePath, []byte("db_password: s3cret"))
	if err != nil {
		t.Fatalf("WriteConfigFile failed: %v", err)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "db_password: s3cret" {
		t.Errorf("unexpected content: %s", string(content))
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %04o", info.Mode().Perm())
	}
}

func TestWriteConfigStringFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "app.json")

	err := WriteConfigStringFile(filePath, `{"secret": "value"}`)
	if err != nil {
		t.Fatalf("WriteConfigStringFile failed: %v", err)
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %04o", info.Mode().Perm())
	}
}

func TestWriteConfigFileOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yml")

	_ = WriteConfigFile(filePath, []byte("old"))
	err := WriteConfigFile(filePath, []byte("new"))
	if err != nil {
		t.Fatalf("WriteConfigFile overwrite failed: %v", err)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "new" {
		t.Errorf("expected 'new', got '%s'", string(content))
	}
}

func TestAppendToLogFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "app.log")

	err := AppendToLogFile(filePath, []byte("line1\n"))
	if err != nil {
		t.Fatalf("AppendToLogFile failed: %v", err)
	}

	err = AppendToLogFile(filePath, []byte("line2\n"))
	if err != nil {
		t.Fatalf("AppendToLogFile second call failed: %v", err)
	}

	content, _ := os.ReadFile(filePath)
	expected := "line1\nline2\n"
	if string(content) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(content))
	}

	info, _ := os.Stat(filePath)
	if info.Mode().Perm() != 0644 {
		t.Errorf("expected perm 0644, got %04o", info.Mode().Perm())
	}
}

func TestAppendLogString(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "syslog")

	_ = AppendLogString(filePath, "msg1\n")
	_ = AppendLogString(filePath, "msg2\n")

	content, _ := os.ReadFile(filePath)
	if string(content) != "msg1\nmsg2\n" {
		t.Errorf("expected 'msg1\\nmsg2\\n', got '%s'", string(content))
	}
}
