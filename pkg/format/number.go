package format

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// digitsOnlyRegex 预编译的正则表达式，用于校验字符串是否只包含数字
var digitsOnlyRegex = regexp.MustCompile(`^\d+$`)

// IsDigitsOnly 判断字符串是否只包含数字
func IsDigitsOnly(s string) bool {
	return digitsOnlyRegex.MatchString(s)
}

// IsNumber 判断字符串是否为有效数字
// 支持整数、浮点数、科学计数法（如 "123", "-45.67", "1.5e-10"）
// 空字符串、纯空格、NaN、Inf 返回 false
func IsNumber(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	// 排除 NaN 和 Inf
	if strings.EqualFold(s, "nan") ||
		strings.EqualFold(s, "inf") ||
		strings.EqualFold(s, "+inf") ||
		strings.EqualFold(s, "-inf") ||
		strings.EqualFold(s, "infinity") ||
		strings.EqualFold(s, "+infinity") ||
		strings.EqualFold(s, "-infinity") {
		return false
	}
	return true
}

// Float64ToPercentFloat 保留两位小数，返回float64
func Float64ToPercentFloat(f float64) float64 {
	res := f * 100
	return math.Round(res*100) / 100
}
