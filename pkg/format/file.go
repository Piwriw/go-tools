package format

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// DirSizeWithDU 获取目录大小（单位：字节）
// 使用du命令获取目录大小，支持macOS和Linux
//
// 注意：macOS上的du命令默认单位是KB，需要转换为字节
// 注意：Linux上的du命令默认单位是字节，无需转换
func DirSizeWithDU(path string) (int64, error) {
	var cmd *exec.Cmd
	// 检查操作系统类型，macOS使用不同的du选项
	if runtime.GOOS == "darwin" {
		// macOS上的du命令使用-k选项表示以KB为单位，需要乘以1024转换为字节
		cmd = exec.Command("du", "-sk", path)
	} else {
		// Linux使用-b选项表示以字节为单位
		cmd = exec.Command("du", "-sb", path)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(string(output))
	if len(parts) < 1 {
		return 0, fmt.Errorf("unexpected output from du")
	}

	size, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	// 如果是macOS，需要将KB转换为字节
	if runtime.GOOS == "darwin" {
		size = size * 1024
	}

	return size, nil
}

// ============================================================================
// File Size Formatting
// ============================================================================

// FormatFileSize 格式化文件大小为人类可读的字符串。
// 专门用于文件大小显示，使用 1024 进制。
// 例如：1536 → "1.50 KB"，1048576 → "1.00 MB"。
func FormatFileSize[T number](bytes T) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
	)

	b := float64(bytes)
	switch {
	case b >= PB:
		return fmt.Sprintf("%.2f PB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2f TB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2f GB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", b/KB)
	default:
		return fmt.Sprintf("%.0f B", b)
	}
}

// FormatFileSizeExact 格式化文件大小，保留指定小数位数。
// decimals 指定小数位数（0-3）。
func FormatFileSizeExact[T number](bytes T, decimals int) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
	)

	// 限制小数位数范围
	if decimals < 0 {
		decimals = 0
	}
	if decimals > 3 {
		decimals = 3
	}

	// 格式化模板
	format := fmt.Sprintf("%%.%df %%s", decimals)

	b := float64(bytes)
	switch {
	case b >= PB:
		return fmt.Sprintf(format, b/PB, "PB")
	case b >= TB:
		return fmt.Sprintf(format, b/TB, "TB")
	case b >= GB:
		return fmt.Sprintf(format, b/GB, "GB")
	case b >= MB:
		return fmt.Sprintf(format, b/MB, "MB")
	case b >= KB:
		return fmt.Sprintf(format, b/KB, "KB")
	default:
		return fmt.Sprintf("%d B", int64(b))
	}
}

// ============================================================================
// Path Formatting
// ============================================================================

// FormatPath 格式化路径用于显示。
// 可选参数：
//   - homeSymbol: 用于替换主目录的符号（如 "~"）
//   - shorten: 是否缩短路径（只显示最后 n 个组件）
func FormatPath(path string, options ...PathFormatOption) string {
	opts := defaultPathFormatOptions()
	for _, opt := range options {
		opt(&opts)
	}

	result := path

	// 展开 ~ 并准备替换
	if opts.homeSymbol != "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" && runtime.GOOS == "windows" {
			homeDir = os.Getenv("USERPROFILE")
		}

		if homeDir != "" && strings.HasPrefix(result, homeDir) {
			result = opts.homeSymbol + strings.TrimPrefix(result, homeDir)
		}
	}

	// 清理路径
	result = filepath.Clean(result)

	// 缩短路径
	if opts.shorten > 0 && opts.shorten < countPathComponents(result) {
		result = shortenPath(result, opts.shorten)
	}

	return result
}

// PathFormatOption 路径格式化选项
type PathFormatOption func(*PathFormatOptions)

// PathFormatOptions 路径格式化选项配置
type PathFormatOptions struct {
	homeSymbol string // 主目录替换符号（如 "~"）
	shorten    int    // 保留的路径组件数，0 表示不缩短
}

// WithHomeSymbol 设置主目录替换符号
func WithHomeSymbol(symbol string) PathFormatOption {
	return func(o *PathFormatOptions) {
		o.homeSymbol = symbol
	}
}

// WithPathShortening 设置路径缩短（保留最后 n 个组件）
func WithPathShortening(components int) PathFormatOption {
	return func(o *PathFormatOptions) {
		o.shorten = components
	}
}

func defaultPathFormatOptions() PathFormatOptions {
	return PathFormatOptions{
		homeSymbol: "~",
		shorten:    0,
	}
}

func countPathComponents(path string) int {
	path = filepath.Clean(path)
	if path == "/" || path == "." {
		return 1
	}

	// 分割路径
	parts := strings.Split(filepath.ToSlash(path), "/")

	// 过滤空字符串
	var nonEmpty []string
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}

	return len(nonEmpty)
}

