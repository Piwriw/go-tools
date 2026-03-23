package format

import (
	"encoding/json"
	"fmt"
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

// ToFloat64 将任意类型转换为 float64
// 支持的类型包括：float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, json.Number, string, bool, nil
// 对于 string 类型，会尝试将其转换为 float64
// 对于 bool 类型，true 转换为 1，false 转换为 0
// 对于 nil，返回 0 和错误
// 对于其他类型，返回 0 和错误
// 对于 uint64 类型，如果值超过 float64 能精准表示的最大整数，返回错误
// 对于 json.Number 类型，会尝试将其转换为 float64
// 对于其他类型，返回 0 和错误
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		const maxExactFloat64 = 1<<53 - 1 // 9007199254740991 (float64 能精准表示的最大整数)
		if v > maxExactFloat64 {
			return 0, fmt.Errorf("value too large to convert without precision loss: %d", v)
		}
		return float64(v), nil
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert json.Number to float64: %v", err)
		}
		return f, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to float64: %v", err)
		}
		return f, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, fmt.Errorf("cannot convert nil to float64")
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
