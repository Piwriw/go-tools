package gocron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

// CronJob 自定义任务类型
type CronJob struct {
	goJob gocron.Job // 底层的 gocron.CronJob
}

// NewCronJob 创建一个新的自定义任务
func NewCronJob(gocronJob gocron.Job) *CronJob {
	return &CronJob{
		goJob: gocronJob,
	}
}

func (j *CronJob) ID() string {
	return j.goJob.ID().String()
}

// Name 获取任务名称
func (j *CronJob) Name() string {
	return j.goJob.Name()
}

// NextRun 获取下次执行时间
func (j *CronJob) NextRun() (time.Time, error) {
	return j.goJob.NextRun()
}

// LastRun 获取上次执行时间
func (j *CronJob) LastRun() (time.Time, error) {
	return j.goJob.LastRun()
}
