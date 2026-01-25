# GORM 审计插件实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标:** 构建一个基于 GORM Plugin 接口的审计插件系统，支持可配置的审计级别、事件驱动架构、多种存储方式

**架构:** 使用 GORM 的 Plugin 接口和 Callback 机制，在数据库操作前后拦截并捕获审计信息，通过事件分发器异步发送到用户定义的处理器

**技术栈:** Go 1.18+, GORM v1.31+, context, worker pool

---

## 目录结构

```
gorm-audit/
├── audit.go              # 核心插件定义（Plugin 接口实现）
├── event.go              # 事件类型定义（Operation 枚举、AuditEvent）
├── config.go             # 配置类型定义
├── callback.go           # GORM Callback 处理器
├── dispatcher.go         # 事件分发器（带 panic 恢复）
├── worker_pool.go        # Worker Pool 实现（可选）
├── panic.go              # Panic 处理器
├── handler/              # 内置事件处理器
│   ├── handler.go        # Handler 接口定义
│   ├── console.go        # 控制台输出
│   └── database.go       # 数据库存储
├── example/
│   └── main.go           # 使用示例
└── README.md             # 文档
```

---

## Task 1: 创建基础类型定义（event.go）

**文件:**
- 创建: `gorm-audit/event.go`

**Step 1: 创建 event.go 文件并定义 Operation 枚举**

```go
package audit

import "time"

// Operation 定义审计操作类型的枚举
type Operation string

const (
    OperationCreate Operation = "create"
    OperationUpdate Operation = "update"
    OperationDelete Operation = "delete"
    OperationQuery  Operation = "query"
)

// String 实现 Stringer 接口
func (o Operation) String() string {
    return string(o)
}

// IsValid 验证操作类型是否有效
func (o Operation) IsValid() bool {
    switch o {
    case OperationCreate, OperationUpdate, OperationDelete, OperationQuery:
        return true
    }
    return false
}

// AuditEvent 审计事件
type AuditEvent struct {
    Timestamp  time.Time
    Operation  Operation
    Table      string
    PrimaryKey string
    OldValues  map[string]any
    NewValues  map[string]any
    SQL        string
    SQLArgs    []any
    UserID     string
    Username   string
    IP         string
    UserAgent  string
    RequestID  string
}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 3: 创建单元测试**

创建: `gorm-audit/event_test.go`

```go
package audit

import "testing"

func TestOperationString(t *testing.T) {
    tests := []struct {
        name      string
        op        Operation
        expected  string
    }{
        {"Create", OperationCreate, "create"},
        {"Update", OperationUpdate, "update"},
        {"Delete", OperationDelete, "delete"},
        {"Query", OperationQuery, "query"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.op.String(); got != tt.expected {
                t.Errorf("String() = %v, want %v", got, tt.expected)
            }
        })
    }
}

func TestOperationIsValid(t *testing.T) {
    tests := []struct {
        name      string
        op        Operation
        expected  bool
    }{
        {"Valid Create", OperationCreate, true},
        {"Valid Update", OperationUpdate, true},
        {"Valid Delete", OperationDelete, true},
        {"Valid Query", OperationQuery, true},
        {"Invalid", Operation("invalid"), false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.op.IsValid(); got != tt.expected {
                t.Errorf("IsValid() = %v, want %v", got, tt.expected)
            }
        })
    }
}

func TestAuditEventCreation(t *testing.T) {
    event := &AuditEvent{
        Timestamp:  time.Now(),
        Operation:  OperationCreate,
        Table:      "users",
        PrimaryKey: "1",
        OldValues:  make(map[string]any),
        NewValues:  map[string]any{"name": "test"},
    }

    if event.Operation != OperationCreate {
        t.Errorf("expected OperationCreate, got %v", event.Operation)
    }
    if event.Table != "users" {
        t.Errorf("expected table 'users', got %v", event.Table)
    }
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v ./event_test.go ./event.go`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/event.go gorm-audit/event_test.go
git commit -m "feat(audit): add Operation enum and AuditEvent types"
```

---

## Task 2: 创建配置类型定义（config.go）

**文件:**
- 创建: `gorm-audit/config.go`

**Step 1: 创建 config.go 文件**

```go
package audit

// AuditLevel 审计级别
type AuditLevel int

const (
    AuditLevelAll         AuditLevel = iota // 记录所有操作
    AuditLevelChangesOnly                   // 仅记录变更操作
    AuditLevelNone                          // 不记录
)

// String 实现 Stringer 接口
func (a AuditLevel) String() string {
    switch a {
    case AuditLevelAll:
        return "all"
    case AuditLevelChangesOnly:
        return "changes_only"
    case AuditLevelNone:
        return "none"
    default:
        return "unknown"
    }
}

// ContextKeyConfig Context 键配置
type ContextKeyConfig struct {
    UserID    any
    Username  any
    IP        any
    UserAgent any
    RequestID any
}

// DefaultContextKeys 返回默认的 context key 配置
func DefaultContextKeys() ContextKeyConfig {
    type contextKey string
    return ContextKeyConfig{
        UserID:    contextKey("user_id"),
        Username:  contextKey("username"),
        IP:        contextKey("ip"),
        UserAgent: contextKey("user_agent"),
        RequestID: contextKey("request_id"),
    }
}

// WorkerPoolConfig worker pool 配置
type WorkerPoolConfig struct {
    WorkerCount int           // worker 数量
    QueueSize   int           // 队列大小
    Timeout     int           // 处理超时（毫秒）
}

// DefaultWorkerPoolConfig 返回默认的 worker pool 配置
func DefaultWorkerPoolConfig() *WorkerPoolConfig {
    return &WorkerPoolConfig{
        WorkerCount: 10,
        QueueSize:   1000,
        Timeout:     5000, // 5秒
    }
}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 3: 创建单元测试**

创建: `gorm-audit/config_test.go`

```go
package audit

import (
    "testing"
    "time"
)

func TestAuditLevelString(t *testing.T) {
    tests := []struct {
        name     string
        level    AuditLevel
        expected string
    }{
        {"All", AuditLevelAll, "all"},
        {"ChangesOnly", AuditLevelChangesOnly, "changes_only"},
        {"None", AuditLevelNone, "none"},
        {"Unknown", AuditLevel(100), "unknown"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.level.String(); got != tt.expected {
                t.Errorf("String() = %v, want %v", got, tt.expected)
            }
        })
    }
}

