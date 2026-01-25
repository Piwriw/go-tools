package audit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

// WorkerPool 工作池，用于异步处理审计事件
type WorkerPool struct {
	queue     chan *handler.Event
	workers   int
	timeout   time.Duration
	wg        sync.WaitGroup
	stopChan   chan struct{}
	handler   handler.EventHandler
	closed    bool
	closedMu  sync.Mutex

	// 新增：批量处理器
	batchProcessor *BatchProcessor
	enableBatch    bool
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(h handler.EventHandler, config *WorkerPoolConfig) *WorkerPool {
	if config == nil {
		config = DefaultWorkerPoolConfig()
	}

	wp := &WorkerPool{
		queue:       make(chan *handler.Event, config.QueueSize),
		workers:     config.WorkerCount,
		timeout:     time.Duration(config.Timeout) * time.Millisecond,
		stopChan:    make(chan struct{}),
		handler:     h,
		enableBatch: config.EnableBatch,
	}

	// 如果启用批量处理，创建 BatchProcessor
	if config.EnableBatch {
		wp.batchProcessor = NewBatchProcessor(h, config)
		wp.batchProcessor.Start()
	} else {
		// 启动 workers
		for i := 0; i < wp.workers; i++ {
			wp.wg.Add(1)
			go wp.worker(i)
		}
	}

	return wp
}

// worker 工作协程
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case event, ok := <-wp.queue:
			if !ok {
				return
			}
			wp.processEventWithRecovery(event)
		case <-wp.stopChan:
			return
		}
	}
}

// processEvent 处理事件
func (wp *WorkerPool) processEvent(event *handler.Event) {
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), wp.timeout)
	defer cancel()

	// 调用处理器
	if err := wp.handler.Handle(ctx, event); err != nil {
		// 记录错误，但不中断 worker
		fmt.Printf("[WorkerPool] worker error: %v\n", err)
	}
}

// processEventWithRecovery 处理事件，带 panic 恢复
func (wp *WorkerPool) processEventWithRecovery(event *handler.Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[WorkerPool] panic recovered in worker: %v\n", r)
		}
	}()

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), wp.timeout)
	defer cancel()

	// 调用处理器
	if err := wp.handler.Handle(ctx, event); err != nil {
		// 记录错误，但不中断 worker
		fmt.Printf("[WorkerPool] worker error: %v\n", err)
	}
}

// Dispatch 分发事件到工作池
func (wp *WorkerPool) Dispatch(event *handler.Event) bool {
	// 如果启用批量处理，使用 BatchProcessor
	if wp.enableBatch && wp.batchProcessor != nil {
		return wp.batchProcessor.Dispatch(event)
	}

	// 否则使用原有的队列机制
	select {
	case wp.queue <- event:
		return true
	default:
		// 队列已满，丢弃事件
		return false
	}
}

// Close 关闭工作池，等待所有 worker 完成
func (wp *WorkerPool) Close() {
	wp.closedMu.Lock()
	defer wp.closedMu.Unlock()

	if wp.closed {
		return
	}
	wp.closed = true

	// 关闭批量处理器或 workers
	if wp.enableBatch && wp.batchProcessor != nil {
		wp.batchProcessor.Close()
	} else {
		close(wp.stopChan)
		close(wp.queue)
		wp.wg.Wait()
	}

	wp.workers = 0
}

// Stats 获取工作池统计信息
func (wp *WorkerPool) Stats() WorkerPoolStats {
	if wp.enableBatch && wp.batchProcessor != nil {
		batchStats := wp.batchProcessor.Stats()
		return WorkerPoolStats{
			QueueLength:   batchStats.BufferSize,
			QueueCapacity: wp.batchProcessor.batchSize,
			WorkerCount:   1, // BatchProcessor 算作 1 个 worker
		}
	}

	return WorkerPoolStats{
		QueueLength:   len(wp.queue),
		QueueCapacity: cap(wp.queue),
		WorkerCount:   wp.workers,
	}
}

// WorkerPoolStats 工作池统计信息
type WorkerPoolStats struct {
	QueueLength  int
	QueueCapacity int
	WorkerCount  int
}
