package audit

import (
	"strings"
	"testing"
)

func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()
	if collector == nil {
		t.Fatal("NewMetricsCollector() returned nil")
	}
}

func TestMetricsCollectorRecordEvent(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some events
	collector.RecordEvent("users", "create", "success", 0.001)
	collector.RecordEvent("users", "update", "success", 0.005)
	collector.RecordEvent("orders", "create", "error", 0.1)

	// Check counters
	if collector.GetTotalEvents() != 3 {
		t.Errorf("expected 3 total events, got %d", collector.GetTotalEvents())
	}
}

func TestMetricsCollectorGauges(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetQueueSize(100)
	if collector.GetQueueSize() != 100 {
		t.Errorf("expected queue size 100, got %d", collector.GetQueueSize())
	}

	collector.SetBufferSize(500)
	if collector.GetBufferSize() != 500 {
		t.Errorf("expected buffer size 500, got %d", collector.GetBufferSize())
	}
}

func TestMetricsCollectorString(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordEvent("users", "create", "success", 0.001)
	collector.SetQueueSize(10)

	output := collector.String()
	if output == "" {
		t.Fatal("String() returned empty string")
	}

	// Check for Prometheus format markers
	if !strings.Contains(output, "# HELP") {
		t.Error("missing HELP comments")
	}
	if !strings.Contains(output, "# TYPE") {
		t.Error("missing TYPE comments")
	}
	if !strings.Contains(output, "gorm_audit_events_total") {
		t.Error("missing gorm_audit_events_total metric")
	}
	if !strings.Contains(output, "gorm_audit_events_duration_seconds") {
		t.Error("missing gorm_audit_events_duration_seconds histogram")
	}
}

func TestLatencyHistogramObserve(t *testing.T) {
	histogram := NewLatencyHistogram()

	// Record some observations
	histogram.Observe(0.001) // 1ms
	histogram.Observe(0.005) // 5ms
	histogram.Observe(0.1)   // 100ms

	histogram.mu.Lock()
	if histogram.count != 3 {
		t.Errorf("expected 3 observations, got %d", histogram.count)
	}
	// Use approximate comparison for floating point
	if histogram.sum < 0.105 || histogram.sum > 0.107 {
		t.Errorf("expected sum ~0.106, got %g", histogram.sum)
	}
	histogram.mu.Unlock()
}