func TestDefaultContextKeys(t *testing.T) {
    keys := DefaultContextKeys()

    if keys.UserID == nil {
        t.Error("UserID should not be nil")
    }
    if keys.Username == nil {
        t.Error("Username should not be nil")
    }
}

func TestDefaultWorkerPoolConfig(t *testing.T) {
    config := DefaultWorkerPoolConfig()

    if config.WorkerCount != 10 {
        t.Errorf("expected WorkerCount 10, got %d", config.WorkerCount)
    }
    if config.QueueSize != 1000 {
        t.Errorf("expected QueueSize 1000, got %d", config.QueueSize)
    }
    if config.Timeout != 5000 {
        t.Errorf("expected Timeout 5000, got %d", config.Timeout)
    }
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/config.go gorm-audit/config_test.go
git commit -m "feat(audit): add Config types with defaults"
```

---

## Task 3: 创建 Handler 接口定义（handler/handler.go）

**文件:**
- 创建: `gorm-audit/handler/handler.go`

**Step 1: 创建 handler 目录和接口文件**

```go
package handler

import "context"

// EventHandler 事件处理器接口
type EventHandler interface {
    Handle(ctx context.Context, event *Event) error
}

// EventHandlerFunc 函数式事件处理器
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
    return f(ctx, event)
}

// Event 事件定义（避免循环导入，重新定义）
type Event struct {
    Timestamp  string
    Operation  string
    Table      string
    PrimaryKey string
    Data       map[string]any
}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 3: 创建单元测试**

创建: `gorm-audit/handler/handler_test.go`

```go
package handler

import (
    "context"
    "testing"
)

func TestEventHandlerFunc(t *testing.T) {
    called := false
    fn := EventHandlerFunc(func(ctx context.Context, event *Event) error {
        called = true
        return nil
    })

    err := fn.Handle(context.Background(), &Event{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if !called {
        t.Error("handler was not called")
    }
}

func TestEventHandlerFuncWithError(t *testing.T) {
    expectedErr := &TestError{}
    fn := EventHandlerFunc(func(ctx context.Context, event *Event) error {
        return expectedErr
    })

    err := fn.Handle(context.Background(), &Event{})
    if err != expectedErr {
        t.Errorf("expected error %v, got %v", expectedErr, err)
    }
}

type TestError struct{}

func (e *TestError) Error() string {
    return "test error"
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v ./handler/...`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/handler/
git commit -m "feat(audit): add EventHandler interface"
```

---

## Task 4: 创建核心 Plugin 实现（audit.go）

**文件:**
- 创建: `gorm-audit/audit.go`

**Step 1: 创建 audit.go 核心插件文件**

```go
package audit

import (
    "gorm.io/gorm"
)

// Config 插件完整配置
type Config struct {
    Level         AuditLevel
    IncludeQuery  bool
    ContextKeys   ContextKeyConfig
    UseWorkerPool bool
    WorkerConfig  *WorkerPoolConfig
}

// Audit GORM 审计插件
type Audit struct {
    config *Config
}

// New 创建新的审计插件实例
func New(config *Config) *Audit {
    if config == nil {
        config = &Config{
            Level:        AuditLevelChangesOnly,
            IncludeQuery: false,
            ContextKeys:  DefaultContextKeys(),
        }
    }
    if config.ContextKeys.UserID == nil {
        config.ContextKeys = DefaultContextKeys()
    }
    return &Audit{config: config}
}

// Name 返回插件名称
func (a *Audit) Name() string {
    return "audit"
}

// Initialize 实现 gorm.Plugin 接口
func (a *Audit) Initialize(db *gorm.DB) error {
    // 注册 Create 回调
    if err := db.Callback().Create().Before("gorm:create").Register("audit:create:before", a.beforeCreate); err != nil {
        return err
    }
    if err := db.Callback().Create().After("gorm:create").Register("audit:create:after", a.afterCreate); err != nil {
        return err
    }

    // 注册 Update 回调
    if err := db.Callback().Update().Before("gorm:update").Register("audit:update:before", a.beforeUpdate); err != nil {
        return err
    }
    if err := db.Callback().Update().After("gorm:update").Register("audit:update:after", a.afterUpdate); err != nil {
        return err
    }

    // 注册 Delete 回调
    if err := db.Callback().Delete().Before("gorm:delete").Register("audit:delete:before", a.beforeDelete); err != nil {
        return err
    }
    if err := db.Callback().Delete().After("gorm:delete").Register("audit:delete:after", a.afterDelete); err != nil {
        return err
    }

    // 可选：注册 Query 回调
    if a.config.IncludeQuery {
        if err := db.Callback().Query().Before("gorm:query").Register("audit:query:before", a.beforeQuery); err != nil {
            return err
        }
        if err := db.Callback().Query().After("gorm:query").Register("audit:query:after", a.afterQuery); err != nil {
            return err
        }
    }

    return nil
}

// Use 添加事件处理器
func (a *Audit) Use(handler interface{}) *Audit {
    // 这个方法将在实现 dispatcher 后补充完整
    return a
}

// ==================== Callback 占位方法 ====================

func (a *Audit) beforeCreate(db *gorm.DB)  {}
func (a *Audit) afterCreate(db *gorm.DB)   {}
func (a *Audit) beforeUpdate(db *gorm.DB)  {}
func (a *Audit) afterUpdate(db *gorm.DB)   {}
func (a *Audit) beforeDelete(db *gorm.DB)  {}
func (a *Audit) afterDelete(db *gorm.DB)   {}
func (a *Audit) beforeQuery(db *gorm.DB)   {}
func (a *Audit) afterQuery(db *gorm.DB)    {}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 3: 创建单元测试**

创建: `gorm-audit/audit_test.go`

