package logger

import (
	"github.com/fatih/color"
)

// 预定义颜色（支持 ANSI 和 fatih/color）
var (
	ColorReset  = Color{"\033[0m", color.New()}
	ColorRed    = Color{"\033[31m", color.New(color.FgRed)}
	ColorGreen  = Color{"\033[32m", color.New(color.FgGreen)}
	ColorYellow = Color{"\033[33m", color.New(color.FgYellow)}
	ColorBlue   = Color{"\033[34m", color.New(color.FgBlue)}
	ColorPurple = Color{"\033[35m", color.New(color.FgMagenta)}
	ColorCyan   = Color{"\033[36m", color.New(color.FgCyan)}
	ColorWhite  = Color{"\033[37m", color.New(color.FgWhite)}
)

// 默认主题配置
var (
	// DefaultANSIColorScheme ANSI 默认主题
	DefaultANSIColorScheme = &ColorScheme{
		CodeType: CodeTypeANSI,
		Debug:    &ColorCyan,
		Info:     &ColorGreen,
		Warn:     &ColorYellow,
		Error:    &ColorRed,
		Fatal:    &ColorPurple,
	}

	// DefaultFatihColorScheme Fatih 默认主题
	DefaultFatihColorScheme = &ColorScheme{
		CodeType: CodeTypeFATIH,
		Debug:    &ColorCyan,
		Info:     &ColorGreen,
		Warn:     &ColorYellow,
		Error:    &ColorRed,
		Fatal:    &ColorPurple,
	}

	// HighContrastColorScheme 高对比度主题
	HighContrastColorScheme = &ColorScheme{
		CodeType: CodeTypeANSI,
		Debug:    &Color{"\033[36;1m", color.New(color.FgCyan, color.Bold)},
		Info:     &Color{"\033[32;1m", color.New(color.FgGreen, color.Bold)},
		Warn:     &Color{"\033[33;1m", color.New(color.FgYellow, color.Bold)},
		Error:    &Color{"\033[31;1;4m", color.New(color.FgRed, color.Bold, color.Underline)},
		Fatal:    &Color{"\033[35;1;7m", color.New(color.FgMagenta, color.Bold, color.BgWhite)},
	}
)

type CodeType string

const (
	CodeTypeANSI  CodeType = "ansi"
	CodeTypeFATIH CodeType = "fatih"
)

type Color struct {
	ansi  string
	color *color.Color
}
type ColorScheme struct {
	CodeType CodeType
	Debug    *Color
	Info     *Color
	Warn     *Color
	Error    *Color
	Fatal    *Color
}

// Colorize 使用 fatih/color 的实现
func (col *ColorScheme) Colorize(level Level, msg string) string {
	var c *Color
	switch level {
	case DebugLevel:
		c = col.Debug
	case InfoLevel:
		c = col.Info
	case WarnLevel:
		c = col.Warn
	case ErrorLevel:
		c = col.Error
	case FatalLevel:
		c = col.Fatal
	default:
		return msg
	}
	// 检查颜色指针是否为nil
	if c == nil {
		return msg
	}
	switch col.CodeType {
	case CodeTypeANSI:
		return c.SprintANSI(msg)
	case CodeTypeFATIH:
		return c.Sprint(msg)
	default:
		return c.Sprint(msg)
	}
}

func (c Color) Sprint(msg string) string {
	color.NoColor = false
	if color.NoColor {
		return msg
	}
	return c.color.Sprint(msg) // 兼容 fatih/color
}

func (c Color) SprintANSI(msg string) string {
	return c.ansi + msg + ColorReset.ansi
}
