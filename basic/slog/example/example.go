package main

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logger := InitLogger("app.log", "warn", 1, 5, 0, true)
	defer logger.Close()
	for i := 0; i < 100000; i++ {
		slog.Info("LOG int arr", slog.Any("intarr", "sss"))
		slog.Info("This is info log")
		slog.Warn("This is warning log")
		slog.Error("This is error log", "err", "x", "ss", "xx")
	}
}

func InitLogger(filePath string, logLevel string, maxSize int, maxBackups int, maxAge int, isCompress bool) *lumberjack.Logger {
	// 设置日志文件轮转
	logFile := &lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件的最大尺寸（MB）
		MaxBackups: maxBackups, // 保留的旧日志文件数
		MaxAge:     maxAge,     // 日志文件保存的天数
		Compress:   isCompress, // 是否压缩/归档旧日志文件
	}

	// 设置控制台和文件输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var level slog.Level
	switch logLevel {
	case slog.LevelInfo.String(), strings.ToLower(slog.LevelInfo.String()):
		level = slog.LevelInfo
	case slog.LevelDebug.String(), strings.ToLower(slog.LevelDebug.String()):
		level = slog.LevelDebug
	case slog.LevelWarn.String(), strings.ToLower(slog.LevelWarn.String()):
		level = slog.LevelWarn
	case slog.LevelError.String(), strings.ToLower(slog.LevelError.String()):
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 创建自定义的 `slog.Handler`，将日志输出到 multiWriter
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{AddSource: true, Level: level})

	// 创建 logger 实例
	logger := slog.New(handler)
	slog.SetDefault(logger)
	// 在程序结束时关闭日志文件
	return logFile
}
