package audit

import (
	"sync"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
)

// BatchStats 批量处理器统计信息
type BatchStats struct {
	BufferSize    int       // 当前缓冲区大小
	TotalEvents   int64     // 总事件数
	TotalBatches  int64     // 总批次数
	TotalFlushes  int64     // 总刷新次数
	TotalErrors   int64     // 总错误数
	AvgBatchSize  float64   // 平均批次大小
	LastFlushTime time.Time // 最后刷新时间
	mu            sync.Mutex
}

// BatchProcessor 批量事件处理器
type BatchProcessor struct {
	buffer        chan *handler.Event
	handler       handler.EventHandler
	batchSize     int
	flushInterval time.Duration
	flushTicker   *time.Ticker
	done          chan struct{}
	wg            sync.WaitGroup
	stats         BatchStats
}

// NewBatchProcessor 创建批量处理器
func NewBatchProcessor(h handler.EventHandler, config *WorkerPoolConfig) *BatchProcessor {
	return &BatchProcessor{
		buffer:        make(chan *handler.Event, config.QueueSize),
		handler:       h,
		batchSize:     config.BatchSize,
		flushInterval: config.FlushInterval,
		flushTicker:   time.NewTicker(config.FlushInterval),
		done:          make(chan struct{}),
	}
}

// Start 启动批量处理器
func (bp *BatchProcessor) Start() {
	bp.wg.Add(1)
	go func() {
		defer bp.wg.Done()
		defer bp.flushTicker.Stop()

		batch := make([]*handler.Event, 0, bp.batchSize)

		for {
			select {
			case event := <-bp.buffer:
				batch = append(batch, event)

				// 更新统计
				bp.stats.mu.Lock()
				bp.stats.BufferSize = len(batch)
				bp.stats.TotalEvents++
				bp.stats.mu.Unlock()

				// 数量达到阈值，触发刷新
				if len(batch) >= bp.batchSize {
					bp.flush(batch)
					batch = make([]*handler.Event, 0, bp.batchSize)
				}

			case <-bp.flushTicker.C:
				// 定时刷新
				if len(batch) > 0 {
					bp.flush(batch)
					batch = make([]*handler.Event, 0, bp.batchSize)
				}

			case <-bp.done:
				// 关闭时刷新剩余事件
				if len(batch) > 0 {
					bp.flush(batch)
				}
				return
			}
		}
	}()
}
