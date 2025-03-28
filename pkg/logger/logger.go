package logger

import (
	"fmt"
	"time"
)

// Level 定义日志级别类型
type Level int

const (
	defaultLevel       = InfoLevel
	defaultJSONFormat  = false
	defaultAddSource   = false
	defaultLogFile     = "./app.log"
	defaultTimeFormat  = time.DateTime
	defaultErrorOutput = "./app_error.log"
)

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger 接口定义
type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)

	Info(args ...any)
	Infof(format string, args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)

	Fatal(args ...any)
	Fatalf(format string, args ...any)

	SetLevel(level Level)
	WithFields(fields map[string]any) Logger
}

// LoggerType 定义支持的Logger类型
type LoggerType string

const (
	SlogLogger   LoggerType = "slog"   // 默认使用slog
	ZapLogger    LoggerType = "zap"    // 可选
	LogrusLogger LoggerType = "logrus" // 可选
)

// NewLogger 创建一个新的Logger实例，默认使用slog
func NewLogger(options ...Option) (Logger, error) {
	return NewLoggerWithType(SlogLogger, options...)
}

// DefaultLogger 创建一个默认的Logger实例，默认使用slog
func DefaultLogger() (Logger, error) {
	opts := applyOptions(
		WithAddSource(),
		WithLevel(defaultLevel),
		WithTimeFormat(defaultTimeFormat),
		WithFileOutput(defaultLogFile),
		WithErrorOutPut(defaultErrorOutput),
	)
	return newSlogLogger(opts)
}

// NewLoggerWithType 创建指定类型的Logger实例
func NewLoggerWithType(loggerType LoggerType, options ...Option) (Logger, error) {
	opts := applyOptions(options...)

	switch loggerType {
	case SlogLogger:
		return newSlogLogger(opts)
	case ZapLogger:
		return newZapLogger(opts)
	case LogrusLogger:
		return newLogrusLogger(opts)
	default:
		return nil, fmt.Errorf("unknown logger type: %s", loggerType)
	}
}

// Option 配置选项
type Option func(*Options)

// Options 包含Logger的配置
type Options struct {
	// 设置日志级别
	Level Level
	// 以JSON格式输出日志
	JSONFormat bool
	// 设置日志文件路径
	FilePath string
	// 是否打印日志函数调用信息，默认不打印
	AddSource bool
	// 格式化日志打印时间
	TimeFormat string
	// ErrorOutput 错误日志输出
	ErrorOutput string
	// 日志轮转配置
	LogRotation *LogRotation
	// 设置颜色输出
	// 开启默认有一个默认的颜色方案，可通过 WithColorScheme 自定义颜色方案
	ColorEnabled bool
	// 主题颜色方案
	ColorScheme *ColorScheme
	// 其他配置项...
}

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithJSONFormat() Option {
	return func(o *Options) {
		o.JSONFormat = true
	}
}

func WithFileOutput(path string) Option {
	return func(o *Options) {
		o.FilePath = path
	}
}

func WithAddSource() Option {
	return func(o *Options) {
		o.AddSource = true
	}
}

// WithTimeFormat 格式化日志打印时间
// 默认为 time.DateTime
func WithTimeFormat(format string) Option {
	return func(o *Options) {
		o.TimeFormat = format
	}
}

// WithErrorOutPut 错误日志输出
// Zap 默认开启Error日志的堆栈打印，Logrus 默认关闭
func WithErrorOutPut(path string) Option {
	return func(o *Options) {
		o.ErrorOutput = path
	}
}

// WithLogRotation 日志轮转配置
// 日志轮转配置，默认不开启
func WithLogRotation(filePath string, maxSize int, maxBackups int, maxAge int, isCompress bool) Option {
	return func(o *Options) {
		o.LogRotation = &LogRotation{
			FilePath:   filePath,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   isCompress,
		}
	}
}

// WithColor 启用颜色输出
// 启用颜色输出，默认不开启
// 注意：颜色输出会影响性能，建议在开发环境中使用
func WithColor() Option {
	return func(o *Options) {
		o.ColorEnabled = true
	}
}

// WithColorScheme 添加主题颜色方案
// 自定义颜色方案，默认使用默认颜色方案
// 注意：颜色输出会影响性能，建议在开发环境中使用
func WithColorScheme(scheme ColorScheme) Option {
	return func(o *Options) {
		o.ColorScheme = &scheme
	}
}

func applyOptions(opts ...Option) Options {
	options := Options{
		Level:      defaultLevel,
		JSONFormat: defaultJSONFormat,
		AddSource:  defaultAddSource,
		TimeFormat: time.DateTime,
	}
	for _, opt := range opts {
		opt(&options)
	}
	if options.ColorEnabled && options.ColorScheme == nil {
		options.ColorScheme = DefaultFatihColorScheme
	}
	return options
}
