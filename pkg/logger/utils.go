package logger

import (
	"io"
	"os"
)

func getOutput(filePath string) io.Writer {
	if filePath == "" {
		return io.Discard
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return io.Discard
	}
	return file
}
