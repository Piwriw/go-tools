package main

import (
	"fmt"
	"github.com/piwriw/gorm/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s",
		"10.0.0.192", "3306", "root", "123456", "joohwan_dev")))
	if err != nil {
		panic(err)
	}
	var ms []model.HestiaInstanceModel
	db.Select("name", "age").Find(&ms)

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
