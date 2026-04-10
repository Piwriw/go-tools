package file

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

// SecureDelete 安全删除文件（覆写后删除）。
// 通过多次覆写文件内容来防止数据恢复。
// passes 是覆写次数，推荐至少 3 次。
func SecureDelete(path string, passes int) error {
	// 检查文件是否存在
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，视为成功
		}
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 只处理文件，不处理目录
	if info.IsDir() {
		return fmt.Errorf("安全删除不支持目录: %s", path)
	}

	fileSize := info.Size()

	// 打开文件用于写入
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 多次覆写
	for i := 0; i < passes; i++ {
		// 生成随机数据或模式数据
		var pattern byte
		switch i % 3 {
		case 0:
			pattern = 0x00 // 全零
		case 1:
			pattern = 0xFF // 全一
		default:
			// 使用随机数据
			buf := make([]byte, 32*1024) // 32KB 缓冲区
			var offset int64
			for offset < fileSize {
				writeSize := int64(len(buf))
				if offset+writeSize > fileSize {
					writeSize = fileSize - offset
				}

				if _, err := rand.Read(buf[:writeSize]); err != nil {
					return fmt.Errorf("生成随机数据失败: %w", err)
				}

				if _, err := file.WriteAt(buf[:writeSize], offset); err != nil {
					return fmt.Errorf("覆写失败: %w", err)
				}

				offset += writeSize
			}

			// 确保数据写入磁盘
			if err := file.Sync(); err != nil {
				return fmt.Errorf("同步数据失败: %w", err)
			}

			continue
		}

		// 使用模式数据覆写
		buf := make([]byte, 32*1024)
		for j := range buf {
			buf[j] = pattern
		}

		var offset int64
		for offset < fileSize {
			writeSize := int64(len(buf))
			if offset+writeSize > fileSize {
				writeSize = fileSize - offset
			}

			if _, err := file.WriteAt(buf[:writeSize], offset); err != nil {
				return fmt.Errorf("覆写失败: %w", err)
			}

			offset += writeSize
		}

		// 确保数据写入磁盘
		if err := file.Sync(); err != nil {
			return fmt.Errorf("同步数据失败: %w", err)
		}
	}

	// 关闭文件后删除
	file.Close()
	return os.Remove(path)
}

// SetExecutable 设置文件为可执行权限。
func SetExecutable(path string) error {
	return os.Chmod(path, 0755)
}

// SetReadOnly 设置文件为只读权限。
func SetReadOnly(path string) error {
	return os.Chmod(path, 0444)
}

// SetReadWrite 设置文件为可读写权限。
func SetReadWrite(path string) error {
	return os.Chmod(path, 0644)
}

// SetPermissions 设置文件权限。
func SetPermissions(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// GetFilePermissions 获取文件权限字符串。
// 返回类似 "rwxr-xr-x" 的权限字符串。
func GetFilePermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	return modeToString(info.Mode().Perm()), nil
}

// modeToString 将 FileMode 转换为权限字符串
func modeToString(mode os.FileMode) string {
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

// GetOctalPermissions 获取文件权限的八进制表示。
// 返回类似 "0644" 的字符串。
func GetOctalPermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	return fmt.Sprintf("%04o", info.Mode().Perm()), nil
}

// SecureDeleteByPattern 按通配符模式安全删除文件。
// pattern 是文件路径模式，支持 * 和 ? 通配符。
// passes 是覆写次数。
func SecureDeleteByPattern(pattern string, passes int) (int, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, fmt.Errorf("解析模式失败: %w", err)
	}

	count := 0
	for _, match := range matches {
		// 检查是否是文件
		info, err := os.Stat(match)
		if err != nil {
			continue
		}

		// 只处理文件
		if info.IsDir() {
			continue
		}

		// 安全删除
		if err := SecureDelete(match, passes); err != nil {
			return count, fmt.Errorf("安全删除文件 %s 失败: %w", match, err)
		}
		count++
	}

	return count, nil
}

// CopyPermissions 复制文件权限。
// 将 src 的权限应用到 dst。
func CopyPermissions(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	return os.Chmod(dst, info.Mode().Perm())
}

// IsExecutable 判断文件是否可执行。
func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	mode := info.Mode()
	return mode&0111 != 0
}

// IsReadOnly 判断文件是否只读。
func IsReadOnly(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	mode := info.Mode().Perm()
	return mode&0222 == 0
}

// MakePrivate 设置文件为私有权限（仅所有者可读写）。
func MakePrivate(path string) error {
	return os.Chmod(path, 0600)
}

// MakePrivateExecutable 设置文件为私有可执行权限。
func MakePrivateExecutable(path string) error {
	return os.Chmod(path, 0700)
}
