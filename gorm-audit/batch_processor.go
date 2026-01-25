package audit

import (
	"sync"
	"time"
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