```go
package audit

import (
    "testing"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestNewAudit(t *testing.T) {
    // 测试 nil 配置
    p := New(nil)
    if p == nil {
        t.Fatal("expected non-nil plugin")
    }
    if p.config.Level != AuditLevelChangesOnly {
        t.Errorf("expected default level ChangesOnly, got %v", p.config.Level)
    }

    // 测试自定义配置
    customConfig := &Config{
        Level:        AuditLevelAll,
        IncludeQuery: true,
    }
    p = New(customConfig)
    if p.config.Level != AuditLevelAll {
        t.Errorf("expected level All, got %v", p.config.Level)
    }
    if !p.config.IncludeQuery {
        t.Error("expected IncludeQuery to be true")
    }
}

func TestAuditName(t *testing.T) {
    p := New(nil)
    if p.Name() != "audit" {
        t.Errorf("expected name 'audit', got %v", p.Name())
    }
}

func TestAuditInitialize(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }

    p := New(nil)
    err = p.Initialize(db)
    if err != nil {
        t.Errorf("failed to initialize plugin: %v", err)
    }
}

func TestAuditInitializeWithQuery(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }

    p := New(&Config{
        IncludeQuery: true,
    })
    err = p.Initialize(db)
    if err != nil {
        t.Errorf("failed to initialize plugin with query: %v", err)
    }
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/audit.go gorm-audit/audit_test.go
git commit -m "feat(audit): add core Audit plugin with GORM interface"
```

---

## Task 5: 实现事件分发器（dispatcher.go）

**文件:**
- 创建: `gorm-audit/dispatcher.go`

**Step 1: 创建 dispatcher.go 文件**

```go
package audit

import (
    "context"
    "fmt"
    "log"
    "runtime"
    "sync"
)

// Dispatcher 事件分发器
type Dispatcher struct {
    handlers []EventHandler
    mu       sync.RWMutex
}

// NewDispatcher 创建新的分发器
func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make([]EventHandler, 0),
    }
}

// Add 添加事件处理器
func (d *Dispatcher) Add(handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers = append(d.handlers, handler)
}

// Dispatch 分发事件到所有处理器
func (d *Dispatcher) Dispatch(ctx context.Context, event *AuditEvent) {
    d.mu.RLock()
    handlers := make([]EventHandler, len(d.handlers))
    copy(handlers, d.handlers)
    d.mu.RUnlock()

    for _, handler := range handlers {
        if handler != nil {
            go d.safeHandle(ctx, handler, event)
        }
    }
}

// safeHandle 安全执行事件处理器，带 panic 恢复
func (d *Dispatcher) safeHandle(ctx context.Context, handler EventHandler, event *AuditEvent) {
    defer func() {
        if r := recover(); r != nil {
            d.handlePanic(r, handler, event)
        }
    }()

    _ = handler.Handle(ctx, event)
}

// handlePanic 处理 panic 情况
func (d *Dispatcher) handlePanic(r any, handler EventHandler, event *AuditEvent) {
    buf := make([]byte, 4096)
    n := runtime.Stack(buf, false)
    stack := string(buf[:n])

    msg := fmt.Sprintf(
        "[AUDIT] panic recovered: %v, handler: %T, table: %s, operation: %s\nstack: %s",
        r, handler, event.Table, event.Operation, stack,
    )

    log.Println(msg)
}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 3: 创建单元测试**

创建: `gorm-audit/dispatcher_test.go`

```go
package audit

import (
    "context"
    "errors"
    "sync/atomic"
    "testing"
    "time"
)

// MockEventHandler 用于测试的模拟处理器
type MockEventHandler struct {
    CallCount int32
    Delay     time.Duration
    Panic     bool
    Error     error
}

func (m *MockEventHandler) Handle(ctx context.Context, event *AuditEvent) error {
    atomic.AddInt32(&m.CallCount, 1)
    if m.Delay > 0 {
        time.Sleep(m.Delay)
    }
    if m.Panic {
        panic("test panic")
    }
    return m.Error
}

func TestDispatcherAdd(t *testing.T) {
    d := NewDispatcher()
    handler := &MockEventHandler{}

    d.Add(handler)

    if len(d.handlers) != 1 {
        t.Errorf("expected 1 handler, got %d", len(d.handlers))
    }
}

func TestDispatcherDispatch(t *testing.T) {
    d := NewDispatcher()
    handler := &MockEventHandler{}
    d.Add(handler)

    event := &AuditEvent{
        Timestamp: time.Now(),
        Operation: OperationCreate,
        Table:     "users",
    }

    d.Dispatch(context.Background(), event)

    // 等待异步处理
    time.Sleep(100 * time.Millisecond)

    if atomic.LoadInt32(&handler.CallCount) != 1 {
        t.Errorf("expected handler to be called once, got %d", handler.CallCount)
    }
}

func TestDispatcherDispatchMultiple(t *testing.T) {
    d := NewDispatcher()
    handler1 := &MockEventHandler{}
    handler2 := &MockEventHandler{}
    d.Add(handler1)
    d.Add(handler2)

    event := &AuditEvent{
        Timestamp: time.Now(),
        Operation: OperationUpdate,
        Table:     "users",
    }

    d.Dispatch(context.Background(), event)

    time.Sleep(100 * time.Millisecond)

    if atomic.LoadInt32(&handler1.CallCount) != 1 {
        t.Error("handler1 should be called once")
    }
    if atomic.LoadInt32(&handler2.CallCount) != 1 {
        t.Error("handler2 should be called once")
    }
}

