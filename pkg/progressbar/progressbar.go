package progressbar

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"log/slog"
	"time"
)

const (
	defaultMetricPort = "0.0.0.0:19999"
)

type ProgressBar struct {
	total   int
	bar     *progressbar.ProgressBar
	options Options
}

type Options struct {
	options []progressbar.Option
}

func (p *ProgressBar) Total(total int) *ProgressBar {
	p.total = total
	return p
}

// Metric starts an HTTP server dedicated to serving progress bar updates. This allows you to
// display the status in various UI elements, such as an OS status bar with an `xbar` extension.
// It is recommended to run this function in a separate goroutine to avoid blocking the main thread.
//
// hostPort specifies the address and port to bind the server to, for example, "0.0.0.0:19999".
func (p *ProgressBar) Metric(hostPort string) {
	if hostPort == "" {
		hostPort = defaultMetricPort
	}
	p.bar.StartHTTPServer(hostPort)
}

func Add(total int, ps *Options) *ProgressBar {
	return &ProgressBar{
		bar: progressbar.NewOptions(total, ps.options...),
	}
}

type ProgressTask struct {
	fn     any
	params []any
}

func NewProgressTask(fn any, params ...any) ProgressTask {
	return ProgressTask{
		fn:     fn,
		params: params,
	}
}

func AutoRun(ps *Options, tasks ...ProgressTask) error {
	bar := &ProgressBar{
		bar: progressbar.NewOptions(len(tasks), ps.options...),
	}
	for _, task := range tasks {
		if err := callFunc(task.fn, task.params...); err != nil {
			slog.Error("AutoRun", "err", err)
			return bar.Exit()
		}
		if err := bar.Next(); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProgressBar) Next() error {
	return p.bar.Add(1)
}

// Finish will fill the bar to full
func (p *ProgressBar) Finish() error {
	return p.bar.Finish()
}

// Exit will exit the bar to keep current state
func (p *ProgressBar) Exit() error {
	return p.bar.Exit()
}

// Clear erases the progress bar from the current line
func (p *ProgressBar) Clear() error {
	return p.bar.Clear()
}

// Set will set the bar to a current number
func (p *ProgressBar) Set(step int) error {
	return p.bar.Set(step)
}

// IsFinished returns true if progress bar is completed
func (p *ProgressBar) IsFinished() bool {
	return p.bar.IsFinished()
}

// IsStarted returns true if progress bar is started
func (p *ProgressBar) IsStarted() bool {
	return p.bar.IsStarted()
}

// State returns the current state
func (p *ProgressBar) State() progressbar.State {
	return p.bar.State()
}

// Describe will change the description shown before the progress, which
// can be changed on the fly (as for a slow running process).
func (p *ProgressBar) Describe(description string) {
	p.bar.Describe(description)
}

func ProgressOptions() *Options {
	return &Options{
		options: make([]progressbar.Option, 0),
	}
}

func (p *Options) Writer(w io.Writer) *Options {
	p.options = append(p.options, progressbar.OptionSetWriter(w))
	return p
}

func (p *Options) Width(width int) *Options {
	p.options = append(p.options, progressbar.OptionSetWidth(width))
	return p
}

func (p *Options) ShowTotalBytes() *Options {
	p.options = append(p.options, progressbar.OptionShowTotalBytes(true))
	return p
}

func (p *Options) SpinnerChangeInterval(interval time.Duration) *Options {
	p.options = append(p.options, progressbar.OptionSetSpinnerChangeInterval(interval))
	return p
}

func (p *Options) SpinnerType(spinnerType int) *Options {
	p.options = append(p.options, progressbar.OptionSpinnerType(spinnerType))
	return p
}

func (p *Options) SpinnerCustom(spinner []string) *Options {
	p.options = append(p.options, progressbar.OptionSpinnerCustom(spinner))
	return p
}

func (p *Options) Theme(t progressbar.Theme) *Options {
	p.options = append(p.options, progressbar.OptionSetTheme(t))
	return p
}

func (p *Options) Visibility(visibility bool) *Options {
	p.options = append(p.options, progressbar.OptionSetVisibility(visibility))
	return p
}

func (p *Options) FullWidth() *Options {
	p.options = append(p.options, progressbar.OptionFullWidth())
	return p
}

func (p *Options) RenderBlankState(r bool) *Options {
	p.options = append(p.options, progressbar.OptionSetRenderBlankState(r))
	return p
}

// Throttle will wait the specified duration before updating again. The default
// duration is 0 seconds.
func (p *Options) Throttle(duration time.Duration) *Options {
	p.options = append(p.options, progressbar.OptionThrottle(duration))
	return p
}

func (p *Options) ShowCount() *Options {
	p.options = append(p.options, progressbar.OptionShowCount())
	return p
}

func (p *Options) ShowIts() *Options {
	p.options = append(p.options, progressbar.OptionShowIts())
	return p
}

func (p *Options) Completion(cmpl func()) *Options {
	p.options = append(p.options, progressbar.OptionOnCompletion(cmpl))
	return p
}

func (p *Options) EnableColorCodes(colorCodes bool) *Options {
	p.options = append(p.options, progressbar.OptionEnableColorCodes(colorCodes))
	return p
}

func (p *Options) ElapsedTime(elapsedTime bool) *Options {
	p.options = append(p.options, progressbar.OptionSetElapsedTime(elapsedTime))
	return p
}

func (p *Options) PredictTime(predictTime bool) *Options {
	p.options = append(p.options, progressbar.OptionSetPredictTime(true))
	return p
}

func (p *Options) OptionShowCount() *Options {
	p.options = append(p.options, progressbar.OptionShowCount())
	return p
}

func (p *Options) ShowElapsedTimeOnFinish() *Options {
	p.options = append(p.options, progressbar.OptionShowElapsedTimeOnFinish())
	return p
}

func (p *Options) SetItsString(iterationString string) *Options {
	p.options = append(p.options, progressbar.OptionSetItsString(iterationString))
	return p
}

func (p *Options) ClearOnFinish() *Options {
	p.options = append(p.options, progressbar.OptionClearOnFinish())
	return p
}

func (p *Options) OptionShowBytes(val bool) *Options {
	p.options = append(p.options, progressbar.OptionShowBytes(val))
	return p
}

func (p *Options) UseANSICodes(val bool) *Options {
	p.options = append(p.options, progressbar.OptionUseANSICodes(val))
	return p
}

func (p *Options) OptionUseIECUnits(val bool) *Options {
	p.options = append(p.options, progressbar.OptionUseIECUnits(val))
	return p
}

func (p *Options) ShowDescriptionAtLineEnd() *Options {
	p.options = append(p.options, progressbar.OptionShowDescriptionAtLineEnd())
	return p
}

func (p *Options) MaxDetailRow(row int) *Options {
	p.options = append(p.options, progressbar.OptionSetMaxDetailRow(row))
	return p
}
