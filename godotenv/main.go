package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}
	if err := godotenv.Load(fmt.Sprintf("./godotenv/.env.%s", env)); err != nil {
		fmt.Println(os.Getenv("env"))
	}
}
