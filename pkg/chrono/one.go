package chrono

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type OnceJob struct {
	Name       string
	WorkTime   []time.Time
	TaskFunc   any
	Parameters []any
	Hooks      []gocron.EventListener
	err        error
}

func NewOnceJob(workTimes ...time.Time) *OnceJob {
	return &OnceJob{
		WorkTime: workTimes,
	}
}

func (c *OnceJob) Names(name string) *OnceJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *OnceJob) Task(task any, parameters ...any) *OnceJob {
	if task == nil {
		c.err = fmt.Errorf("%w: %s", c.err, ErrTaskFuncNil)
	}
	c.TaskFunc = task
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *OnceJob) addHooks(hook ...gocron.EventListener) *OnceJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *OnceJob) DefaultHooks() *OnceJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *OnceJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *OnceJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *OnceJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *OnceJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *OnceJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *OnceJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *OnceJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *OnceJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}
