package format

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//	DirSizeWithDU 获取目录大小（单位：字节）
//	获取目录大小（单位：字节）
//	使用du命令获取目录大小，支持macOS和Linux
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
