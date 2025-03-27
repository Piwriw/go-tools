package logger

import (
	"io"
	"os"
)

func getFilePathOutputs(filePaths ...string) io.Writer {
	if len(filePaths) == 0 {
		return nil
	}
	iowriters := make([]io.Writer, 0, len(filePaths))
	for _, path := range filePaths {
		output := getOutput(path)
		if output != nil && output != io.Discard {
			iowriters = append(iowriters, output)
		}
	}
	return io.MultiWriter(iowriters...)
}

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