func TestDispatcherPanicRecovery(t *testing.T) {
    d := NewDispatcher()
    handler := &MockEventHandler{Panic: true}
    d.Add(handler)

    event := &AuditEvent{
        Timestamp: time.Now(),
        Operation: OperationDelete,
        Table:     "users",
    }

    // 不应该 panic
    d.Dispatch(context.Background(), event)

    time.Sleep(100 * time.Millisecond)

    if atomic.LoadInt32(&handler.CallCount) != 1 {
        t.Error("handler should be called despite panic")
    }
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/dispatcher.go gorm-audit/dispatcher_test.go
git commit -m "feat(audit): add event dispatcher with panic recovery"
```

---

## Task 6: 更新 Audit 集成 Dispatcher（audit.go）

**文件:**
- 修改: `gorm-audit/audit.go:14-18, 38-50`

**Step 1: 更新 Audit 结构体和 New 函数**

```go
// Audit GORM 审计插件
type Audit struct {
    config     *Config
    dispatcher *Dispatcher
}

// New 创建新的审计插件实例
func New(config *Config) *Audit {
    if config == nil {
        config = &Config{
            Level:        AuditLevelChangesOnly,
            IncludeQuery: false,
            ContextKeys:  DefaultContextKeys(),
        }
    }
    if config.ContextKeys.UserID == nil {
        config.ContextKeys = DefaultContextKeys()
    }

    return &Audit{
        config:     config,
        dispatcher: NewDispatcher(),
    }
}
```

**Step 2: 添加 Use 方法**

在 `audit.go` 中 `Name()` 方法后添加：

```go
// Use 添加事件处理器
func (a *Audit) Use(handler EventHandler) *Audit {
    a.dispatcher.Add(handler)
    return a
}
```

**Step 3: 更新单元测试**

修改: `gorm-audit/audit_test.go:8-16`

```go
func TestNewAudit(t *testing.T) {
    // 测试 nil 配置
    p := New(nil)
    if p == nil {
        t.Fatal("expected non-nil plugin")
    }
    if p.config.Level != AuditLevelChangesOnly {
        t.Errorf("expected default level ChangesOnly, got %v", p.config.Level)
    }
    if p.dispatcher == nil {
        t.Error("expected dispatcher to be initialized")
    }

    // 测试自定义配置
    customConfig := &Config{
        Level:        AuditLevelAll,
        IncludeQuery: true,
    }
    p = New(customConfig)
    if p.config.Level != AuditLevelAll {
        t.Errorf("expected level All, got %v", p.config.Level)
    }
    if !p.config.IncludeQuery {
        t.Error("expected IncludeQuery to be true")
    }
}
```

**Step 4: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v`
Expected: PASS

**Step 5: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/audit.go gorm-audit/audit_test.go
git commit -m "feat(audit): integrate Dispatcher into Audit plugin"
```

---

## Task 7: 实现 Callback 处理器（callback.go）

**文件:**
- 创建: `gorm-audit/callback.go`

**Step 1: 创建 callback.go 文件**

```go
package audit

import (
    "context"
    "gorm.io/gorm"
    "reflect"
    "strings"
    "time"
)

// auditContextKey 存储在 DB 实例中的临时数据键
type auditContextKey struct{}

// auditData 存储审计过程中的临时数据
type auditData struct {
    startTime time.Time
    oldValues map[string]any
    newValues map[string]any
}

// SkipAudit 在 db 实例上设置跳过审计标记
func SkipAudit(db *gorm.DB) *gorm.DB {
    return db.InstanceSet("audit:skip", true)
}

// shouldSkip 检查是否应该跳过审计
func (a *Audit) shouldSkip(db *gorm.DB) bool {
    if db == nil || db.Statement == nil {
        return true
    }
    // 检查是否有跳过标记
    if _, ok := db.InstanceGet("audit:skip"); ok {
        return true
    }
    return false
}

// shouldAuditForLevel 根据级别检查是否应该审计
func (a *Audit) shouldAuditForLevel(op Operation) bool {
    if a.config.Level == AuditLevelNone {
        return false
    }
    if a.config.Level == AuditLevelChangesOnly && op == OperationQuery {
        return false
    }
    return true
}

// ==================== Create ====================

func (a *Audit) beforeCreate(db *gorm.DB) {
    if a.shouldSkip(db) || !a.shouldAuditForLevel(OperationCreate) {
        return
    }

    data := &auditData{
        startTime: time.Now(),
        newValues: make(map[string]any),
    }

    // 获取即将插入的值
    if db.Statement != nil && db.Statement.Dest != nil {
        data.newValues = a.extractValues(db.Statement.Dest)
    }

    db.InstanceSet(auditContextKey{}, data)
}

func (a *Audit) afterCreate(db *gorm.DB) {
    a.processAudit(db, OperationCreate)
}

// ==================== Update ====================

func (a *Audit) beforeUpdate(db *gorm.DB) {
    if a.shouldSkip(db) || !a.shouldAuditForLevel(OperationUpdate) {
        return
    }

    data := &auditData{
        startTime: time.Now(),
        oldValues: make(map[string]any),
    }

    // 查询旧值（需要跳过审计）
    if db.Statement != nil && db.Statement.Model != nil {
        data.oldValues = a.queryOldValues(db)
    }

    db.InstanceSet(auditContextKey{}, data)
}

func (a *Audit) afterUpdate(db *gorm.DB) {
    a.processAudit(db, OperationUpdate)
}

// ==================== Delete ====================

func (a *Audit) beforeDelete(db *gorm.DB) {
    if a.shouldSkip(db) || !a.shouldAuditForLevel(OperationDelete) {
        return
    }

    data := &auditData{
        startTime: time.Now(),
        oldValues: make(map[string]any),
    }

    // 删除前记录旧值
    if db.Statement != nil && db.Statement.Model != nil {
        data.oldValues = a.queryOldValues(db)
    }

    db.InstanceSet(auditContextKey{}, data)
}

func (a *Audit) afterDelete(db *gorm.DB) {
    a.processAudit(db, OperationDelete)
}

// ==================== Query ====================

func (a *Audit) beforeQuery(db *gorm.DB) {
    if a.shouldSkip(db) || !a.shouldAuditForLevel(OperationQuery) {
        return
    }

    data := &auditData{
        startTime: time.Now(),
    }

    db.InstanceSet(auditContextKey{}, data)
}

func (a *Audit) afterQuery(db *gorm.DB) {
    a.processAudit(db, OperationQuery)
}

// ==================== 核心处理逻辑 ====================

// processAudit 处理审计事件
func (a *Audit) processAudit(db *gorm.DB, op Operation) {
    if a.shouldSkip(db) || !a.shouldAuditForLevel(op) {
        return
    }

    // 获取审计上下文数据
    data, ok := db.InstanceGet(auditContextKey{})
    if !ok {
        return
    }
    auditCtx := data.(*auditData)

    // 获取上下文信息
    ctx := db.Statement.Context
    if ctx == nil {
        ctx = context.Background()
    }

    // 构建事件
    event := &AuditEvent{
        Timestamp:  auditCtx.startTime,
        Operation:  op,
        Table:      db.Statement.Table,
        PrimaryKey: a.extractPrimaryKey(db),
        OldValues:  auditCtx.oldValues,
        NewValues:  a.extractValues(db.Statement.Dest),
        SQL:        db.Statement.SQL.String(),
        SQLArgs:    db.StatementVars,
        UserID:     a.getContextValue(ctx, a.config.ContextKeys.UserID),
        Username:   a.getContextValue(ctx, a.config.ContextKeys.Username),
        IP:         a.getContextValue(ctx, a.config.ContextKeys.IP),
        UserAgent:  a.getContextValue(ctx, a.config.ContextKeys.UserAgent),
        RequestID:  a.getContextValue(ctx, a.config.ContextKeys.RequestID),
    }

    // 异步分发事件
    a.dispatcher.Dispatch(ctx, event)
}

// extractValues 从模型中提取值
func (a *Audit) extractValues(dest any) map[string]any {
    values := make(map[string]any)

    if dest == nil {
        return values
    }

    v := reflect.ValueOf(dest)
    for v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    switch v.Kind() {
    case reflect.Struct:
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            field := t.Field(i)
            // 跳过非导出字段
            if !field.IsExported() {
                continue
            }
            // 获取 gorm 标签
            tag := field.Tag.Get("gorm")
            if tag != "" && !strings.Contains(tag, "-") {
                // 使用 gorm 标签中的列名，或使用字段名的蛇形命名
                columnName := a.getColumnName(tag, field.Name)
                values[columnName] = v.Field(i).Interface()
            }
        }
    case reflect.Map:
        for _, key := range v.MapKeys() {
            if key.Kind() == reflect.String {
                values[key.String()] = v.MapIndex(key).Interface()
            }
        }
    }

    return values
}

