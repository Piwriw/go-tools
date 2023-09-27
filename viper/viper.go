package main

import (
	"fmt"
	"github.com/spf13/viper"
)

/*
go get github.com/spf13/viper
viper是一个配置管理的解决方案
URL Blog：https://piwriw.github.io/2023/03/02/web/go/%E6%95%B4%E5%90%88%E7%AE%A1%E7%90%86%E9%85%8D%E7%BD%AE%20-%20viper/?highlight=viper
*/
func main() {
	viper.SetConfigFile("./config.yaml")  // 指定配置文件路径
	viper.SetConfigName("config")         // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")           // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")              // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()           // 查找并读取配置文件
	if err != nil {                       // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 设置默认值
	viper.SetDefault("ContentDir", "content")
}
