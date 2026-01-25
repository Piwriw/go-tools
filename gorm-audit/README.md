# GORM Audit Plugin

[English](#gorm-audit-plugin) | [中文文档](#gorm-审计插件-中文文档)

A comprehensive audit logging plugin for GORM that tracks all database operations with flexible event handlers and detailed context tracking.

## Features

- **Complete Operation Tracking**: Audit Create, Update, Delete, and Query operations
- **Flexible Event Handlers**: Chain multiple handlers with middleware support
- **Context-Aware**: Track user information, request IDs, IP addresses, and more
- **Async Processing**: Non-blocking audit event dispatching
- **Skip Functionality**: Selectively skip audit for specific operations
- **Multiple Output Formats**: Built-in console handler with colored/text/JSON output
- **Panic Recovery**: Safe handlers with automatic panic recovery
- **Configurable Levels**: Audit all operations, changes only, or disable

## Installation

```bash
go get github.com/piwriw/gorm/gorm-audit
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/piwriw/gorm/gorm-audit"
    "github.com/piwriw/gorm/gorm-audit/handler"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:100"`
    Email string `gorm:"size:100;uniqueIndex"`
}

func main() {
    // Connect to database
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Create audit plugin
    auditPlugin := audit.New(&audit.Config{
        Level:        audit.AuditLevelAll, // Record all operations
        IncludeQuery: true,                // Include Query operations
    })

    // Add console handler
    consoleHandler := handler.NewConsoleHandler()
    consoleHandler.SetColor(true)
    auditPlugin.Use(consoleHandler)

    // Initialize plugin
    if err := db.Use(auditPlugin); err != nil {
        log.Fatal(err)
    }

    // Use context to track user information
    ctx := context.WithValue(context.Background(), "user_id", "12345")
    ctx = context.WithValue(ctx, "username", "admin")

    // Create user - will be audited
    user := User{Name: "John", Email: "john@example.com"}
    db.WithContext(ctx).Create(&user)
}
```

## Configuration

### Audit Levels

```go
auditPlugin := audit.New(&audit.Config{
    Level: audit.AuditLevelAll,          // Record all operations (Create, Update, Delete, Query)
    // Level: audit.AuditLevelChangesOnly, // Only record changes (Create, Update, Delete)
    // Level: audit.AuditLevelNone,        // Disable auditing
})
```

### Event Filtering

Support flexible event filtering mechanisms:

#### Table Name Filtering

```go
// Only audit users and orders tables
audit.New(&audit.Config{
    Filters: []audit.Filter{
        audit.NewTableFilter(audit.FilterModeWhitelist, []string{"users", "orders"}),
    },
})
```

#### Operation Type Filtering

```go
// Only audit create and update operations
audit.NewOperationFilter([]types.Operation{
    types.OperationCreate,
    types.OperationUpdate,
})
```

#### User Filtering

```go
// Exclude test users
audit.NewUserFilter(audit.FilterModeBlacklist, []string{"test_user"})
```

#### Composite Filters

```go
audit.NewCompositeFilter(audit.FilterLogicAnd,
    audit.NewTableFilter(audit.FilterModeWhitelist, []string{"users"}),
    audit.NewOperationFilter([]types.Operation{types.OperationCreate}),
)
```

### Config Hot Reload

Support runtime configuration reload:

```go
// Control via environment variables
// export GORM_AUDIT_LEVEL=all

audit := audit.New(&audit.Config{
    Level: audit.AuditLevelChangesOnly,
})

// Runtime reload
audit.Reload()
```

Supported environment variables:
- `GORM_AUDIT_LEVEL`: `all`, `changes_only`, `none`

### Context Keys

Customize the keys used to extract user information from context:

```go
type contextKey string