// getColumnName 从 gorm 标签或字段名获取列名
func (a *Audit) getColumnName(tag, fieldName string) string {
    // 解析 gorm 标签获取 column 值
    parts := strings.Split(tag, ";")
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if strings.HasPrefix(part, "column:") {
            return strings.TrimPrefix(part, "column:")
        }
    }
    // 默认返回字段名
    return fieldName
}

// extractPrimaryKey 提取主键值
func (a *Audit) extractPrimaryKey(db *gorm.DB) string {
    if db.Statement == nil || db.Statement.Schema == nil {
        return ""
    }

    // 获取主键字段名
    primaryFields := db.Statement.Schema.PrimaryFields
    if len(primaryFields) == 0 {
        return ""
    }

    var keyParts []string
    for _, pf := range primaryFields {
        if db.Statement.Dest != nil {
            v := reflect.ValueOf(db.Statement.Dest)
            for v.Kind() == reflect.Ptr {
                v = v.Elem()
            }
            if v.Kind() == reflect.Struct {
                f := v.FieldByName(pf.Name)
                if f.IsValid() {
                    keyParts = append(keyParts, fmt.Sprintf("%v", f.Interface()))
                }
            }
        }
    }

    return strings.Join(keyParts, ",")
}

// queryOldValues 查询数据库中的旧值
func (a *Audit) queryOldValues(db *gorm.DB) map[string]any {
    oldValues := make(map[string]any)

    // 克隆查询以避免修改原始查询
    query := db.Session(&gorm.Session{SkipHooks: true})

    // 创建一个空的模型来接收数据
    dest := reflect.New(reflect.TypeOf(db.Statement.Model)).Interface()

    // 执行查询（跳过审计）
    err := SkipAudit(query).First(dest).Error
    if err != nil {
        return oldValues
    }

    return a.extractValues(dest)
}

// getContextValue 从 context 获取值
func (a *Audit) getContextValue(ctx context.Context, key any) string {
    if key == nil || ctx == nil {
        return ""
    }

    var val any
    if strKey, ok := key.(string); ok {
        val = ctx.Value(strKey)
    } else {
        val = ctx.Value(key)
    }

    if val == nil {
        return ""
    }

    if str, ok := val.(string); ok {
        return str
    }
    return fmt.Sprintf("%v", val)
}
```

**Step 2: 运行 go build 验证语法**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误（注意需要添加 fmt 导入）

**Step 3: 修复 import**

在 `callback.go` 顶部添加 `"fmt"`:

```go
import (
    "context"
    "fmt"
    "gorm.io/gorm"
    "reflect"
    "strings"
    "time"
)
```

**Step 4: 再次运行 go build 验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 无错误

**Step 5: 创建单元测试**

创建: `gorm-audit/callback_test.go`

```go
package audit

