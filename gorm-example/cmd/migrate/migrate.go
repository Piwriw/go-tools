package main

import (
	"github.com/piwriw/gorm/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// 迁移数据库
func Migrate(db *gorm.DB) {
	db.AutoMigrate(model.HestiaInstanceModel{})

}
func main() {

	dsn := "root:123456@tcp(10.0.0.197:3306)/joohwan_dev?parseTime=true"

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	pgDb, err := gorm.Open(postgres.Open("postgres://yunqu:YunquTech01*@@10.0.0.195:15432/devops"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Migrate(mysqlDb)
	migrateData(pgDb, mysqlDb)
}
func migrateData(dgPg, dbMySql *gorm.DB) {
	tx := dbMySql.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// 查询 PostgreSQL 中的数据
	var acms []model.HestiaInstanceModel

	if err := dgPg.Find(&acms).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to retrieve data from PostgreSQL: %v", err)
	}

	// 将数据插入到 MySQL 中
	for _, result := range acms {
		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert data into MySQL: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to commit transaction: %v", err)
	}

}
