# PromQL Validator

独立的 PromQL 语法和语义校验库。

## 功能

- ✅ **语法检查** - 验证 PromQL 查询的词法和语法
- ✅ **语义检查** - 验证表达式类型、函数参数等
- ✅ **错误定位** - 精确显示错误位置
- ✅ **元数据提取** - 提取查询中的指标和函数信息

## 安装

```bash
go get github.com/prometheus/promql-validator
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/prometheus/promql-validator/validator"
)

func main() {
    // 校验一个简单查询
    result := validator.Validate("http_requests_total")

    if result.Valid {
        fmt.Println("✅ 查询有效!")
        fmt.Printf("类型: %s\n", result.ExprType)
        fmt.Printf("指标: %v\n", result.Metrics)
    } else {
        fmt.Println("❌ 查询无效:")
        for _, err := range result.Errors {
            fmt.Printf("  - %s\n", err.Error())
        }
    }
}
```

## 使用示例

### 基本校验

```go
result := validator.Validate("sum by (job) (rate(http_requests_total[5m]))")
fmt.Println(result.String())
```

### 输出示例

**有效查询：**
```
✅ Query is VALID
   Expression type: Vector
   Metrics: [http_requests_total]
   Functions: [rate]
```

**无效查询：**
```
❌ Query is INVALID

Errors:
  - position 20: unclosed parenthesis
```

## API 文档

### Validate 函数

```go
func Validate(query string) *ValidationResult
```

校验 PromQL 查询字符串。

### ValidationResult 结构体

```go
type ValidationResult struct {
    Valid     bool              // 查询是否有效
    Errors    []ParseError      // 错误列表
    Warnings  []Warning         // 警告列表
    ExprType  parser.ValueType  // 表达式类型
    Metrics   []string          // 涉及的指标名称
    Functions []string          // 使用的函数名
}
```

## 开发说明

此模块使用 `replace` 指令引用主仓库的代码进行开发：

```go
// go.mod
replace github.com/prometheus/prometheus => ../
```

## 许可证

Apache License 2.0
