package chrono

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type WeeklyJob struct {
	Name          string
	Interval      uint
	DaysOfTheWeek gocron.Weekdays
	AtTimes       gocron.AtTimes
	TaskFunc      any
	Parameters    []any
	Hooks         []gocron.EventListener
	err           error
}

func NewWeeklyJob(interval uint, days gocron.Weekdays, atTime gocron.AtTimes) *WeeklyJob {
	return &WeeklyJob{
		Interval:      interval,
		DaysOfTheWeek: days,
		AtTimes:       atTime,
	}
}

func (c *WeeklyJob) Names(name string) *WeeklyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *WeeklyJob) Task(task any, parameters ...any) *WeeklyJob {
	if task == nil {
		c.err = fmt.Errorf("%w: %s", c.err, ErrTaskFuncNil)
	}
	c.TaskFunc = task
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *WeeklyJob) addHooks(hook ...gocron.EventListener) *WeeklyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *WeeklyJob) DefaultHooks() *WeeklyJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *WeeklyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *WeeklyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *WeeklyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *WeeklyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *WeeklyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *WeeklyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}
