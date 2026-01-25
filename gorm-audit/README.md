# GORM Audit Plugin

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
