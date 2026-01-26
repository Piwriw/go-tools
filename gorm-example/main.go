package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func main() {
	// 构建 MySQL DSN（Data Source Name）
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// 使用 GORM 连接 MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	s, err := db.DB()
	s.Stats() // 自动迁移：创建表
	//err = db.AutoMigrate(&User{})
	//if err != nil {
	//	log.Fatalf("迁移表失败: %v", err)
	//}
	// 查询数据
	user := &User{ID: 1,
		Name: "user1",
		Age:  18,
	}

	if err := db.Where(user).First(user).Error; err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	//// 查询数据
	//user = &User{}
	//if err := db.Where("id = ?", 1).First(user); err != nil {
	//	log.Fatalf("Failed to find user: %v", err)
	//}
}
