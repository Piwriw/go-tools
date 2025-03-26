package logger

import "os"

func getOutput(filePath string) *os.File {
	if filePath == "" {
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	return file
}
