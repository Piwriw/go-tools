package file

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// ReadLines 按行读取文件内容，返回字符串切片。
// 每行包含换行符，使用时可根据需要 TrimSpace。
// 空文件返回空切片而非 nil。
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 确保空文件返回空切片而非 nil
	if lines == nil {
		lines = []string{}
	}

	return lines, nil
}

// ReadFileAsString 读取文件内容并返回字符串。
// os.ReadFile 返回 []byte，此函数简化了类型转换。
func ReadFileAsString(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(data), nil
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts int           // 最大尝试次数，必须 >= 1
	WaitTime    time.Duration // 每次重试前的等待时间
}

// DefaultRetryConfig 默认重试配置：最多重试 3 次，每次间隔 100ms
var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 3,
	WaitTime:    100 * time.Millisecond,
}

// ReadFileWithRetry 带重试机制的文件读取。
// 适用于处理临时文件锁定、网络驱动器等可能出现瞬时错误的场景。
// 如果所有重试都失败，返回最后一次的错误。
func ReadFileWithRetry(filePath string, config *RetryConfig) ([]byte, error) {
	if config == nil {
		config = &DefaultRetryConfig
	}
	if config.MaxAttempts < 1 {
		config.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		data, err := os.ReadFile(filePath)
		if err == nil {
			return data, nil
		}

		lastErr = err

		// 最后一次尝试后不再等待
		if attempt < config.MaxAttempts-1 {
			time.Sleep(config.WaitTime)
		}
	}

	return nil, fmt.Errorf("读取文件失败（已重试 %d 次）: %w", config.MaxAttempts, lastErr)
}

// ReadFileAsStringWithRetry 带重试机制的文件读取（返回字符串）。
func ReadFileAsStringWithRetry(filePath string, config *RetryConfig) (string, error) {
	data, err := ReadFileWithRetry(filePath, config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
