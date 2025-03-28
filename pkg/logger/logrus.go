package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var defaultCallerPrettyfierFunc = func(f *runtime.Frame) (string, string) {
	// 自定义文件名 + 行号格式
	filename := path.Base(f.File)
	return "", fmt.Sprintf("%s:%d", filename, f.Line)
}

// logrusLogger 实现

type logrusLogger struct {
	logger      *logrus.Logger
	errorLogger *logrus.Logger
	level       Level
	fields      logrus.Fields
	colorScheme *ColorScheme
}

var _ Logger = (*logrusLogger)(nil)

func newLogrusLogger(opts Options) (Logger, error) {
	logger := logrus.New()
	errorLogger := logrus.New()
	if opts.AddSource {
		logger.SetReportCaller(true)
		errorLogger.SetReportCaller(true)
	}

	if opts.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: defaultCallerPrettyfierFunc,
			TimestampFormat:  opts.TimeFormat,
		})
		errorLogger.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: defaultCallerPrettyfierFunc,
			TimestampFormat:  opts.TimeFormat,
		}
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: defaultCallerPrettyfierFunc,
			TimestampFormat:  opts.TimeFormat,
		})
		errorLogger.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: defaultCallerPrettyfierFunc,
			TimestampFormat:  opts.TimeFormat,
		}
	}
	logger.SetLevel(ToLogrusLoggerLevel(opts.Level))
	errorLogger.SetLevel(ToLogrusLoggerLevel(ErrorLevel))
	// 设置控制台和文件输出
	multiWriter := io.MultiWriter(os.Stdout, getOutput(opts.FilePath))
	logger.SetOutput(multiWriter)
	errorLogger.SetOutput(getOutput(opts.ErrorOutput))
	logrusLogger := &logrusLogger{
		logger:      logger,
		errorLogger: errorLogger,
		level:       opts.Level,
	}
	// 设置颜色输出
	if opts.ColorEnabled {
		logrusLogger.colorScheme = opts.ColorScheme
	}
	return logrusLogger, nil
}

func ToLogrusLoggerLevel(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func FromLogrusLoggerLevel(level logrus.Level) Level {
	switch level {
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.FatalLevel:
		return FatalLevel
	default:
		return defaultLevel
	}
}

func (l *logrusLogger) log(level logrus.Level, args ...any) {
	msg := fmt.Sprint(args...)
	if l.colorScheme != nil {
		msg = l.colorScheme.Colorize(FromLogrusLoggerLevel(level), msg)
	}
	if l.fields != nil {
		l.logger.WithFields(l.fields).Print(msg)
		return
	}
	l.logger.Log(level, msg)
	if l.errorLogger != nil && level >= logrus.ErrorLevel {
		l.errorLogger.Log(level, fmt.Sprint(msg))
	}
}

func (l *logrusLogger) Debug(args ...any) { l.log(logrus.DebugLevel, args...) }
func (l *logrusLogger) Debugf(format string, args ...any) {
	l.log(logrus.DebugLevel, fmt.Sprintf(format, args...))
}
func (l *logrusLogger) Info(args ...any) { l.log(logrus.InfoLevel, args...) }
func (l *logrusLogger) Infof(format string, args ...any) {
	l.log(logrus.InfoLevel, fmt.Sprintf(format, args...))
}
func (l *logrusLogger) Warn(args ...any) {
	l.log(logrus.WarnLevel, args...)
}
func (l *logrusLogger) Warnf(format string, args ...any) {
	l.log(logrus.WarnLevel, fmt.Sprintf(format, args...))
}
func (l *logrusLogger) Error(args ...any) { l.log(logrus.ErrorLevel, args...) }
func (l *logrusLogger) Errorf(format string, args ...any) {
	l.log(logrus.ErrorLevel, fmt.Sprintf(format, args...))
}
func (l *logrusLogger) Fatal(args ...any) { l.log(logrus.FatalLevel, args...) }
func (l *logrusLogger) Fatalf(format string, args ...any) {
	l.log(logrus.FatalLevel, fmt.Sprintf(format, args...))
}
func (l *logrusLogger) SetLevel(level Level) { l.logger.SetLevel(logrus.Level(level)) }

func (l *logrusLogger) WithFields(fields map[string]any) Logger {
	return &logrusLogger{
		logger: l.logger,
		level:  l.level,
		fields: logrus.Fields(fields),
	}
}
