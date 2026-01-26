package audit

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

const (
	// MinSampleRate is the minimum valid sampling rate (0%)
	MinSampleRate = 0.0
	// MaxSampleRate is the maximum valid sampling rate (100%)
	MaxSampleRate = 1.0
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

// RandomSampler 随机采样器（固定概率）
type RandomSampler struct {
	rate float64
	rng  *rand.Rand
	mu   sync.RWMutex
}

// NewRandomSampler 创建随机采样器
// The rate parameter is clamped to [0.0, 1.0] range for safety.
// Note: Using math/rand with time-based seeding is acceptable for non-critical
// audit log sampling. For production systems requiring cryptographically secure
// randomness, consider using crypto/rand instead.
func NewRandomSampler(rate float64) *RandomSampler {
	// Clamp rate to valid range [0.0, 1.0]
	if rate < MinSampleRate {
		rate = MinSampleRate
	}
	if rate > MaxSampleRate {
		rate = MaxSampleRate
	}

	return &RandomSampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldSample 实现 SamplingStrategy 接口
func (r *RandomSampler) ShouldSample(ctx context.Context, event *handler.Event) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.rate >= MaxSampleRate {
		return true
	}
	if r.rate <= MinSampleRate {
		return false
	}
	return r.rng.Float64() < r.rate
}

// UpdateRate 更新采样率
func (r *RandomSampler) UpdateRate(rate float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clamp rate to valid range [0.0, 1.0]
	if rate < MinSampleRate {
		rate = MinSampleRate
	}
	if rate > MaxSampleRate {
		rate = MaxSampleRate
	}

	r.rate = rate
}

// GetEffectiveRate 返回有效采样率
func (r *RandomSampler) GetEffectiveRate() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rate
}

// String 返回策略描述
func (r *RandomSampler) String() string {
	return "random"
}
