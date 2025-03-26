package logger

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestName(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", "root", "root", "10.0.0.195", 3306,
		"mysql")
	// dsn := "root:123456@tcp(10.0.0.197:3306)/joohwan_dev?parseTime=true"
	conf := GormLogger{
		SQLConfig: &SQLConfig{
			SQLFile:      "./sql.log",
			MaxSQLLength: 5,
		},
	}
	gormLogger := SetGormLogger(conf)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		ID int `gorm:"column:id"`
	}

	// 使用Raw执行查询并Scan结果
	if err := db.Raw("SELECT 1 as id").Scan(&result).Error; err != nil {
		t.Fatalf("Query failed: %v", err)
	}
}
