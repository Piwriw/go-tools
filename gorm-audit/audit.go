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

func (a *Audit) beforeCreate(db *gorm.DB) {}
func (a *Audit) afterCreate(db *gorm.DB)  {}
func (a *Audit) beforeUpdate(db *gorm.DB) {}
func (a *Audit) afterUpdate(db *gorm.DB)  {}
func (a *Audit) beforeDelete(db *gorm.DB) {}
func (a *Audit) afterDelete(db *gorm.DB)  {}
func (a *Audit) beforeQuery(db *gorm.DB)  {}
func (a *Audit) afterQuery(db *gorm.DB)   {}
