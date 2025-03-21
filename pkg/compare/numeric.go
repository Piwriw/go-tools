package compare

import (
	"strconv"
)

// IsNumeric 判断字符串 s 是否可以转换为 T 类型，并返回转换后的值
func IsNumeric[T int | float64](s string) (bool, T) {
	// 尝试转换为整数
	if intValue, err := strconv.Atoi(s); err == nil {
		return convertToType[T](intValue)
	}

	// 尝试转换为浮点数
	if floatValue, err := strconv.ParseFloat(s, 64); err == nil {
		return convertToType[T](floatValue)
	}

	// 转换失败，返回 false 和零值
	var zero T
	return false, zero
}

// convertToType 将值转换为 T 类型
func convertToType[T int | float64](value any) (bool, T) {
	var zero T
	switch any(zero).(type) {
	case int:
		// 如果 T 是 int 类型，直接转换 int 值
		if v, ok := value.(int); ok {
			return true, T(v)
		}
	case float64:
		// 如果 T 是 float64 类型，转换 int 到 float64
		if v, ok := value.(int); ok {
			return true, T(float64(v)) // 将 int 转换为 float64
		} else if v, ok := value.(float64); ok {
			return true, T(v)
		}
	}
	return false, zero
}
