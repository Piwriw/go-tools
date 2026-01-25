# Prom Package Test Summary

## 概述 / Overview

本测试套件为 `prom` 包提供全面的单元测试覆盖，包括 Prometheus 客户端查询、PromQL 验证和数据降采样功能。
This test suite provides comprehensive unit test coverage for the `prom` package, including Prometheus client queries, PromQL validation, and data downsampling functionality.

## 测试统计 / Test Statistics

| 指标 / Metric | 值 / Value |
|--------------|------------|
| **测试用例总数 / Total Test Cases** | 173 |
| **代码覆盖率 / Code Coverage** | 87.8% |
| **并行测试 / Parallel Tests** | 支持 / Supported |
| **竞态检测 / Race Detection** | 通过 / Passed |
| **Mock 服务 / Mock Servers** | httptest.Server |

## 测试覆盖的函数 / Functions Covered

### `promql.go` - Prometheus 客户端

| 函数 / Function | 覆盖率 / Coverage | 说明 / Notes |
|----------------|------------------|-------------|
| `InitPromClient` | 80.0% | 客户端初始化 |
| `Client` | 100.0% | 获取原生客户端 |
| `HTTPClient` | 100.0% | 获取 HTTP 客户端 |
| `NewPrometheusClient` | 91.7% | 创建新客户端 |
| `WithToken` | 100.0% | 设置 Token 选项 |
| `WithTimeout` | 100.0% | 设置超时选项 |
| `RoundTrip` | 100.0% | HTTP 请求拦截 |
| `Query` | 90.0% | 执行即时查询 |
| `QueryResPromQL` | 66.7% | 查询并返回结果 |
| `QueryOneValue` | 81.8% | 查询单个值 |
| `QueryVector` | 85.7% | 查询向量类型 |
| `QueryRange` | 90.9% | 执行范围查询 |
| `QueryRangeMatrix` | 86.7% | 查询范围矩阵 |
| `Validate` | 100.0% | 验证 PromQL |
| `Prettify` | 100.0% | 美化 PromQL |
| `QueryAllMetrics` | 83.3% | 查询所有指标 |
| `QueryMetric` | 83.3% | 查询特定指标 |
| `Reload` | 80.0% | 重新加载配置 |
| `ValidateMetric` | 92.9% | 验证指标名称 |
| `extractPromQLQuery` | 100.0% | 提取 PromQL |
| `PushGateway` | 0.0% | 推送到 Pushgateway |

### `promres.go` - 结果处理与降采样

| 函数 / Function | 覆盖率 / Coverage | 说明 / Notes |
|----------------|------------------|-------------|
| `String` | 100.0% | 字符串表示 |
| `Rows` | 75.0% | 转换为行 |
| `DownSample` (multiple) | 100.0% | 降采样方法 |
| `DownSampleWithOptions` | 100.0% | 带选项降采样 |
| `setSampler` | 85.7% | 设置采样器 |
| `ensureLastPoint` | 77.8% | 确保最后点 |
| `Sample` (multiple) | 90.0%-100.0% | 采样方法 |
| `aggregationSample` | 90.2% | 聚合采样 |
| `Len` | 75.0% | 获取长度 |
| `Data` | 82.6% | 获取数据 |
| `ResultType` | 100.0% | 获取结果类型 |
| `Metric` | 80.0% | 获取指标信息 |
| `Values` | 53.3% | 获取值 |
| `getPromQLValue` | 75.0% | 获取 PromQL 值 |
| `GetVal` | 100.0% | 获取值 |
| `GetValStr` | 100.0% | 获取字符串值 |
| `GetValue` | 100.0% | 获取值 |
| `convertToFieldMap` | 92.0% | 转换为字段映射 |

## 测试场景 / Test Scenarios

### 1. 客户端初始化测试 / Client Initialization Tests

**TestNewPrometheusClient** (6 个场景)
- ✅ 有效地址无选项
- ✅ 带 Token 的地址
- ✅ 带超时选项
- ✅ 零超时（使用默认值）
- ✅ 负超时（应被忽略）
- ✅ 多个选项组合

