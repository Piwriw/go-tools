package progressbar

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestAutoRun(t *testing.T) {
	tasks := make([]ProgressTask, 0)
	for i := 0; i < 10; i++ {
		tasks = append(tasks, NewProgressTask(TaskTime))
	}
	tasks = append(tasks, NewProgressTask(TaskTimeErr, 1))

	tasks = append(tasks, NewProgressTask(Task, 1))
	AutoRun(ProgressOptions().
		Writer(os.Stderr).
		Width(10).
		Throttle(65*time.Millisecond).
		ShowCount().
		ShowIts().
		FullWidth(),
		tasks...)
}

func TestDescribe(t *testing.T) {
	bar := Add(10,
		ProgressOptions().
			Writer(os.Stderr).
			Width(10).
			Throttle(65*time.Millisecond).
			ShowCount().
			ShowIts().
			FullWidth().Completion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}))
	for i := 0; i < 10; i++ {
		bar.Describe(fmt.Sprintf("Task %d", i))
		Task(i)
		bar.Next()
	}
}

func TestAdd(t *testing.T) {
	bar := Add(10,
		ProgressOptions().
			Writer(os.Stderr).
			Width(10).
			ShowTotalBytes().
			Throttle(65*time.Millisecond).
			ShowCount().ShowIts().
			SpinnerType(14).
			FullWidth())
	for i := 0; i < 10; i++ {
		bar.Next()
		Task(i)
	}
}
func TaskTimeErr(num int) error {
	slog.Info("Task Done")
	time.Sleep(time.Duration(1) * time.Second)
	return errors.New("Task Validate Fail")
}

func TaskTime() error {
	slog.Info("Task Done")
	time.Sleep(time.Duration(1) * time.Second)
	return nil
}

func Task(num int) {
	slog.Info("Task Done")
	time.Sleep(time.Duration(num) * time.Second)
}
