package main

import (
	"log/slog"
	"os"
)

func main() {
	name := "sss"
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
