package logger

import "fmt"

// Level 定义日志级别类型
type Level int

const (
	defaultLevel      = InfoLevel
	defaultJSONFormat = false
	defaultAddSource  = false
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
	Debug(args ...interface{})
	Debugf(format string, args ...any)

	Info(args ...interface{})
	Infof(format string, args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)

	Fatal(args ...interface{})
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
	Level Level
	// 以JSON格式输出日志
	JSONFormat bool
	FilePath   string
	// 是否记录源码位置(slog特有)
	AddSource bool
	// 格式化日志打印时间
	TimeFormat string
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

func WithTimeFormat(format string) Option {
	return func(o *Options) {
		o.TimeFormat = format
	}
}

func applyOptions(opts ...Option) Options {
	options := Options{
		Level:      defaultLevel,
		JSONFormat: defaultJSONFormat,
		AddSource:  defaultAddSource,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
