package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/piwriw/gorm/gorm-audit"
	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/piwriw/gorm/gorm-audit/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Name    string `gorm:"size:100" json:"name"`
	Email   string `gorm:"size:100;uniqueIndex" json:"email"`
	Age     int    `json:"age"`
	Country string `gorm:"size:50" json:"country"`
	gorm.Model
}

// Product 产品模型
type Product struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Name  string  `gorm:"size:200" json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	gorm.Model
}

// 自定义 context key 类型
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	usernameKey  contextKey = "username"
	ipKey        contextKey = "ip"
	userAgentKey contextKey = "user_agent"
	requestIDKey contextKey = "request_id"
)

func main() {
	// 创建数据库连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}, &Product{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// 创建审计插件 - 演示事件过滤功能
	auditPlugin := audit.New(&audit.Config{
		Level:        audit.AuditLevelAll, // 记录所有操作
		IncludeQuery: true,                // 包含查询操作
		ContextKeys: audit.ContextKeyConfig{
			UserID:    userIDKey,
			Username:  usernameKey,
			IP:        ipKey,
			UserAgent: userAgentKey,
			RequestID: requestIDKey,
		},
		// 添加事件过滤器
		Filters: []audit.Filter{
			// 只审计 users 和 products 表（白名单模式）
			audit.NewTableFilter(audit.FilterModeWhitelist, []string{"users", "products"}),

			// 只审计 create 和 update 操作
			audit.NewOperationFilter([]types.Operation{
				types.OperationCreate,
				types.OperationUpdate,
			}),

			// 排除测试用户（黑名单模式）
			audit.NewUserFilter(audit.FilterModeBlacklist, []string{"test_user"}),
		},
	})

	// 添加控制台处理器（带颜色）
	consoleHandler := handler.NewConsoleHandler()
	consoleHandler.SetColor(true) // 启用颜色
	auditPlugin.Use(consoleHandler)

	// 添加自定义处理器
	customHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
		// 自定义处理逻辑
		if event.Operation == handler.OperationDelete {
			log.Printf("[CUSTOM] Delete operation detected on table: %s, key: %s", event.Table, event.PrimaryKey)
		}
		return nil
	})
	auditPlugin.Use(customHandler)

	// 初始化插件
	if err := db.Use(auditPlugin); err != nil {
		log.Fatalf("Failed to initialize audit plugin: %v", err)
	}

	fmt.Println("=== GORM Audit Plugin Demo ===\n")

	// 创建 context
	ctx := context.WithValue(context.Background(), userIDKey, "12345")
	ctx = context.WithValue(ctx, usernameKey, "admin")
	ctx = context.WithValue(ctx, ipKey, "192.168.1.100")
	ctx = context.WithValue(ctx, userAgentKey, "Mozilla/5.0")
	ctx = context.WithValue(ctx, requestIDKey, "req-001")

	// 1. Create 操作
	fmt.Println("\n--- Create Operations ---")
	createUser(ctx, db, "Alice", "alice@example.com", 25, "USA")
	createUser(ctx, db, "Bob", "bob@example.com", 30, "Canada")
	createProduct(ctx, db, "Laptop", 999.99, 50)
	createProduct(ctx, db, "Mouse", 29.99, 200)

	// 2. Read 操作
	fmt.Println("\n--- Read Operations ---")
	var users []User
	db.WithContext(ctx).Where("age > ?", 20).Find(&users)
	fmt.Printf("Found %d users\n", len(users))

	var products []Product
	db.WithContext(ctx).Where("price > ?", 100).Find(&products)
	fmt.Printf("Found %d products\n", len(products))

	// 3. Update 操作
	fmt.Println("\n--- Update Operations ---")
	updateUser(ctx, db, 1, map[string]any{
		"age":     26,
		"country": "UK",
	})

	// 4. Delete 操作
	fmt.Println("\n--- Delete Operations ---")
	deleteUser(ctx, db, 2)

	// 5. Skip Audit 示例
	fmt.Println("\n--- Skip Audit Example ---")
	skipAuditUser(ctx, db, "Charlie", "charlie@example.com", 35, "Australia")

	// 6. 配置热更新示例
	fmt.Println("\n--- Config Reload Example ---")
	fmt.Printf("Current audit level: %v\n", auditPlugin.GetLevel())

	// 演示配置热更新
	fmt.Println("Attempting to reload config...")
	if err := auditPlugin.Reload(); err != nil {
		fmt.Printf("Reload failed (expected if no env vars set): %v\n", err)
	} else {
		fmt.Printf("Reload successful! New audit level: %v\n", auditPlugin.GetLevel())
	}

	// 提示如何使用环境变量
	fmt.Println("\n💡 Tip: You can control audit level via environment variables:")
	fmt.Println("   export GORM_AUDIT_LEVEL=all          # Audit all operations")
	fmt.Println("   export GORM_AUDIT_LEVEL=changes_only # Audit only changes")
	fmt.Println("   export GORM_AUDIT_LEVEL=none         # Disable auditing")

	// 等待异步处理器完成
	fmt.Println("\nWaiting for async handlers to complete...")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n=== Demo Complete ===")
}

// createUser 创建用户
func createUser(ctx context.Context, db *gorm.DB, name, email string, age int, country string) {
	user := User{
		Name:    name,
		Email:   email,
		Age:     age,
		Country: country,
	}
	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		fmt.Printf("Created user: %s (ID: %d)\n", name, user.ID)
	}
}

// createProduct 创建产品
func createProduct(ctx context.Context, db *gorm.DB, name string, price float64, stock int) {
	product := Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}
	if err := db.WithContext(ctx).Create(&product).Error; err != nil {
		log.Printf("Failed to create product: %v", err)
	} else {
		fmt.Printf("Created product: %s (ID: %d)\n", name, product.ID)
	}
}

// updateUser 更新用户
func updateUser(ctx context.Context, db *gorm.DB, id uint, updates map[string]any) {
	if err := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Printf("Updated user ID: %d\n", id)
	}
}

// deleteUser 删除用户
func deleteUser(ctx context.Context, db *gorm.DB, id uint) {
	if err := db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		log.Printf("Failed to delete user: %v", err)
	} else {
		fmt.Printf("Deleted user ID: %d\n", id)
	}
}

// skipAuditUser 跳过审计创建用户
func skipAuditUser(ctx context.Context, db *gorm.DB, name, email string, age int, country string) {
	user := User{
		Name:    name,
		Email:   email,
		Age:     age,
		Country: country,
	}
	if err := db.WithContext(ctx).Scopes(audit.SkipAudit).Create(&user).Error; err != nil {
		log.Printf("Failed to create user (skip audit): %v", err)
	} else {
		fmt.Printf("Created user (audit skipped): %s (ID: %d)\n", name, user.ID)
	}
}