**TestInitPromClient** (2 个场景)
- ✅ 有效初始化
- ✅ 带 Token 初始化

**TestPrometheusClientMethods** (2 个场景)
- ✅ Client() 方法
- ✅ HTTPClient() 方法

### 2. 认证测试 / Authentication Tests

**TestAuthTransport** (2 个场景)
- ✅ 基本认证头设置
- ✅ 空 Token 处理

### 3. PromQL 操作测试 / PromQL Operation Tests

**TestExtractPromQLQuery** (9 个场景)
- ✅ 简单查询
- ✅ 带管道的查询
- ✅ 复杂查询
- ✅ 带空格的查询
- ✅ 嵌套引号查询
- ✅ 无查询模式
- ✅ 空输入
- ✅ 多个查询（提取第一个）
- ✅ 带换行符的查询

**TestValidateMetric** (10 个场景)
- ✅ 有效的指标名和标签
- ✅ 带下划线的有效指标
- ❌ 以数字开头的无效指标名
- ❌ 包含特殊字符的无效指标名
- ❌ 空指标名
- ❌ 无效标签名
- ✅ 带下划线的有效标签名
- ✅ 带 PromQL 查询的标签（有效）
- ❌ 带 PromQL 查询的标签（无效）
- ✅ 所有有效标签名的指标

**TestPrettify** (7 个场景)
- ✅ 简单指标
- ✅ 带聚合的查询
- ✅ 复杂查询
- ✅ 二元操作
- ❌ 无效查询（语法错误）
- ✅ 空查询
- ✅ 带注释的查询

### 4. 查询执行测试 / Query Execution Tests

**TestQuery** (3 个场景)
- ✅ 成功查询
- ❌ 查询错误响应
- ❌ HTTP 错误

**TestQueryVector** (3 个场景)
- ✅ 成功的向量查询
- ❌ 查询返回 Scalar（应报错）
- ❌ 查询返回 Matrix（应报错）

**TestQueryRangeMatrix** (2 个场景)
- ✅ 成功的范围查询
- ❌ 查询返回 Vector（应报错）

**TestQueryOneValue** (3 个场景)
- ✅ 单值查询
- ✅ 无结果
- ❌ 多个值（应报错）

**TestQueryAllMetrics** (3 个场景)
- ✅ 成功查询
- ✅ 空结果
- ❌ 错误响应

**TestQueryRange** (3 个场景) - **新增**
- ✅ 成功的范围查询
- ✅ 空结果
- ❌ 查询错误响应

**TestQueryMetric** (3 个场景) - **新增**
- ✅ 成功的指标查询
- ✅ 指标未找到
- ❌ 查询错误

### 5. PromQL 验证测试 / PromQL Validation Tests

**TestValidate** (9 个场景)
- ✅ 有效的简单查询
- ✅ 带函数的有效查询
- ✅ 带聚合的有效查询
- ✅ 带二元操作的有效查询
- ❌ 无效查询（语法错误）
- ❌ 无效查询（未闭合括号）
- ✅ 空查询
- ✅ 带正则表达式的有效查询
- ✅ 带范围的有效查询

### 6. 结果处理测试 / Result Processing Tests

**TestResPromQLRows** (3 个场景) - **新增**
- ✅ Vector 值
- ✅ Matrix 值
- ✅ 空 Vector

**TestResPromQLLen** (3 个场景) - **新增**
- ✅ 带两个样本的 Vector
- ✅ 带一个序列两个点的 Matrix
- ✅ 空 Vector

**TestResPromQLData** (3 个场景) - **新增**
- ✅ Vector 值
- ✅ Matrix 值
- ✅ 空 Vector

**TestResPromQLResultType** (3 个场景) - **新增**
- ✅ Vector 类型
- ✅ Matrix 类型
- ❌ Scalar 类型（不支持）

**TestResPromQLMetric** (4 个场景) - **新增**
- ✅ Vector 单样本
- ✅ Vector 多样本
- ✅ Matrix 单系列
- ✅ 空 Vector

### 7. 降采样测试 / Downsampling Tests

