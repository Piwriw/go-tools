package chrono

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// Scheduler 封装 gocron 的调度器
type Scheduler struct {
	scheduler gocron.Scheduler
	monitor   SchedulerMonitor
}

// NewScheduler creates a new scheduler.
func NewScheduler(monitor SchedulerMonitor) (*Scheduler, error) {
	// 根据 monitor 是否为空来决定如何创建调度器
	if monitor == nil {
		monitor = newDefaultSchedulerMonitor()
	}
	s, err := gocron.NewScheduler(gocron.WithMonitorStatus(monitor), gocron.WithMonitor(monitor))
	// 错误处理
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	// 返回 Scheduler
	return &Scheduler{
		scheduler: s,
		monitor:   monitor,
	}, nil
}

// Start Starts the scheduler.
func (s *Scheduler) Start() {
	s.scheduler.Start()
}

// Stop Stops the scheduler.
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// RemoveJob Removes a job.
func (s *Scheduler) RemoveJob(job gocron.Job) error {
	return s.scheduler.RemoveJob(job.ID())
}

func (s *Scheduler) Watch() {
	s.monitor.Watch()
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

// GetJobByID get job BY ID
func (s *Scheduler) GetJobByID(jobID string) (gocron.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return job, nil
		}
	}
	return nil, fmt.Errorf("job %s not found", jobID)
}

// AddCronJob adds a new cron job.
func (s *Scheduler) AddCronJob(job *CronJob) (gocron.Job, error) {
	if job == nil {
		return nil, fmt.Errorf("job cannot be nil")
	}
	if job.err != nil {
		return nil, job.err
	}
	// check if job has a task function
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}

	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false),                 // 使用 cron 表达式
		gocron.NewTask(job.TaskFunc, job.Parameters...), // 任务函数
		gocron.WithEventListeners(job.Hooks...),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add cron job: %w", err)
	}

	return jobInstance, nil
}

// AddCronJobs adds list new cron job.
func (s *Scheduler) AddCronJobs(jobs ...*CronJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, cronJob := range jobs {
		cronJobInstance, err := s.AddCronJob(cronJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add cron jobs: %v", errs)
	}
	return jobList, nil
}

func (s *Scheduler) AddCronJobWithOptions(job *CronJob, options ...gocron.JobOption) (gocron.Job, error) {
	if job == nil {
		return nil, fmt.Errorf("job cannot be nil")
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false),
		gocron.NewTask(job.TaskFunc),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return jobInstance, nil
}

// AddOnceJob adds a new once job.
func (s *Scheduler) AddOnceJob(job *OnceJob) (gocron.Job, error) {
	jobInstance, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(job.WorkTime...)),
		gocron.NewTask(job.TaskFunc, job.Parameters...), // 任务函数
		gocron.WithEventListeners(job.Hooks...),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add once job: %w", err)
	}
	return jobInstance, nil
}

// AddOnceJobs adds list new once job.
func (s *Scheduler) AddOnceJobs(jobs ...*OnceJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, onceJob := range jobs {
		cronJobInstance, err := s.AddOnceJob(onceJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add once jobs: %v", errs)
	}
	return jobList, nil
}

// AddIntervalJob 添加一个间隔任务
func (s *Scheduler) AddIntervalJob(job *IntervalJob) (gocron.Job, error) {
	jobInstance, err := s.scheduler.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(job.TaskFunc, job.Parameters...),
		gocron.WithEventListeners(job.Hooks...),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return jobInstance, nil
}

// AddIntervalJobs adds list new interval job.
func (s *Scheduler) AddIntervalJobs(jobs ...*IntervalJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, intervalJob := range jobs {
		intervalJobInstance, err := s.AddIntervalJob(intervalJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, intervalJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add interval jobs: %v", errs)
	}
	return jobList, nil
}

// AddIntervalJobWithOptions 添加一个间隔任务,支持原生拓展方法
func (s *Scheduler) AddIntervalJobWithOptions(interval time.Duration, task func(), options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval), // 使用时间间隔
		gocron.NewTask(task),         // 任务函数
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add interval job: %w", err)
	}
	return job, nil
}

// AddDailyJob add daily Job
func (s *Scheduler) AddDailyJob(job *DailyJob) (gocron.Job, error) {
	jobInstance, err := s.scheduler.NewJob(
		gocron.DailyJob(job.Interval, job.AtTimes),
		gocron.NewTask(job.TaskFunc, job.Parameters...),
		gocron.WithEventListeners(job.Hooks...),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return jobInstance, nil
}

// AddDailyJobs adds list new interval job.
func (s *Scheduler) AddDailyJobs(jobs ...*DailyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, dailyJob := range jobs {
		dailyJobInstance, err := s.AddDailyJob(dailyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add daily jobs: %v", errs)
	}
	return jobList, nil
}

// AddDailyJobWithOptions 添加一个间隔任务,支持原生拓展方法
func (s *Scheduler) AddDailyJobWithOptions(interval time.Duration, task func(), options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval), // 使用时间间隔
		gocron.NewTask(task),         // 任务函数
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add interval job: %w", err)
	}
	return job, nil
}
