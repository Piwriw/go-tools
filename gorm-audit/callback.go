package audit

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"gorm.io/gorm"
)

// auditContextKey 用于在 GORM context 中存储审计数据
const auditContextKey = "audit_context"

// auditData 存储审计过程中的临时数据
type auditData struct {
	startTime string
	oldValues map[string]any
}

// SkipAudit 跳过当前操作的审计
func SkipAudit(db *gorm.DB) *gorm.DB {
	return db.InstanceSet(auditContextKey, &auditData{})
}

// shouldSkip 检查是否应该跳过审计
func (a *Audit) shouldSkip(db *gorm.DB) bool {
	_, ok := db.InstanceGet(auditContextKey)
	return ok
}

// shouldAuditForLevel 检查是否应该审计该操作
func (a *Audit) shouldAuditForLevel(op Operation, db *gorm.DB) bool {
	switch a.config.Level {
	case AuditLevelAll:
		return true
	case AuditLevelChangesOnly:
		return op != OperationQuery
	default:
		return false
	}
}

// ==================== Create Callbacks ====================

func (a *Audit) beforeCreate(db *gorm.DB) {
	if a.shouldSkip(db) {
		return
	}

	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	auditCtx := &auditData{
		startTime: db.Statement.DB.NowFunc().Format("2006-01-02T15:04:05.000"),
		oldValues: make(map[string]any),
	}

	// 对于 Create 操作，oldValues 为空（因为是新记录）
	_ = db.InstanceSet(auditContextKey, auditCtx)
}

func (a *Audit) afterCreate(db *gorm.DB) {
	a.processAudit(db, OperationCreate)
}

// ==================== Update Callbacks ====================

func (a *Audit) beforeUpdate(db *gorm.DB) {
	if a.shouldSkip(db) {
		return
	}

	oldValues := a.queryOldValues(db)

	auditCtx := &auditData{
		startTime: db.Statement.DB.NowFunc().Format("2006-01-02T15:04:05.000"),
		oldValues: oldValues,
	}

	_ = db.InstanceSet(auditContextKey, auditCtx)
}

func (a *Audit) afterUpdate(db *gorm.DB) {
	a.processAudit(db, OperationUpdate)
}

// ==================== Delete Callbacks ====================

func (a *Audit) beforeDelete(db *gorm.DB) {
	if a.shouldSkip(db) {
		return
	}

	oldValues := a.queryOldValues(db)

	auditCtx := &auditData{
		startTime: db.Statement.DB.NowFunc().Format("2006-01-02T15:04:05.000"),
		oldValues: oldValues,
	}

	_ = db.InstanceSet(auditContextKey, auditCtx)
}

func (a *Audit) afterDelete(db *gorm.DB) {
	a.processAudit(db, OperationDelete)
}

// ==================== Query Callbacks ====================

func (a *Audit) beforeQuery(db *gorm.DB) {
	if a.shouldSkip(db) {
		return
	}

	auditCtx := &auditData{
		startTime: db.Statement.DB.NowFunc().Format("2006-01-02T15:04:05.000"),
		oldValues: make(map[string]any),
	}

	_ = db.InstanceSet(auditContextKey, auditCtx)
}

func (a *Audit) afterQuery(db *gorm.DB) {
	a.processAudit(db, OperationQuery)
}

// ==================== Core Processing ====================

// processAudit 核心审计处理逻辑
func (a *Audit) processAudit(db *gorm.DB, op Operation) {
	if a.shouldSkip(db) {
		return
	}

	if !a.shouldAuditForLevel(op, db) {
		return
	}

	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// 获取审计上下文
	auditCtxVal, ok := db.InstanceGet(auditContextKey)
	if !ok {
		return
	}

	auditCtx, ok := auditCtxVal.(*auditData)
	if !ok || auditCtx == nil {
		return
	}

	// 构建 handler.Event 对象
	event := &handler.Event{
		Timestamp:  auditCtx.startTime,
		Operation:  handler.Operation(op),
		Table:      db.Statement.Table,
		PrimaryKey: a.extractPrimaryKey(db),
		OldValues:  auditCtx.oldValues,
		NewValues:  a.extractValues(db.Statement.Dest),
		SQL:        db.Statement.SQL.String(),
		SQLArgs:    db.Statement.Vars,
		UserID:     a.getContextValue(ctx, a.config.ContextKeys.UserID),
		Username:   a.getContextValue(ctx, a.config.ContextKeys.Username),
		IP:         a.getContextValue(ctx, a.config.ContextKeys.IP),
		UserAgent:  a.getContextValue(ctx, a.config.ContextKeys.UserAgent),
		RequestID:  a.getContextValue(ctx, a.config.ContextKeys.RequestID),
	}

	// 分发事件
	a.dispatcher.DispatchHandler(ctx, event)
}

// ==================== Helper Methods ====================

// extractValues 从对象中提取值
func (a *Audit) extractValues(dest any) map[string]any {
	if dest == nil {
		return nil
	}

	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	values := make(map[string]any)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		// 获取 GORM 列名
		tag := field.Tag.Get("gorm")
		columnName := field.Name
		if tag != "" {
			for _, part := range strings.Split(tag, ";") {
				if strings.HasPrefix(part, "column:") {
					columnName = strings.TrimPrefix(part, "column:")
					break
				}
			}
		}

		values[columnName] = fieldValue.Interface()
	}

	return values
}

// extractPrimaryKey 提取主键值
func (a *Audit) extractPrimaryKey(db *gorm.DB) string {
	if db.Statement.SQL.Len() == 0 {
		return ""
	}

	// 尝试从 Statement.ReflectValue 中获取主键
	if db.Statement.ReflectValue.IsValid() {
		primaryFields := db.Statement.Schema.PrimaryFields
		if len(primaryFields) > 0 {
			var keys []string
			for _, pf := range primaryFields {
				if v, ok := pf.ValueOf(db.Statement.Context, db.Statement.ReflectValue); ok {
					keys = append(keys, fmt.Sprintf("%v", v))
				}
			}
			if len(keys) > 0 {
				return strings.Join(keys, ",")
			}
		}
	}

	return ""
}

// queryOldValues 查询旧值（用于 Update 和 Delete）
func (a *Audit) queryOldValues(db *gorm.DB) map[string]any {
	if db.Statement.ReflectValue.IsValid() && db.Statement.Schema != nil {
		primaryFields := db.Statement.Schema.PrimaryFields
		if len(primaryFields) > 0 {
			// 构建查询条件
			conds := []any{}
			for _, pf := range primaryFields {
				if v, ok := pf.ValueOf(db.Statement.Context, db.Statement.ReflectValue); ok {
					conds = append(conds, pf.DBName, v)
				}
			}

			if len(conds) > 0 {
				// 执行查询获取旧值
				var oldValues map[string]any
				err := db.Session(&gorm.Session{NewDB: true}).
					Table(db.Statement.Table).
					Where(conds[0], conds[1:]...).
					Scan(&oldValues).Error

				if err == nil && len(oldValues) > 0 {
					return oldValues
				}
			}
		}
	}

	// 如果无法查询，返回当前值作为旧值（适用于 Soft Delete）
	return a.extractValues(db.Statement.Dest)
}

// getContextValue 从 context 中获取值
func (a *Audit) getContextValue(ctx context.Context, key any) string {
	if key == nil {
		return ""
	}

	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	if str, ok := val.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", val)
}