auditPlugin := audit.New(&audit.Config{
    ContextKeys: audit.ContextKeyConfig{
        UserID:    contextKey("user_id"),
        Username:  contextKey("username"),
        IP:        contextKey("ip"),
        UserAgent: contextKey("user_agent"),
        RequestID: contextKey("request_id"),
    },
})
```

## Event Handlers

### Console Handler

```go
consoleHandler := handler.NewConsoleHandler()
consoleHandler.SetColor(true)  // Enable colored output
consoleHandler.SetJSON(false)  // Use text format (true for JSON)
auditPlugin.Use(consoleHandler)
```

### Custom Handler

```go
// Using function handler
customHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
    log.Printf("Custom handler: %+v", event)
    return nil
})
auditPlugin.Use(customHandler)

// Using struct handler
type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, event *handler.Event) error {
    // Process event
    return nil
}

auditPlugin.Use(&MyHandler{})
```

### Handler Middleware

#### Filter Middleware

```go
// Only log delete operations
filterHandler := handler.NewFilterMiddleware(
    func(event *handler.Event) bool {
        return event.Operation == handler.OperationDelete
    },
    consoleHandler,
)
auditPlugin.Use(filterHandler)
```

#### Retry Middleware

```go
// Retry up to 3 times with 100ms delay
retryHandler := handler.NewRetryMiddleware(
    consoleHandler,
    3,              // max retries
    100*time.Millisecond,
)
auditPlugin.Use(retryHandler)
```

#### Chain Middleware

```go
// Chain multiple handlers
chain := handler.NewChainMiddleware(handler1)
chain.Then(handler2).Then(handler3)
auditPlugin.Use(chain)
```

## Worker Pool

For high-concurrency scenarios, use a worker pool to control resource usage:

```go
auditPlugin := audit.New(&audit.Config{
    UseWorkerPool: true,
    WorkerConfig: &audit.WorkerPoolConfig{
        WorkerCount: 10,         // Number of worker goroutines
        QueueSize:   10000,      // Event queue buffer size
        Timeout:     5000,       // Event processing timeout (ms)
    },
})
```

**Benefits**:
- Controls goroutine count to prevent resource exhaustion
- Buffers events in memory during traffic spikes
- Prevents database connection pool depletion

## Batch Processing

Batch processing can improve performance in high-concurrency scenarios:

```go
audit.New(&audit.Config{
    UseWorkerPool: true,
    WorkerConfig: &audit.WorkerPoolConfig{
        WorkerCount:   10,
        QueueSize:     10000,
        Timeout:       5000,
        EnableBatch:   true,                 // Enable batch processing
        BatchSize:     1000,                 // Trigger on 1000 events
        FlushInterval: 10 * time.Second,     // Trigger on 10 seconds
    },
})
```

**How it works**:
- Events first enter an in-memory buffer
- Buffer flushes when reaching 1000 events or after 10 seconds
- Batched events are processed individually through Handler
- Single event failure doesn't affect other events

**Use cases**:
- High-concurrency write scenarios
- Handler processing overhead is high
- Millisecond-level latency is acceptable

**Performance improvements**:
- Reduces Handler call frequency
- Batch flushing reduces I/O overhead
- Suitable for log-type Handlers

## Context Tracking

Pass user information through context:

```go
ctx := context.Background()
ctx = context.WithValue(ctx, "user_id", "12345")
ctx = context.WithValue(ctx, "username", "admin")
ctx = context.WithValue(ctx, "ip", "192.168.1.100")
ctx = context.WithValue(ctx, "user_agent", "Mozilla/5.0")
ctx = context.WithValue(ctx, "request_id", "req-001")

// All operations with this context will include the above information
db.WithContext(ctx).Create(&user)
```

## Skip Audit

Skip auditing for specific operations:

```go
// Using SkipAudit function
db.Scopes(audit.SkipAudit).Create(&sensitiveUser)

