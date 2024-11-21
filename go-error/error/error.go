package main

import (
	orlerror "errors"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

var zeroErr = orlerror.New("zero error")

func doTask() error {
	//slog.Error("task is failed", slog.Any("error", zeroErr))
	return zeroErr
}

func doTaskWarp() error {
	//slog.Error("task is failed", slog.Any("error", zeroErr))
	return errors.Wrapf(zeroErr, "doTaskWarp")
}

func doTaskErrorf() error {
	//slog.Error("task is failed", slog.Any("error", zeroErr))
	return errors.Errorf("doTaskErrof:%v", zeroErr)
}

func doTaskFmtErrorf() error {
	//slog.Error("task is failed", slog.Any("error", zeroErr))
	return fmt.Errorf("doTaskErrof:%v", zeroErr)
}

func doTaskErrorw() error {
	//slog.Error("task is failed", slog.Any("error", zeroErr))
	return fmt.Errorf("doTaskErrow:%w", zeroErr)
}

func main() {
	if err := doTask(); err != nil {
		slog.Error("doTask err", slog.Any("error", err))
	}
	if err := doTaskWarp(); err != nil {
		slog.Error("doTask err", slog.Any("error", err))
	}
	if err := doTaskErrorf(); err != nil {
		slog.Error("doTask err", slog.Any("error", err))
	}
	if err := doTaskErrorw(); err != nil {
		slog.Error("doTask err", slog.Any("error", err))
		slog.Error("doTask UnWarp err", slog.Any("error", errors.Cause(err)))
	}
	if err := doTaskFmtErrorf(); err != nil {
		slog.Error("doTask err", slog.Any("error", err))
	}

}
