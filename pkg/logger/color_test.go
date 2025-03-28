package logger

import (
	"fmt"
	"testing"
)

// 在 macOS 上测试这个示例
func TestColorize(t *testing.T) {
	msg := "Hello, macOS Terminal!"

	fmt.Println(DefaultFatihColorScheme.Colorize(DebugLevel, msg))
	fmt.Println(DefaultFatihColorScheme.Colorize(InfoLevel, msg))
	fmt.Println(DefaultFatihColorScheme.Colorize(WarnLevel, msg))
	fmt.Println(DefaultFatihColorScheme.Colorize(ErrorLevel, msg))
	fmt.Println(DefaultFatihColorScheme.Colorize(FatalLevel, msg))
}
