package audit

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// defaultBuckets 默认的延迟 bucket 边界（秒）
var defaultBuckets = []float64{
	0.001, // 1ms
	0.005, // 5ms
	0.01,  // 10ms
	0.025, // 25ms
	0.05,  // 50ms
	0.1,   // 100ms
	0.25,  // 250ms
	0.5,   // 500ms
	1.0,   // 1s
	2.5,   // 2.5s
	5.0,   // 5s
	10.0,  // 10s
}

// LatencyHistogram 延迟直方图
type LatencyHistogram struct {
	buckets []float64 // bucket 边界
	counts  []int64   // 每个 bucket 的计数
	sum     float64   // 延迟总和
	count   int64     // 总计数
	mu      sync.Mutex
}

// NewLatencyHistogram 创建延迟直方图
func NewLatencyHistogram() *LatencyHistogram {
	return &LatencyHistogram{
		buckets: defaultBuckets,
		counts:  make([]int64, len(defaultBuckets)+1),
	}
}

// Observe 记录一个观测值
func (h *LatencyHistogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += value
	h.count++

	// 找到合适的 bucket
	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
			return
		}
	}
	// 超过所有 bucket，计入最后一个
	h.counts[len(h.counts)-1]++
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	// Counters (使用 atomic)
	totalEvents   int64
	successEvents int64
	failedEvents  int64

	// Gauges
	queueSize  int64
	bufferSize int64

	// Histogram
	latency *LatencyHistogram

	// 按维度分组 (使用 map + mutex)
	byTable     map[string]int64
	byOperation map[string]int64
	byStatus    map[string]int64
	dimMu       sync.Mutex
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		latency:     NewLatencyHistogram(),
		byTable:     make(map[string]int64),
		byOperation: make(map[string]int64),
		byStatus:    make(map[string]int64),
	}
}

// RecordEvent 记录一个事件
func (m *MetricsCollector) RecordEvent(table, operation, status string, latency float64) {
	// 更新计数器
	atomic.AddInt64(&m.totalEvents, 1)

	switch status {
	case "success":
		atomic.AddInt64(&m.successEvents, 1)
	case "error":
		atomic.AddInt64(&m.failedEvents, 1)
	}

	// 更新延迟直方图
	m.latency.Observe(latency)

	// 更新维度统计
	m.dimMu.Lock()
	m.byTable[table]++
	m.byOperation[operation]++
	m.byStatus[status]++
	m.dimMu.Unlock()
}

// SetQueueSize 设置队列大小
func (m *MetricsCollector) SetQueueSize(size int64) {
	atomic.StoreInt64(&m.queueSize, size)
}

// SetBufferSize 设置缓冲区大小
func (m *MetricsCollector) SetBufferSize(size int64) {
	atomic.StoreInt64(&m.bufferSize, size)
}

// GetQueueSize 获取队列大小
func (m *MetricsCollector) GetQueueSize() int64 {
	return atomic.LoadInt64(&m.queueSize)
}

// GetBufferSize 获取缓冲区大小
func (m *MetricsCollector) GetBufferSize() int64 {
	return atomic.LoadInt64(&m.bufferSize)
}

// String 返回 Prometheus 格式的指标
func (m *MetricsCollector) String() string {
	var sb strings.Builder

	// 1. 总事件数 counter
	sb.WriteString("# HELP gorm_audit_events_total Total number of audit events\n")
	sb.WriteString("# TYPE gorm_audit_events_total counter\n")
	total := atomic.LoadInt64(&m.totalEvents)
	sb.WriteString(fmt.Sprintf("gorm_audit_events_total %d\n", total))

	// 2. 成功事件 counter
	success := atomic.LoadInt64(&m.successEvents)
	sb.WriteString(fmt.Sprintf("gorm_audit_events_success_total %d\n", success))

	// 3. 失败事件 counter
	failed := atomic.LoadInt64(&m.failedEvents)
	sb.WriteString(fmt.Sprintf("gorm_audit_events_error_total %d\n", failed))

	// 4. 按表名分组
	m.dimMu.Lock()
	for table, count := range m.byTable {
		sb.WriteString(fmt.Sprintf("gorm_audit_events_total{table=\"%s\"} %d\n", table, count))
	}
	m.dimMu.Unlock()

	// 5. 按操作类型分组
	m.dimMu.Lock()
	for op, count := range m.byOperation {
		sb.WriteString(fmt.Sprintf("gorm_audit_events_total{operation=\"%s\"} %d\n", op, count))
	}
	m.dimMu.Unlock()

	// 6. 按状态分组
	m.dimMu.Lock()
	for status, count := range m.byStatus {
		sb.WriteString(fmt.Sprintf("gorm_audit_events_total{status=\"%s\"} %d\n", status, count))
	}
	m.dimMu.Unlock()

	// 7. 队列大小 gauge
	queueSize := atomic.LoadInt64(&m.queueSize)
	sb.WriteString("# HELP gorm_audit_queue_size Current queue size\n")
	sb.WriteString("# TYPE gorm_audit_queue_size gauge\n")
	sb.WriteString(fmt.Sprintf("gorm_audit_queue_size %d\n", queueSize))

	// 8. 缓冲区大小 gauge
	bufferSize := atomic.LoadInt64(&m.bufferSize)
	sb.WriteString("# HELP gorm_audit_buffer_size Current buffer size\n")
	sb.WriteString("# TYPE gorm_audit_buffer_size gauge\n")
	sb.WriteString(fmt.Sprintf("gorm_audit_buffer_size %d\n", bufferSize))

	// 9. 延迟直方图
	sb.WriteString("# HELP gorm_audit_events_duration_seconds Event processing duration\n")
	sb.WriteString("# TYPE gorm_audit_events_duration_seconds histogram\n")
	m.latency.mu.Lock()
	for i, count := range m.latency.counts {
		if i < len(m.latency.buckets) {
			sb.WriteString(fmt.Sprintf("gorm_audit_events_duration_seconds_bucket{le=\"%g\"} %d\n",
				m.latency.buckets[i], count))
		} else {
			sb.WriteString(fmt.Sprintf("gorm_audit_events_duration_seconds_bucket{le=\"+Inf\"} %d\n", count))
		}
	}
	sb.WriteString(fmt.Sprintf("gorm_audit_events_duration_seconds_sum %g\n", m.latency.sum))
	sb.WriteString(fmt.Sprintf("gorm_audit_events_duration_seconds_count %d\n", m.latency.count))
	m.latency.mu.Unlock()

	return sb.String()
}
