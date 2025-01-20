package chrono

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// Scheduler 封装 gocron 的调度器
type Scheduler struct {
	scheduler gocron.Scheduler
}

// NewScheduler creates a new scheduler.
func NewScheduler() (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scheduler{
		scheduler: s,
	}, nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.scheduler.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// RemoveJob 移除任务
func (s *Scheduler) RemoveJob(job gocron.Job) error {
	return s.scheduler.RemoveJob(job.ID())
}

// GetJobs add all Jobs
func (s *Scheduler) GetJobs() ([]gocron.Job, error) {
	return s.scheduler.Jobs(), nil
}

func (s *Scheduler) GetJobByName(jobName string) (gocron.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == jobName {
			return job, nil
		}
	}
	return nil, fmt.Errorf("job %s not found", jobName)
}

func (s *Scheduler) GetJobByID(jobID string) (gocron.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return job, nil
		}
	}
	return nil, fmt.Errorf("job %s not found", jobID)
}

// AddCronJob adds a new cron job.
func (s *Scheduler) AddCronJob(cronExpr string, task func()) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

func (s *Scheduler) AddCronJobWithName(cronExpr string, task func(), taskName string) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
		gocron.WithName(taskName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

func (s *Scheduler) AddCronJobWithOptions(cronExpr string, task func(), options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

// AddOnceJob adds a new cron job.
func (s *Scheduler) AddOnceJob(task func(), times ...time.Time) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)),
		gocron.NewTask(task),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

// AddIntervalJob 添加一个间隔任务
// func (s *Scheduler) AddIntervalJob(interval time.Duration, task func()) (gocron.CronJob, error) {
// 	job, err := s.scheduler.NewCronJob(
// 		gocron.DurationJob(interval), // 使用时间间隔
// 		gocron.NewTask(task),         // 任务函数
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("添加任务失败: %w", err)
// 	}
// 	return job, nil
// }
