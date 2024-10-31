package main

import (
	"github.com/piwriw/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID int
	// 其他字段...
}

func main() {

	// 打印 JSON
	//fmt.Println(string(jsonData))
	pgDb, err := gorm.Open(postgres.Open("postgres://yunqu:YunquTech01*@@10.0.0.195:15432/devops"), &gorm.Config{
		Logger: logger.SetGormLogger(),
	})
	if err != nil {
		panic(err)
	}
	sql := `select id from athena_config_alertsddd `

	for i := 0; i < 1000; i++ {
		sql += "union all select id from athena_config_alertsddd"
		//var user User
		//pgDb.First(&user)
		//fmt.Println(user)
		//fmt.Println(user.ID)
		//fmt.Println(user.Name)
		//fmt.Println(user.Age)
	}

	pgDb.Exec(sql)
	//var id int
}
