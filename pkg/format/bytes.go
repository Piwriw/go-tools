package format

import (
	"fmt"
	"math"
	"strings"
)

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// BytesToMBExact 将字节（Byte）转换为MB，保留小数
func BytesToMBExact[T number](bytes T) float64 {
	return float64(bytes) / (1024 * 1024)
}

// BytesToGBExact 将字节（Byte）转换为GB，保留小数
func BytesToGBExact[T number](bytes T) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

// BytesToMB 将字节（Byte）转换为MB，保留两位小数，返回与输入相同类型
func BytesToMB[T number](bytes T) T {
	mb := float64(bytes) / (1024 * 1024)
	return T(math.Round(mb*100) / 100)
}

// BytesToGB 将字节转换为GB，保留两位小数，返回与输入相同类型
func BytesToGB[T number](bytes T) T {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return T(math.Round(gb*100) / 100)
}

// KBToMB 将KB转换为MB，保留两位小数
func KBToMB(kb float64) float64 {
	mb := kb / 1024
	return math.Round(mb*100) / 100
}

// KBToGB 将KB转换为GB，保留两位小数
func KBToGB(kb float64) float64 {
	gb := kb / (1024 * 1024)
	return math.Round(gb*100) / 100
}

// BytesToKB 将字节转换为KB，保留两位小数
func BytesToKB[T number](bytes T) float64 {
	return math.Round(float64(bytes)/1024*100) / 100
}

// BytesToTB 将字节转换为TB，保留两位小数
func BytesToTB[T number](bytes T) float64 {
	return math.Round(float64(bytes)/(1024*1024*1024*1024)*100) / 100
}

// MBToBytes 将MB转换为字节数
func MBToBytes(mb float64) uint64 {
	return uint64(mb * 1024 * 1024)
}

// GBToBytes 将GB转换为字节数
func GBToBytes(gb float64) uint64 {
	return uint64(gb * 1024 * 1024 * 1024)
}

// TBToBytes 将TB转换为字节数
func TBToBytes(tb float64) uint64 {
	return uint64(tb * 1024 * 1024 * 1024 * 1024)
}

// HumanReadableBytes 将字节数格式化为人类可读的字符串，如 "1.50 MB"
// 采用二进制单位（1024进制）
func HumanReadableBytes[T number](bytes T) string {
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

// ParseHumanReadableBytes 将人类可读的字节字符串（如 "1.5MB"、"2 GB"、"1024KB"）解析为字节数
// 支持的单位：B, KB, MB, GB, TB, PB（不区分大小写，支持有无空格）
func ParseHumanReadableBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty input")
	}

	// 提取数字部分
	var numStr string
	var i int
	for i = 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || c == '.' || c == '-' || c == '+' {
			numStr += string(c)
		} else {
			break
		}
	}

	if numStr == "" {
		return 0, fmt.Errorf("no number found in %q", s)
	}

	var val float64
	if _, err := fmt.Sscanf(numStr, "%f", &val); err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", numStr, err)
	}

	if val < 0 {
		return 0, fmt.Errorf("negative value not allowed: %f", val)
	}

	unit := strings.TrimSpace(s[i:])
	unit = strings.ToUpper(unit)

	multiplier := 1.0
	switch unit {
	case "", "B":
		multiplier = 1
	case "KB", "K", "KIB":
		multiplier = 1024
	case "MB", "M", "MIB":
		multiplier = 1024 * 1024
	case "GB", "G", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "TB", "T", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "PB", "P", "PIB":
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown unit %q", unit)
	}

	return uint64(val * multiplier), nil
}
