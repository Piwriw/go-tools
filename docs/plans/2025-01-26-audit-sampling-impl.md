# GORM Audit: 采样和降级实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标:** 为 GORM Audit 插件添加可配置的采样和降级功能，在系统负载高时自动降低审计开销。

**架构:** 通过 SamplingStrategy 接口实现可扩展的采样策略，通过 DegradationController 实现基于队列深度和 CPU 的智能降级，集成到 Dispatcher 的事件处理流程中。

**技术栈:** Go 1.21+, sync/atomic, context, GORM

---

## 文件结构

```
gorm-audit/
├── sampling.go           # 采样策略接口和实现
├── sampling_test.go      # 采样测试
├── degradation.go        # 降级控制器
├── degradation_test.go   # 降级测试
├── config.go             # 添加采样和降级配置（修改）
├── dispatcher.go         # 集成采样和降级（修改）
└── metrics.go            # 添加采样相关指标（修改）
```

---

### Task 1: 添加采样和降级配置到 config.go

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/config.go`

**Step 1: 在文件末尾添加采样配置结构**

```go
// StrategyType 采样策略类型
type StrategyType string

const (
	StrategyRandom  StrategyType = "random"  // 固定概率采样
	StrategyUniform StrategyType = "uniform" // 均匀采样
	StrategySmart   StrategyType = "smart"   // 智能采样
	StrategyCustom  StrategyType = "custom"  // 自定义
)

// SamplingConfig 采样配置
type SamplingConfig struct {
	Enabled bool         // 是否启用采样
	Strategy StrategyType // 采样策略类型
	Rate     float64     // 采样率 0.0-1.0

	// 策略特定配置
	WindowSize  int            // uniform: 时间窗口（秒）
	PriorityMap map[string]int // smart: 操作优先级
}

// DegradationLevel 降级级别定义
type DegradationLevel struct {
	Name         string
	TriggerCPU   float64 // 触发 CPU 阈值 (0.0-1.0)
	TriggerQueue int     // 触发队列阈值（绝对值）
	Action       DegradationAction
}

// DegradationAction 降级行为
type DegradationAction struct {
	AuditLevel AuditLevel // 降级后的审计级别
	SampleRate float64    // 降级后的采样率
}

// DegradationConfig 降级配置
type DegradationConfig struct {
	Enabled bool

	// 触发阈值
	CPUThreshold   float64 // CPU 使用率阈值
	QueueThreshold int     // 队列深度阈值

	// 降级级别（按严重程度排序）
	Levels []DegradationLevel

	// 恢复配置
	RecoveryCooldown time.Duration // 冷却期
}

// DefaultSamplingConfig 返回默认采样配置
func DefaultSamplingConfig() *SamplingConfig {
	return &SamplingConfig{
		Enabled:       false,
		Strategy:      StrategyRandom,
		Rate:          1.0,
		WindowSize:    10,
		PriorityMap:   map[string]int{"create": 8, "delete": 10, "update": 5, "query": 1},
	}
}

// DefaultDegradationConfig 返回默认降级配置
func DefaultDegradationConfig(queueSize int) *DegradationConfig {
	return &DegradationConfig{
		Enabled:         false,
		CPUThreshold:    0.8,
		QueueThreshold:  queueSize * 90 / 100,
		Levels:          DefaultDegradationLevels(queueSize),
		RecoveryCooldown: 30 * time.Second,
	}
}

