package chrono

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type IntervalJob struct {
	Name       string
	Interval   time.Duration
	TaskFunc   any
	Parameters []any
	Hooks      []gocron.EventListener
	WatchFunc  func(event MonitorJobSpec)
	err        error
}

func NewIntervalJob(interval time.Duration) *IntervalJob {
	return &IntervalJob{
		Interval: interval,
	}
}

func (c *IntervalJob) Names(name string) *IntervalJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *IntervalJob) Task(task any, parameters ...any) *IntervalJob {
	if task == nil {
		c.err = errors.Join(c.err, ErrTaskFuncNil)
		return c
	}
	c.TaskFunc = task
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *IntervalJob) Watch(watch func(event MonitorJobSpec)) *IntervalJob {
	c.WatchFunc = watch
	return c
}

func (c *IntervalJob) addHooks(hook ...gocron.EventListener) *IntervalJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *IntervalJob) DefaultHooks() *IntervalJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *IntervalJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *IntervalJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *IntervalJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *IntervalJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *IntervalJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *IntervalJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}
