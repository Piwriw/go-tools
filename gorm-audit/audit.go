package audit

import (
	"context"
	"sync"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"gorm.io/gorm"
)

// Config 插件完整配置
type Config struct {
	Level         AuditLevel
	IncludeQuery  bool
	ContextKeys   ContextKeyConfig
	UseWorkerPool bool
	WorkerConfig  *WorkerPoolConfig
	Filters       []Filter // 事件过滤器列表
	EnableMetrics bool     // 是否启用指标收集
}

// Audit GORM 审计插件
type Audit struct {
	config     *Config
	dispatcher *Dispatcher
	configMu   sync.RWMutex // 保护 config 的并发访问
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

	var dispatcher *Dispatcher
	if config.UseWorkerPool {
		// Worker Pool 需要一个默认的处理器
		// 这里暂时使用空实现，用户需要通过 Use 方法添加实际的处理器
		defaultHandler := handler.EventHandlerFunc(func(ctx context.Context, event *handler.Event) error {
			return nil
		})
		dispatcher = NewDispatcherWithWorkerPool(defaultHandler, config.WorkerConfig)
	} else {
		dispatcher = NewDispatcher()
	}

	return &Audit{
		config:     config,
		dispatcher: dispatcher,
	}
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
func (a *Audit) Use(h handler.EventHandler) *Audit {
	a.dispatcher.AddHandler(h)
	return a
}

// Reload 重新加载配置（从环境变量）
func (a *Audit) Reload() error {
	level := getAuditLevelFromEnv()

	a.configMu.Lock()
	a.config.Level = level
	a.configMu.Unlock()

	return nil
}

// GetLevel 线程安全地获取当前审计级别
func (a *Audit) GetLevel() AuditLevel {
	a.configMu.RLock()
	defer a.configMu.RUnlock()
	return a.config.Level
}

// Metrics 返回 Prometheus 格式的指标
func (a *Audit) Metrics() string {
	if a.dispatcher != nil {
		return a.dispatcher.Metrics()
	}
	return ""
}
