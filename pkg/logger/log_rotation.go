package logger

import "gopkg.in/natefinch/lumberjack.v2"

type LogRotation struct {
	logger     *lumberjack.Logger
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件的最大大小（单位：MB）
	MaxBackups int    // 保留的旧日志文件的最大数量
	MaxAge     int    // 保留的旧日志文件的最大天数
	Compress   bool   // 是否压缩/归档旧日志文件

}

func initLogRotation(filePath string, maxSize int, maxBackups int, maxAge int, isCompress bool) *LogRotation {
	// 设置日志文件轮转
	logFile := &lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件的最大尺寸（MB）
		MaxBackups: maxBackups, // 保留的旧日志文件数
		MaxAge:     maxAge,     // 日志文件保存的天数
		Compress:   isCompress, // 是否压缩/归档旧日志文件
	}
	return &LogRotation{
		logger:     logFile,
		FilePath:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   isCompress,
	}
}