// DefaultDegradationLevels 返回默认降级级别
func DefaultDegradationLevels(queueSize int) []DegradationLevel {
	return []DegradationLevel{
		{
			Name:         "mild",
			TriggerCPU:   0.7,
			TriggerQueue: queueSize * 70 / 100,
			Action: DegradationAction{
				AuditLevel: AuditLevelChangesOnly,
				SampleRate: 0.5,
			},
		},
		{
			Name:         "severe",
			TriggerCPU:   0.85,
			TriggerQueue: queueSize * 85 / 100,
			Action: DegradationAction{
				AuditLevel: AuditLevelNone,
				SampleRate: 0.1,
			},
		},
		{
			Name:         "critical",
			TriggerCPU:   0.95,
			TriggerQueue: queueSize * 95 / 100,
			Action: DegradationAction{
				AuditLevel: AuditLevelNone,
				SampleRate: 0.0,
			},
		},
	}
}
```

**Step 2: 修改 DefaultConfig 函数，添加采样和降级配置**

找到 `DefaultConfig()` 函数，修改为：

```go
// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:         AuditLevelAll,
		IncludeQuery:  true,
		ContextKeys:   DefaultContextKeys(),
		UseWorkerPool: false,
		WorkerConfig:  nil,
		Filters:       nil,
		EnableMetrics: false,
		Sampling:      nil,
		Degradation:   nil,
	}
}
```

**Step 3: 修改 Config 结构体，添加采样和降级字段**

找到 `type Config struct`，添加字段：

```go
type Config struct {
	Level         AuditLevel
	IncludeQuery  bool
	ContextKeys   ContextKeyConfig
	UseWorkerPool bool
	WorkerConfig  *WorkerPoolConfig
	Filters       []Filter // 事件过滤器列表
	EnableMetrics bool     // 是否启用指标收集

	// 采样和降级
	Sampling    *SamplingConfig
	Degradation *DegradationConfig
}
```

**Step 4: 运行测试验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v -run TestDefault`
Expected: PASS (或没有相关测试失败)

**Step 5: 提交**

```bash
git add gorm-audit/config.go
git commit -m "feat(audit): add sampling and degradation config structures"
```

---

### Task 2: 创建采样策略接口和基础结构

**Files:**
- Create: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/sampling.go`

**Step 1: 创建文件并添加采样策略接口**

```go
package audit

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

// SamplingStrategy 采样策略接口
type SamplingStrategy interface {
	// ShouldSample 决定是否采样该事件
	ShouldSample(ctx context.Context, event *handler.Event) bool

	// UpdateRate 动态更新采样率
	UpdateRate(rate float64)

	// GetEffectiveRate 返回有效采样率
	GetEffectiveRate() float64

	// String 返回策略描述
	String() string
}

// Sampler 采样器别名
type Sampler = SamplingStrategy
```

**Step 2: 添加随机采样器实现**

```go
// RandomSampler 随机采样器（固定概率）
type RandomSampler struct {
	rate float64
	rng  *rand.Rand
	mu   sync.RWMutex
}

// NewRandomSampler 创建随机采样器
func NewRandomSampler(rate float64) *RandomSampler {
	return &RandomSampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldSample 实现 SamplingStrategy 接口
func (r *RandomSampler) ShouldSample(ctx context.Context, event *handler.Event) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.rate >= 1.0 {
		return true
	}
	if r.rate <= 0.0 {
		return false
	}
	return r.rng.Float64() < r.rate
}

// UpdateRate 更新采样率
func (r *RandomSampler) UpdateRate(rate float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rate = rate
}

// GetEffectiveRate 返回有效采样率
func (r *RandomSampler) GetEffectiveRate() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rate
}

// String 返回策略描述
func (r *RandomSampler) String() string {
	return "random"
}
```

**Step 3: 运行测试验证编译**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 成功编译，无错误

**Step 4: 提交**

```bash
git add gorm-audit/sampling.go
git commit -m "feat(audit): add SamplingStrategy interface and RandomSampler"
```

---

### Task 3: 添加采样测试

**Files:**
- Create: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/sampling_test.go`

**Step 1: 创建测试文件**

```go
package audit

import (
	"context"
	"testing"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewRandomSampler(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	assert.NotNil(t, sampler)
	assert.Equal(t, 0.5, sampler.GetEffectiveRate())
	assert.Equal(t, "random", sampler.String())
}

func TestRandomSamplerShouldSample(t *testing.T) {
	tests := []struct {
		name     string
		rate     float64
		minCount int
		maxCount int
	}{
		{"always sample", 1.0, 1000, 1000},
		{"never sample", 0.0, 0, 0},
		{"50% sample", 0.5, 400, 600},
		{"10% sample", 0.1, 50, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := NewRandomSampler(tt.rate)
			event := &handler.Event{}

			count := 0
			trials := 1000
			for i := 0; i < trials; i++ {
				if sampler.ShouldSample(context.Background(), event) {
					count++
				}
			}

			assert.GreaterOrEqual(t, count, tt.minCount)
			assert.LessOrEqual(t, count, tt.maxCount)
		})
	}
}

func TestRandomSamplerUpdateRate(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	assert.Equal(t, 0.5, sampler.GetEffectiveRate())

	sampler.UpdateRate(0.1)
	assert.Equal(t, 0.1, sampler.GetEffectiveRate())
}

func TestRandomSamplerThreadSafe(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	event := &handler.Event{}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				sampler.ShouldSample(context.Background(), event)
				sampler.UpdateRate(0.3)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证没有 panic 且采样率在合理范围内
	rate := sampler.GetEffectiveRate()
	assert.GreaterOrEqual(t, rate, 0.0)
	assert.LessOrEqual(t, rate, 1.0)
}
```

