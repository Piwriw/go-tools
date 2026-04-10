package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDeleted 确保文件或目录被删除。
// 如果路径不存在，返回 nil（不报错），这与 os.Remove 不同。
// 如果是目录，递归删除其所有内容。
func EnsureDeleted(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 路径不存在，直接返回成功
		return nil
	}
	if err != nil {
		return fmt.Errorf("检查路径失败: %w", err)
	}

	// 使用 RemoveAll 处理文件和目录
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	return nil
}

// ForceDeleteDir 强制删除目录，处理只读文件。
// 在删除前会尝试修改文件权限为可写，然后再删除。
// 适用于 Windows 或权限受限的场景。
func ForceDeleteDir(dir string) error {
	// 检查目录是否存在
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("检查目录失败: %w", err)
	}

	// 如果不是目录，尝试直接删除
	if !info.IsDir() {
		return EnsureDeleted(dir)
	}

	// 遍历目录，修改权限后删除
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 修改文件权限为可写
		if !info.IsDir() {
			// 添加写权限（用户、组、其他）
			mode := info.Mode()
			newMode := mode | 0200 // 添加用户写权限
			if err := os.Chmod(path, newMode); err != nil {
				// 权限修改失败不是致命错误，继续尝试删除
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("修改权限失败: %w", err)
	}

	// 删除整个目录
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("删除目录失败: %w", err)
	}

	return nil
}

// DeleteByPattern 按通配符模式删除文件。
// pattern 是文件路径模式，支持 * 和 ? 通配符。
// 只删除匹配的文件，不删除目录。
// 返回删除的文件数量和可能的错误。
func DeleteByPattern(pattern string) (int, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, fmt.Errorf("解析模式失败: %w", err)
	}

	count := 0
	for _, match := range matches {
		// 检查是否是文件
		info, err := os.Stat(match)
		if err != nil {
			continue // 跳过无法访问的文件
		}

		// 只删除文件，跳过目录
		if info.IsDir() {
			continue
		}

		// 尝试删除文件
		if err := os.Remove(match); err != nil {
			return count, fmt.Errorf("删除文件 %s 失败: %w", match, err)
		}
		count++
	}

	return count, nil
}

// DeleteByPatternInDir 在指定目录中按通配符模式删除文件。
// dir 是要搜索的目录，pattern 是文件名模式（不含路径）。
// 递归搜索子目录。
func DeleteByPatternInDir(dir, pattern string) (int, error) {
	fullPattern := filepath.Join(dir, pattern)
	return DeleteByPattern(fullPattern)
}
