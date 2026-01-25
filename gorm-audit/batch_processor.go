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
