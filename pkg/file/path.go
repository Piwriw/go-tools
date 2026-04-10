package file

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFileExt 获取文件扩展名（包含点号）。
// 例如："test.txt" → ".txt"，"archive.tar.gz" → ".gz"。
// 如果没有扩展名，返回空字符串。
func GetFileExt(path string) string {
	return filepath.Ext(path)
}

// GetFileName 获取文件名（不含扩展名）。
// 例如："/path/to/test.txt" → "test"，"archive.tar.gz" → "archive.tar"。
// 如果路径以 / 结尾，返回空字符串。
func GetFileName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	if len(ext) > len(base) {
		return base // 扩展名比基本名长（不可能，但安全检查）
	}
	return base[:len(base)-len(ext)]
}

// GetFileNameWithExt 获取文件名（含扩展名）。
// 例如："/path/to/test.txt" → "test.txt"。
func GetFileNameWithExt(path string) string {
	return filepath.Base(path)
}

// GetDirName 获取文件所在目录的名称。
// 例如："/path/to/file.txt" → "to"。
// 返回最后一个目录组件的名称。
func GetDirName(path string) string {
	dir := filepath.Dir(path)
	return filepath.Base(dir)
}

// GetDirPath 获取文件所在目录的完整路径。
// 例如："/path/to/file.txt" → "/path/to"。
func GetDirPath(path string) string {
	return filepath.Dir(path)
}

// JoinPath 安全地拼接路径组件。
// 自动跳过空字符串组件。
func JoinPath(parts ...string) string {
	// 过滤空组件
	var nonEmpty []string
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}

	return filepath.Join(nonEmpty...)
}

// AbsPath 获取绝对路径（含错误处理）。
// 如果路径已经是绝对路径，直接返回规范化后的路径。
// 如果是相对路径，基于当前工作目录转换为绝对路径。
func AbsPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %w", err)
	}

	return abs, nil
}

// MustAbsPath 获取绝对路径，失败时返回空字符串。
func MustAbsPath(path string) string {
	abs, err := AbsPath(path)
	if err != nil {
		return ""
	}
	return abs
}

// ResolveHome 展开 ~ 和 ~user 路径。
// 例如："~/Documents" → "/home/user/Documents"。
// 如果路径不包含 ~，原样返回。
func ResolveHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	// 获取用户主目录
	var homeDir string
	if path == "~" || strings.HasPrefix(path, "~/") {
		// 当前用户主目录
		homeDir = os.Getenv("HOME")
		if homeDir == "" && runtime.GOOS == "windows" {
			homeDir = os.Getenv("USERPROFILE")
		}
		if homeDir == "" {
			return path // 无法获取主目录，原样返回
		}
	} else {
		// 其他用户主目录（~user）
		// 简化处理：Unix 系统通常在 /home/user 或 /Users/user
		parts := strings.SplitN(path, "/", 2)
		if len(parts) < 2 {
			return path
		}
		username := parts[0][1:] // 去掉 ~
		if runtime.GOOS == "darwin" {
			homeDir = "/Users/" + username
		} else {
			homeDir = "/home/" + username
		}
	}

	// 替换 ~ 为主目录
	if path == "~" {
		return homeDir
	}

	// ~/path 或 ~user/path
	remainingPath := strings.TrimPrefix(path, "~")
	if strings.HasPrefix(remainingPath, "/") {
		remainingPath = remainingPath[1:]
	}

	return filepath.Join(homeDir, remainingPath)
}

// IsAbs 判断路径是否为绝对路径。
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// RelPath 获取相对于基准路径的相对路径。
// 如果无法转换为相对路径，返回原路径。
func RelPath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", fmt.Errorf("计算相对路径失败: %w", err)
	}
	return rel, nil
}

// CleanPath 规范化路径（清理冗余分隔符、.、.. 等）。
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// SplitPath 分割路径为 (目录, 文件) 元组。
// 例如："/path/to/file.txt" → ("/path/to", "file.txt")。
func SplitPath(path string) (dir, file string) {
	return filepath.Split(path)
}

// GetCommonPath 获取多个路径的公共前缀路径。
func GetCommonPath(paths ...string) string {
	if len(paths) == 0 {
		return ""
	}
	if len(paths) == 1 {
		return CleanPath(paths[0])
	}

	// 转换为绝对路径并规范化
	var absPaths []string
	for _, p := range paths {
		abs, err := AbsPath(p)
		if err != nil {
			continue
		}
		absPaths = append(absPaths, CleanPath(abs))
	}

	if len(absPaths) == 0 {
		return ""
	}

	// 以第一个路径为基准
	common := absPaths[0]

	for _, p := range absPaths[1:] {
		common = findCommonPrefix(common, p)
		if common == "" {
			break
		}
	}

	return common
}

// findCommonPrefix 查找两个路径的公共前缀
func findCommonPrefix(a, b string) string {
	// 确保使用系统的路径分隔符
	a = filepath.ToSlash(a)
	b = filepath.ToSlash(b)

	aParts := strings.Split(a, "/")
	bParts := strings.Split(b, "/")

	var common []string
	minLen := len(aParts)
	if len(bParts) < minLen {
		minLen = len(bParts)
	}

	for i := 0; i < minLen; i++ {
		if aParts[i] == bParts[i] {
			common = append(common, aParts[i])
		} else {
			break
		}
	}

	if len(common) == 0 {
		return ""
	}

	return filepath.Join(common...)
}

// HasParent 判断 child 路径是否在 parent 路径下。
func HasParent(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	// 相对路径不能以 .. 开头（表示在父目录之外）
	return !strings.HasPrefix(rel, "..")
}

// ChangeExt 更改文件扩展名。
// 例如："/path/to/file.txt", ".log" → "/path/to/file.log"。
func ChangeExt(path, newExt string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path + newExt
	}
	return path[:len(path)-len(ext)] + newExt
}

// EnsureExt 确保路径以指定扩展名结尾。
// 如果已有扩展名且不匹配，则替换。
func EnsureExt(path, ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	if filepath.Ext(path) == "" {
		return path + ext
	}

	return ChangeExt(path, ext)
}
