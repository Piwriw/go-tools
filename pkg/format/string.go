package format

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func ConvertMapKeysToUpper(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		upperKey := strings.ToUpper(key)
		result[upperKey] = value
	}

	return result
}

func CovertToUpperSingle(str string) string {
	return strings.ToUpper(str)
}

func CovertToUpperMultiple(strArr ...string) []string {
	for i, s := range strArr {
		strArr[i] = strings.ToUpper(s)
	}
	return strArr
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%f", v)
	case float32:
		return fmt.Sprintf("%f", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case json.Number:
		return v.String()
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	case time.Time:
		return v.Format(time.RFC3339) // 格式化时间
	case []byte:
		return string(v) // 将 []byte 转换为字符串
	default:
		// 处理指针类型
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem() // 解引用指针
		}
		// 处理结构体、map、slice 等
		if val.Kind() == reflect.Struct {
			bytes, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("unknown type: %T", v)
			}
			return string(bytes)
		}
		if val.Kind() == reflect.Map || val.Kind() == reflect.Slice {
			bytes, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("unknown type: %T", v)
			}
			return string(bytes)
		}
		// 处理实现了 Stringer 接口的类型
		if stringer, ok := v.(fmt.Stringer); ok {
			return stringer.String()
		}
		return fmt.Sprintf("unknown type: %T", v)
	}
}
