package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger 实现 Logger 接口

type zapLogger struct {
	logger      *zap.SugaredLogger
	errorLogger *zap.SugaredLogger
	level       Level
	colorScheme *ColorScheme
}

var _ Logger = (*zapLogger)(nil)

func newZapLogger(opts Options) (Logger, error) {
	// 构建通用的 zap.Config
	buildConfig := func(level zap.AtomicLevel, disableCaller bool) zap.Config {
		cfg := zap.NewProductionConfig()
		cfg.Level = level
		cfg.DisableCaller = disableCaller
		cfg.DisableStacktrace = true // 关闭 error 级别的堆栈打印（zap 默认会打印）
		if !opts.JSONFormat {
			cfg.Encoding = "console"
		}
		return cfg
	}

	// 创建主日志配置
	mainCfg := buildConfig(ToZapLevel(opts.Level), !opts.AddSource)
	if opts.FilePath != "" {
		mainCfg.OutputPaths = []string{"stderr", opts.FilePath}
	}
	// 创建 error 日志配置
	errorCfg := buildConfig(ToZapLevel(ErrorLevel), !opts.AddSource)
	if opts.FilePath != "" {
		mainCfg.OutputPaths = []string{opts.FilePath}
	}
	if opts.ErrorOutput != "" {
		errorCfg.OutputPaths = []string{opts.ErrorOutput}
	}
	errorCfg.DisableStacktrace = false
	logger, err := mainCfg.Build()
	if err != nil {
		return nil, err
	}
	errorLogger, err := errorCfg.Build()
	if err != nil {
		return nil, err
	}
	zapLogger := &zapLogger{
		logger:      logger.Sugar(),
		errorLogger: errorLogger.Sugar(),
		level:       opts.Level,
	}
	// 设置颜色输出
	if opts.ColorEnabled {
		zapLogger.colorScheme = opts.ColorScheme
	}
	return zapLogger, nil
}

func ToZapLevel(level Level) zap.AtomicLevel {
	switch level {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevel, FatalLevel:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func FromZapLevel(zapLevel zapcore.Level) Level {
	switch zapLevel {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.FatalLevel:
		return FatalLevel
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return ErrorLevel
	default:
		return defaultLevel
	}
}

func (l *zapLogger) log(level zapcore.Level, msg string, args ...any) {
	if l.colorScheme != nil {
		msg = l.colorScheme.Colorize(FromZapLevel(level), msg)
	}
	switch level {
	case zap.DebugLevel:
		l.logger.Debugw(msg, args...)
	case zap.InfoLevel:
		l.logger.Infow(msg, args...)
	case zap.WarnLevel:
		l.logger.Warnw(msg, args...)
	case zap.ErrorLevel:
		l.logger.Errorw(msg, args...)
		if l.errorLogger != nil {
			l.errorLogger.Errorw(msg, args...)
		}
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
