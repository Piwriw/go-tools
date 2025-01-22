package chrono

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type defaultSchedulerMonitor struct {
	mu       sync.Mutex
	counter  map[string]int
	time     map[string][]time.Duration
	taskList map[uuid.UUID]MonitorTaskSpec
}

type MonitorTaskSpec struct {
	TaskID    uuid.UUID
	TaskName  string
	StartTime time.Time
	EndTime   time.Time
	Status    gocron.JobStatus
	Tags      []string
	Err       error
}

func NewMonitorTaskSpec(id uuid.UUID, name string, startTime, endTime time.Time, tags []string, status gocron.JobStatus, err error) MonitorTaskSpec {
	return MonitorTaskSpec{
		TaskID:    id,
		TaskName:  name,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		Tags:      tags,
		Err:       err,
	}
}

func newDefaultSchedulerMonitor() *defaultSchedulerMonitor {
	return &defaultSchedulerMonitor{
		counter:  make(map[string]int),
		time:     make(map[string][]time.Duration),
		taskList: make(map[uuid.UUID]MonitorTaskSpec),
	}
}

// IncrementJob 增加任务的执行次数
func (s *defaultSchedulerMonitor) IncrementJob(id uuid.UUID, name string, tags []string, status gocron.JobStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("IncrementJob", "JobID", id, "JobName", name, "tags", tags, "status", status)
	_, ok := s.counter[name]
	if !ok {
		s.counter[name] = 0
	}
	s.counter[name]++
}

// RecordJobTiming 记录任务的执行时间
func (s *defaultSchedulerMonitor) RecordJobTiming(startTime, endTime time.Time, id uuid.UUID, name string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("RecordJobTiming", "JobID", id, "JobName", name, "startTime", startTime.Format("2006-01-02 15:04:05"),
		"endTime", endTime.Format("2006-01-02 15:04:05"), "duration", endTime.Sub(startTime), "tags", tags)
	_, ok := s.time[name]
	if !ok {
		s.time[name] = make([]time.Duration, 0)
	}
	s.time[name] = append(s.time[name], endTime.Sub(startTime))
}

// RecordJobTimingWithStatus 记录任务的执行时间、状态
func (s *defaultSchedulerMonitor) RecordJobTimingWithStatus(startTime, endTime time.Time, id uuid.UUID, name string, tags []string, status gocron.JobStatus, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("RecordJobTimingWithStatus", "JobID", id, "JobName", name, "startTime", startTime.Format("2006-01-02 15:04:05"),
		"endTime", endTime.Format("2006-01-02 15:04:05"), "duration", endTime.Sub(startTime), "status", status, "err", err)
	s.taskList[id] = NewMonitorTaskSpec(id, name, startTime, endTime, tags, status, err)
}

// Scheduler 封装 gocron 的调度器
type Scheduler struct {
	scheduler gocron.Scheduler
	monitor   gocron.MonitorStatus
}

type Options struct {
	Monitor bool
}

// NewScheduler creates a new scheduler.
func NewScheduler(monitor gocron.MonitorStatus) (*Scheduler, error) {
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
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}

	// 创建一个新的定时任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false),                 // 使用 cron 表达式
		gocron.NewTask(job.TaskFunc, job.Parameters...), // 任务函数
		gocron.WithEventListeners(job.Hooks...),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}

	return jobInstance, nil
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
func (s *Scheduler) AddIntervalJob(interval time.Duration, task func()) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval), // 使用时间间隔
		gocron.NewTask(task),         // 任务函数
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

// AddIntervalJobWithName 添加一个间隔任务,支持任务名称
func (s *Scheduler) AddIntervalJobWithName(interval time.Duration, task func(), name string) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval), // 使用时间间隔
		gocron.NewTask(task),         // 任务函数
		gocron.WithName(name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}

// AddIntervalJobWithOptions 添加一个间隔任务,支持原生拓展方法
func (s *Scheduler) AddIntervalJobWithOptions(interval time.Duration, task func(), options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval), // 使用时间间隔
		gocron.NewTask(task),         // 任务函数
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	return job, nil
}
