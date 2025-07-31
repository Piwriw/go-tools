package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	defaultMaxSQLLength  = 1024
	defaultSlowThreshold = 100
	defaultLogLevel      = 3
)

// GormLoggerOption 定义 GORM Logger 配置选项类型
type GormLoggerOption func(*gormLoggerOptions)

// gormLoggerOptions 存储所有可配置选项
type gormLoggerOptions struct {
	logLevel                  int
	slowThreshold             int
	maxSQLLength              int
	colorful                  bool
	ignoreRecordNotFoundError bool
	logFile                   string
}

// WithGormLogLevel 设置日志级别,Gorm日志级别
//
//	1 Silent silent log level
//	2 Error error log level
//	3 Warn warn log level
//	4 Info info log level
func WithGormLogLevel(level int) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.logLevel = level
	}
}

// WithGormSlowThreshold 设置慢查询阈值(毫秒)
func WithGormSlowThreshold(threshold int) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.slowThreshold = threshold
	}
}

// WithGormMaxSQLLength 设置SQL最大长度
// 默认为1024
// <0打印全部SQL
func WithGormMaxSQLLength(length int) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.maxSQLLength = length
	}
}

// WithDisGormMaxSQLLength 不设置SQL最大长度限制
// 默认为1024
// <0打印全部SQL
func WithDisGormMaxSQLLength() GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.maxSQLLength = -1
	}
}

// WithGormColorful 设置是否启用颜色输出
func WithGormColorful(colorful bool) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.colorful = colorful
	}
}

// WithIgnoreRecordNotFoundError 设置是否忽略记录未找到错误
func WithIgnoreRecordNotFoundError(ignore bool) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.ignoreRecordNotFoundError = ignore
	}
}

// WithGormLogFile 设置日志文件路径
func WithGormLogFile(logFile string) GormLoggerOption {
	return func(o *gormLoggerOptions) {
		o.logFile = logFile
	}
}

// NewGormLogger 使用 Options 模式创建 GORM 自定义日志记录器
func NewGormLogger(opts ...GormLoggerOption) *GormCustomLogger {
	// 默认配置
	options := &gormLoggerOptions{
		logLevel:                  defaultLogLevel,
		slowThreshold:             defaultSlowThreshold,
		maxSQLLength:              defaultMaxSQLLength,
		colorful:                  !color.NoColor,
		ignoreRecordNotFoundError: false,
	}

	// 应用传入的选项
	for _, opt := range opts {
		opt(options)
	}

	// 验证和处理配置
	if options.logLevel <= 0 {
		slog.Warn("set Gorm LogLevel 0,will use default LogLevel=3(warn)")
		options.logLevel = defaultLogLevel
	}

	if options.slowThreshold <= 0 {
		slog.Warn("set Gorm SlowThreshold 0,will use default SlowThreshold=100ms")
		options.slowThreshold = defaultSlowThreshold
	}

	if options.maxSQLLength == 0 {
		slog.Warn("set Gorm maxSQLLength 0,will use default maxSQLLength=1024")
		options.maxSQLLength = defaultMaxSQLLength
	} else if options.maxSQLLength < 0 {
		slog.Warn("set Gorm maxSQLLength < 0,will use print all SQL")
		options.maxSQLLength = -1
	}

	// 配置 GORM 日志输出
	logConfig := gormlogger.Config{
		SlowThreshold:             time.Duration(options.slowThreshold) * time.Millisecond,
		LogLevel:                  gormlogger.LogLevel(options.logLevel),
		Colorful:                  options.colorful,
		IgnoreRecordNotFoundError: options.ignoreRecordNotFoundError,
	}

	newLog := gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig)
	return newGormCustomLogger(newLog, logConfig, options.maxSQLLength, options.logFile)
}

type GormCustomLogger struct {
	gormlogger.Interface // 嵌入现成的 GORM logger 接口
	config               gormlogger.Config
	maxSQLLength         int
	logFile              string
	logout               *slog.Logger
}

func newGormCustomLogger(loggerIn gormlogger.Interface, config gormlogger.Config, maxSQLLength int, logFile string) *GormCustomLogger {
	if maxSQLLength == 0 {
		maxSQLLength = defaultMaxSQLLength
	}
	/*
	 SQL不应该被因为slog而被截断
	*/
	handler := NewLogger(WithLogPath(logFile), WithCloseSource(), WithDisMaxMessageLength())
	logout := slog.New(handler)
	return &GormCustomLogger{
		Interface:    loggerIn,
		config:       config,
		maxSQLLength: maxSQLLength,
		logFile:      logFile,
		logout:       logout,
	}
}

func (l *GormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= gormlogger.Silent {
		return
	}

	var (
		elapsed     = time.Since(begin)
		splitString = func(sql string) string {
			sql = strings.ReplaceAll(sql, "\n", " ")
			sql = strings.ReplaceAll(sql, "\t", " ")
			if l.maxSQLLength <= 0 {
				return sql
			}
			if len(sql) > l.maxSQLLength {
				return sql[:l.maxSQLLength] + " ... (SQL text truncated)" + "maxSQLLength:" + strconv.Itoa(l.maxSQLLength)
			}
			return sql
		}
	)
	switch {
	case err != nil && l.config.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.logout.ErrorContext(ctx, fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			l.logout.ErrorContext(ctx, fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		if rows == -1 {
			l.logout.WarnContext(ctx, fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms]  %s", l.config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			l.logout.WarnContext(ctx, fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms] [rows:%d]  %s", l.config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	case l.config.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.logout.InfoContext(ctx, fmt.Sprintf("GORM SQL  %s  [%fms]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			l.logout.InfoContext(ctx, fmt.Sprintf("GORM SQL  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
	}
}

// ParamsFilter filter params
func (l *GormCustomLogger) ParamsFilter(_ context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
