package chrono

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func TestWatchJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	name := "Joohwan"
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Names("TestWatchJob").
		Task(task, 1, 2).Watch(func(event MonitorJobSpec) {
		fmt.Println("watchFunc", event, name)
	})

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestMonthlyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	monthlyJob := NewMonthJob(1, gocron.NewDaysOfTheMonth(23), gocron.NewAtTimes(gocron.NewAtTime(11, 43, 3))).
		Names("TestMonthlyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddMonthlyJob(monthlyJob)
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

func TestWeeklyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	weeklyJob := NewWeeklyJob(1, gocron.NewWeekdays(time.Thursday), gocron.NewAtTimes(gocron.NewAtTime(11, 34, 3))).
		Names("TestWeeklyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddWeeklyJob(weeklyJob)
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
func TestDailyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	dailyJob := NewDailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(19, 40, 3))).
		Names("TestDailyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddDailyJob(dailyJob)
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

func add(a, b int) {
	fmt.Println(a + b)
}

func countdis(a, b int) {
	fmt.Println(a - b)
}

func TestIntervalJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	taskErr := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return errors.New("error")
	}
	// 定义 funcName
	a := 1
	b := 2
	// 定义 watchFunc
	watchFunc := func(event MonitorJobSpec) {
		fmt.Println("watchFunc", event)
		add(a, b)
	}
	watchFunc2 := func(event MonitorJobSpec) {
		fmt.Println("watchFunc2", event)
		countdis(a, b)
	}
	intervalJob := NewIntervalJob(15*time.Second).
		Names("TestDurationJob").
		Watch(watchFunc).
		Task(task, 1, 2)

	job, err := scheduler.AddIntervalJob(intervalJob)
	if err != nil {
		t.Fatal(err)
	}

	intervalJob2 := NewIntervalJob(15*time.Second).
		Names("TestDurationJob2").
		Watch(watchFunc2).
		Task(taskErr, 3, 5)

	_, err = scheduler.AddIntervalJob(intervalJob2)
	if err != nil {
		t.Fatal(err)
	}

	scheduler.Start()
	go scheduler.Watch()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestOnceJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	onceJob := NewOnceJob(time.Now().Add(time.Minute)).
		Names("TestOnceJob").
		Task(task, 1, 2)

	job, err := scheduler.AddOnceJob(onceJob)

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

func TestMonitor(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Names("TestMonitor").
		Task(task, 1, 2)

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestDefaultHooks(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return errors.New("some error")
	}
	cronJob := NewCronJob(
		DayTimeToCron(time.Now().Add(time.Minute*1))).
		Task(task, 1, 2).
		DefaultHooks()

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

func TestDayTimeToCron(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	// Task with a parameter using a closure
	task := func(a, b int) {
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
		fmt.Println("Task executed with parameters:", a, b)

	}
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Task(task, 1, 2).
		AfterJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("AfterJobRuns")
		}).
		BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("BeforeJobRuns")
		}).AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
		fmt.Println("AfterJobRuns")
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