**Step 2: 运行测试**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v -run TestRandom`
Expected: 全部 PASS

**Step 3: 提交**

```bash
git add gorm-audit/sampling_test.go
git commit -m "test(audit): add RandomSampler tests"
```

---

### Task 4: 创建降级控制器

**Files:**
- Create: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/degradation.go`

**Step 1: 创建降级控制器文件**

```go
package audit

import (
	"context"
	"fmt"
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
	cpu := d.getCPUUsage()
	queue := d.queueChecker.GetQueueDepth()

	d.mu.Lock()
	defer d.mu.Unlock()

	currentLevel := atomic.LoadInt32(&d.currentLevel)

	// 检查是否需要降级
	for i := len(d.config.Levels) - 1; i >= 0; i-- {
		level := d.config.Levels[i]
		if cpu >= level.TriggerCPU || queue >= level.TriggerQueue {
			if int(currentLevel) != i {
				d.setLevel(i)
			}
			return
		}
	}

	// 检查是否可以恢复
	if currentLevel > 0 {
		lastTime := d.lastDegraded.Load().(time.Time)
		if time.Since(lastTime) > d.config.RecoveryCooldown {
			d.setLevel(0)
		}
	}
}

// setLevel 设置降级级别
func (d *DegradationController) setLevel(level int) {
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
```

**Step 2: 修改 DegradationConfig 添加对 SamplingConfig 的引用**

在 config.go 中修改 DegradationConfig：

```go
// DegradationConfig 降级配置
type DegradationConfig struct {
	Enabled bool

	// 触发阈值
	CPUThreshold   float64 // CPU 使用率阈值
	QueueThreshold int     // 队列深度阈值

	// 降级级别（按严重程度排序）
	Levels []DegradationLevel

	// 恢复配置
	RecoveryCooldown time.Duration // 冷却期

	// 采样配置引用（用于恢复时重置采样率）
	Sampling *SamplingConfig
}
```

**Step 3: 编译验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 成功编译

**Step 4: 提交**

```bash
git add gorm-audit/degradation.go
git commit -m "feat(audit): add DegradationController"
```

---

### Task 5: 集成采样和降级到 Dispatcher

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/dispatcher.go`

**Step 1: 在 Dispatcher 结构体中添加采样和降级字段**

找到 `type Dispatcher struct`，添加字段：

```go
type Dispatcher struct {
	handlers        []EventHandler
	handlerHandlers []handler.EventHandler
	mu              sync.RWMutex
	workerPool      *WorkerPool
	useWorkerPool   bool
	metrics         *MetricsCollector

	// 采样和降级
	sampler              Sampler
	degradationController *DegradationController
}
```

**Step 2: 修改 NewDispatcher 函数**

```go
// NewDispatcher 创建新的分发器
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:        make([]EventHandler, 0),
		handlerHandlers: make([]handler.EventHandler, 0),
		useWorkerPool:   false,
		metrics:         NewMetricsCollector(),
	}
}
```

**Step 3: 修改 NewDispatcherWithWorkerPool 函数**

```go
// NewDispatcherWithWorkerPool 创建带 Worker Pool 的分发器
func NewDispatcherWithWorkerPool(h handler.EventHandler, config *WorkerPoolConfig) *Dispatcher {
	return &Dispatcher{
		handlers:        make([]EventHandler, 0),
		handlerHandlers: make([]handler.EventHandler, 0),
		workerPool:      NewWorkerPool(h, config),
		useWorkerPool:   true,
		metrics:         NewMetricsCollector(),
	}
}
```

**Step 4: 添加 SetupSamplingAndDegradation 方法**

在 dispatcher.go 中添加：

```go
// SetupSamplingAndDegradation 设置采样和降级
func (d *Dispatcher) SetupSamplingAndDegradation(
	samplingConfig *SamplingConfig,
	degradationConfig *DegradationConfig,
) {
	if samplingConfig != nil && samplingConfig.Enabled {
		switch samplingConfig.Strategy {
		case StrategyRandom, "":
			d.sampler = NewRandomSampler(samplingConfig.Rate)
		default:
			d.sampler = NewRandomSampler(samplingConfig.Rate)
		}
	}

	if degradationConfig != nil && degradationConfig.Enabled && d.workerPool != nil {
		// 设置采样配置引用
		degradationConfig.Sampling = samplingConfig

		checker := &workerPoolQueueChecker{wp: d.workerPool}
		d.degradationController = NewDegradationController(
			degradationConfig,
			d.sampler,
			checker,
		)

		// 启动降级控制器
		go d.degradationController.Start(context.Background())
	}
}

