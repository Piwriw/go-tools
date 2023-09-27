package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

/*
go get github.com/joho/godotenv

godotenv库从.env文件中读取配置 (默认根目录的.env)

	 _ "github.com/joho/godotenv/autoload"
	默认会帮你进行load
*/
func test() {
	if err := godotenv.Load("./godotenv/.env"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("author", os.Getenv("author"))
}
