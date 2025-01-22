package chrono

import (
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type SchedulerMonitor interface {
	gocron.MonitorStatus
	Watch()
}
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

func (s *defaultSchedulerMonitor) Watch() {
}
