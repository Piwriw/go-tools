package file

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// FileExists 检查文件或目录是否存在。
// 路径不存在时返回 false（不报错），这比 os.Stat 更方便。
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir 判断路径是否为目录。
// 路径不存在时返回 false。
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 判断路径是否为文件。
// 路径不存在时返回 false。
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsEmptyDir 判断目录是否为空。
// 空目录是指不包含任何文件或子目录。
// 路径不存在或不是目录时返回 false。
func IsEmptyDir(dir string) (bool, error) {
	// 打开目录
	f, err := os.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("打开目录失败: %w", err)
	}
	defer f.Close()

	// 读取目录条目（最多读1个）
	_, err = f.Readdir(1)
	if err == io.EOF {
		// 没有读取到任何条目，目录为空
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("读取目录失败: %w", err)
	}

	// 读取到至少一个条目，目录不为空
	return false, nil
}

// GetFileSize 获取文件大小（字节）。
// 如果是目录或路径不存在，返回错误。
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if info.IsDir() {
		return 0, fmt.Errorf("路径是目录，不是文件: %s", path)
	}

	return info.Size(), nil
}

// GetFileModTime 获取文件修改时间。
func GetFileModTime(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return info.ModTime().Unix(), nil
}

// HashType 哈希类型
type HashType string

const (
	HashMD5    HashType = "md5"
	HashSHA256 HashType = "sha256"
)

// GetFileHash 计算文件的哈希值。
// 支持 MD5 和 SHA256 算法。
// 返回十六进制编码的哈希字符串。
func GetFileHash(path string, hashType HashType) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var hashHasher interface {
		io.Writer
		Sum([]byte) []byte
	}

	switch hashType {
	case HashMD5:
		hashHasher = md5.New()
	case HashSHA256:
		hashHasher = sha256.New()
	default:
		return "", fmt.Errorf("不支持的哈希类型: %s", hashType)
	}

	if _, err := io.Copy(hashHasher, file); err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return hex.EncodeToString(hashHasher.Sum(nil)), nil
}

// GetFileMD5 计算文件的 MD5 哈希值。
func GetFileMD5(path string) (string, error) {
	return GetFileHash(path, HashMD5)
}

// GetFileSHA256 计算文件的 SHA256 哈希值。
func GetFileSHA256(path string) (string, error) {
	return GetFileHash(path, HashSHA256)
}

// FilesEqual 比较两个文件的内容是否相同（通过 MD5）。
// 只比较内容，不比较元数据（如修改时间）。
func FilesEqual(path1, path2 string) (bool, error) {
	hash1, err := GetFileMD5(path1)
	if err != nil {
		return false, fmt.Errorf("获取文件1哈希失败: %w", err)
	}

	hash2, err := GetFileMD5(path2)
	if err != nil {
		return false, fmt.Errorf("获取文件2哈希失败: %w", err)
	}

	return hash1 == hash2, nil
}

// GetFileInfo 获取文件的完整信息。
type FileInfoDetail struct {
	Path    string // 文件路径
	Size    int64  // 文件大小（字节）
	ModTime int64  // 修改时间（Unix 时间戳）
	IsDir   bool   // 是否为目录
	Mode    string // 权限模式（八进制字符串，如 "0644"）
}

// GetFileInfo 获取文件的详细信息。
func GetFileInfo(path string) (*FileInfoDetail, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return &FileInfoDetail{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime().Unix(),
		IsDir:   info.IsDir(),
		Mode:    fmt.Sprintf("%04o", info.Mode().Perm()),
	}, nil
}