**TestDownSample_Uniform** (8 个场景)
- ✅ 空 Vector
- ✅ 单点 Vector
- ✅ 无需降采样
- ✅ 降采样到 2 点
- ✅ 降采样到 1 点
- ✅ 带 ResultPair
- ✅ Matrix 降采样

**TestDownSample_Max** (6 个场景)
- ✅ 降采样到 1 点（取最大值）
- ✅ 多点取最大值
- ✅ 降采样到 3 点
- ✅ 带 ResultPair
- ✅ 带选项降采样

**TestDownSample_Min** (6 个场景)
- ✅ 降采样到 1 点（取最小值）
- ✅ 多点取最小值
- ✅ 降采样到 3 点
- ✅ 带 ResultPair
- ✅ 带选项降采样

**TestDownSample_Average** (6 个场景)
- ✅ 降采样到 1 点（取平均值）
- ✅ 多点取平均值
- ✅ 降采样到 3 点
- ✅ 带 ResultPair
- ✅ 带选项降采样

**TestDownSample_LTTB** (4 个场景)
- ✅ 基本 LTTB 降采样
- ✅ 少于目标点数
- ✅ 带 ResultPair
- ✅ 带选项降采样

**TestDownSampleWithOptions_Validation** (2 个场景)
- ❌ 零目标点数
- ❌ 负目标点数

### 8. 配置重载测试 / Configuration Reload Tests

**TestReload** (2 个场景)
- ✅ 成功重载
- ❌ 重载失败

### 9. 辅助函数测试 / Helper Function Tests

**TestSetSampler** (7 个场景)
- ✅ Uniform 采样器
- ✅ Max 采样器
- ✅ Min 采样器
- ✅ Average 采样器
- ✅ LTTB 采样器
- ❌ 无效算法
- ❌ 空 Vector

**TestEnsureLastPoint** (5 个场景)
- ✅ 单点
- ✅ 两点（无需修改）
- ✅ 最后点时间不同（添加新点）
- ✅ 多点保持最后点
- ✅ 带 ResultPair

**TestConvertToFieldMap** (10 个场景)
- ✅ Vector 单样本
- ✅ Vector 多样本
- ✅ Matrix 单系列
- ✅ Matrix 多系列
- ✅ 带标签的样本
- ❌ Scalar 类型（不支持）
- ❌ 空结果
- ❌ String 类型（已弃用）

### 10. 示例函数 / Example Functions

| 示例 / Example | 描述 / Description |
|----------------|-------------------|
| `ExampleNewPrometheusClient` | 创建 Prometheus 客户端 |
| `ExamplePrometheusClient_Query` | 执行 PromQL 查询 |
| `ExamplePrometheusClient_QueryVector` | 查询 Vector 类型结果 |
| `ExamplePrometheusClient_QueryRange` | 范围查询 |
| `ExamplePrometheusClient_Validate` | 验证 PromQL |
| `ExamplePrometheusClient_Prettify` | 美化 PromQL |
| `ExamplePrometheusClient_QueryMetric` | 查询特定指标 |
| `ExampleResultPair_DownSample` | Uniform 降采样 |
| `ExampleResultPair_DownSampleWithOptions_Max` | Max 降采样 |
| `ExampleResultPair_DownSampleWithOptions_Min` | Min 降采样 |
| `ExampleResultPair_DownSampleWithOptions_Average` | Average 降采样 |
| `ExampleResultPair_DownSampleWithOptions_LTTB` | LTTB 降采样 |

## 运行命令 / Execution Commands

### 运行所有测试
```bash
go test ./prom/...
```

### 运行测试并显示详细输出
```bash
go test -v ./prom/...
```

### 运行测试并检测竞态条件
```bash
go test -race ./prom/...
```

### 生成覆盖率报告
```bash
go test -coverprofile=/tmp/prom.cover.out ./prom/...
go tool cover -html=/tmp/prom.cover.out -o prom/coverage.html
go tool cover -func=/tmp/prom.cover.out
```

### 运行基准测试
```bash
go test -bench=. -benchmem ./prom/...
```

