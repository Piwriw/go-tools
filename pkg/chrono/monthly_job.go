package chrono

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type MonthJob struct {
	Name           string
	Interval       uint
	DaysOfTheMonth gocron.DaysOfTheMonth
	AtTimes        gocron.AtTimes
	TaskFunc       any
	Parameters     []any
	Hooks          []gocron.EventListener
	err            error
}

func NewMonthJob(interval uint, days gocron.DaysOfTheMonth, atTime gocron.AtTimes) *MonthJob {
	return &MonthJob{
		Interval:       interval,
		DaysOfTheMonth: days,
		AtTimes:        atTime,
	}
}

func (c *MonthJob) Names(name string) *MonthJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *MonthJob) Task(task any, parameters ...any) *MonthJob {
	if task == nil {
		c.err = fmt.Errorf("%w: %s", c.err, ErrTaskFuncNil)
	}
	c.TaskFunc = task
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *MonthJob) addHooks(hook ...gocron.EventListener) *MonthJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *MonthJob) DefaultHooks() *MonthJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *MonthJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *MonthJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *MonthJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *MonthJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *MonthJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *MonthJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *MonthJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *MonthJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}
