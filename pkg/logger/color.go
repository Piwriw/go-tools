package logger

import (
	"github.com/fatih/color"
)

// Color 类型定义（兼容原有代码）
type Color string

// 保留原有颜色常量定义（可选）
const (
	ColorReset  Color = "\033[0m"
	ColorRed    Color = "\033[31m"
	ColorGreen  Color = "\033[32m"
	ColorYellow Color = "\033[33m"
	ColorBlue   Color = "\033[34m"
	ColorPurple Color = "\033[35m"
	ColorCyan   Color = "\033[36m"
	ColorWhite  Color = "\033[37m"
	ColorGray   Color = "\033[90m"
)

// ColorScheme 使用 fatih/color 的颜色属性
type ColorScheme struct {
	Debug *color.Color
	Info  *color.Color
	Warn  *color.Color
	Error *color.Color
	Fatal *color.Color
}

// 默认颜色方案
var defaultColorScheme = ColorScheme{
	Debug: color.New(color.FgCyan),
	Info:  color.New(color.FgGreen),
	Warn:  color.New(color.FgYellow),
	Error: color.New(color.FgRed),
	Fatal: color.New(color.FgMagenta),
}

// 无颜色方案
var noColorScheme = ColorScheme{
	Debug: color.New(),
	Info:  color.New(),
	Warn:  color.New(),
	Error: color.New(),
	Fatal: color.New(),
}

// Colorize 使用 fatih/color 的实现
func Colorize(level Level, msg string, scheme ColorScheme) string {
	var c *color.Color
	color.NoColor = false

	switch level {
	case DebugLevel:
		c = scheme.Debug
	case InfoLevel:
		c = scheme.Info
	case WarnLevel:
		c = scheme.Warn
	case ErrorLevel:
		c = scheme.Error
	case FatalLevel:
		c = scheme.Fatal
	default:
		return msg
	}

	return c.Sprint(msg)
}
