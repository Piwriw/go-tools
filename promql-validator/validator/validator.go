// Copyright 2025 The Prometheus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package validator 提供 PromQL 查询的语法和语义校验功能。
//
// 该包可以验证 PromQL 查询字符串的有效性，并返回详细的验证结果，
// 包括：
//   - 语法错误及其位置
//   - 表达式类型（标量、瞬时向量、范围向量）
//   - 查询中使用的指标名称
//   - 查询中调用的函数列表
//
// 基本使用：
//
//	result := validator.Validate("http_requests_total")
//	if result.Valid {
//	    fmt.Println("查询有效")
//	}
package validator

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
)

const (
	// emptyQueryMsg 是空查询的错误消息
	emptyQueryMsg = "empty query"

	// maxQueryLength 是查询字符串的最大允许长度
	maxQueryLength = 100000 // 100KB

	// positionFormat 是带位置信息的错误格式
	positionFormat = "position %d: %s"
)

// ParseError 表示 PromQL 解析错误
type ParseError struct {
	Pos     int
	Message string
}

func (e ParseError) Error() string {
	// 防止空消息导致格式问题
	if e.Message == "" {
		e.Message = "unknown error"
	}

	if e.Pos <= 0 {
		return e.Message
	}
	return fmt.Sprintf(positionFormat, e.Pos, e.Message)
}

// newParseError 创建新的 ParseError
func newParseError(pos int, message string) ParseError {
	return ParseError{
		Pos:     pos,
		Message: message,
	}
}

// newParseErrorWithoutPos 创建无位置信息的 ParseError
func newParseErrorWithoutPos(message string) ParseError {
	return ParseError{
		Message: message,
	}
}

// Warning 表示非致命的警告信息
type Warning struct {
	Message string
}

func (w Warning) String() string {
	return w.Message
}

// ValidationResult 表示 PromQL 校验结果
type ValidationResult struct {
	Valid     bool
	Errors    []ParseError
	Warnings  []Warning
	ExprType  parser.ValueType
	Metrics   []string
	Functions []string
}

// String returns a human-readable representation of the validation result
func (r *ValidationResult) String() string {
	var sb strings.Builder

	r.writeStatus(&sb)
	r.writeErrors(&sb)
	r.writeWarnings(&sb)

	return sb.String()
}

// writeStatus 写入验证状态和元数据
func (r *ValidationResult) writeStatus(sb *strings.Builder) {
	if r.Valid {
		sb.WriteString("✅ Query is VALID\n")
		sb.WriteString(fmt.Sprintf("   Expression type: %s\n", r.ExprType))
		r.writeMetadata(sb)
	} else {
		sb.WriteString("❌ Query is INVALID\n")
	}
}

// writeMetadata 写入指标和函数信息
func (r *ValidationResult) writeMetadata(sb *strings.Builder) {
	if len(r.Metrics) > 0 {
		sb.WriteString(fmt.Sprintf("   Metrics: %v\n", r.Metrics))
	}
	if len(r.Functions) > 0 {
		sb.WriteString(fmt.Sprintf("   Functions: %v\n", r.Functions))
	}
}

// writeErrors 写入错误信息
func (r *ValidationResult) writeErrors(sb *strings.Builder) {
	if len(r.Errors) == 0 {
		return
	}
	sb.WriteString("\nErrors:\n")
	for _, e := range r.Errors {
		sb.WriteString(fmt.Sprintf("  - %s\n", e.Error()))
	}
}

// writeWarnings 写入警告信息
func (r *ValidationResult) writeWarnings(sb *strings.Builder) {
	if len(r.Warnings) == 0 {
		return
	}
	sb.WriteString("\nWarnings:\n")
	for _, w := range r.Warnings {
		sb.WriteString(fmt.Sprintf("  - %s\n", w.String()))
	}
}

// newValidationResult 创建初始化好的 ValidationResult
func newValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:     false,
		Errors:    make([]ParseError, 0, 4),
		Warnings:  make([]Warning, 0, 2),
		Metrics:   make([]string, 0, 4),
		Functions: make([]string, 0, 2),
	}
}

// addError 添加单个错误
func (r *ValidationResult) addError(err ParseError) {
	r.Errors = append(r.Errors, err)
}