func shortenPath(path string, keepComponents int) string {
	path = filepath.Clean(path)
	parts := strings.Split(filepath.ToSlash(path), "/")

	// 过滤空字符串
	var nonEmpty []string
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}

	// 如果组件数足够，直接返回
	if len(nonEmpty) <= keepComponents {
		return path
	}

	// 保留最后几个组件，前面用 "..." 表示
	kept := nonEmpty[len(nonEmpty)-keepComponents:]
	if filepath.IsAbs(path) {
		return "/" + strings.Join(kept, "/")
	}

	result := strings.Join(kept, string(filepath.Separator))
	return "..." + string(filepath.Separator) + result
}

// FormatPathForDisplay 格式化路径用于终端显示。
// 自动替换主目录为 ~，并使用系统路径分隔符。
func FormatPathForDisplay(path string) string {
	return FormatPath(path,
		WithHomeSymbol("~"),
		WithPathShortening(0),
	)
}

// ============================================================================
// Permissions Formatting
// ============================================================================

// FormatPermissions 格式化文件权限为字符串（rwx 格式）。
// 例如：0644 → "rw-r--r--"，0755 → "rwxr-xr-x"。
func FormatPermissions(mode os.FileMode) string {
	var str string

	// 用户权限
	if mode&0400 != 0 {
		str += "r"
	} else {
		str += "-"
	}
	if mode&0200 != 0 {
		str += "w"
	} else {
		str += "-"
	}
	if mode&0100 != 0 {
		str += "x"
	} else {
		str += "-"
	}

	// 组权限
	if mode&0040 != 0 {
		str += "r"
	} else {
		str += "-"
	}
	if mode&0020 != 0 {
		str += "w"
	} else {
		str += "-"
	}
	if mode&0010 != 0 {
		str += "x"
	} else {
		str += "-"
	}

	// 其他权限
	if mode&0004 != 0 {
		str += "r"
	} else {
		str += "-"
	}
	if mode&0002 != 0 {
		str += "w"
	} else {
		str += "-"
	}
	if mode&0001 != 0 {
		str += "x"
	} else {
		str += "-"
	}

	return str
}

// FormatPermissionsOctal 格式化文件权限为八进制字符串。
// 例如：0644 → "0644"，0755 → "0755"。
func FormatPermissionsOctal(mode os.FileMode) string {
	return fmt.Sprintf("%04o", mode.Perm())
}

// FormatPathPermissions 获取路径的权限并格式化为 rwx 字符串。
func FormatPathPermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}
	return FormatPermissions(info.Mode()), nil
}

// FormatPathPermissionsOctal 获取路径的权限并格式化为八进制字符串。
func FormatPathPermissionsOctal(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}
	return FormatPermissionsOctal(info.Mode()), nil
}

// FormatFileInfo 格式化文件信息为易读字符串。
// 返回格式如："-rw-r--r-- 1.5 KB /path/to/file.txt"
func FormatFileInfo(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 文件类型
	var typeChar string
	switch {
	case info.IsDir():
		typeChar = "d"
	case info.Mode()&os.ModeSymlink != 0:
		typeChar = "l"
	default:
		typeChar = "-"
	}

	// 权限
	perms := FormatPermissions(info.Mode())

	// 大小
	size := FormatFileSize(info.Size())

	// 路径
	displayPath := FormatPathForDisplay(path)

	return fmt.Sprintf("%s%s %s %s", typeChar, perms, size, displayPath), nil
}

// ============================================================================
// File Type Detection
// ============================================================================

// GetFileType 获取文件类型的描述字符串。
// 返回：file, directory, symlink, 或其他特殊类型。
func GetFileType(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return "unknown"
	}

	mode := info.Mode()

	switch {
	case mode&os.ModeSymlink != 0:
		return "symlink"
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			return "char device"
		}
		return "block device"
	case mode&os.ModeNamedPipe != 0:
		return "named pipe"
	case mode&os.ModeSocket != 0:
		return "socket"
	case info.IsDir():
		return "directory"
	default:
		return "file"
	}
}

// FormatFileType 格式化文件类型为单字符表示。
// 返回：- (file), d (directory), l (symlink), 等。
func FormatFileType(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return "?"
	}

	mode := info.Mode()

	switch {
	case mode&os.ModeSymlink != 0:
		return "l"
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			return "c"
		}
		return "b"
	case mode&os.ModeNamedPipe != 0:
		return "p"
	case mode&os.ModeSocket != 0:
		return "s"
	case info.IsDir():
		return "d"
	default:
		return "-"
	}
}
