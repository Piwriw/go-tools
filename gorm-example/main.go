package main

import (
	"fmt"
	"github.com/piwriw/gorm/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", "root", "123456", "10.0.0.197", 3308, "joohwan_dev")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	var ms []model.HestiaInstanceModel
	db.Select("name", "age").Find(&ms)

	db.Where("name = ?", "jinzhu").Session(&gorm.Session{NewDB: true})

	//if err != nil {
	//	panic(err)
	//}
	//result, err := dao.Use(db).Emp.WithContext(context.Background()).GetByID(1)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("User:%v", result)
	//
	//name, err := dao.Use(db).Emp.WithContext(context.Background()).GetByName("joohwan")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("User:%v", result)
	//all, err := dao.Use(db).Emp.WithContext(context.Background()).GetEmps(20)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("User:%v", len(all))
}
