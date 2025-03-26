package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type slogLogger struct {
	logger          *slog.Logger
	filePath        string
	addSource       bool
	level           Level
	replaceAttrFunc func(groups []string, a slog.Attr) slog.Attr
}

var _ Logger = (*slogLogger)(nil)

var defaultReplaceAttrFunc = func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Time()
		a.Value = slog.StringValue(t.Format(time.DateTime))
	}
	return a
}

func newSlogLogger(opts Options) (Logger, error) {
	if opts.TimeFormat != "" {
		defaultReplaceAttrFunc = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(opts.TimeFormat))
			}
			return a
		}
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource:   opts.AddSource,
		Level:       getSlogLoggerLevel(opts.Level),
		ReplaceAttr: defaultReplaceAttrFunc,
	}
	// 设置控制台和文件输出
	multiWriter := io.MultiWriter(os.Stdout, getOutput(opts.FilePath))
	var handler slog.Handler
	if opts.JSONFormat {
		handler = slog.NewJSONHandler(multiWriter, handlerOpts)
	} else {
		handler = slog.NewTextHandler(multiWriter, handlerOpts)
	}

	return &slogLogger{
		filePath:        opts.FilePath,
		addSource:       opts.AddSource,
		logger:          slog.New(handler),
		level:           opts.Level,
		replaceAttrFunc: defaultReplaceAttrFunc,
	}, nil
}

// getSlogLoggerLevel 将自定义的 Level 转换为 slog.Level
func getSlogLoggerLevel(level Level) slog.Level {
	switch level {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case FatalLevel:
		return slog.LevelError // slog 没有 Fatal 级别
	default:
		return slog.LevelInfo
	}
}

func (l *slogLogger) log(level slog.Level, msg string, args ...any) {
	if !l.logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // 跳过 3 层调用栈
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])

	if len(args) > 0 {
		if len(args)%2 != 0 {
			args = append(args, "!MISSING!")
		}
		for i := 0; i < len(args); i += 2 {
			if key, ok := args[i].(string); ok {
				r.AddAttrs(slog.Any(key, args[i+1]))
			}
		}
	}

	_ = l.logger.Handler().Handle(context.Background(), r)
}

func (l *slogLogger) Debug(args ...any) {
	l.log(slog.LevelDebug, fmt.Sprint(args...))
}

func (l *slogLogger) Debugf(format string, args ...any) {
	l.log(slog.LevelDebug, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Info(args ...any) {
	l.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (l *slogLogger) Infof(format string, args ...any) {
	l.log(slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Warn(args ...any) {
	l.log(slog.LevelWarn, fmt.Sprint(args...))
}

func (l *slogLogger) Warnf(format string, args ...any) {
	l.log(slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Error(args ...any) {
	l.log(slog.LevelError, fmt.Sprint(args...))
}

func (l *slogLogger) Errorf(format string, args ...any) {
	l.log(slog.LevelError, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Fatal(args ...any) {
	l.log(slog.LevelError, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *slogLogger) Fatalf(format string, args ...any) {
	l.log(slog.LevelError, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *slogLogger) SetLevel(level Level) {
	// 重新创建 Logger 实例以应用新的日志级别
	handlerOpts := &slog.HandlerOptions{
		AddSource:   l.addSource,
		Level:       getSlogLoggerLevel(level),
		ReplaceAttr: l.replaceAttrFunc,
	}
	newHandler := slog.NewTextHandler(getOutput(l.filePath), handlerOpts)
	l.logger = slog.New(newHandler)
	l.level = level
}

func (l *slogLogger) WithFields(fields map[string]any) Logger {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Attr{Key: k, Value: slog.AnyValue(v)})
	}
	//  解决 []slog.Attr 转换为 []any 的问题
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return &slogLogger{
		logger:          l.logger.With(args...),
		level:           l.level,
		addSource:       l.addSource,
		replaceAttrFunc: l.replaceAttrFunc,
	}
}
