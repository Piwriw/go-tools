package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	defaultMaxSQLLength  = 1024
	defaultSlowThreshold = 100
	defaultLogLevel      = 3
)

func SetGormLogger() *GormCustomLogger {
	logLevel := setting.CFG.GetLogLevel()
	if logLevel <= 0 {
		slog.Warn("set Gorm LogLevel 0,will use default LogLevel=3(warn)")
		logLevel = defaultLogLevel
	}
	slowThreshold := setting.CFG.GetSlowThreshold()
	if slowThreshold <= 0 {
		slog.Warn("set Gorm SlowThreshold 0,will use default SlowThreshold=100ms")
		slowThreshold = defaultSlowThreshold
	}
	maxSQLLength := setting.CFG.GetMaxSQLLength()
	if maxSQLLength == 0 {
		slog.Warn("set Gorm MaxSQLLength 0,will use default MaxSQLLength=1024")
		maxSQLLength = defaultMaxSQLLength
	} else if maxSQLLength < 0 {
		slog.Warn("set Gorm MaxSQLLength < 0,will use print all SQL")
	}
	// 配置 GORM 日志输出
	logConfig := logger.Config{
		SlowThreshold:             time.Duration(slowThreshold) * time.Millisecond,
		LogLevel:                  logger.LogLevel(logLevel),
		Colorful:                  !color.NoColor,
		IgnoreRecordNotFoundError: false}
	newLog := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig)
	return NewGormCustomLogger(newLog, logConfig, maxSQLLength)
}

type GormCustomLogger struct {
	logger.Interface // 嵌入现成的 GORM logger 接口
	config           logger.Config
	MaxSQLLength     int
}

func NewGormCustomLogger(logger logger.Interface, config logger.Config, maxSQLLength int) *GormCustomLogger {
	if maxSQLLength == 0 {
		maxSQLLength = defaultMaxSQLLength
	}
	return &GormCustomLogger{
		Interface:    logger,
		config:       config,
		MaxSQLLength: maxSQLLength,
	}
}

func (l GormCustomLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= logger.Silent {
		return
	}

	var (
		elapsed     = time.Since(begin)
		splitString = func(sql string) string {
			sql = strings.ReplaceAll(sql, "\n", " ")
			sql = strings.ReplaceAll(sql, "\t", " ")
			if l.MaxSQLLength < 0 {
				return sql
			}
			if len(sql) > l.MaxSQLLength {
				return sql[:l.MaxSQLLength] + " ... (SQL text truncated)"
			}
			return sql
		}
		logout = slog.Default()
	)

	switch {
	case err != nil && l.config.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			logout.Error(fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Error(fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= logger.Warn:
		sql, rows := fc()
		if rows == -1 {
			logout.Warn(fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms]  %s", l.config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Warn(fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms] [rows:%d]  %s", l.config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	case l.config.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			logout.Info(fmt.Sprintf("GORM SQL  %s  [%fms]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Info(fmt.Sprintf("GORM SQL  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	}
}
