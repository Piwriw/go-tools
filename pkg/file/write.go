package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteOptions 写入选项
type WriteOptions struct {
	Perm    os.FileMode // 文件权限，默认 0644
	DirPerm os.FileMode // 目录权限，默认 0755
	Append  bool        // 是否追加模式，默认 false（覆盖）
	Create  bool        // 是否创建文件，默认 true
}

// DefaultWriteOptions 默认写入选项
var DefaultWriteOptions = WriteOptions{
	Perm:    0644,
	DirPerm: 0755,
	Append:  false,
	Create:  true,
}

// WriteToFile 将数据写入指定文件
// 根据提供的文件路径、数据和写入选项，将数据写入文件。如果目录不存在则自动创建，
// 并根据选项决定是否追加写入或覆盖写入。
//
// 参数:
//   - filePath: string 类型，表示要写入的文件路径
//   - data: []byte 类型，表示要写入的字节数据
//   - opts: *WriteOptions 类型，表示写入时的选项配置，如果为 nil 则使用默认配置
//
// 返回值:
//   - error: 如果在创建目录、打开文件或写入数据过程中发生错误，则返回相应的错误信息，否则返回 nil
func WriteToFile(filePath string, data []byte, opts *WriteOptions) error {
	if opts == nil {
		opts = &DefaultWriteOptions
	}

	// 获取文件所在目录
	dir := filepath.Dir(filePath)

	// 创建目录（如果不存在）
	if err := os.MkdirAll(dir, opts.DirPerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 设置文件打开标志
	flags := os.O_WRONLY
	if opts.Create {
		flags |= os.O_CREATE
	}
	if opts.Append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	// 创建并打开文件
	file, err := os.OpenFile(filePath, flags, opts.Perm)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// WriteStringToFile 写入字符串内容的便捷函数
func WriteStringToFile(filePath string, content string, opts *WriteOptions) error {
	return WriteToFile(filePath, []byte(content), opts)
}
