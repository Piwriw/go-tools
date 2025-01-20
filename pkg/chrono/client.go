package chrono

import (
	"time"
)

// Scheduler 调度器接口
// type Scheduler interface {
// 	Start()                                               // 启动调度器
// 	Stop() error                                          // 停止调度器
// 	AddCronJob(cronExpr string, task func()) (Job, error) // 添加 Cron 任务
// 	// AddCronJobWithName(cronExpr string, task func(), name string) (Job, error)             // 添加带名称的 Cron 任务
// 	// AddCronJobWithOptions(cronExpr string, task func(), options ...JobOption) (Job, error) // 添加带选项的 Cron 任务
// 	// AddOnceJob(task func(), times ...time.Time) (Job error)                                // 添加一次性任务
// 	// RemoveJob(job Job) error                                                               // 移除任务
// 	// GetJobs() ([]Job, error)                                                               // 获取所有任务
// 	// GetJobByName(name string) (Job, error)                                                 // 根据名称获取任务
// 	// GetJobByID(id string) (Job, error)                                                     // 根据 ID 获取任务
// }

type CronScheduler interface {
	Start()                                                                                // 启动调度器
	Stop() error                                                                           // 停止调度器
	AddCronJob(cronExpr string, task func()) (Job, error)                                  // 添加 Cron 任务
	AddCronJobWithName(cronExpr string, task func(), name string) (Job, error)             // 添加带名称的 Cron 任务
	AddCronJobWithOptions(cronExpr string, task func(), options ...JobOption) (Job, error) // 添加带选项的 Cron 任务
	GetJobs() ([]Job, error)                                                               // 获取所有任务
	GetJobByName(name string) (Job, error)                                                 // 根据名称获取任务
	GetJobByID(id string) (Job, error)                                                     // 根据 ID 获取任务
}

// Job 任务接口
type Job interface {
	ID() string                  // 获取任务 ID
	Name() string                // 获取任务名称
	NextRun() (time.Time, error) // 获取下次执行时间
	LastRun() (time.Time, error) // 获取上次执行时间
}

// JobOption 任务选项
type JobOption interface{}
