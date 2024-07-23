package main

import (
	"log/slog"
	"os"
)

func main() {
	name := "sss"
	slog.Info("This is info log")
	slog.Warn("This is warning log")
	slog.Error("This is error log")

	// Group 日志
	g1 := slog.Group("g1", slog.String("module", "authentication"), slog.String("method", "login"))
	g2 := slog.Group("g2", slog.String("module", "authentication"), slog.String("method", "login"))

	slog.Info("User login attempt",
		slog.String("service", "login-service"),
		g1,
	)
	slog.Info("User login attempt",
		slog.String("service", "login-service"),
		g2,
	)

	slog.Info("msg", slog.String("name", name))
	slog.Error("ERROR: value is empty", slog.Any("name", name))
	// 设置日志格式为JSON
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})))
	slog.Error("ERROR: value is empty", slog.Any("name", name))
	// 设置日志为Text
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})))
	slog.Error("ERROR: value is empty", slog.Any("name", name))
	// 设置日志里面有代码行
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true})))
	slog.Error("ERROR: value is empty", slog.Any("name", name))

	// 设置日志写入本地文件
	f, err := os.Create("foo.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	slog.SetDefault(slog.New(slog.NewJSONHandler(f, nil)))
	slog.Error("ERROR: value is empty", slog.Any("name", name))
}
