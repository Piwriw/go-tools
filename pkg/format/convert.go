package format

import (
	"encoding/json"
	"math"
	"strconv"
)

// ToInt 将任意类型转换为 int。
// 转换失败时返回 defaults[0]，未提供默认值时返回 0。
// 提供多个默认值时仅使用第一个。
func ToInt(value any, defaults ...int) int {
	fallback := defaultInt(defaults...)

	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		if v < int64(math.MinInt) || v > int64(math.MaxInt) {
			return fallback
		}
		return int(v)
	case uint:
		if uint64(v) > uint64(math.MaxInt) {
			return fallback
		}
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		if uint64(v) > uint64(math.MaxInt) {
			return fallback
		}
		return int(v)
	case uint64:
		if v > uint64(math.MaxInt) {
			return fallback
		}
		return int(v)
	case float32:
		return floatToInt(float64(v), fallback)
	case float64:
		return floatToInt(v, fallback)
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return fallback
		}
		return floatToInt(f, fallback)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fallback
		}
		return floatToInt(f, fallback)
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return fallback
	}
}

func defaultInt(defaults ...int) int {
	if len(defaults) == 0 {
		return 0
	}
	return defaults[0]
}

func floatToInt(v float64, fallback int) int {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fallback
	}
	if math.Trunc(v) != v {
		return fallback
	}
	if v < float64(math.MinInt) || v > float64(math.MaxInt) {
		return fallback
	}
	return int(v)
}
