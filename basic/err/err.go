package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

func Cause(err error) error {
	for err != nil {
		err = errors.Unwrap(err)
	}
	return err
}

func main() {
	bytes, err := do("")
	if err != nil {
		slog.Error("do is failed ", slog.Any("err", err))
		slog.Error("原始错误打印", slog.Any("err", Cause(err)))
		return
	}
	slog.Info("do is success ", slog.Any("bytes", bytes))
}

func do(str string) ([]byte, error) {
	if str == "" {
		err := errors.New("原始错误")
		w := fmt.Errorf("外面包了一个错误%w", err)
		errDouble := fmt.Errorf("外面包了一个错误2层%w", w)
		return nil, errDouble
	}
	marshal, err := json.Marshal(str)
	if err != nil {
		slog.Error("marshal is failed ", slog.Any("err", err))
		return nil, err
	}
	return marshal, nil

}
