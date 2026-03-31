package format

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToInt 将任意类型转换为 int。
// 转换失败时返回 defaults[0]，未提供默认值时返回 0。
// 提供多个默认值时仅使用第一个。
func ToInt(value any, defaults ...int) int {
	result, err := MustToInt(value)
	if err == nil {
		return result
	}

	fallback := 0
	if len(defaults) > 0 {
		fallback = defaults[0]
	}
	return fallback
}

// MustToInt 将任意类型转换为 int。
// 转换失败时返回 error。
func MustToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		if v < int64(math.MinInt) || v > int64(math.MaxInt) {
			return 0, fmt.Errorf("cannot convert %T to int: out of range", value)
		}
		return int(v), nil
	case uint:
		if uint64(v) > uint64(math.MaxInt) {
			return 0, fmt.Errorf("cannot convert %T to int: out of range", value)
		}
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		if uint64(v) > uint64(math.MaxInt) {
			return 0, fmt.Errorf("cannot convert %T to int: out of range", value)
		}
		return int(v), nil
	case uint64:
		if v > uint64(math.MaxInt) {
			return 0, fmt.Errorf("cannot convert %T to int: out of range", value)
		}
		return int(v), nil
	case float32:
		return floatToInt(float64(v))
	case float64:
		return floatToInt(v)
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert %T to int: %w", value, err)
		}
		return floatToInt(f)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %T to int: %w", value, err)
		}
		return floatToInt(f)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// floatToInt 将浮点数转换为 int。
// 仅接受有限且无小数部分的值，超出 int 范围时返回 error。
func floatToInt(v float64) (int, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("cannot convert float64 to int: invalid number")
	}
	if math.Trunc(v) != v {
		return 0, fmt.Errorf("cannot convert float64 to int: not an integer")
	}
	if v < float64(math.MinInt) || v > float64(math.MaxInt) {
		return 0, fmt.Errorf("cannot convert float64 to int: out of range")
	}
	return int(v), nil
}

// ToFloat64 将任意类型转换为 float64。
// 转换失败时返回 defaults[0]，未提供默认值时返回 0。
// 提供多个默认值时仅使用第一个。
func ToFloat64(value any, defaults ...float64) float64 {
	result, err := MustToFloat64(value)
	if err == nil {
		return result
	}

	fallback := 0.0
	if len(defaults) > 0 {
		fallback = defaults[0]
	}
	return fallback
}

// MustToFloat64 将任意类型转换为 float64。
// 转换失败时返回 error。
func MustToFloat64(value any) (float64, error) {
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
		const maxExactFloat64 = 1<<53 - 1
		if v > maxExactFloat64 {
			return 0, fmt.Errorf("cannot convert %T to float64: value too large", value)
		}
		return float64(v), nil
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert %T to float64: %w", value, err)
		}
		return f, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %T to float64: %w", value, err)
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
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToString 将任意类型转换为 string。
// 转换失败时返回 defaults[0]，未提供默认值时返回空字符串。
// 提供多个默认值时仅使用第一个。
func ToString(value any, defaults ...string) string {
	result, err := MustToString(value)
	if err == nil {
		return result
	}

	fallback := ""
	if len(defaults) > 0 {
		fallback = defaults[0]
	}
	return fallback
}

// MustToString 将任意类型转换为 string。
// 转换失败时返回 error。
func MustToString(value any) (string, error) {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%f", v), nil
	case float32:
		return fmt.Sprintf("%f", v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case json.Number:
		return v.String(), nil
	case string:
		return v, nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case nil:
		return "null", nil
	case time.Time:
		return v.Format(time.RFC3339), nil
	case []byte:
		return string(v), nil
	default:
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return "", fmt.Errorf("cannot convert %T to string", value)
			}
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct || val.Kind() == reflect.Map || val.Kind() == reflect.Slice {
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("cannot convert %T to string: %w", value, err)
			}
			return string(bytes), nil
		}
		if stringer, ok := v.(fmt.Stringer); ok {
			return stringer.String(), nil
		}
		return "", fmt.Errorf("cannot convert %T to string", value)
	}
}

// ToBool 将任意类型转换为 bool。
// 转换失败时返回 defaults[0]，未提供默认值时返回 false。
// 提供多个默认值时仅使用第一个。
func ToBool(value any, defaults ...bool) bool {
	result, err := MustToBool(value)
	if err == nil {
		return result
	}

	fallback := false
	if len(defaults) > 0 {
		fallback = defaults[0]
	}
	return fallback
}

// MustToBool 将任意类型转换为 bool。
// 转换失败时返回 error。
func MustToBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case int8:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case int16:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case int32:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case int64:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case uint:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case uint8:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case uint16:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case uint32:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case uint64:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case float32:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case float64:
		if v == 0 {
			return false, nil
		}
		if v == 1 {
			return true, nil
		}
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return false, fmt.Errorf("cannot convert %T to bool: %w", value, err)
		}
		if f == 0 {
			return false, nil
		}
		if f == 1 {
			return true, nil
		}
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		switch s {
		case "true", "1", "yes", "y", "on":
			return true, nil
		case "false", "0", "no", "n", "off":
			return false, nil
		}
	case nil:
		return false, fmt.Errorf("cannot convert nil to bool")
	}
	return false, fmt.Errorf("cannot convert %T to bool", value)
}
