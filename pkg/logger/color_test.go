package logger

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

// 在 macOS 上测试这个示例
func TestColorize(t *testing.T) {
	msg := "Hello, macOS Terminal!"
	defaultColorScheme := &ColorScheme{
		Debug: color.New(color.FgCyan),
		Info:  color.New(color.FgGreen),
		Warn:  color.New(color.FgYellow),
		Error: color.New(color.FgRed),
		Fatal: color.New(color.FgMagenta),
	}
	fmt.Println(defaultColorScheme.Colorize(DebugLevel, msg))
	fmt.Println(defaultColorScheme.Colorize(InfoLevel, msg))
	fmt.Println(defaultColorScheme.Colorize(WarnLevel, msg))
	fmt.Println(defaultColorScheme.Colorize(ErrorLevel, msg))
	fmt.Println(defaultColorScheme.Colorize(FatalLevel, msg))
}