import (
    "context"
    "testing"
    "time"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type TestUser struct {
    ID   uint `gorm:"primaryKey"`
    Name string
    Age  int
}

func TestSkipAudit(t *testing.T) {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    p := New(nil)
    p.Initialize(db)

    db2 := SkipAudit(db)
    if _, ok := db2.InstanceGet("audit:skip"); !ok {
        t.Error("skip flag should be set")
    }
}

func TestExtractValues(t *testing.T) {
    p := New(nil)

    user := &TestUser{ID: 1, Name: "test", Age: 25}
    values := p.extractValues(user)

    if values["ID"] != uint(1) {
        t.Errorf("expected ID 1, got %v", values["ID"])
    }
    if values["Name"] != "test" {
        t.Errorf("expected Name 'test', got %v", values["Name"])
    }
    if values["Age"] != 25 {
        t.Errorf("expected Age 25, got %v", values["Age"])
    }
}

func TestExtractValuesNil(t *testing.T) {
    p := New(nil)
    values := p.extractValues(nil)

    if len(values) != 0 {
        t.Errorf("expected empty values for nil, got %d", len(values))
    }
}

func TestGetContextValue(t *testing.T) {
    p := New(nil)
    ctx := context.Background()

    // 空 context
    val := p.getContextValue(ctx, "test_key")
    if val != "" {
        t.Errorf("expected empty string, got %v", val)
    }

    // 带 value 的 context
    ctx = context.WithValue(ctx, "user_id", "12345")
    val = p.getContextValue(ctx, "user_id")
    if val != "12345" {
        t.Errorf("expected '12345', got %v", val)
    }
}
```

**Step 6: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v`
Expected: PASS

**Step 7: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/callback.go gorm-audit/callback_test.go
git commit -m "feat(audit): implement GORM callback handlers"
```

---

## Task 8: 创建内置 Console Handler（handler/console.go）

**文件:**
- 创建: `gorm-audit/handler/console.go`

**Step 1: 创建 console.go 文件**

```go
package handler

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// ConsoleHandler 控制台输出处理器
type ConsoleHandler struct {
    PrettyPrint bool
}

// NewConsoleHandler 创建控制台处理器
func NewConsoleHandler() *ConsoleHandler {
    return &ConsoleHandler{
        PrettyPrint: true,
    }
}

// Handle 实现 EventHandler 接口
func (h *ConsoleHandler) Handle(ctx context.Context, event *Event) error {
    if h.PrettyPrint {
        h.printPretty(event)
    } else {
        h.printJSON(event)
    }
    return nil
}

// printPretty 美化打印
func (h *ConsoleHandler) printPretty(event *Event) {
    fmt.Printf("===== Audit Event =====\n")
    fmt.Printf("Time:      %s\n", event.Timestamp)
    fmt.Printf("Operation: %s\n", event.Operation)
    fmt.Printf("Table:     %s\n", event.Table)
    fmt.Printf("PrimaryKey: %s\n", event.PrimaryKey)
    if event.Data["user_id"] != nil {
        fmt.Printf("UserID:    %v\n", event.Data["user_id"])
    }
    if event.Data["username"] != nil {
        fmt.Printf("Username:  %v\n", event.Data["username"])
    }
    if event.Data["ip"] != nil {
        fmt.Printf("IP:        %v\n", event.Data["ip"])
    }
    fmt.Printf("=======================\n")
}

// printJSON JSON 格式打印
func (h *ConsoleHandler) printJSON(event *Event) {
    data, _ := json.MarshalIndent(event, "", "  ")
    fmt.Println(string(data))
}
```

**Step 2: 更新 handler/handler.go 中的 Event 定义**

修改: `gorm-audit/handler/handler.go:19-25`

```go
// Event 事件定义
type Event struct {
    Timestamp  string
    Operation  string
    Table      string
    PrimaryKey string
    Data       map[string]any
}
```

**Step 3: 更新 callback.go 中的事件转换**

修改: `gorm-audit/callback.go:256-272` (processAudit 方法中)

需要更新 dispatcher.Dispatch 调用，将 AuditEvent 转换为 handler.Event:

```go
// 构建发送给 handler 的事件
    handlerEvent := &handler.Event{
        Timestamp:  auditCtx.startTime.Format(time.RFC3339),
        Operation:  string(op),
        Table:      db.Statement.Table,
        PrimaryKey: a.extractPrimaryKey(db),
        Data:       make(map[string]any),
    }

    // 合并旧值和新值到 Data
    if len(auditCtx.oldValues) > 0 {
        handlerEvent.Data["old_values"] = auditCtx.oldValues
    }
    if len(auditCtx.newValues) > 0 {
        handlerEvent.Data["new_values"] = auditCtx.newValues
    }
    handlerEvent.Data["user_id"] = a.getContextValue(ctx, a.config.ContextKeys.UserID)
    handlerEvent.Data["username"] = a.getContextValue(ctx, a.config.ContextKeys.Username)
    handlerEvent.Data["ip"] = a.getContextValue(ctx, a.config.ContextKeys.IP)
    handlerEvent.Data["user_agent"] = a.getContextValue(ctx, a.config.ContextKeys.UserAgent)
    handlerEvent.Data["request_id"] = a.getContextValue(ctx, a.config.ContextKeys.RequestID)
    handlerEvent.Data["sql"] = db.Statement.SQL.String()
    handlerEvent.Data["sql_args"] = db.StatementVars

    // 异步分发事件
    a.dispatcher.Dispatch(ctx, handlerEvent)
```

等等，这样设计不太合理。让我重新思考一下架构...

实际上，我们应该让 EventHandler 接口直接接收 AuditEvent，而不是创建一个新的事件类型。让我修改设计：

**Step 2 (修正): 更新 handler/handler.go 修改 EventHandler 接口**

修改: `gorm-audit/handler/handler.go:5-26`

```go
package handler

import (
    "context"
    "time"
)

// Operation 操作类型（与 audit.Operation 相同）
type Operation string

const (
    OperationCreate Operation = "create"
    OperationUpdate Operation = "update"
    OperationDelete Operation = "delete"
    OperationQuery  Operation = "query"
)

// Event 事件定义（完整版本）
type Event struct {
    Timestamp  time.Time
    Operation  Operation
    Table      string
    PrimaryKey string
    OldValues  map[string]any
    NewValues  map[string]any
    SQL        string
    SQLArgs    []any
    UserID     string
    Username   string
    IP         string
    UserAgent  string
    RequestID  string
}

// EventHandler 事件处理器接口
type EventHandler interface {
    Handle(ctx context.Context, event *Event) error
}

// EventHandlerFunc 函数式事件处理器
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
    return f(ctx, event)
}
```

**Step 3: 更新 callback.go 使用正确的类型**

修改: `gorm-audit/callback.go:1-7` 添加 handler 包导入

```go
import (
    "context"
    "fmt"
    "gorm.io/gorm"
    "reflect"
    "strings"
    "time"

    "github.com/piwriw/gorm/gorm-audit/handler"
)
```

修改: `gorm-audit/callback.go:256-272` 的 dispatcher.Dispatch 调用：

```go
// 构建事件
    event := &handler.Event{
        Timestamp:  auditCtx.startTime,
        Operation:  handler.Operation(op),
        Table:      db.Statement.Table,
        PrimaryKey: a.extractPrimaryKey(db),
        OldValues:  auditCtx.oldValues,
        NewValues:  a.extractValues(db.Statement.Dest),
        SQL:        db.Statement.SQL.String(),
        SQLArgs:    db.StatementVars,
        UserID:     a.getContextValue(ctx, a.config.ContextKeys.UserID),
        Username:   a.getContextValue(ctx, a.config.ContextKeys.Username),
        IP:         a.getContextValue(ctx, a.config.ContextKeys.IP),
        UserAgent:  a.getContextValue(ctx, a.config.ContextKeys.UserAgent),
        RequestID:  a.getContextValue(ctx, a.config.ContextKeys.RequestID),
    }

    // 异步分发事件
    a.dispatcher.Dispatch(ctx, event)
