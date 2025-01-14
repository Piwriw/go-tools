package format

import (
	"encoding/json"
	"fmt"
	"strings"
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
		return fmt.Sprintf("%f", v) // 转换为字符串并保留浮点数格式
	case int:
		return fmt.Sprintf("%d", v) // 转换为整数字符串
	case json.Number:
		return v.String() // json.Number 自带 String 方法
	case string:
		return v // 字符串直接返回
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null" // nil 转换为 "null"
	default:
		return fmt.Sprintf("unknown type: %T", v) // 未知类型
	}
}
