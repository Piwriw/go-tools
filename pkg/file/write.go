package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// WriteOptions configures file write behavior.
type WriteOptions struct {
	Perm    os.FileMode // File permission, default 0644
	DirPerm os.FileMode // Directory permission, default 0755
	Append  bool        // Append mode, default false (overwrite)
	Create  bool        // Create file if not exists, default true
}

// DefaultWriteOptions provides sensible defaults for file writes.
var DefaultWriteOptions = WriteOptions{
	Perm:    0644,
	DirPerm: 0755,
	Append:  false,
	Create:  true,
}

// WriteToFile writes byte data to the specified file.
// It creates the parent directory if it does not exist, and supports both
// overwrite and append modes controlled by WriteOptions.
func WriteToFile(filePath string, data []byte, opts *WriteOptions) error {
	if opts == nil {
		opts = &DefaultWriteOptions
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, opts.DirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	flags := os.O_WRONLY
	if opts.Create {
		flags |= os.O_CREATE
	}
	if opts.Append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(filePath, flags, opts.Perm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WriteStringToFile is a convenience wrapper that writes a string to a file.
func WriteStringToFile(filePath string, content string, opts *WriteOptions) error {
	return WriteToFile(filePath, []byte(content), opts)
}

// AppendToFile appends byte data to the end of a file.
// Creates the file if it does not exist.
func AppendToFile(filePath string, data []byte) error {
	opts := DefaultWriteOptions
	opts.Append = true
	return WriteToFile(filePath, data, &opts)
}

// AppendStringToFile appends a string to the end of a file.
// Creates the file if it does not exist.
func AppendStringToFile(filePath string, content string) error {
	return AppendToFile(filePath, []byte(content))
}

// WriteLinesToFile writes a string slice to a file, one line at a time.
// A newline character is appended after each line.
// Creates an empty file if lines is empty.
func WriteLinesToFile(filePath string, lines []string, opts *WriteOptions) error {
	total := len(lines)
	for _, line := range lines {
		total += len(line)
	}
	buf := make([]byte, 0, total)
	for _, line := range lines {
		buf = append(buf, line...)
		buf = append(buf, '\n')
	}
	return WriteToFile(filePath, buf, opts)
}

// CopyFile copies a file's content and permissions from src to dst.
// The destination directory is created automatically if it does not exist.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	return nil
}

// MoveFile moves a file from src to dst.
// It tries os.Rename first (atomic on the same filesystem), and falls back
// to copy+delete for cross-filesystem moves.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	if !errors.Is(err, syscall.EXDEV) {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	if err := CopyFile(src, dst); err != nil {
		return fmt.Errorf("failed to move across filesystems: %w", err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

// Touch creates an empty file or updates the file's access and modification times.
// If the file does not exist, it is created; if it already exists, its timestamps
// are updated to the current time.
func Touch(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	file.Close()

	now := time.Now()
	return os.Chtimes(filePath, now, now)
}

// AtomicWriteToFile writes data to a file atomically.
// It first writes to a temporary file in the same directory, then renames it
// to the target path. This prevents half-written states and is suitable for
// configuration file updates and other consistency-critical scenarios.
func AtomicWriteToFile(filePath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	success := false
	defer func() {
		if !success {
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	tmpFile.Close()

	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	success = true
	return nil
}

// --- Convenience writers with common permissions ---
// These functions write content first, then explicitly os.Chmod to guarantee
// the target permission regardless of the process umask.

// WritePrivateFile writes data to a file with 0600 permissions (owner read/write only).
// Suitable for credentials, tokens, SSH keys, and other sensitive data.
func WritePrivateFile(filePath string, data []byte) error {
	opts := WriteOptions{Perm: 0600, DirPerm: 0700, Create: true}
	if err := WriteToFile(filePath, data, &opts); err != nil {
		return err
	}
	return os.Chmod(filePath, 0600)
}

// WritePrivateStringFile writes a string to a file with 0600 permissions.
func WritePrivateStringFile(filePath string, content string) error {
	return WritePrivateFile(filePath, []byte(content))
}

// WriteExecutableFile writes data to a file with 0755 permissions (executable).
// Suitable for scripts, binaries, and other executable content.
func WriteExecutableFile(filePath string, data []byte) error {
	if err := WriteToFile(filePath, data, nil); err != nil {
		return err
	}
	return os.Chmod(filePath, 0755)
}

// WriteSharedFile writes data to a file with 0666 permissions (group/other writable).
// Suitable for collaborative files, shared logs, and temp files.
func WriteSharedFile(filePath string, data []byte) error {
	if err := WriteToFile(filePath, data, nil); err != nil {
		return err
	}
	return os.Chmod(filePath, 0666)
}

// WriteSharedStringFile writes a string to a file with 0666 permissions.
func WriteSharedStringFile(filePath string, content string) error {
	return WriteSharedFile(filePath, []byte(content))
}

// WriteConfigFile atomically writes data to a file with 0600 permissions.
// Uses atomic write to prevent half-written config states.
// Suitable for YAML/JSON/TOML config files containing secrets.
func WriteConfigFile(filePath string, data []byte) error {
	return AtomicWriteToFile(filePath, data, 0600)
}

// WriteConfigStringFile atomically writes a string to a file with 0600 permissions.
func WriteConfigStringFile(filePath string, content string) error {
	return WriteConfigFile(filePath, []byte(content))
}

// AppendToLogFile appends data to a file with 0644 permissions.
// Creates the file if it does not exist. A convenience wrapper for log writing.
func AppendToLogFile(filePath string, data []byte) error {
	return AppendToFile(filePath, data)
}

// AppendLogString appends a string to a log file with 0644 permissions.
func AppendLogString(filePath string, content string) error {
	return AppendStringToFile(filePath, content)
}