```

**Step 4: 更新 dispatcher.go 中的类型引用**

修改: `gorm-audit/dispatcher.go:1-8` 添加 handler 导入，移除 AuditEvent

```go
package audit

import (
    "context"
    "fmt"
    "log"
    "runtime"
    "sync"

    "github.com/piwriw/gorm/gorm-audit/handler"
)
```

修改: `gorm-audit/dispatcher.go:23-24, 32, 48, 54, 62`

将所有 `AuditEvent` 替换为 `handler.Event`:

```go
// Dispatcher 事件分发器
type Dispatcher struct {
    handlers []handler.EventHandler
    mu       sync.RWMutex
}

// NewDispatcher 创建新的分发器
func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make([]handler.EventHandler, 0),
    }
}

// Add 添加事件处理器
func (d *Dispatcher) Add(h handler.EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers = append(d.handlers, h)
}

// Dispatch 分发事件到所有处理器
func (d *Dispatcher) Dispatch(ctx context.Context, event *handler.Event) {
    // ... rest of code
}

func (d *Dispatcher) safeHandle(ctx context.Context, h handler.EventHandler, event *handler.Event) {
    // ... rest of code
}

func (d *Dispatcher) handlePanic(r any, h handler.EventHandler, event *handler.Event) {
    // ... rest of code
}
```

**Step 5: 更新 audit.go 中的 Use 方法**

修改: `gorm-audit/audit.go:48-51`

```go
// Use 添加事件处理器
func (a *Audit) Use(h handler.EventHandler) *Audit {
    a.dispatcher.Add(h)
    return a
}
```

**Step 6: 运行 go build 验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go mod tidy && go build`
Expected: 无错误

**Step 7: 更新 console.go 适配新接口**

修改: `gorm-audit/handler/console.go:29-45`

```go
// Handle 实现 EventHandler 接口
func (h *ConsoleHandler) Handle(ctx context.Context, event *Event) error {
    if h.PrettyPrint {
        h.printPretty(event)
    } else {
        h.printJSON(event)
    }
    return nil
}

// printPretty 美化打印
func (h *ConsoleHandler) printPretty(event *Event) {
    fmt.Printf("===== Audit Event =====\n")
    fmt.Printf("Time:      %s\n", event.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("Operation: %s\n", event.Operation)
    fmt.Printf("Table:     %s\n", event.Table)
    fmt.Printf("PrimaryKey: %s\n", event.PrimaryKey)
    if event.UserID != "" {
        fmt.Printf("UserID:    %s\n", event.UserID)
    }
    if event.Username != "" {
        fmt.Printf("Username:  %s\n", event.Username)
    }
    if event.IP != "" {
        fmt.Printf("IP:        %s\n", event.IP)
    }
    if len(event.OldValues) > 0 {
        fmt.Printf("OldValues: %+v\n", event.OldValues)
    }
    if len(event.NewValues) > 0 {
        fmt.Printf("NewValues: %+v\n", event.NewValues)
    }
    fmt.Printf("=======================\n")
}

// printJSON JSON 格式打印
func (h *ConsoleHandler) printJSON(event *Event) {
    data, _ := json.MarshalIndent(event, "", "  ")
    fmt.Println(string(data))
}
```

**Step 8: 创建 console handler 测试**

创建: `gorm-audit/handler/console_test.go`

```go
package handler

import (
    "context"
    "testing"
    "time"
)

func TestConsoleHandlerHandle(t *testing.T) {
    h := NewConsoleHandler()

    event := &Event{
        Timestamp:  time.Now(),
        Operation:  OperationCreate,
        Table:      "users",
        PrimaryKey: "1",
        UserID:     "123",
        Username:   "admin",
        NewValues:  map[string]any{"name": "test"},
    }

    err := h.Handle(context.Background(), event)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestConsoleHandlerJSON(t *testing.T) {
    h := &ConsoleHandler{PrettyPrint: false}

    event := &Event{
        Timestamp: time.Now(),
        Operation: OperationUpdate,
        Table:     "users",
    }

    err := h.Handle(context.Background(), event)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

**Step 9: 运行测试验证通过**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v ./...`
Expected: PASS

**Step 10: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/
git commit -m "feat(audit): add console handler and update event types"
```

---

## Task 9: 创建使用示例（example/main.go）

**文件:**
- 创建: `gorm-audit/example/main.go`

**Step 1: 创建示例程序**

```go
package main

import (
    "context"
    "log"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "github.com/piwriw/gorm/gorm-audit"
    "github.com/piwriw/gorm/gorm-audit/handler"
)

// 定义 Context 键
type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    UsernameKey  contextKey = "username"
    IPKey        contextKey = "ip"
    UserAgentKey contextKey = "user_agent"
    RequestIDKey contextKey = "request_id"
)

// User 模型
type User struct {
    ID   uint `gorm:"primaryKey"`
    Name string
    Age  int
}

