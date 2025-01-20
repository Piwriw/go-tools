package gocron

import (
	"fmt"
	"time"

	"github.piwriw.go-tools/pkg/chrono"

	"github.com/google/uuid"

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
func (s *Scheduler) RemoveJob(job chrono.Job) error {
	jobID, err := uuid.Parse(job.ID())
	if err != nil {
		return fmt.Errorf("failed to parse job ID: %w", err)
	}

	// 调用底层的 RemoveJob 方法
	return s.scheduler.RemoveJob(jobID)
}

// GetJobs add all Jobs
func (s *Scheduler) GetJobs() ([]chrono.Job, error) {
	jobs := make([]chrono.Job, 0)
	for _, job := range s.scheduler.Jobs() {
		jobs = append(jobs, NewCronJob(job))
	}
	return jobs, nil
}

func (s *Scheduler) GetJobByName(jobName string) (chrono.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == jobName {
			return NewCronJob(job), nil
		}
	}
	return nil, fmt.Errorf("job %s not found", jobName)
}

func (s *Scheduler) GetJobByID(jobID string) (chrono.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return NewCronJob(job), nil
		}
	}
	return nil, fmt.Errorf("job %s not found", jobID)
}

// AddCronJob adds a new cron job.
func (s *Scheduler) AddCronJob(cronExpr string, task func()) (chrono.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return NewCronJob(job), nil
}

func (s *Scheduler) AddCronJobWithName(cronExpr string, task func(), taskName string) (chrono.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
		gocron.WithName(taskName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return &CronJob{
		goJob: job,
	}, nil
}

// convertJobOption 将 gron.JobOption 转换为 gocron.JobOption
func convertJobOption(option ...chrono.JobOption) (gocron.JobOption, error) {

	return nil, fmt.Errorf("not implemented")
}
func (s *Scheduler) AddCronJobWithOptions(cronExpr string, task func(), options ...chrono.JobOption) (chrono.Job, error) {
	convertJobOption(options...)
	job, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false), // Use cron expression
		gocron.NewTask(task),            // Task function
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return &CronJob{
		goJob: job,
	}, nil
}

// AddOnceJob adds a new cron job.
func (s *Scheduler) AddOnceJob(task func(), times ...time.Time) (chrono.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)),
		gocron.NewTask(task),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return &CronJob{
		goJob: job,
	}, nil
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
