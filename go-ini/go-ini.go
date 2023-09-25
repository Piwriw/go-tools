package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

/*
install  go get gopkg.in/ini.v1
*/
func main() {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("app_name:%s\n mysql[port]:%s\n", cfg.Section("").Key("app_name").String(),
		cfg.Section("mysql").Key("port").String())
}
