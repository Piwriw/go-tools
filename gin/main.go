package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		go func(ctx *gin.Context) {
			//defer func() {
			//	if err := recover(); err != nil {
			//		slog.Error("panic", slog.Any("err", err))
			//	}
			//}()
			arr := []string{"111", "222"}
			fmt.Println(arr[3])
		}(c)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
