package format

import (
	"strings"
)

// ConvertMapKeysToUpper 将 map 的所有 key 转换为大写。
// 注意：值为 nil 时返回空 map，不会 panic。
func ConvertMapKeysToUpper(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		upperKey := strings.ToUpper(key)
		result[upperKey] = value
	}

	return result
}

// ConvertToUpperSingle 将单个字符串转为大写。
func ConvertToUpperSingle(str string) string {
	return strings.ToUpper(str)
}

// ConvertToUpperMultiple 将多个字符串批量转为大写，直接修改原切片并返回。
func ConvertToUpperMultiple(strArr ...string) []string {
	for i, s := range strArr {
		strArr[i] = strings.ToUpper(s)
	}
	return strArr
}

// ============================================================================
// Blank / Default
// ============================================================================

// IsBlank 判断字符串是否为空或仅包含空白字符。
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// DefaultIfEmpty 当字符串为空时返回默认值，否则返回原值。
func DefaultIfEmpty(s, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

// ============================================================================
// Truncate / Substring
// ============================================================================

// Truncate 截断字符串，超过 maxLen 时追加 suffix（如 "..."）。
// maxLen 为最终返回字符串的最大字节长度（含 suffix）。
func Truncate(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= len(suffix) {
		return suffix[:maxLen]
	}
	return s[:maxLen-len(suffix)] + suffix
}

// Substring 按 rune 安全截取子串（支持中文等多字节字符）。
// start 和 end 均为 rune 索引，end 不包含。
func Substring(s string, start, end int) string {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}

// ============================================================================
// Case Conversion: Camel / Snake / Kebab
// ============================================================================

// ToCamelCase 将 snake_case 或 kebab-case 字符串转为 camelCase。
func ToCamelCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// ToSnakeCase 将 camelCase 或 PascalCase 字符串转为 snake_case。
func ToSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

// ToKebabCase 将 camelCase 或 PascalCase 字符串转为 kebab-case。
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// ============================================================================
// Masking (脱敏)
// ============================================================================

// MaskString 对字符串进行脱敏，保留前后各 keepLen 个字符，中间用 mask 替换。
func MaskString(s string, keepLen int, mask rune) string {
	runes := []rune(s)
	length := len(runes)
	if length <= keepLen*2 {
		return strings.Repeat(string(mask), length)
	}
	left := string(runes[:keepLen])
	right := string(runes[length-keepLen:])
	middle := strings.Repeat(string(mask), length-keepLen*2)
	return left + middle + right
}

// MaskPhone 对手机号脱敏，保留前 3 位和后 4 位，如 "138****1234"。
func MaskPhone(phone string) string {
	return MaskString(phone, 3, '*')
}

// MaskEmail 对邮箱脱敏，保留首字符和 @ 后的域名，如 "t***@example.com"。
func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return email
	}
	return string(email[0]) + strings.Repeat("*", at-1) + email[at:]
}

// ============================================================================
// Padding
// ============================================================================

// PadLeft 左填充字符串到指定长度。
func PadLeft(s string, pad rune, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(pad), length-len(s)) + s
}

// PadRight 右填充字符串到指定长度。
func PadRight(s string, pad rune, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(pad), length-len(s))
}

// ============================================================================
// Reverse / Remove
// ============================================================================

// Reverse 反转字符串（按 rune，支持中文）。
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// RemoveSpaces 移除字符串中的所有空白字符。
func RemoveSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