func main() {
    // 1. 初始化内存数据库
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 自动迁移
    if err := db.AutoMigrate(&User{}); err != nil {
        log.Fatal(err)
    }

    // 2. 创建审计插件
    auditPlugin := audit.New(&audit.Config{
        Level:        audit.AuditLevelAll,
        IncludeQuery: false,
        ContextKeys: audit.ContextKeyConfig{
            UserID:    UserIDKey,
            Username:  UsernameKey,
            IP:        IPKey,
            UserAgent: UserAgentKey,
            RequestID: RequestIDKey,
        },
    })

    // 3. 添加控制台处理器
    auditPlugin.Use(handler.NewConsoleHandler())

    // 4. 注册插件到 GORM
    if err := db.Use(auditPlugin); err != nil {
        log.Fatal(err)
    }

    log.Println("Database initialized with audit plugin")

    // 5. 使用带审计的数据库操作
    ctx := context.Background()
    ctx = context.WithValue(ctx, UserIDKey, "12345")
    ctx = context.WithValue(ctx, UsernameKey, "admin")
    ctx = context.WithValue(ctx, IPKey, "192.168.1.100")
    ctx = context.WithValue(ctx, UserAgentKey, "Mozilla/5.0")
    ctx = context.WithValue(ctx, RequestIDKey, "req-001")

    // Create 操作
    log.Println("\n=== Creating user ===")
    user := &User{Name: "Alice", Age: 25}
    if err := db.WithContext(ctx).Create(user).Error; err != nil {
        log.Fatal(err)
    }

    // Update 操作
    log.Println("\n=== Updating user ===")
    if err := db.WithContext(ctx).Model(user).Update("Age", 26).Error; err != nil {
        log.Fatal(err)
    }

    // Delete 操作
    log.Println("\n=== Deleting user ===")
    if err := db.WithContext(ctx).Delete(user).Error; err != nil {
        log.Fatal(err)
    }

    log.Println("\n=== Demo completed ===")
}
```

**Step 2: 创建 go.mod**

创建: `gorm-audit/example/go.mod`

```go
module github.com/piwriw/gorm/gorm-audit/example

go 1.22

require (
    github.com/piwriw/gorm/gorm-audit v0.0.0
    gorm.io/driver/sqlite v1.5.6
    gorm.io/gorm v1.31.1
)

replace github.com/piwriw/gorm/gorm-audit => ../
```

**Step 3: 运行示例验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit/example && go mod tidy && go run main.go`
Expected: 看到审计事件输出

**Step 4: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/example/
git commit -m "feat(audit): add usage example with console handler"
```

---

## Task 10: 创建 go.mod 和 README

**文件:**
- 创建: `gorm-audit/go.mod`
- 创建: `gorm-audit/README.md`

**Step 1: 创建 go.mod**

```go
module github.com/piwriw/gorm/gorm-audit

go 1.22

require (
    gorm.io/gorm v1.31.1
)
```

**Step 2: 运行 go mod tidy**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go mod tidy`
Expected: 成功下载依赖

**Step 3: 创建 README.md**

```markdown
# GORM Audit Plugin

基于 GORM Plugin 接口的审计插件，支持可配置审计级别、事件驱动架构。

## 特性

- **可配置审计级别**: 支持记录所有操作 / 仅变更 / 不记录
- **事件驱动架构**: 通过 EventHandler 接口灵活处理审计事件
- **自动上下文提取**: 从 context 自动获取用户、IP 等信息
- **变更前后记录**: 自动记录数据变更前后的值
- **Panic 恢复**: 事件处理器异常不影响主业务
- **跳过审计**: 支持临时跳过审计（如批量操作）

## 安装

```bash
go get github.com/piwriw/gorm/gorm-audit
```

## 快速开始

```go
package main

import (
    "context"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/piwriw/gorm/gorm-audit"
    "github.com/piwriw/gorm/gorm-audit/handler"
)

func main() {
    db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    // 创建审计插件
    auditPlugin := audit.New(&audit.Config{
        Level: audit.AuditLevelChangesOnly,
    })

    // 添加处理器
    auditPlugin.Use(handler.NewConsoleHandler())

    // 注册插件
    db.Use(auditPlugin)

    // 使用带审计的操作
    ctx := context.WithValue(context.Background(), "user_id", "123")
    db.WithContext(ctx).Create(&User{Name: "test"})
}
```

## 配置

### 审计级别

```go
type AuditLevel int

const (
    AuditLevelAll         // 记录所有操作（包括查询）
    AuditLevelChangesOnly // 仅记录变更操作（Create/Update/Delete）
    AuditLevelNone        // 不记录
)
```

### Context 配置

```go
ContextKeys: audit.ContextKeyConfig{
    UserID:    "user_id",    // context key
    Username:  "username",
    IP:        "ip",
    UserAgent: "user_agent",
    RequestID: "request_id",
}
```

## 内置处理器

### ConsoleHandler

控制台输出审计事件：

```go
auditPlugin.Use(handler.NewConsoleHandler())
```

## 自定义处理器

实现 EventHandler 接口：

```go
type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, event *audit.Event) error {
    // 处理审计事件
    return nil
}

auditPlugin.Use(&MyHandler{})
```

## 跳过审计

```go
import "github.com/piwriw/gorm/gorm-audit"

// 批量操作时跳过审计
audit.SkipAudit(db).CreateInBatches(users, 100)
```

## 事件结构

```go
type Event struct {
    Timestamp  time.Time
    Operation  Operation     // create/update/delete/query
    Table      string
    PrimaryKey string
    OldValues  map[string]any
    NewValues  map[string]any
    SQL        string
    SQLArgs    []any
    UserID     string
    Username   string
    IP         string
    UserAgent  string
    RequestID  string
}
```

## License

MIT
```

**Step 4: 提交**

```bash
cd /Users/joohwan/GolandProjects/go-tools
git add gorm-audit/go.mod gorm-audit/README.md
git commit -m "feat(audit): add go.mod and README documentation"
```

---

## 完成检查清单

- [ ] 所有单元测试通过
- [ ] 示例程序正常运行
- [ ] 代码无编译警告
- [ ] README 文档完整

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v ./...`
Expected: 全部 PASS

---

**计划完成！**
