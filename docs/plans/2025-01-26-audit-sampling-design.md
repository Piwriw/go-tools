# GORM Audit: 采样和降级设计

## 概述

为 GORM Audit 插件添加采样和降级功能，在系统负载高时自动降低审计开销，防止审计功能影响主业务性能。

## 设计目标

- **可配置**: 采样策略、降级阈值完全可配置
- **可扩展**: 支持用户自定义采样和降级策略
- **低开销**: 采样判断性能开销 < 1 微秒
- **智能降级**: 根据队列深度和 CPU 自动降级
- **平滑恢复**: 支持冷却期，避免频繁切换

## 核心功能

### 1. 采样

按比例记录事件，而非全部记录：

| 采样率 | 记录 | 跳过 |
|--------|------|------|
| 100% | 全部 | 无 |
| 10% | 100条中记10条 | 90条 |
| 1% | 1000条中记10条 | 990条 |

### 2. 降级

根据系统压力自动调整审计级别：

```go
正常情况  → AuditLevelAll      // 记录所有操作
CPU > 70% → AuditLevelChangesOnly // 只记录变更，采样 50%
CPU > 85% → AuditLevelNone        // 暂停审计，采样 10%
CPU > 95% → 完全停止               // 采样 0%
```

## 配置结构

### 完整配置

```go
// 采样配置
type SamplingConfig struct {
    Enabled   bool
    Strategy  StrategyType
    Rate      float64              // 默认采样率 0.0-1.0

    // 策略特定配置
    WindowSize     int             // uniform: 时间窗口（秒）
    PriorityMap    map[string]int  // smart: 操作优先级

    // 自定义策略
    CustomStrategy SamplingStrategy
}

// 降级配置
type DegradationConfig struct {
    Enabled bool

    // 触发阈值
    CPUThreshold    float64  // CPU 使用率阈值 (0.0-1.0)，默认 0.8
    QueueThreshold  int      // 队列深度阈值，默认 90% 队列容量

    // 降级级别（按严重程度排序）
    Levels []DegradationLevel

    // 恢复配置
    RecoveryCooldown time.Duration  // 冷却期，默认 30 秒
}

type DegradationLevel struct {
    Name         string
    TriggerCPU   float64     // 触发 CPU 阈值
    TriggerQueue int         // 触发队列阈值（百分比或绝对值）
    Action       DegradationAction
}

type DegradationAction struct {
    AuditLevel    AuditLevel  // 降级后的审计级别
    SampleRate    float64     // 降级后的采样率（覆盖采样配置）
}

// 集成到 Config
type Config struct {
    Level            AuditLevel
    IncludeQuery     bool
    ContextKeys      ContextKeyConfig
    UseWorkerPool    bool
    WorkerConfig     *WorkerPoolConfig
    Filters          []Filter
    EnableMetrics    bool

    // 新增
    Sampling     *SamplingConfig
    Degradation  *DegradationConfig
}
```

### 策略类型

```go
type StrategyType string

const (
    StrategyRandom     StrategyType = "random"      // 固定概率采样
    StrategyUniform    StrategyType = "uniform"     // 均匀采样
    StrategySmart      StrategyType = "smart"       // 智能采样
    StrategyCustom     StrategyType = "custom"      // 自定义
)
```

## 组件架构

### 采样策略接口

```go
// 采样策略接口（用户可实现）
type SamplingStrategy interface {
    // ShouldSample 决定是否采样该事件
    ShouldSample(ctx context.Context, event *handler.Event) bool

    // UpdateConfig 动态更新配置
    UpdateConfig(config interface{}) error

    // GetEffectiveRate 返回有效采样率
    GetEffectiveRate() float64

    // String 返回策略描述
    String() string
}
```

### 内置采样器

```go
// 随机采样器
type RandomSampler struct {
    rate    float64
    rng     *rand.Rand
    mu      sync.RWMutex
}

// 均匀采样器
type UniformSampler struct {
    rate       float64
    counter    int64
    windowSize int
    mu         sync.Mutex
}

// 智能采样器
type SmartSampler struct {
    baseRate    float64
    priorities  map[string]int  // operation -> priority (1-10)
    rng         *rand.Rand
    mu          sync.RWMutex
}
```

### 降级控制器

