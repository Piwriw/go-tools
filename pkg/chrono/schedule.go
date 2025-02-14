package chrono

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// Scheduler base gocron scheduler
type Scheduler struct {
	ctx          context.Context
	scheduler    gocron.Scheduler
	monitor      SchedulerMonitor
	aliasMap     map[string]string
	watchFuncMap map[string]func(event JobWatchInterface)
	mu           sync.Mutex // 用于保护 watchFuncMap
	schOptions   *SchedulerOptions
}

type SchedulerOptions struct {
	aliasEnable ChronoOption
}

// Enable 用于查询某个选项是否启用
func (s *Scheduler) Enable(option string) bool {
	switch option {
	case s.schOptions.aliasEnable.Name():
		if s.schOptions.aliasEnable != nil {
			return s.schOptions.aliasEnable.Enable()
		}
		return false
	}
	return false
}

type SchedulerOption func(*SchedulerOptions)

func WithAliasMode(enabled bool) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.aliasEnable = &AliasOption{enabled: enabled}
	}
}

func WithWatch(enabled bool) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.aliasEnable = &AliasOption{enabled: enabled}
	}
}

type Event struct {
	JobID       string
	JobName     string
	NextRunTime time.Time
	LastTime    time.Time
	Err         error
}

func (s *Scheduler) Watch() {
	event := s.monitor.Watch()
	for {
		select {
		case <-s.ctx.Done():
			return
		case e := <-event:
			fn, ok := s.watchFuncMap[e.GetJobID()]
			if !ok {
				slog.Error("job not found", "jobID", e.GetJobID())
				continue
			}
			fn(e)
		}
	}
}

// NewScheduler creates a new scheduler.
func NewScheduler(ctx context.Context, monitor SchedulerMonitor, options ...SchedulerOption) (*Scheduler, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 根据 monitor 是否为空来决定如何创建调度器
	if monitor == nil {
		monitor = newDefaultSchedulerMonitor()
	}
	s, err := gocron.NewScheduler(gocron.WithMonitorStatus(monitor), gocron.WithMonitor(monitor))
	// 错误处理
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}
	schOptions := &SchedulerOptions{}
	for _, option := range options {
		option(schOptions)
	}

	return &Scheduler{
		scheduler:    s,
		monitor:      monitor,
		ctx:          ctx,
		watchFuncMap: make(map[string]func(event JobWatchInterface)),
		aliasMap:     make(map[string]string),
		schOptions:   schOptions,
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
func (s *Scheduler) RemoveJob(jobID string) error {
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID %s: %w", jobID, err)
	}
	s.removeAlias(jobID)
	s.removeWatchFunc(jobID)
	return s.scheduler.RemoveJob(jobUUID)
}

// RemoveJobByAlias Removes a job by alias.
func (s *Scheduler) RemoveJobByAlias(alias string) error {
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return fmt.Errorf("alias %s not found", alias)
	}
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID %s: %w", jobID, err)
	}
	s.removeAlias(jobID)
	s.removeWatchFunc(jobID)
	return s.scheduler.RemoveJob(jobUUID)
}

func (s *Scheduler) GetAlias(taskID string) (string, error) {
	for alias, readTaskID := range s.aliasMap {
		if taskID == readTaskID {
			return alias, nil
		}
	}
	return "", ErrFoundAlias
}

// addAlias add alias
func (s *Scheduler) addAlias(alias string, jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if alias == "" {
		slog.Warn("alias is empty", "alias", alias)
		return
	}
	if jobID == "" {
		slog.Warn("jobID is empty", "jobID", jobID)
		return
	}
	s.aliasMap[alias] = jobID
}

// removeWatchFunc  remove alias
func (s *Scheduler) removeAlias(alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.aliasMap[alias]; exists {
		delete(s.aliasMap, alias)
		slog.Info("alias  removed", "alias", alias)
	} else {
		slog.Warn("alias not found in aliasMap", "alias", alias)
	}
}

// addWatchFunc add watch Func
func (s *Scheduler) addWatchFunc(jobID string, fn func(event JobWatchInterface)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jobID == "" {
		slog.Warn("jobID is empty", "jobID", jobID)
		return
	}
	if fn == nil {
		slog.Warn("watchFunc is empty", "jobID", jobID)
		return
	}
	s.watchFuncMap[jobID] = fn
}

// removeWatchFunc  remove watch Func
func (s *Scheduler) removeWatchFunc(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.watchFuncMap[jobID]; exists {
		delete(s.watchFuncMap, jobID)
		slog.Info("Watch function removed", "jobID", jobID)
	} else {
		slog.Warn("Job not found in watchFuncMap", "jobID", jobID)
	}
}

