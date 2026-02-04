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

// Package validator 提供 PromQL 语法和语义校验功能
package validator

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
)

// ParseError 表示 PromQL 解析错误
type ParseError struct {
	Pos     int
	Message string
}

func (e ParseError) Error() string {
	if e.Pos <= 0 {
		return e.Message
	}
	return fmt.Sprintf("position %d: %s", e.Pos, e.Message)
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

	if r.Valid {
		sb.WriteString("✅ Query is VALID\n")
		sb.WriteString(fmt.Sprintf("   Expression type: %s\n", r.ExprType))
		if len(r.Metrics) > 0 {
			sb.WriteString(fmt.Sprintf("   Metrics: %v\n", r.Metrics))
		}
		if len(r.Functions) > 0 {
			sb.WriteString(fmt.Sprintf("   Functions: %v\n", r.Functions))
		}
	} else {
		sb.WriteString("❌ Query is INVALID\n")
	}

	if len(r.Errors) > 0 {
		sb.WriteString("\nErrors:\n")
		for _, e := range r.Errors {
			sb.WriteString(fmt.Sprintf("  - %s\n", e.Error()))
		}
	}

	if len(r.Warnings) > 0 {
		sb.WriteString("\nWarnings:\n")
		for _, w := range r.Warnings {
			sb.WriteString(fmt.Sprintf("  - %s\n", w.String()))
		}
	}

	return sb.String()
}

// Validate 检查 PromQL 查询在语法和语义上是否正确
func Validate(query string) *ValidationResult {
	result := &ValidationResult{
		Valid:     false,
		Errors:    make([]ParseError, 0),
		Warnings:  make([]Warning, 0),
		Metrics:   make([]string, 0),
		Functions: make([]string, 0),
	}

	query = strings.TrimSpace(query)
	if query == "" {
		result.Errors = append(result.Errors, ParseError{
			Message: "empty query",
		})
		return result
	}

	p := parser.NewParser(query)
	expr, err := p.ParseExpr()
	defer p.Close()

	if err != nil {
		if parseErrs, ok := err.(parser.ParseErrors); ok {
			for _, pe := range parseErrs {
				result.Errors = append(result.Errors, ParseError{
					Pos:     int(pe.PositionRange.Start),
					Message: pe.Err.Error(),
				})
			}
		} else {
			result.Errors = append(result.Errors, ParseError{
				Message: err.Error(),
			})
		}
		return result
	}

	result.ExprType = expr.Type()
	collectMetadata(expr, result)
	result.Valid = len(result.Errors) == 0

	return result
}

func collectMetadata(expr parser.Expr, result *ValidationResult) {
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			if n.Name != "" {
				result.Metrics = append(result.Metrics, n.Name)
			}
		case *parser.MatrixSelector:
			if vs, ok := n.VectorSelector.(*parser.VectorSelector); ok {
				if vs.Name != "" {
					result.Metrics = append(result.Metrics, vs.Name)
				}
			}
		case *parser.Call:
			result.Functions = append(result.Functions, n.Func.Name)
		}
		return nil
	})

	result.Metrics = uniqueStrings(result.Metrics)
	result.Functions = uniqueStrings(result.Functions)
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
