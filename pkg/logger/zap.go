package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger 实现 Logger 接口

type zapLogger struct {
	logger *zap.SugaredLogger
	level  Level
}

var _ Logger = (*zapLogger)(nil)

func newZapLogger(opts Options) (Logger, error) {
	var zapLevel zap.AtomicLevel
	switch opts.Level {
	case DebugLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevel, FatalLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zapLevel
	cfg.DisableCaller = !opts.AddSource
	// 关闭打印堆栈（error级别，zap默认打印）
	cfg.DisableStacktrace = true
	if !opts.JSONFormat {
		cfg.Encoding = "console"
	}

	if opts.FilePath != "" {
		cfg.OutputPaths = []string{opts.FilePath, "stderr"}
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: logger.Sugar(),
		level:  opts.Level,
	}, nil
}

func (l *zapLogger) log(level zapcore.Level, msg string, args ...any) {
	switch level {
	case zap.DebugLevel:
		l.logger.Debugw(msg, args...)
	case zap.InfoLevel:
		l.logger.Infow(msg, args...)
	case zap.WarnLevel:
		l.logger.Warnw(msg, args...)
	case zap.ErrorLevel:
		l.logger.Errorw(msg, args...)
	}
}

func (l *zapLogger) Debug(args ...any) { l.log(zap.DebugLevel, fmt.Sprint(args...)) }
func (l *zapLogger) Debugf(format string, args ...any) {
	l.log(zap.DebugLevel, fmt.Sprintf(format, args...))
}
func (l *zapLogger) Info(args ...any) { l.log(zap.InfoLevel, fmt.Sprint(args...)) }
func (l *zapLogger) Infof(format string, args ...any) {
	l.log(zap.InfoLevel, fmt.Sprintf(format, args...))
}
func (l *zapLogger) Warn(args ...any) { l.log(zap.WarnLevel, fmt.Sprint(args...)) }
func (l *zapLogger) Warnf(format string, args ...any) {
	l.log(zap.WarnLevel, fmt.Sprintf(format, args...))
}
func (l *zapLogger) Error(args ...any) { l.log(zap.ErrorLevel, fmt.Sprint(args...)) }
func (l *zapLogger) Errorf(format string, args ...any) {
	l.log(zap.ErrorLevel, fmt.Sprintf(format, args...))
}
func (l *zapLogger) Fatal(args ...any) {
	l.log(zap.ErrorLevel, fmt.Sprint(args...))
	os.Exit(1)
}
func (l *zapLogger) Fatalf(format string, args ...any) {
	l.log(zap.ErrorLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}
func (l *zapLogger) SetLevel(level Level) { /* Zap Logger 的 Level 不能动态修改 */ }
func (l *zapLogger) WithFields(fields map[string]any) Logger {
	return &zapLogger{
		logger: l.logger.With(fields),
		level:  l.level,
	}
}