// Skip multiple operations
tx := db.Scopes(audit.SkipAudit)
tx.Create(&user1)
tx.Create(&user2)
```

## Event Structure

Each audit event contains:

```go
type Event struct {
    Timestamp  string              // Operation timestamp
    Operation  Operation           // CREATE, UPDATE, DELETE, QUERY
    Table      string              // Table name
    PrimaryKey string              // Primary key value
    OldValues  map[string]any      // Values before change (Update/Delete)
    NewValues  map[string]any      // Values after change (Create/Update)
    SQL        string              // Executed SQL
    SQLArgs    []any               // SQL arguments
    UserID     string              // User ID from context
    Username   string              // Username from context
    IP         string              // IP address from context
    UserAgent  string              // User agent from context
    RequestID  string              // Request ID from context
}
```

## Examples

See the [example](example/) directory for a complete working example demonstrating:

- Create, Read, Update, Delete operations
- Context tracking
- Custom handlers
- Skip audit functionality

Run the example:

```bash
cd example
go run main.go
```

## Testing

Run tests:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## Architecture

```
gorm-audit/
├── audit.go           # Main plugin and GORM integration
├── callback.go        # GORM callback implementations
├── config.go          # Configuration types and defaults
├── dispatcher.go      # Event dispatching with panic recovery
├── event.go           # Internal event types
└── handler/
    ├── handler.go     # Public event interface and types
    └── console.go     # Console output handler
```

## Performance Considerations

- Event handlers are executed asynchronously using goroutines
- Handlers should be non-blocking to avoid performance impact
- Consider using a worker pool for high-volume scenarios (configurable via `UseWorkerPool`)

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

**[⬆ Back to top](#gorm-audit-plugin)** | **[中文文档 ↓](#gorm-审计插件-中文文档)**

# GORM 审计插件 [中文文档]

[English](#gorm-audit-plugin) | [中文文档](#gorm-审计插件-中文文档)

一个全面的 GORM 审计日志插件，跟踪所有数据库操作，支持灵活的事件处理器和详细的上下文跟踪。

## 特性

- **完整的操作跟踪**: 审计创建、更新、删除和查询操作
- **灵活的事件处理器**: 支持多个处理器链，包含中间件支持
- **上下文感知**: 跟踪用户信息、请求 ID、IP 地址等
- **异步处理**: 非阻塞的审计事件分发
- **跳过功能**: 可选择性地跳过特定操作的审计
- **多种输出格式**: 内置控制台处理器，支持彩色/文本/JSON 输出
- **Panic 恢复**: 安全的处理器，自动恢复 panic
- **可配置级别**: 审计所有操作、仅变更操作或禁用

## 安装

```bash
go get github.com/piwriw/gorm/gorm-audit
```

## 快速开始

```go
package main

import (
    "context"
    "log"

    "github.com/piwriw/gorm/gorm-audit"
    "github.com/piwriw/gorm/gorm-audit/handler"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:100"`
    Email string `gorm:"size:100;uniqueIndex"`
}

