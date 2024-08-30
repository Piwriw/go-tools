package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Age      int
	GameUser GameUsers `json:"game_user" gorm:"type:json"`
}
type GameUser struct {
	ID       int
	Name     string
	PassWord string
}
type GameUsers []GameUser

func (e GameUsers) Value() (driver.Value, error) {
	str, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(str), nil
}

func (e *GameUsers) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Unmarshal the json.RawMessage into an InstanceData struct
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, e); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), e); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(v.([]byte), e); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	dsn := "root:123456@tcp(10.0.0.197:3303)/joohwan_dev?parseTime=true"

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//mysqlDb.AutoMigrate(&User{})
	//insert(mysqlDb)
	fmt.Println(getById(mysqlDb))
}

func insert(db *gorm.DB) {
	games := []GameUser{{ID: 1, Name: "gameuser1"}, {ID: 1, Name: "gameuser1"}}
	user := User{
		Name:     "user1",
		Age:      18,
		GameUser: games,
	}
	if err := db.Debug().Create(&user).Error; err != nil {
		panic(err)
	}
}

func getById(db *gorm.DB) User {
	user := User{}
	if err := db.Debug().First(&user).Where("id", 1).Error; err != nil {
		panic(err)
	}
	return user
}
