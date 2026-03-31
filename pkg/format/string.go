package format

import (
	"strings"
)

// ConvertMapKeysToUpper 将 map 的所有 key 转换为大写
func ConvertMapKeysToUpper(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		upperKey := strings.ToUpper(key)
		result[upperKey] = value
	}

	return result
}

func ConvertToUpperSingle(str string) string {
	return strings.ToUpper(str)
}

func ConvertToUpperMultiple(strArr ...string) []string {
	for i, s := range strArr {
		strArr[i] = strings.ToUpper(s)
	}
	return strArr
}