func main() {
    // 连接数据库
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 创建审计插件
    auditPlugin := audit.New(&audit.Config{
        Level:        audit.AuditLevelAll, // 记录所有操作
        IncludeQuery: true,                // 包含查询操作
    })

    // 添加控制台处理器
    consoleHandler := handler.NewConsoleHandler()
    consoleHandler.SetColor(true)
    auditPlugin.Use(consoleHandler)

    // 初始化插件
    if err := db.Use(auditPlugin); err != nil {
        log.Fatal(err)
    }

    // 使用 context 跟踪用户信息
    ctx := context.WithValue(context.Background(), "user_id", "12345")
    ctx = context.WithValue(ctx, "username", "admin")

    // 创建用户 - 将被审计
    user := User{Name: "张三", Email: "zhangsan@example.com"}
    db.WithContext(ctx).Create(&user)
}
```

## 配置

### 审计级别

```go
auditPlugin := audit.New(&audit.Config{
    Level: audit.AuditLevelAll,          // 记录所有操作（创建、更新、删除、查询）
    // Level: audit.AuditLevelChangesOnly, // 仅记录变更操作（创建、更新、删除）
    // Level: audit.AuditLevelNone,        // 禁用审计
})
```

### 事件过滤

支持灵活的事件过滤机制：

#### 表名过滤

```go
// 只审计 users 和 orders 表
audit.New(&audit.Config{
    Filters: []audit.Filter{
        audit.NewTableFilter(audit.FilterModeWhitelist, []string{"users", "orders"}),
    },
})
```

#### 操作类型过滤

```go
// 只审计 create 和 update 操作
audit.NewOperationFilter([]types.Operation{
    types.OperationCreate,
    types.OperationUpdate,
})
```

#### 用户过滤

```go
// 排除测试用户
audit.NewUserFilter(audit.FilterModeBlacklist, []string{"test_user"})
```

#### 组合过滤器

```go
audit.NewCompositeFilter(audit.FilterLogicAnd,
    audit.NewTableFilter(audit.FilterModeWhitelist, []string{"users"}),
    audit.NewOperationFilter([]types.Operation{types.OperationCreate}),
)
```

### 配置热更新

支持运行时重新加载配置：

```go
// 通过环境变量控制
// export GORM_AUDIT_LEVEL=all

audit := audit.New(&audit.Config{
    Level: audit.AuditLevelChangesOnly,
})

// 运行时重新加载
audit.Reload()
```

支持的环境变量：
- `GORM_AUDIT_LEVEL`: `all`, `changes_only`, `none`

### 上下文键

自定义从 context 中提取用户信息的键：

```go
type contextKey string

auditPlugin := audit.New(&audit.Config{
    ContextKeys: audit.ContextKeyConfig{
        UserID:    contextKey("user_id"),
        Username:  contextKey("username"),
        IP:        contextKey("ip"),
        UserAgent: contextKey("user_agent"),
        RequestID: contextKey("request_id"),
    },
})
```

## 事件处理器

### 控制台处理器

```go
consoleHandler := handler.NewConsoleHandler()
consoleHandler.SetColor(true)  // 启用彩色输出
consoleHandler.SetJSON(false)  // 使用文本格式（true 为 JSON）
auditPlugin.Use(consoleHandler)
```

### 自定义处理器

```go
// 使用函数处理器
customHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
    log.Printf("自定义处理器: %+v", event)
    return nil
})
auditPlugin.Use(customHandler)

// 使用结构体处理器
type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, event *handler.Event) error {
    // 处理事件
    return nil
}

