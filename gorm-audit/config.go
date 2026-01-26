package audit

import "time"

// AuditLevel 审计级别
type AuditLevel int

const (
	AuditLevelAll         AuditLevel = iota // 记录所有操作
	AuditLevelChangesOnly                   // 仅记录变更操作
	AuditLevelNone                          // 不记录
)

// String 实现 Stringer 接口
func (a AuditLevel) String() string {
	switch a {
	case AuditLevelAll:
		return "all"
	case AuditLevelChangesOnly:
		return "changes_only"
	case AuditLevelNone:
		return "none"
	default:
		return "unknown"
	}
}

// ContextKeyConfig Context 键配置
type ContextKeyConfig struct {
	UserID    any
	Username  any
	IP        any
	UserAgent any
	RequestID any
}

// DefaultContextKeys 返回默认的 context key 配置
func DefaultContextKeys() ContextKeyConfig {
	type contextKey string
	return ContextKeyConfig{
		UserID:    contextKey("user_id"),
		Username:  contextKey("username"),
		IP:        contextKey("ip"),
		UserAgent: contextKey("user_agent"),
		RequestID: contextKey("request_id"),
	}
}

// WorkerPoolConfig worker pool 配置
type WorkerPoolConfig struct {
	WorkerCount int // worker 数量
	QueueSize   int // 队列大小
	Timeout     int // 处理超时（毫秒）

	// 批量处理配置
	EnableBatch   bool          // 是否启用批量处理
	BatchSize     int           // 批量大小，默认 1000
	FlushInterval time.Duration // 刷新间隔，默认 10秒
}

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

// DefaultWorkerPoolConfig 返回默认的 worker pool 配置
func DefaultWorkerPoolConfig() *WorkerPoolConfig {
	return &WorkerPoolConfig{
		WorkerCount:   10,
		QueueSize:     1000,
		Timeout:       5000, // 5秒
		EnableBatch:   false,
		BatchSize:     1000,
		FlushInterval: 10 * time.Second,
	}
}

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

	// 采样配置引用（用于恢复时重置采样率）
	Sampling *SamplingConfig
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
