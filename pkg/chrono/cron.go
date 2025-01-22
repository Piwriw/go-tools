package chrono

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

const (
	DayTimeType   string = "%d %d * * * "
	WeekTimeType  string = "%d %d * * %d "
	MonthTimeType string = "%d %d * %d * "
)

type CronJob struct {
	Name       string
	Expr       string
	TaskFunc   any
	Parameters []any
	Hooks      gocron.EventListener
	err        error
}

func NewCronJob(expr string) *CronJob {
	return &CronJob{
		Expr: expr,
	}
}

func (c *CronJob) Error() string {
	return c.err.Error()
}

func (c *CronJob) Names(name string) *CronJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *CronJob) Task(task any, parameters ...any) *CronJob {
	if task == nil {
		c.err = fmt.Errorf("%w; wrong task: task function cannot be nil", c.err)
	}
	c.TaskFunc = task
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *CronJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	c.Hooks = gocron.BeforeJobRuns(eventListenerFunc)
	return c
}

func (c *CronJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
	c.Hooks = gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

func (c *CronJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	c.Hooks = gocron.AfterJobRuns(eventListenerFunc)
	return c
}

func (c *CronJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	c.Hooks = gocron.AfterJobRunsWithError(eventListenerFunc)
	return c
}

func (c *CronJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
	c.Hooks = gocron.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

func (c *CronJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	c.Hooks = gocron.AfterLockError(eventListenerFunc)
	return c
}

type TimeType string

// DayTimeToCron 将 time.Time 转换为 Cron 表达式
func DayTimeToCron(t time.Time) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(DayTimeType, minute, hour)
}

func WeekTimeToCron(t time.Time, week time.Weekday) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(WeekTimeType, minute, hour, week)
}

func MonthTimeToCron(t time.Time, month time.Month) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(MonthTimeType, minute, hour, month)
}