### 运行特定测试
```bash
# 只运行客户端初始化测试
go test -run TestNewPrometheusClient ./prom/...

# 只运行降采样测试
go test -run TestDownSample ./prom/...

# 只运行查询测试
go test -run TestQuery ./prom/...
```

## 覆盖率详情 / Coverage Details

```
github.piwriw.go-tools/pkg/prom/promql.go:38:   InitPromClient         80.0%
github.piwriw.go-tools/pkg/prom/promql.go:56:   Client                 100.0%
github.piwriw.go-tools/pkg/prom/promql.go:60:   HTTPClient             100.0%
github.piwriw.go-tools/pkg/prom/promql.go:65:   NewPrometheusClient    91.7%
github.piwriw.go-tools/pkg/prom/promql.go:99:   WithToken              100.0%
github.piwriw.go-tools/pkg/prom/promql.go:107:  WithTimeout            100.0%
github.piwriw.go-tools/pkg/prom/promql.go:122:  RoundTrip              100.0%
github.piwriw.go-tools/pkg/prom/promql.go:128:  Query                  90.0%
github.piwriw.go-tools/pkg/prom/promql.go:146:  QueryResPromQL         66.7%
github.piwriw.go-tools/pkg/prom/promql.go:160:  QueryOneValue          81.8%
github.piwriw.go-tools/pkg/prom/promql.go:178:  QueryVector            85.7%
github.piwriw.go-tools/pkg/prom/promql.go:202:  QueryRange             90.9%
github.piwriw.go-tools/pkg/prom/promql.go:226:  QueryRangeMatrix       86.7%
github.piwriw.go-tools/pkg/prom/promql.go:254:  Validate               100.0%
github.piwriw.go-tools/pkg/prom/promql.go:274:  Prettify               100.0%
github.piwriw.go-tools/pkg/prom/promql.go:283:  PushGateway            0.0%
github.piwriw.go-tools/pkg/prom/promql.go:293:  QueryAllMetrics        83.3%
github.piwriw.go-tools/pkg/prom/promql.go:305:  QueryMetric            83.3%
github.piwriw.go-tools/pkg/prom/promql.go:317:  Reload                 80.0%
github.piwriw.go-tools/pkg/prom/promql.go:342:  ValidateMetric         92.9%
github.piwriw.go-tools/pkg/prom/promql.go:367:  extractPromQLQuery      100.0%
github.piwriw.go-tools/pkg/prom/promres.go:47:   String                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:51:   Rows                   75.0%
github.piwriw.go-tools/pkg/prom/promres.go:68:   DownSample             100.0%
github.piwriw.go-tools/pkg/prom/promres.go:75:   DownSampleWithOptions   100.0%
github.piwriw.go-tools/pkg/prom/promres.go:88:   DownSample             100.0%
github.piwriw.go-tools/pkg/prom/promres.go:108:  DownSampleWithOptions   100.0%
github.piwriw.go-tools/pkg/prom/promres.go:143:  setSampler              85.7%
github.piwriw.go-tools/pkg/prom/promres.go:161:  ensureLastPoint         77.8%
github.piwriw.go-tools/pkg/prom/promres.go:183:  Sample                 90.9%
github.piwriw.go-tools/pkg/prom/promres.go:218:  Sample                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:225:  Sample                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:232:  Sample                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:237:  aggregationSample      90.2%
github.piwriw.go-tools/pkg/prom/promres.go:321:  Sample                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:398:  Len                    75.0%
github.piwriw.go-tools/pkg/prom/promres.go:411:  Data                   82.6%
github.piwriw.go-tools/pkg/prom/promres.go:461:  ResultType             100.0%
github.piwriw.go-tools/pkg/prom/promres.go:471:  Metric                 80.0%
github.piwriw.go-tools/pkg/prom/promres.go:498:  Values                 53.3%
github.piwriw.go-tools/pkg/prom/promres.go:540:  getPromQLValue          75.0%
github.piwriw.go-tools/pkg/prom/promres.go:551:  GetVal                 100.0%
github.piwriw.go-tools/pkg/prom/promres.go:560:  GetValStr              100.0%
github.piwriw.go-tools/pkg/prom/promres.go:573:  GetValue               100.0%
github.piwriw.go-tools/pkg/prom/promres.go:585:  convertToFieldMap       92.0%
total:                                                      (statements)    87.8%
```

