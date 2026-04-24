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

// --- Permission bit operations ---

// HasPermission checks whether a specific permission bit is set on a file.
// Use os.FileMode bit masks like 0400 (owner read), 0200 (owner write), 0100 (owner exec),
// 0040 (group read), 0004 (other read), etc.
func HasPermission(path string, bit os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return info.Mode().Perm()&bit != 0, nil
}

// AddPermission adds permission bits to a file (preserving existing bits).
// Example: AddPermission("file.txt", 0111) adds execute for all.
func AddPermission(path string, bits os.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	return os.Chmod(path, info.Mode().Perm()|bits)
}

// RemovePermission removes permission bits from a file (preserving existing bits).
// Example: RemovePermission("file.txt", 0222) removes all write bits.
func RemovePermission(path string, bits os.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	return os.Chmod(path, info.Mode().Perm()&^bits)
}

// --- Readable / Writable checks ---

// IsReadable checks whether the owner has read permission.
func IsReadable(path string) bool {
	ok, err := HasPermission(path, 0400)
	return err == nil && ok
}

// IsWritable checks whether the owner has write permission.
func IsWritable(path string) bool {
	ok, err := HasPermission(path, 0200)
	return err == nil && ok
}

// --- Recursive operations ---

// SetPermissionsRecursive applies a permission mode to all files and/or directories
// under the given root path (inclusive).
// applyTo controls what is affected: "file", "dir", or "all".
func SetPermissionsRecursive(root string, mode os.FileMode, applyTo string) error {
	validTargets := map[string]bool{"file": true, "dir": true, "all": true}
	if !validTargets[applyTo] {
		return fmt.Errorf("invalid applyTo value %q: must be \"file\", \"dir\", or \"all\"", applyTo)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}

		shouldApply := applyTo == "all" ||
			(applyTo == "file" && !info.IsDir()) ||
			(applyTo == "dir" && info.IsDir())

		if shouldApply {
			if chErr := os.Chmod(path, mode); chErr != nil {
				return fmt.Errorf("failed to chmod %s: %w", path, chErr)
			}
		}
		return nil
	})
}

// SetDefaultPermissions sets files under root to filePerm and directories to dirPerm.
// This is the common "web-safe" pattern (0644 for files, 0755 for dirs).
func SetDefaultPermissions(root string, filePerm, dirPerm os.FileMode) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return os.Chmod(path, dirPerm)
		}
		return os.Chmod(path, filePerm)
	})
}

// CopyPermissionsRecursive copies file permissions from src to dst directory trees.
// Files and directories are matched by their relative path under srcRoot/dstRoot.
func CopyPermissionsRecursive(srcRoot, dstRoot string) error {
	return filepath.Walk(srcRoot, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(srcRoot, srcPath)
		if err != nil {
			return nil
		}
		dstPath := filepath.Join(dstRoot, rel)

		if !FileExists(dstPath) {
			return nil // skip if dst counterpart doesn't exist
		}

		return os.Chmod(dstPath, info.Mode().Perm())
	})
}

// AddPermissionRecursive adds permission bits to all files and/or directories
// under the given root path. applyTo follows the same rules as SetPermissionsRecursive.
func AddPermissionRecursive(root string, bits os.FileMode, applyTo string) error {
	validTargets := map[string]bool{"file": true, "dir": true, "all": true}
	if !validTargets[applyTo] {
		return fmt.Errorf("invalid applyTo value %q: must be \"file\", \"dir\", or \"all\"", applyTo)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		shouldApply := applyTo == "all" ||
			(applyTo == "file" && !info.IsDir()) ||
			(applyTo == "dir" && info.IsDir())

		if shouldApply {
			newMode := info.Mode().Perm() | bits
			if chErr := os.Chmod(path, newMode); chErr != nil {
				return fmt.Errorf("failed to chmod %s: %w", path, chErr)
			}
		}
		return nil
	})
}

// RemovePermissionRecursive removes permission bits from all files and/or directories
// under the given root path. applyTo follows the same rules as SetPermissionsRecursive.
func RemovePermissionRecursive(root string, bits os.FileMode, applyTo string) error {
	validTargets := map[string]bool{"file": true, "dir": true, "all": true}
	if !validTargets[applyTo] {
		return fmt.Errorf("invalid applyTo value %q: must be \"file\", \"dir\", or \"all\"", applyTo)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		shouldApply := applyTo == "all" ||
			(applyTo == "file" && !info.IsDir()) ||
			(applyTo == "dir" && info.IsDir())

		if shouldApply {
			newMode := info.Mode().Perm() &^ bits
			if chErr := os.Chmod(path, newMode); chErr != nil {
				return fmt.Errorf("failed to chmod %s: %w", path, chErr)
			}
		}
		return nil
	})
}
