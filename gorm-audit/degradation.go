package audit

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

// DegradationController 降级控制器
type DegradationController struct {
	config        *DegradationConfig
	currentLevel  int32          // 当前降级级别索引（使用 atomic）
	lastDegraded  atomic.Value   // time.Time
	sampler       Sampler
	mu            sync.RWMutex
	stopped       atomic.Value   // bool
	cancel        context.CancelFunc

	// 队列检查器
	queueChecker QueueDepthChecker
}

// QueueDepthChecker 队列深度检查接口
type QueueDepthChecker interface {
	GetQueueDepth() int
	GetQueueCapacity() int
}

// NewDegradationController 创建降级控制器
func NewDegradationController(config *DegradationConfig, sampler Sampler, checker QueueDepthChecker) *DegradationController {
	dc := &DegradationController{
		config:       config,
		sampler:      sampler,
		queueChecker: checker,
	}
	dc.lastDegraded.Store(time.Time{})
	dc.stopped.Store(false)
	return dc
}

// Start 启动降级控制器
func (d *DegradationController) Start(ctx context.Context) {
	ctx, d.cancel = context.WithCancel(ctx)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Add panic recovery to prevent goroutine crashes
	defer func() {
		if r := recover(); r != nil {
			// Log the panic and allow graceful shutdown
			log.Printf("[DegradationController] panic recovered: %v", r)
		}
	}()

	for {
		select {
		case <-ticker.C:
			d.evaluate()
		case <-ctx.Done():
			return
		}
	}
}

// Stop 停止降级控制器
func (d *DegradationController) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
	d.stopped.Store(true)
}

// evaluate 评估是否需要降级
func (d *DegradationController) evaluate() {
	// Add nil check for queueChecker (I5)
	queue := 0
	if d.queueChecker != nil {
		queue = d.queueChecker.GetQueueDepth()
	}
	cpu := d.getCPUUsage()

	d.mu.Lock()
	currentLevel := atomic.LoadInt32(&d.currentLevel)

	// Determine the new level (holding lock only for state check)
	var newLevel int = -1
	// 检查是否需要降级
	for i := len(d.config.Levels) - 1; i >= 0; i-- {
		level := d.config.Levels[i]
		if cpu >= level.TriggerCPU || queue >= level.TriggerQueue {
			newLevel = i
			break
		}
	}

	// 检查是否可以恢复
	shouldRecover := false
	if newLevel == -1 && currentLevel > 0 {
		// Safe type assertion with check (I1)
		if lastTime, ok := d.lastDegraded.Load().(time.Time); ok && !lastTime.IsZero() {
			if time.Since(lastTime) > d.config.RecoveryCooldown {
				shouldRecover = true
			}
		}
	}

	d.mu.Unlock() // Release lock before calling setLevel() (C2)

	// Call setLevel() outside the lock to prevent deadlock
	if newLevel != -1 && int(currentLevel) != newLevel {
		d.setLevel(newLevel)
	} else if shouldRecover {
		d.setLevel(0)
	}
}

// setLevel 设置降级级别
func (d *DegradationController) setLevel(level int) {
	// Add bounds checking (I3)
	if level < 0 || level >= len(d.config.Levels) {
		return
	}

	atomic.StoreInt32(&d.currentLevel, int32(level))
	d.lastDegraded.Store(time.Now())

	// 更新采样器
	if level > 0 && d.sampler != nil {
		action := d.config.Levels[level].Action
		d.sampler.UpdateRate(action.SampleRate)
	} else if d.sampler != nil {
		// 恢复正常采样率
		if d.config.Sampling != nil {
			d.sampler.UpdateRate(d.config.Sampling.Rate)
		}
	}
}

// ShouldSkip 判断是否应该跳过该事件
func (d *DegradationController) ShouldSkip(event *handler.Event) bool {
	level := atomic.LoadInt32(&d.currentLevel)
	if level == 0 {
		return false
	}

	// Add bounds checking (M5)
	if level < 0 || int(level) >= len(d.config.Levels) {
		return false
	}

	action := d.config.Levels[level].Action

	// 检查审计级别
	switch action.AuditLevel {
	case AuditLevelNone:
		return true
	case AuditLevelChangesOnly:
		return event.Operation == "query"
	case AuditLevelAll:
		return false
	}

	return false
}

// GetCurrentLevel 获取当前降级级别
func (d *DegradationController) GetCurrentLevel() int {
	return int(atomic.LoadInt32(&d.currentLevel))
}

// getCPUUsage 获取 CPU 使用率（简化版）
func (d *DegradationController) getCPUUsage() float64 {
	// 使用 goroutine 数量作为近似指标
	// 生产环境可以使用 gopsutil 或读取 /proc/stat
	numGoroutines := float64(runtime.NumGoroutine())
	// 假设 10000 个 goroutine 为 100% CPU（简化）
	return numGoroutines / 10000.0
}
