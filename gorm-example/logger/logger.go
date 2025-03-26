package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	defaultMaxSQLLength  = 1024
	defaultSlowThreshold = 100
	defaultLogLevel      = 3
)

type GormLogger struct {
	LogConfig *LogConfig
	SQLConfig *SQLConfig
}

type LogConfig struct {
	LogLevel      int  `mapstructure:"LogLevel" yaml:"LogLevel"`           // 日志等级 1 silent 2 error 3 warn 4 info 默认是 3 warn
	SlowThreshold int  `mapstructure:"SlowThreshold" yaml:"SlowThreshold"` //  默认100ms 设置的是毫秒记录
	Colorful      bool `mapstructure:"Colorful" yaml:"Colorful"`
}

type SQLConfig struct {
	SQLFile      string `mapstructure:"SQLFile" yaml:"SQLFile"`
	MaxSQLLength int    `mapstructure:"MaxSQLLength" yaml:"MaxSQLLength"`
}

func SetGormLogger(conf GormLogger) *GormCustomLogger {
	// 自动设置默认值
	conf.setDefaults()
	// 配置 GORM 日志输出
	logConfig := logger.Config{
		SlowThreshold:             time.Duration(conf.LogConfig.SlowThreshold) * time.Millisecond,
		LogLevel:                  logger.LogLevel(conf.LogConfig.LogLevel),
		Colorful:                  conf.LogConfig.Colorful,
		IgnoreRecordNotFoundError: false,
	}
	newLog := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig)
	return NewGormCustomLogger(newLog, logConfig, *conf.SQLConfig)
}

// 为 GormLogger 设置默认值的方法
func (conf *GormLogger) setDefaults() {
	// LogConfig
	if conf.LogConfig == nil {
		conf.LogConfig = &LogConfig{}
	}
	if conf.LogConfig.LogLevel <= 0 {
		conf.LogConfig.LogLevel = defaultLogLevel
		slog.Warn("Invalid LogLevel, using default LogLevel=3 (warn)")
	}
	if conf.LogConfig.SlowThreshold <= 0 {
		conf.LogConfig.SlowThreshold = defaultSlowThreshold
		slog.Warn("Invalid SlowThreshold, using default SlowThreshold=100ms")
	}

	// SQLConfig
	if conf.SQLConfig == nil {
		conf.SQLConfig = &SQLConfig{}
	}
	if conf.SQLConfig.MaxSQLLength == 0 {
		conf.SQLConfig.MaxSQLLength = defaultMaxSQLLength
		slog.Warn("Invalid MaxSQLLength, using default MaxSQLLength=1024")
	} else if conf.SQLConfig.MaxSQLLength < 0 {
		slog.Warn("MaxSQLLength < 0, printing all SQL")
	}
}

type GormCustomLogger struct {
	Logger       logger.Interface // 嵌入现成的 GORM logger 接口
	Config       logger.Config
	MaxSQLLength int
	SQLFile      string
}

func NewGormCustomLogger(logger logger.Interface, config logger.Config, sqlConf SQLConfig) *GormCustomLogger {
	return &GormCustomLogger{
		Logger:       logger,
		Config:       config,
		MaxSQLLength: sqlConf.MaxSQLLength,
		SQLFile:      sqlConf.SQLFile,
	}
}

func (l *GormCustomLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Config.LogLevel <= logger.Silent {
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

	// Open SQL file if it's provided
	file, err := os.OpenFile(l.SQLFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logout.Error(fmt.Sprintf("Failed to open SQL log file %s: %v", l.SQLFile, err))
		return
	}
	defer file.Close()

	// Create a log writer
	sqlLogger := log.New(file, "", log.LstdFlags)
	// 打印SQL日志
	switch {
	case l.Config.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			logout.Error(fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Error(fmt.Sprintf("GORM EXECUTE ERROR  %s  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
		sqlLogger.Printf("GORM EXECUTE SQL  %s  %v  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Config.LogLevel >= logger.Warn:
		sql, rows := fc()
		if rows == -1 {
			logout.Warn(fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms]  %s", l.Config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Warn(fmt.Sprintf("GORM SLOW SQL >= %v  %s  [%fms] [rows:%d]  %s", l.Config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
		sqlLogger.Printf("GORM SLOW SQL >= %v  %s  [%fms] [rows:%d]  %s", l.Config.SlowThreshold, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql))
	case l.Config.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			logout.Info(fmt.Sprintf("GORM SQL  %s  [%fms]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, splitString(sql)))
		} else {
			logout.Info(fmt.Sprintf("GORM SQL  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, splitString(sql)))
		}
		sqlLogger.Printf("GORM SQL  %s  [%fms] [rows:%d]  %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
