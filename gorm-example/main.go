package main

import (
	"context"
	"github.com/piwriw/gorm/dao"
	"github.com/piwriw/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
)

func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:abc123@47.107.113.111:55433/postgres"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		err = dao.Use(db).Emp.WithContext(context.Background()).Create(&model.Emp{
			ID:       int32(3 + i),
			Name:     "joohwan",
			Addresss: "hangzhou",
			Age:      rand.Int31(),
		})
	}

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