```go
type DegradationController struct {
    config       *DegradationConfig
    currentLevel int           // 当前降级级别索引
    lastDegraded time.Time     // 上次降级时间
    cooldown     time.Duration // 冷却期
    mu           sync.RWMutex

    // 依赖
    cpuCollector  *CPUCollector
    queueChecker  QueueDepthChecker
    sampler       Sampler  // 用于动态调整采样率

    // 自定义触发器
    customTrigger DegradationTrigger
}

// CPU 采集器（轻量级）
type CPUCollector struct {
    interval    time.Duration
    lastCPUTime time.Time
    lastCPU     float64
    mu          sync.Mutex
}

func (c *CPUCollector) GetCPUUsage() float64 {
    // 使用 runtime.GetNumGoroutine() + 时间计算近似值
    // 或使用 /proc/stat (Linux)
}
```

### 自定义降级触发器

```go
// 降级判断接口（用户可实现）
type DegradationTrigger interface {
    // ShouldDegrade 判断是否应该降级
    ShouldDegrade(ctx context.Context) (shouldDegrade bool, level int)

    // GetMetrics 获取当前指标
    GetMetrics() map[string]interface{}
}
```

## 数据流

### 完整的事件处理流程

```go
func (d *Dispatcher) DispatchHandler(ctx context.Context, event *handler.Event) {
    // 1. 检查降级状态
    if d.degradationController != nil {
        level := d.degradationController.GetCurrentLevel()
        if level.ShouldSkip(event) {
            return  // 降级，跳过此事件
        }
    }

    // 2. 采样判断
    if d.sampler != nil {
        if !d.sampler.ShouldSample(ctx, event) {
            return  // 未被采样，跳过
        }
    }

    // 3. 记录指标
    start := time.Now()

    // 4. 分发到 handler
    for _, h := range handlers {
        go func(handler handler.EventHandler) {
            status := "success"
            if err := handler.Handle(ctx, event); err != nil {
                status = "error"
            }

            // 5. 记录指标
            if d.metrics != nil {
                latency := time.Since(start).Seconds()
                d.metrics.RecordEvent(...)
            }
        }(h)
    }
}
```

### 降级控制器工作流程

```go
// 定期检查系统状态
func (d *DegradationController) Start(ctx context.Context) {
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

func (d *DegradationController) evaluate() {
    d.mu.Lock()
    defer d.mu.Unlock()

    cpu := d.cpuCollector.GetCPUUsage()
    queue := d.queueChecker.GetDepth()

    // 找到匹配的降级级别
    for i, level := range d.config.Levels {
        if cpu >= level.TriggerCPU || queue >= level.TriggerQueue {
            d.setLevel(i)
            return
        }
    }

    // 检查是否可以恢复
    if d.currentLevel > 0 {
        if time.Since(d.lastDegraded) > d.cooldown {
            d.setLevel(0)  // 恢复正常
        }
    }
}

func (d *DegradationController) setLevel(level int) {
    d.currentLevel = level
    d.lastDegraded = time.Now()

    // 通知采样器更新采样率
    action := d.config.Levels[level].Action
    if d.sampler != nil {
        d.sampler.UpdateRate(action.SampleRate)
    }
}
```

## 默认配置

### 默认降级级别

```go
func DefaultDegradationLevels(workerQueueSize int) []DegradationLevel {
    return []DegradationLevel{
        {
            Name:         "mild",
            TriggerCPU:   0.7,    // CPU 70%
            TriggerQueue: workerQueueSize * 70 / 100,
            Action: DegradationAction{
                AuditLevel: AuditLevelChangesOnly,
                SampleRate: 0.5,
            },
        },
        {
            Name:         "severe",
            TriggerCPU:   0.85,
            TriggerQueue: workerQueueSize * 85 / 100,
            Action: DegradationAction{
                AuditLevel: AuditLevelNone,
                SampleRate: 0.1,
            },
        },
        {
            Name:         "critical",
            TriggerCPU:   0.95,
            TriggerQueue: workerQueueSize * 95 / 100,
            Action: DegradationAction{
                AuditLevel: AuditLevelNone,
                SampleRate: 0.0,
            },
        },
    }
}
```

## 错误处理

### 采样器容错

```go
type SafeSampler struct {
    inner SamplingStrategy
    logger Logger
}

func (s *SafeSampler) ShouldSample(ctx context.Context, event *handler.Event) bool {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Error("sampler panic", "error", r)
            // panic 时默认记录
        }
    }()
    return s.inner.ShouldSample(ctx, event)
}
```

### 降级控制器容错

