package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	ShanghaiLoc *time.Location
)

// 限制消息的最大长度
const defaultMaxMessageLength = 1024

func init() {
	var err error
	// 加载上海时区（如果 location 为 nil，默认用上海时区）
	ShanghaiLoc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		slog.Error("ReloadConfig(), fail to load location Asia/Shanghai. ", slog.Any("err", err))
		ShanghaiLoc = time.FixedZone("CST", 8*60*60) // 手动设置 UTC+8
	}
}

// LevelSplitHandler 是自定义的 Handler，按级别分发日志
type LevelSplitHandler struct {
	defaultHandler slog.Handler
	errorHandler   slog.Handler
}

func (h *LevelSplitHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.defaultHandler.Enabled(ctx, level) || h.errorHandler.Enabled(ctx, level)
}

func (h *LevelSplitHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelError {
		if h.errorHandler != nil {
			if err := h.errorHandler.Handle(ctx, record); err != nil {
				return err
			}
		}
		return h.defaultHandler.Handle(ctx, record)
	}
	return h.defaultHandler.Handle(ctx, record)
}

func (h *LevelSplitHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.errorHandler != nil {
		h.errorHandler = h.errorHandler.WithAttrs(attrs)
	}
	if h.defaultHandler != nil {
		h.defaultHandler = h.defaultHandler.WithAttrs(attrs)
	}
	return h
}

func (h *LevelSplitHandler) WithGroup(name string) slog.Handler {
	if h.errorHandler != nil {
		h.errorHandler = h.errorHandler.WithGroup(name)
	}
	if h.defaultHandler != nil {
		h.defaultHandler = h.defaultHandler.WithGroup(name)
	}
	return h
}

// Option 定义 Logger 配置选项类型
type Option func(*loggerOptions)

// loggerOptions 存储所有可配置选项
type loggerOptions struct {
	logLevel         string
	isAddSource      bool
	errorLogPath     string
	logPath          string
	replaceAttrFunc  func(groups []string, a slog.Attr) slog.Attr
	maxMessageLength int
}

// WithLogLevel 设置日志级别
func WithLogLevel(level string) Option {
	return func(o *loggerOptions) {
		o.logLevel = level
	}
}

// WithCloseSource 关闭添加源码信息
func WithCloseSource() Option {
	return func(o *loggerOptions) {
		o.isAddSource = false
	}
}

// WithErrorLogPath 设置错误日志文件路径
func WithErrorLogPath(path string) Option {
	return func(o *loggerOptions) {
		if path == "" {
			return
		}
		o.errorLogPath = path
	}
}

// WithLogPath 添加日志文件路径
func WithLogPath(path string) Option {
	return func(o *loggerOptions) {
		if path == "" {
			return
		}
		o.logPath = path
	}
}

// WithMaxMessageLength 设置最大消息长度
// 设置<=0 不截断
func WithMaxMessageLength(length int) Option {
	return func(o *loggerOptions) {
		o.maxMessageLength = length
	}
}

// WithDisMaxMessageLength 关闭最大消息长度，默认1024
// 设置<=0 不截断
func WithDisMaxMessageLength() Option {
	return func(o *loggerOptions) {
		o.maxMessageLength = -1
	}
}

// InitLogger 使用 Options 模式初始化日志记录器
func InitLogger(opts ...Option) {
	handler := NewLogger(opts...)
	slog.SetDefault(slog.New(handler))
}

// NewLogger 使用 Options 模式初始化日志记录器
func NewLogger(opts ...Option) *LevelSplitHandler {
	// 默认配置
	options := &loggerOptions{
		logLevel:         slog.LevelInfo.String(),
		isAddSource:      true,
		maxMessageLength: defaultMaxMessageLength,
	}
	// 应用传入的选项
	for _, opt := range opts {
		opt(options)
	}

	options.replaceAttrFunc = func(_ []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			// 强制转换为上海时间
			t := a.Value.Time().In(ShanghaiLoc)
			a.Value = slog.StringValue(t.Format(time.DateTime))
		case slog.MessageKey:
			message := a.Value.String()
			if options.maxMessageLength <= 0 {
				return a
			}
			if len(message) > options.maxMessageLength {
				// 截断消息并在末尾添加省略号
				truncated := message[:options.maxMessageLength] + "......"
				a.Value = slog.StringValue(truncated)
			}
		default:
			// 处理日志参数中的长字符串
			if options.maxMessageLength <= 0 {
				return a
			}
			if a.Value.Kind() == slog.KindString && options.maxMessageLength > 0 {
				strValue := a.Value.String()
				if len(strValue) > options.maxMessageLength {
					// 截断并在末尾添加省略号
					truncated := strValue[:options.maxMessageLength] + "......"
					a.Value = slog.StringValue(truncated)
				}
			}
		}
		return a
	}

	// 使用标准输出作为日志目标
	var multiWriter []io.Writer
	multiWriter = append(multiWriter, os.Stdout)

	var level slog.Level
	switch options.logLevel {
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
	var errorHandler slog.Handler
	if options.errorLogPath != "" {
		// 创建 error 日志文件（你可以换成日志轮转方式）
		errorFile, err := os.OpenFile(options.errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("open error log file failed: " + err.Error())
		}
		errorHandler = slog.NewTextHandler(io.MultiWriter(errorFile), &slog.HandlerOptions{
			AddSource: options.isAddSource,
			Level:     slog.LevelError,
		})
	}
	if options.logPath != "" {
		logFile, err := os.OpenFile(options.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("open log  file failed: " + err.Error())
		}
		multiWriter = append(multiWriter, logFile)
	}
	// 创建两个 handler
	defaultHandler := slog.NewTextHandler(io.MultiWriter(multiWriter...), &slog.HandlerOptions{
		AddSource:   options.isAddSource,
		Level:       level,
		ReplaceAttr: options.replaceAttrFunc,
	})

	// 创建分发 handler
	handler := &LevelSplitHandler{
		defaultHandler: defaultHandler,
		errorHandler:   errorHandler,
	}
	return handler
}
