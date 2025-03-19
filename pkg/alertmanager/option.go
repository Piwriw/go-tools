package alertmanager

import (
	"sync"
	"time"
)

type Option func(*AlertStrategyManager)

type hookOption interface {
	Allow() bool
}

type LimitOption interface {
	Allow() bool
}

// TokenBucket 结构体实现令牌桶算法
type TokenBucket struct {
	mu         sync.Mutex // 互斥锁，确保并发安全
	rate       int        // 令牌生成速率（每秒生成多少个令牌）
	capacity   int        // 令牌桶的容量（最多存多少个令牌）
	tokens     int        // 当前桶中的令牌数
	lastRefill time.Time  // 上次补充令牌的时间
}

// NewTokenBucket 创建一个新的令牌桶
// rate: 每秒生成的令牌数
// capacity: 令牌桶的最大容量
func NewTokenBucket(rate, capacity int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // 初始时令牌桶满
		lastRefill: time.Now(),
	}
}

// refill 根据时间流逝计算需要补充的令牌
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds() // 计算自上次补充以来的时间间隔
	tb.lastRefill = now

	// 计算补充的令牌数量
	tb.tokens += int(elapsed * float64(tb.rate))

	// 限制令牌数量不能超过桶的容量
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
}

// Allow 判断是否允许请求通过
// 允许：消耗一个令牌并返回 true
// 拒绝：令牌不足，返回 false
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock() // 加锁，确保并发安全
	defer tb.mu.Unlock()

	tb.refill() // 先补充令牌

	if tb.tokens > 0 {
		tb.tokens-- // 消耗一个令牌
		return true
	}
	return false // 令牌不足，拒绝请求
}
