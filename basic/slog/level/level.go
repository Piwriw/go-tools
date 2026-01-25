package main

import (
	"fmt"
	"log/slog"
	"os"
)

// LevelHandler 是一个自定义的处理器，用于根据日志级别将日志写入不同的文件
type LevelHandler struct {
	level   slog.Level
	handler slog.Handler
}

// NewLevelHandler 创建一个指定级别的日志处理器
func NewLevelHandler(level slog.Level, file *os.File) *LevelHandler {
	return &LevelHandler{
		level:   level,
		handler: slog.NewJSONHandler(file), // 使用 JSON 格式的日志
	}
}

// Handle 处理日志记录，仅当日志级别匹配时才记录日志
func (h *LevelHandler) Handle(r slog.Record) error {
	if r.Level == h.level {
		return h.handler.Handle(r)
	}
	return nil
}

// WithAttrs 增加日志属性
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.handler = h.handler.WithAttrs(attrs)
	return h
}

// WithGroup 增加日志分组
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	h.handler = h.handler.WithGroup(name)
	return h
}

func main() {
	// 创建两个不同级别的日志文件
	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("failed to open info.log: %v\n", err)
		return
	}
	defer infoFile.Close()

	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("failed to open error.log: %v\n", err)
		return
	}
	defer errorFile.Close()

	// 创建针对不同日志级别的处理器
	infoHandler := NewLevelHandler(slog.LevelInfo, infoFile)
	errorHandler := NewLevelHandler(slog.LevelError, errorFile)

	// 创建日志记录器，将多个处理器组合
	logger := slog.New(slog.NewMultiHandler(infoHandler, errorHandler))

	// 记录不同级别的日志
	logger.Info("This is an info message", "module", "main")
	logger.Error("This is an error message", "module", "main")
}