## 未覆盖函数说明 / Uncovered Functions Notes

### `PushGateway` (0.0% 覆盖率)
- **原因**: 该函数需要实际的 Pushgateway HTTP 端点进行集成测试
- **建议**: 在集成测试环境中使用真实的 Pushgateway 进行测试

### `Values` (53.3% 覆盖率)
- **原因**: 部分错误处理分支未被测试覆盖
- **建议**: 添加更多错误场景测试

## 关键测试结果 / Key Test Results

```
=== RUN   TestNewPrometheusClient
--- PASS: TestNewPrometheusClient (0.00s)
=== RUN   TestAuthTransport
--- PASS: TestAuthTransport (0.00s)
=== RUN   TestExtractPromQLQuery
--- PASS: TestExtractPromQLQuery (0.00s)
=== RUN   TestValidateMetric
--- PASS: TestValidateMetric (0.00s)
=== RUN   TestPrettify
--- PASS: TestPrettify (0.00s)
=== RUN   TestQuery
--- PASS: TestQuery (0.00s)
=== RUN   TestQueryVector
--- PASS: TestQueryVector (0.00s)
=== RUN   TestQueryRangeMatrix
--- PASS: TestQueryRangeMatrix (0.00s)
=== RUN   TestValidate
--- PASS: TestValidate (0.07s)
=== RUN   TestQueryOneValue
--- PASS: TestQueryOneValue (0.00s)
=== RUN   TestQueryAllMetrics
--- PASS: TestQueryAllMetrics (0.00s)
=== RUN   TestInitPromClient
--- PASS: TestInitPromClient (0.00s)
=== RUN   TestPrometheusClientMethods
--- PASS: TestPrometheusClientMethods (0.00s)
=== RUN   TestReload
--- PASS: TestReload (0.00s)
=== RUN   TestQueryRange
--- PASS: TestQueryRange (0.00s)
=== RUN   TestQueryMetric
--- PASS: TestQueryMetric (0.00s)
=== RUN   TestResPromQLRows
--- PASS: TestResPromQLRows (0.00s)
=== RUN   TestResPromQLLen
--- PASS: TestResPromQLLen (0.00s)
=== RUN   TestResPromQLData
--- PASS: TestResPromQLData (0.00s)
=== RUN   TestResPromQLResultType
--- PASS: TestResPromQLResultType (0.00s)
=== RUN   TestResPromQLMetric
--- PASS: TestResPromQLMetric (0.00s)
```

## 测试覆盖要点 / Testing Highlights

1. **Mock 服务器**: 使用 `httptest.Server` 模拟 Prometheus API，隔离外部依赖
2. **表格驱动测试**: 所有测试都采用结构化的表格驱动方式
3. **降采样算法**: 全面测试 5 种降采样算法（Uniform, Max, Min, Average, LTTB）
4. **PromQL 验证**: 测试 PromQL 语法验证和美化功能
5. **错误处理**: 覆盖各种错误场景（HTTP 错误、无效查询、类型不匹配等）
6. **结果类型**: 测试不同 PromQL 结果类型（Vector, Matrix, Scalar）的处理

## 注意事项 / Notes

1. **Pushgateway 测试**: `PushGateway` 函数需要集成测试环境，不在单元测试范围内
2. **并发安全**: 所有测试通过竞态检测器验证
3. **外部依赖**: 使用 Mock 服务器避免依赖真实的 Prometheus 实例
4. **时间处理**: 范围查询测试使用相对时间，确保测试稳定性

## 结论 / Conclusion

`prom` 包的测试套件实现了 87.8% 的代码覆盖率，超过 80% 的目标要求。测试全面覆盖了：
- Prometheus 客户端的核心功能（查询、验证、配置重载）
- 五种降采样算法的完整实现
- PromQL 语法验证和美化
- 结果处理和类型转换

所有测试均通过竞态检测器验证，确保代码质量和并发安全性。未覆盖的函数主要是需要外部集成环境的场景（如 Pushgateway）。
