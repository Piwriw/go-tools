package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接 MySQL 数据库
	db, err := sql.Open("mysql", "root:123456@tcp(10.0.0.197:3303)/joohwan_dev")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 设置一个事务超时时间（这里设置为 10 秒）
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback() // 事务结束时回滚

	// 执行一个阻塞的查询
	var id int
	query := "SELECT 1 FROM users "
	err = tx.QueryRow(query).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	// 模拟长时间运行的操作
	time.Sleep(15 * time.Minute)

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("阻塞会话已完成")
}