auditPlugin.Use(&MyHandler{})
```

### 处理器中间件

#### 过滤中间件

```go
// 仅记录删除操作
filterHandler := handler.NewFilterMiddleware(
    func(event *handler.Event) bool {
        return event.Operation == handler.OperationDelete
    },
    consoleHandler,
)
auditPlugin.Use(filterHandler)
```

#### 重试中间件

```go
// 最多重试 3 次，间隔 100ms
retryHandler := handler.NewRetryMiddleware(
    consoleHandler,
    3,                   // 最大重试次数
    100*time.Millisecond, // 重试间隔
)
auditPlugin.Use(retryHandler)
```

#### 链式中间件

```go
// 链接多个处理器
chain := handler.NewChainMiddleware(handler1)
chain.Then(handler2).Then(handler3)
auditPlugin.Use(chain)
```

## 工作池

对于高并发场景，使用工作池来控制资源使用：

```go
auditPlugin := audit.New(&audit.Config{
    UseWorkerPool: true,
    WorkerConfig: &audit.WorkerPoolConfig{
        WorkerCount: 10,         // Worker goroutine 数量
        QueueSize:   10000,      // 事件队列缓冲区大小
        Timeout:     5000,       // 事件处理超时时间（毫秒）
    },
})
```

**优势**：
- 控制 goroutine 数量，防止资源耗尽
- 在流量高峰时在内存中缓冲事件
- 防止数据库连接池耗尽

## 批量处理

批量处理可以提高高并发场景下的性能：

```go
audit.New(&audit.Config{
    UseWorkerPool: true,
    WorkerConfig: &audit.WorkerPoolConfig{
        WorkerCount:   10,
        QueueSize:     10000,
        Timeout:       5000,
        EnableBatch:   true,                 // 启用批量处理
        BatchSize:     1000,                 // 1000 条事件触发
        FlushInterval: 10 * time.Second,     // 10 秒触发
    },
})
```

**工作原理**：
- 事件先进入内存缓冲区
- 缓冲区达到 1000 条或超过 10 秒时触发刷新
- 批量收集后逐个调用 Handler 处理
- 单个事件失败不影响其他事件

**适用场景**：
- 高并发写入场景
- Handler 处理开销较大
- 可以容忍毫秒级延迟

**性能提升**：
- 减少 Handler 调用频率
- 批量刷新可降低 I/O 开销
- 适合日志类 Handler

## 上下文跟踪

通过 context 传递用户信息：

```go
ctx := context.Background()
ctx = context.WithValue(ctx, "user_id", "12345")
ctx = context.WithValue(ctx, "username", "admin")
ctx = context.WithValue(ctx, "ip", "192.168.1.100")
ctx = context.WithValue(ctx, "user_agent", "Mozilla/5.0")
ctx = context.WithValue(ctx, "request_id", "req-001")

// 使用此 context 的所有操作都将包含上述信息
db.WithContext(ctx).Create(&user)
```

## 跳过审计

为特定操作跳过审计：

```go
// 使用 SkipAudit 函数
audit.SkipAudit(db).Create(&sensitiveUser)

// 跳过多个操作
tx := audit.SkipAudit(db)
tx.Create(&user1)
tx.Create(&user2)
```

## 事件结构

每个审计事件包含：

```go
type Event struct {
    Timestamp  string              // 操作时间戳
    Operation  Operation           // CREATE, UPDATE, DELETE, QUERY
    Table      string              // 表名
    PrimaryKey string              // 主键值
    OldValues  map[string]any      // 变更前的值（更新/删除）
    NewValues  map[string]any      // 变更后的值（创建/更新）
    SQL        string              // 执行的 SQL
    SQLArgs    []any               // SQL 参数
    UserID     string              // 用户 ID（来自 context）
    Username   string              // 用户名（来自 context）
    IP         string              // IP 地址（来自 context）
    UserAgent  string              // 用户代理（来自 context）
    RequestID  string              // 请求 ID（来自 context）
}
```

## 示例

查看 [example](example/) 目录中的完整工作示例，演示：

- 创建、读取、更新、删除操作
- 上下文跟踪
- 自定义处理器
- 跳过审计功能

运行示例：

```bash
cd example
go run main.go
```

## 测试

运行测试：

```bash
go test -v ./...
```

运行测试并查看覆盖率：

```bash
go test -cover ./...
```

## 架构

```
gorm-audit/
├── audit.go           # 主插件和 GORM 集成
├── callback.go        # GORM 回调实现
├── config.go          # 配置类型和默认值
├── dispatcher.go      # 事件分发（带 panic 恢复）
├── event.go           # 内部事件类型
└── handler/
    ├── handler.go     # 公共事件接口和类型
    └── console.go     # 控制台输出处理器
```

## 性能考虑

- 事件处理器使用 goroutine 异步执行
- 处理器应该是非阻塞的，以避免性能影响
- 对于高并发场景，建议使用工作池（通过 `UseWorkerPool` 配置）

## 许可证

MIT License

## 贡献

欢迎贡献！请随时提交 Pull Request。

---

**[English ↑](#gorm-audit-plugin)** | **[⬆ 返回顶部](#gorm-审计插件-中文文档)**
