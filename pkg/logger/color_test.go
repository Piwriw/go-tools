package logger

import (
	"fmt"
	"testing"
)

// 在 macOS 上测试这个示例
func TestColorize(t *testing.T) {
	msg := "Hello, macOS Terminal!"
	fmt.Println(Colorize(DebugLevel, msg, defaultColorScheme))
	fmt.Println(Colorize(InfoLevel, msg, defaultColorScheme))
	fmt.Println(Colorize(WarnLevel, msg, defaultColorScheme))
	fmt.Println(Colorize(ErrorLevel, msg, defaultColorScheme))
	fmt.Println(Colorize(FatalLevel, msg, defaultColorScheme))
}