// workerPoolQueueChecker WorkerPool 队列检查器
type workerPoolQueueChecker struct {
	wp *WorkerPool
}

func (w *workerPoolQueueChecker) GetQueueDepth() int {
	if w.wp == nil {
		return 0
	}
	// 假设 WorkerPool 有 GetQueueSize 方法
	return w.wp.GetQueueSize()
}

func (w *workerPoolQueueChecker) GetQueueCapacity() int {
	if w.wp == nil {
		return 1000
	}
	return w.wp.GetQueueCapacity()
}
```

**Step 5: 修改 DispatchHandler 方法，集成采样和降级**

找到 `func (d *Dispatcher) DispatchHandler`，在方法开头添加采样和降级检查：

```go
// DispatchHandler 分发 handler.Event 到 handler 包的处理器
func (d *Dispatcher) DispatchHandler(ctx context.Context, event *handler.Event) {
	d.mu.RLock()
	useWorkerPool := d.useWorkerPool
	wp := d.workerPool
	handlers := make([]handler.EventHandler, len(d.handlerHandlers))
	copy(handlers, d.handlerHandlers)

	// 采样和降级检查
	sampler := d.sampler
	degradation := d.degradationController
	d.mu.RUnlock()

	// 1. 检查降级状态
	if degradation != nil {
		if degradation.ShouldSkip(event) {
			return // 降级，跳过此事件
		}
	}

	// 2. 采样判断
	if sampler != nil {
		if !sampler.ShouldSample(ctx, event) {
			return // 未被采样，跳过
		}
	}

	// 如果启用了 Worker Pool，直接分发到 Worker Pool
	if useWorkerPool && wp != nil {
		wp.Dispatch(event)
		return
	}

	// 记录开始时间
	start := time.Now()

	// 使用 goroutine 分发
	for _, h := range handlers {
		if h != nil {
			go func(handler handler.EventHandler) {
				status := "success"
				if err := handler.Handle(ctx, event); err != nil {
					status = "error"
				}

				// 记录指标（如果启用了指标收集）
				if d.metrics != nil {
					latency := time.Since(start).Seconds()
					d.metrics.RecordEvent(event.Table, string(event.Operation), status, latency)
				}
			}(h)
		}
}
```

**Step 6: 修改 Close 方法，停止降级控制器**

```go
// Close 关闭分发器及 Worker Pool
func (d *Dispatcher) Close() {
	d.mu.Lock()
	wp := d.workerPool
	d.workerPool = nil
	d.useWorkerPool = false

	// 停止降级控制器
	if d.degradationController != nil {
		d.degradationController.Stop()
		d.degradationController = nil
	}
	d.mu.Unlock()

	if wp != nil {
		wp.Close()
	}
}
```

**Step 7: 编译验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 成功编译

**Step 8: 提交**

```bash
git add gorm-audit/dispatcher.go
git commit -m "feat(audit): integrate sampling and degradation into Dispatcher"
```

---

### Task 6: 在 Audit 中初始化采样和降级

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/audit.go`

**Step 1: 修改 New 函数，初始化采样和降级**

找到 `func New(config *Config) *Audit`，在返回前添加：

```go
	// 初始化采样和降级
	if dispatcher != nil {
		dispatcher.SetupSamplingAndDegradation(config.Sampling, config.Degradation)
	}

	return &Audit{
		config:     config,
		dispatcher: dispatcher,
	}
```

**Step 2: 编译验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 成功编译

**Step 3: 提交**

```bash
git add gorm-audit/audit.go
git commit -m "feat(audit): initialize sampling and degradation in Audit"
```

---

### Task 7: 添加降级测试

**Files:**
- Create: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/degradation_test.go`

**Step 1: 创建降级测试文件**

```go
package audit

import (
	"context"
	"testing"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/stretchr/testify/assert"
)

// mockQueueChecker 模拟队列检查器
type mockQueueChecker struct {
	depth     int
	capacity  int
}

func (m *mockQueueChecker) GetQueueDepth() int {
	return m.depth
}

func (m *mockQueueChecker) GetQueueCapacity() int {
	return m.capacity
}

func TestNewDegradationController(t *testing.T) {
	config := &DegradationConfig{
		Enabled:         true,
		Levels:          DefaultDegradationLevels(1000),
		RecoveryCooldown: 100 * time.Millisecond,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)
	assert.NotNil(t, controller)
	assert.Equal(t, 0, controller.GetCurrentLevel())
}

func TestDegradationControllerShouldSkip(t *testing.T) {
	config := &DegradationConfig{
		Enabled: true,
		Levels: []DegradationLevel{
			{
				Name:       "test",
				TriggerCPU: 0.5,
				TriggerQueue: 500,
				Action: DegradationAction{
					AuditLevel: AuditLevelNone,
					SampleRate: 0.0,
				},
			},
		},
		RecoveryCooldown: 30 * time.Second,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)

	// 正常状态，不跳过
	queryEvent := &handler.Event{Operation: "query"}
	createEvent := &handler.Event{Operation: "create"}
	assert.False(t, controller.ShouldSkip(queryEvent))
	assert.False(t, controller.ShouldSkip(createEvent))

	// 设置降级级别
	controller.setLevel(1)

	// 跳过所有事件
	assert.True(t, controller.ShouldSkip(queryEvent))
	assert.True(t, controller.ShouldSkip(createEvent))
}

func TestDegradationControllerSetLevel(t *testing.T) {
	config := &DegradationConfig{
		Enabled: true,
		Levels: []DegradationLevel{
			{
				Action: DegradationAction{SampleRate: 0.5},
			},
			{
				Action: DegradationAction{SampleRate: 0.1},
			},
		},
		RecoveryCooldown: 30 * time.Second,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)

	// 设置降级级别
	controller.setLevel(1)
	assert.Equal(t, 1, controller.GetCurrentLevel())
	assert.InDelta(t, 0.1, sampler.GetEffectiveRate(), 0.01)

	// 恢复正常
	controller.setLevel(0)
	assert.Equal(t, 0, controller.GetCurrentLevel())
}
```

**Step 2: 运行测试**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go test -v -run TestDegradation`
Expected: 全部 PASS

**Step 3: 提交**

```bash
git add gorm-audit/degradation_test.go
git commit -m "test(audit): add DegradationController tests"
```

---

### Task 8: 更新 WorkerPool 添加队列查询方法

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/worker_pool.go`

**Step 1: 在 WorkerPool 结构体中添加方法**

找到 worker_pool.go 文件，添加以下方法：

```go
// GetQueueSize 返回当前队列大小
func (w *WorkerPool) GetQueueSize() int {
	if w.queue == nil {
		return 0
	}
	return w.queue.Len()
}

// GetQueueCapacity 返回队列容量
func (w *WorkerPool) GetQueueCapacity() int {
	if w.config == nil {
		return 1000
	}
	return w.config.QueueSize
}
```

**Step 2: 编译验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit && go build`
Expected: 成功编译

**Step 3: 提交**

```bash
git add gorm-audit/worker_pool.go
git commit -m "feat(audit): add GetQueueSize and GetQueueCapacity to WorkerPool"
```

---

### Task 9: 更新示例代码

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/example/main.go`

**Step 1: 添加采样和降级配置**

找到 `audit.New` 调用，修改配置：

```go
	// 初始化审计插件
	auditPlugin := audit.New(&audit.Config{
		Level:      audit.AuditLevelAll,
		EnableMetrics: true,

		// 采样配置
		Sampling: &audit.SamplingConfig{
			Enabled:  true,
			Strategy: audit.StrategyRandom,
			Rate:     1.0, // 初始 100% 采样
		},

		// 降级配置（仅在启用 Worker Pool 时生效）
		Degradation: nil, // 示例中不启用 Worker Pool，所以设为 nil
	})
```

**Step 2: 编译验证**

Run: `cd /Users/joohwan/GolandProjects/go-tools/gorm-audit/example && go build`
Expected: 成功编译

**Step 3: 提交**

```bash
git add gorm-audit/example/main.go
git commit -m "example(audit): add sampling configuration"
```

---

### Task 10: 更新 README 文档

**Files:**
- Modify: `/Users/joohwan/GolandProjects/go-tools/gorm-audit/README.md`

**Step 1: 在 README 中添加采样和降级章节**

在文件末尾添加：

```markdown
## 采样和降级

### 采样配置

在高并发场景下，可以通过采样降低审计开销：

\`\`\`go
audit := audit.New(&audit.Config{
    Sampling: &audit.SamplingConfig{
        Enabled:  true,
        Strategy: audit.StrategyRandom,
        Rate:     0.1, // 10% 采样
    },
})
\`\`\`

**采样策略类型：**

- `StrategyRandom`: 随机采样（默认）
- `StrategyUniform`: 均匀采样（计划中）
- `StrategySmart`: 智能采样（计划中）

### 降级配置

当系统负载高时自动降低审计级别（需要启用 Worker Pool）：

\`\`\`go
audit := audit.New(&audit.Config{
    UseWorkerPool: true,
    WorkerConfig:  audit.DefaultWorkerPoolConfig(),
    Degradation: &audit.DegradationConfig{
        Enabled:         true,
        CPUThreshold:    0.8,  // CPU 80%
        QueueThreshold:  900,  // 队列 90%
        Levels:          audit.DefaultDegradationLevels(1000),
        RecoveryCooldown: 30 * time.Second,
    },
})
\`\`\`

**默认降级级别：**

| 级别 | CPU 阈值 | 队列阈值 | 审计级别 | 采样率 |
|------|----------|----------|----------|--------|
| mild | 70% | 70% | ChangesOnly | 50% |
| severe | 85% | 85% | None | 10% |
| critical | 95% | 95% | None | 0% |

### 自定义采样策略

实现 `SamplingStrategy` 接口：

\`\`\`go
type MySampler struct {
    blacklist map[string]bool
}

func (m *MySampler) ShouldSample(ctx context.Context, event *handler.Event) bool {
    if m.blacklist[event.Table] {
        return true // 黑名单表总是记录
    }
    return rand.Float64() < 0.1
}

// ... 实现其他接口方法

audit := audit.New(&audit.Config{
    Sampling: &audit.SamplingConfig{
        Enabled:        true,
        CustomStrategy: &MySampler{...},
    },
})
\`\`\`
```

**Step 2: 提交**

```bash
git add gorm-audit/README.md
git commit -m "docs(audit): add sampling and degradation documentation"
```

---

## 验证步骤

### 运行所有测试

```bash
cd /Users/joohwan/GolandProjects/go-tools/gorm-audit
go test -v ./...
```

### 编译示例

```bash
cd /Users/joohwan/GolandProjects/go-tools/gorm-audit/example
go build -o main
```

### 运行示例

```bash
./main
```

---

## 任务清单

- [x] Task 1: 添加采样和降级配置到 config.go
- [x] Task 2: 创建采样策略接口和基础结构
- [x] Task 3: 添加采样测试
- [x] Task 4: 创建降级控制器
- [x] Task 5: 集成采样和降级到 Dispatcher
- [x] Task 6: 在 Audit 中初始化采样和降级
- [x] Task 7: 添加降级测试
- [x] Task 8: 更新 WorkerPool 添加队列查询方法
- [x] Task 9: 更新示例代码
- [x] Task 10: 更新 README 文档
