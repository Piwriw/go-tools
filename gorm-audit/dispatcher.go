package audit

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

const defaultQueueCapacity = 1000

// EventHandler 事件处理器接口（audit 包内部使用）
type EventHandler interface {
	Handle(ctx context.Context, event *AuditEvent) error
}

// Dispatcher 事件分发器
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

// NewDispatcher 创建新的分发器
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:        make([]EventHandler, 0),
		handlerHandlers: make([]handler.EventHandler, 0),
		useWorkerPool:   false,
		metrics:         NewMetricsCollector(),
	}
}

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

// Add 添加事件处理器（内部接口）
func (d *Dispatcher) Add(h EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, h)
}

// AddHandler 添加 handler.EventHandler
func (d *Dispatcher) AddHandler(h handler.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlerHandlers = append(d.handlerHandlers, h)
}

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
	mu sync.RWMutex
}

func (w *workerPoolQueueChecker) GetQueueDepth() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.wp == nil {
		return 0
	}
	return w.wp.GetQueueSize()
}

func (w *workerPoolQueueChecker) GetQueueCapacity() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.wp == nil {
		return defaultQueueCapacity
	}
	return w.wp.GetQueueCapacity()
}

// Dispatch 分发 AuditEvent 到内部处理器
func (d *Dispatcher) Dispatch(ctx context.Context, event *AuditEvent) {
	d.mu.RLock()
	handlers := make([]EventHandler, len(d.handlers))
	copy(handlers, d.handlers)
	d.mu.RUnlock()

	for _, h := range handlers {
		if h != nil {
			go d.safeHandle(ctx, h, event)
		}
	}
}

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
}

// safeHandle 安全执行事件处理器，带 panic 恢复
func (d *Dispatcher) safeHandle(ctx context.Context, h EventHandler, event *AuditEvent) {
	defer func() {
		if r := recover(); r != nil {
			d.handlePanic(r, h, event.Table, string(event.Operation))
		}
	}()

	_ = h.Handle(ctx, event)
}

// safeHandleHandler 安全执行 handler.EventHandler，带 panic 恢复
func (d *Dispatcher) safeHandleHandler(ctx context.Context, h handler.EventHandler, event *handler.Event) {
	defer func() {
		if r := recover(); r != nil {
			d.handlePanic(r, h, event.Table, string(event.Operation))
		}
	}()

	_ = h.Handle(ctx, event)
}

// handlePanic 处理 panic 情况
func (d *Dispatcher) handlePanic(r any, h any, table, operation string) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	msg := fmt.Sprintf(
		"[AUDIT] panic recovered: %v, handler: %T, table: %s, operation: %s\nstack: %s",
		r, h, table, operation, stack,
	)

	log.Println(msg)
}

// Close 关闭分发器及 Worker Pool
func (d *Dispatcher) Close() {
	d.mu.Lock()
	wp := d.workerPool
	d.workerPool = nil
	d.useWorkerPool = false
	d.mu.Unlock()

	// 停止降级控制器
	if d.degradationController != nil {
		d.degradationController.Stop()
		d.degradationController = nil
	}

	if wp != nil {
		wp.Close()
	}
}

// Metrics 返回 Prometheus 格式的指标
func (d *Dispatcher) Metrics() string {
	if d.metrics != nil {
		return d.metrics.String()
	}
	return ""
}