// TODO 批量移除任务
// RemoveJobs Removes job list with rollback support.
// func (s *Scheduler) RemoveJobs(jobIDS ...string) error {
// 	// 获取所有需要删除的任务
// 	jobs, err := s.GetJobByIDS(jobIDS...)
// 	if err != nil {
// 		return fmt.Errorf("failed to get jobs: %w", err)
// 	}
//
// 	// 记录成功删除的任务
// 	removedJobs := make([]gocron.Job, 0, len(jobs))
//
// 	// 遍历任务列表，逐个删除任务
// 	for _, job := range jobs {
// 		if err := s.RemoveJob(job.ID().String()); err != nil {
// 			// 如果删除失败，回滚已删除的任务
// 			if rollbackErr := s.rollbackRemovedJobs(removedJobs); rollbackErr != nil {
// 				return fmt.Errorf("failed to remove job %s: %w; rollback failed: %v", job.ID(), err, rollbackErr)
// 			}
// 			return fmt.Errorf("failed to remove job %s: %w", job.ID(), err)
// 		}
// 		// 记录成功删除的任务
// 		removedJobs = append(removedJobs, job)
// 	}
//
// 	return nil
// }
//
// // rollbackRemovedJobs 回滚已删除的任务
// func (s *Scheduler) rollbackRemovedJobs(jobs []gocron.Job) error {
// 	var rollbackErrors []error
//
// 	// 遍历已删除的任务，逐个重新添加
// 	for _, job := range jobs {
// 		if err := s.scheduler.AddJob(job); err != nil {
// 			rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to re-add job %s: %w", job.ID(), err))
// 		}
// 	}
//
// 	// 如果有回滚错误，返回合并后的错误
// 	if len(rollbackErrors) > 0 {
// 		return fmt.Errorf("rollback errors: %v", rollbackErrors)
// 	}
//
// 	return nil
// }

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

// GetJobByAlias get job BY Alias
func (s *Scheduler) GetJobByAlias(alias string) (gocron.Job, error) {
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, fmt.Errorf("alias %s not found", alias)
	}
	return s.GetJobByID(jobID)
}

// GetJobByIDOrAlias get job BY ID or Alias,first by ID, then by Alias
func (s *Scheduler) GetJobByIDOrAlias(identifier string) (gocron.Job, error) {
	// 优先通过ID查找
	if jobID, err := s.GetJobByID(identifier); err == nil {
		return jobID, nil
	}

	// 如果没有找到ID，尝试通过别名查找
	if jobID, exists := s.aliasMap[identifier]; exists {
		return s.GetJobByAlias(jobID)
	}

	return nil, fmt.Errorf("job with identifier %s not found", identifier)
}

// GetJobByIDS gets jobs by IDs.
func (s *Scheduler) GetJobByIDS(jobIDS ...string) ([]gocron.Job, error) {
	// 创建一个切片用于存储找到的任务
	jobs := make([]gocron.Job, 0, len(jobIDS))

	// 遍历 jobIDS，逐个查找任务
	for _, jobID := range jobIDS {
		job, err := s.GetJobByID(jobID)
		if err != nil {
			return nil, fmt.Errorf("failed to get job %s: %w", jobID, err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// AddCronJob adds a new cron job.
func (s *Scheduler) AddCronJob(job *CronJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// check if job has a task function
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false), // 使用 cron 表达式
		gocron.NewTask(job.TaskFunc),    // 任务函数
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add cron job: %w", err)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
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
		return nil, ErrInvalidJob
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
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(job.WorkTime...)),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add once job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
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

// AddOnceJobWithOptions add a once job, support native extension method
func (s *Scheduler) AddOnceJobWithOptions(startAt gocron.OneTimeJobStartAtOption, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(startAt),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add once job: %w", err)
	}
	return job, nil
}

// AddIntervalJob  add interval job
func (s *Scheduler) AddIntervalJob(job *IntervalJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	// Job options
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}

	jobInstance, err := s.scheduler.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
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

// AddIntervalJobWithOptions add a interval job, support native extension method
func (s *Scheduler) AddIntervalJobWithOptions(interval time.Duration, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add interval job: %w", err)
	}
	return job, nil
}

// AddDailyJob add daily Job
func (s *Scheduler) AddDailyJob(job *DailyJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrTaskFuncNil
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.DailyJob(job.Interval, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
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

// AddDailyJobWithOptions add a daily job, support native extension method
func (s *Scheduler) AddDailyJobWithOptions(interval uint, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DailyJob(interval, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add daily job: %w", err)
	}
	return job, nil
}

// AddWeeklyJob add weekly Job
func (s *Scheduler) AddWeeklyJob(job *WeeklyJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.WeeklyJob(job.Interval, job.DaysOfTheWeek, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add weekly job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddWeeklyJobs adds list new Weekly job.
func (s *Scheduler) AddWeeklyJobs(jobs ...*WeeklyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, weeklyJob := range jobs {
		dailyJobInstance, err := s.AddWeeklyJob(weeklyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add weekly jobs: %v", errs)
	}
	return jobList, nil
}

// AddWeeklyJobWithOptions add a weekly job, support native extension method
func (s *Scheduler) AddWeeklyJobWithOptions(interval uint, daysOfTheWeek gocron.Weekdays, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.WeeklyJob(interval, daysOfTheWeek, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add weekly job: %w", err)
	}
	return job, nil
}

// AddMonthlyJob add monthly Job
func (s *Scheduler) AddMonthlyJob(job *MonthJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("job %s has no task function", job.Name)
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.MonthlyJob(job.Interval, job.DaysOfTheMonth, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add monthly job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddMonthlyJobs adds list new Monthly job.
func (s *Scheduler) AddMonthlyJobs(jobs ...*MonthJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, monthlyJob := range jobs {
		dailyJobInstance, err := s.AddMonthlyJob(monthlyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("failed to add monthly jobs: %v", errs)
	}
	return jobList, nil
}

// AddMonthlyJobWithOptions add a monthly job, support native extension method
func (s *Scheduler) AddMonthlyJobWithOptions(interval uint, daysOfTheMonth gocron.DaysOfTheMonth, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.MonthlyJob(interval, daysOfTheMonth, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add monthly job: %w", err)
	}
	return job, nil
}
