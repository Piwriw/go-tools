package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FindFilesByExt 按扩展名递归查找文件。
// ext 应包含点号（如 ".txt"），比较不区分大小写。
// 返回匹配文件的完整路径列表。
func FindFilesByExt(dir string, ext string) ([]string, error) {
	var matches []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 跳过无法访问的文件/目录
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查扩展名（不区分大小写）
		if filepath.Ext(path) == ext {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindFilesByExts 按多个扩展名递归查找文件。
// 任意一个扩展名匹配即返回该文件。
func FindFilesByExts(dir string, exts []string) ([]string, error) {
	var matches []string

	// 构建扩展名映射（不区分大小写）
	extMap := make(map[string]bool)
	for _, ext := range exts {
		extMap[ext] = true
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if extMap[filepath.Ext(path)] {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindFilesByName 按名称模式查找文件。
// pattern 支持 filepath.Match 的通配符语法（*, ?, []）。
// 返回匹配文件的完整路径列表。
func FindFilesByName(dir string, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 获取文件名
		filename := info.Name()

		// 检查是否匹配模式
		matched, err := filepath.Match(pattern, filename)
		if err != nil {
			return fmt.Errorf("模式匹配失败: %w", err)
		}

		if matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindRecentFiles 查找指定时间内修改过的文件。
// within 指定时间范围（如最近 24 小时：time.Hour * 24）。
// 返回修改时间在范围内的文件路径列表。
func FindRecentFiles(dir string, within time.Duration) ([]string, error) {
	var matches []string
	cutoff := time.Now().Add(-within)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查修改时间
		if info.ModTime().After(cutoff) {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindFilesModifiedAfter 查找在指定时间之后修改过的文件。
func FindFilesModifiedAfter(dir string, after time.Time) ([]string, error) {
	var matches []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().After(after) {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindLargeFiles 查找超过指定大小的文件。
// minSize 是最小文件大小（字节）。
// 返回文件路径和文件大小的列表。
func FindLargeFiles(dir string, minSize int64) ([]FileInfo, error) {
	var matches []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件大小
		if info.Size() >= minSize {
			matches = append(matches, FileInfo{
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FindFilesBySizeRange 查找在指定大小范围内的文件。
// minSize 和 maxSize 是文件大小的上下限（字节）。
// maxSize 为 0 表示无上限。
func FindFilesBySizeRange(dir string, minSize, maxSize int64) ([]FileInfo, error) {
	var matches []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		size := info.Size()

		// 检查大小范围
		if size >= minSize && (maxSize == 0 || size <= maxSize) {
			matches = append(matches, FileInfo{
				Path:    path,
				Size:    size,
				ModTime: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// FileInfo 文件信息
type FileInfo struct {
	Path    string    // 文件完整路径
	Size    int64     // 文件大小（字节）
	ModTime time.Time // 修改时间
}

// FindFilesByPredicate 使用自定义谓词函数查找文件。
// predicate 返回 true 表示该文件符合条件。
func FindFilesByPredicate(dir string, predicate func(path string, info os.FileInfo) bool) ([]string, error) {
	var matches []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if predicate(path, info) {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return matches, nil
}

// CountFilesByExt 统计各扩展名文件的数量。
// 返回扩展名到文件数量的映射。
func CountFilesByExt(dir string) (map[string]int, error) {
	counts := make(map[string]int)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		counts[ext]++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return counts, nil
}