// addErrors 批量添加错误
func (r *ValidationResult) addErrors(errs []ParseError) {
	r.Errors = append(r.Errors, errs...)
}

// addWarning 添加警告
func (r *ValidationResult) addWarning(warning Warning) {
	r.Warnings = append(r.Warnings, warning)
}

// IsEmpty 检查结果是否为空（无错误、无警告、无元数据）
func (r *ValidationResult) IsEmpty() bool {
	return len(r.Errors) == 0 &&
		len(r.Warnings) == 0 &&
		len(r.Metrics) == 0 &&
		len(r.Functions) == 0
}

// HasWarnings 检查是否有警告
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// ErrorCount 返回错误数量
func (r *ValidationResult) ErrorCount() int {
	return len(r.Errors)
}

// WarningCount 返回警告数量
func (r *ValidationResult) WarningCount() int {
	return len(r.Warnings)
}

// Validate 检查 PromQL 查询在语法和语义上是否正确
func Validate(query string) *ValidationResult {
	result := newValidationResult()

	// 1. 预处理和基础验证
	normalizedQuery, parseErr := preprocessQuery(query)
	if parseErr.Message != "" {
		result.addError(parseErr)
		return result
	}

	// 2. 解析查询
	expr, parseErrs := parseQuery(normalizedQuery)
	if parseErrs != nil {
		result.addErrors(parseErrs)
		return result
	}

	// 3. 收集元数据并返回
	return finalizeResult(expr, result)
}

// preprocessQuery 预处理查询字符串
func preprocessQuery(query string) (string, ParseError) {
	if query == "" {
		return "", newParseErrorWithoutPos(emptyQueryMsg)
	}

	// 检查长度限制
	if len(query) > maxQueryLength {
		return "", newParseErrorWithoutPos(
			fmt.Sprintf("query length exceeds maximum allowed size of %d bytes", maxQueryLength),
		)
	}

	normalized := strings.TrimSpace(query)
	if normalized == "" {
		return "", newParseErrorWithoutPos(emptyQueryMsg)
	}

	return normalized, ParseError{}
}

// parseQuery 解析 PromQL 查询
func parseQuery(query string) (parser.Expr, []ParseError) {
	p := parser.NewParser(query)

	// 使用 defer 确保清理，即使发生 panic
	defer func() {
		if p != nil {
			p.Close()
		}
	}()

	expr, err := p.ParseExpr()
	if err != nil {
		return nil, convertParseErrors(err)
	}

	return expr, nil
}

// convertParseErrors 转换解析错误
func convertParseErrors(err error) []ParseError {
	if parseErrs, ok := err.(parser.ParseErrors); ok {
		errors := make([]ParseError, len(parseErrs))
		for i, pe := range parseErrs {
			errors[i] = newParseError(
				int(pe.PositionRange.Start),
				pe.Err.Error(),
			)
		}
		return errors
	}
	return []ParseError{newParseErrorWithoutPos(err.Error())}
}

// finalizeResult 完成验证结果
func finalizeResult(expr parser.Expr, result *ValidationResult) *ValidationResult {
	result.ExprType = expr.Type()
	collectMetadata(expr, result)
	result.Valid = true
	return result
}

func collectMetadata(expr parser.Expr, result *ValidationResult) {
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			// n 不可能为 nil，这是类型断言保证的
			if n.Name != "" {
				result.Metrics = append(result.Metrics, n.Name)
			}

		case *parser.MatrixSelector:
			// 安全地获取 VectorSelector，添加 nil 检查
			if n != nil && n.VectorSelector != nil {
				if vs, ok := n.VectorSelector.(*parser.VectorSelector); ok && vs != nil {
					if vs.Name != "" {
						result.Metrics = append(result.Metrics, vs.Name)
					}
				}
			}

		case *parser.Call:
			// n 不可能为 nil
			if n.Func != nil {
				result.Functions = append(result.Functions, n.Func.Name)
			}
		}
		return nil
	})

	result.Metrics = uniqueStrings(result.Metrics)
	result.Functions = uniqueStrings(result.Functions)
}

// uniqueStrings 去重字符串切片，保持顺序
func uniqueStrings(slice []string) []string {
	if len(slice) == 0 {
		return slice // 空切片直接返回
	}

	seen := make(map[string]struct{}, len(slice))
	result := make([]string, 0, len(slice))

	for _, s := range slice {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result
}
