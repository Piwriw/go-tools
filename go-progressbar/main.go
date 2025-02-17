package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log/slog"
	"os"
	"time"
)

func Task(bar *progressbar.ProgressBar, num int) {
	time.Sleep(time.Duration(num) * time.Second)
	slog.Info("current", slog.Any("xxx", bar.State().CurrentNum))

}
func main() {
	add()
	// finish()
	// exit()
	//addDetail()
	//startHTTPServer()
}
func add() {
	bar := progressbar.Default(10)
	for i := 0; i < 10; i++ {
		bar.Add(1)
		Task(bar, i)
	}
}
func finish() {
	bar := progressbar.Default(10, "Down", "xx")

	for i := 0; i < 10; i++ {
		Task(bar, i)
		bar.Add(1)
		if i == 5 {
			bar.Finish()
		}
	}
}

func exit() {
	bar := progressbar.Default(10, "Down", "xx")

	for i := 0; i < 10; i++ {
		Task(bar, i)
		bar.Add(1)
		if i == 5 {
			bar.Exit()
		}
	}
	bar.State()
}

func addDetail() {
	bar := progressbar.NewOptions64(
		10,
		progressbar.OptionSetDescription("xxx"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowTotalBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetMaxDetailRow(5),
	)

	for i := 0; i < 10; i++ {
		Task(bar, i)
		bar.Add(1)

	}
}

func startHTTPServer() {
	bar := progressbar.Default(10, "Down", "xx")
	bar.StartHTTPServer("0.0.0.0:19999")

	for i := 0; i < 10; i++ {
		Task(bar, i)
		bar.Add(1)
	}
}
