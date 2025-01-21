package chrono

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/go-co-op/gocron/v2"
)

func TestDayTimeToCron(t *testing.T) {
	scheduler, err := NewScheduler(nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	// Task with a parameter using a closure
	task := func() {
		jobs, err := scheduler.GetJobs()
		if err != nil {
			fmt.Println(err)
		}
		for _, job := range jobs {
			// 获取当前任务信息
			jobID := job.ID()
			jobName := job.Name()
			nextRun, err := job.NextRun()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("TASKID", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
			lastRun, err := job.LastRun()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("TASKID", job.ID(), "TASK NAME", job.Name(), "lastRunTime", lastRun.Format("2006-01-02 15:04:05"))
			// 使用 Monitor 记录任务信息
			scheduler.monitor.IncrementJob(jobID, jobName, nil, gocron.Success)

			// 模拟任务执行
			startTime := time.Now()
			time.Sleep(2 * time.Second) // 模拟任务运行 2 秒
			endTime := time.Now()

			// 记录任务执行时间
			scheduler.monitor.RecordJobTiming(startTime, endTime, jobID, jobName, nil)
		}

	}
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute * 1))).
		Task(task).
		AfterJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("AfterJobRuns")
		}).
		BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("BeforeJobRuns")
		})
	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
