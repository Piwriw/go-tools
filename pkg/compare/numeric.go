package compare

import (
	"strconv"
)

// Numeric 支持的数字类型
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// IsNumeric 判断字符串 s 是否可以转换为 T 类型，并返回转换后的值
// 各类型行为：
//   - int/uint 系列：只接受整数格式
//   - float32/float64：接受整数和浮点数格式
func IsNumeric[T Numeric](s string) (bool, T) {
	var zero T

	switch any(zero).(type) {
	case int:
		if i, err := strconv.Atoi(s); err == nil {
			return true, T(i)
		}
	case int8:
		if i, err := strconv.ParseInt(s, 10, 8); err == nil {
			return true, T(i)
		}
	case int16:
		if i, err := strconv.ParseInt(s, 10, 16); err == nil {
			return true, T(i)
		}
	case int32:
		if i, err := strconv.ParseInt(s, 10, 32); err == nil {
			return true, T(i)
		}
	case int64:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return true, T(i)
		}
	case uint:
		if i, err := strconv.ParseUint(s, 10, 64); err == nil {
			return true, T(i)
		}
	case uint8:
		if i, err := strconv.ParseUint(s, 10, 8); err == nil {
			return true, T(i)
		}
	case uint16:
		if i, err := strconv.ParseUint(s, 10, 16); err == nil {
			return true, T(i)
		}
	case uint32:
		if i, err := strconv.ParseUint(s, 10, 32); err == nil {
			return true, T(i)
		}
	case uint64:
		if i, err := strconv.ParseUint(s, 10, 64); err == nil {
			return true, T(i)
		}
	case float32:
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			return true, T(f)
		}
	case float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return true, T(f)
		}
	}

	return false, zero
}