```go
func (d *DegradationController) evaluate() {
    cpu, err := d.cpuCollector.GetCPUUsage()
    if err != nil {
        // 获取失败时，保持当前状态，不降级
        return
    }
    // ...
}
```

## 指标扩展

新增 Prometheus 指标：

- `gorm_audit_sampling_rate` - 当前采样率（gauge）
- `gorm_audit_degradation_level` - 当前降级级别（gauge，0=正常）
- `gorm_audit_events_sampled_total` - 被采样跳过的事件数（counter）
- `gorm_audit_degradations_total` - 降级次数（counter，按级别标签）

## 使用示例

### 基础使用

```go
audit := audit.New(&audit.Config{
    Sampling: &audit.SamplingConfig{
        Enabled:  true,
        Strategy: audit.StrategyRandom,
        Rate:     0.1,  // 10% 采样
    },
    Degradation: &audit.DegradationConfig{
        Enabled:        true,
        CPUThreshold:   0.8,
        QueueThreshold: 900,  // 假设队列大小 1000
        Levels:         audit.DefaultDegradationLevels(1000),
        RecoveryCooldown: 30 * time.Second,
    },
})
```

### 自定义采样策略

```go
type MySampler struct {
    blacklist map[string]bool
}

func (m *MySampler) ShouldSample(ctx context.Context, event *handler.Event) bool {
    if m.blacklist[event.Table] {
        return true  // 黑名单表总是记录
    }
    return rand.Float64() < 0.1
}

audit := audit.New(&audit.Config{
    Sampling: &audit.SamplingConfig{
        Enabled:        true,
        CustomStrategy: &MySampler{...},
    },
})
```

### 自定义降级触发器

```go
type MyTrigger struct {
    httpClient *http.Client
    endpoint   string
}

func (m *MyTrigger) ShouldDegrade(ctx context.Context) (bool, int) {
    resp, _ := m.httpClient.Get(m.endpoint + "/load")
    // 从外部服务获取负载信息
    return false, 0
}

controller := audit.NewDegradationController(config)
controller.SetCustomTrigger(&MyTrigger{...})
```

## 测试策略

### 单元测试

1. **采样策略测试**
   - 验证随机采样分布符合预期
   - 验证均匀采样时间窗口正确
   - 验证智能采样优先级排序

2. **降级逻辑测试**
   - 验证降级级别切换
   - 验证冷却期恢复逻辑
   - 验证采样率动态调整

3. **自定义策略集成测试**
   - 验证用户自定义采样器
   - 验证用户自定义降级触发器

### 集成测试

1. **与 Dispatcher 集成**
   - 验证采样跳过不影响其他事件
   - 验证降级时事件正确过滤

2. **与 Worker Pool 集成**
   - 验证队列深度指标正确获取
   - 验证降级不影响 Worker Pool 正常运行

3. **压力测试**
   - 高负载下的降级行为
   - 恢复后的行为验证

### 测试用例示例

```go
func TestRandomSampler(t *testing.T) {
    sampler := NewRandomSampler(0.5)
    trials := 10000
    count := 0
    for i := 0; i < trials; i++ {
        if sampler.ShouldSample(context.Background(), &handler.Event{}) {
            count++
        }
    }
    ratio := float64(count) / float64(trials)
    assert.InDelta(t, 0.5, ratio, 0.02)
}

func TestDegradationRecovery(t *testing.T) {
    controller := NewDegradationController(&DegradationConfig{
        RecoveryCooldown: 100 * time.Millisecond,
        Levels: []DegradationLevel{...},
    })

    // 触发降级
    controller.setCurrentLevel(2)

    // 立即检查，未到冷却期
    assert.Equal(t, 2, controller.GetCurrentLevel())

    // 等待冷却期
    time.Sleep(150 * time.Millisecond)

    // CPU 恢复正常
    controller.evaluate()
    assert.Equal(t, 0, controller.GetCurrentLevel())
}
```

## 实现文件

- `sampling.go` - 采样策略接口和内置实现
- `sampling_test.go` - 采样测试
- `degradation.go` - 降级控制器
- `degradation_test.go` - 降级测试
- `config.go` - 添加采样和降级配置
- `dispatcher.go` - 集成采样和降级
- `metrics.go` - 添加采样相关指标

## 优势总结

1. **性能保护**: 防止审计功能拖累主业务
2. **灵活配置**: 支持多种策略，用户可扩展
3. **智能降级**: 自动感知负载并调整
4. **平滑恢复**: 冷却期避免抖动
5. **可观测性**: 新增指标监控降级状态
