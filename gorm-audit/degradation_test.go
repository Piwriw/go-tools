package audit

import (
	"testing"
	"time"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/stretchr/testify/assert"
)

// mockQueueChecker 模拟队列检查器
type mockQueueChecker struct {
	depth    int
	capacity int
}

func (m *mockQueueChecker) GetQueueDepth() int {
	return m.depth
}

func (m *mockQueueChecker) GetQueueCapacity() int {
	return m.capacity
}

func TestNewDegradationController(t *testing.T) {
	config := &DegradationConfig{
		Enabled:         true,
		Levels:          DefaultDegradationLevels(1000),
		RecoveryCooldown: 100 * time.Millisecond,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)
	assert.NotNil(t, controller)
	assert.Equal(t, 0, controller.GetCurrentLevel())
}

func TestDegradationControllerShouldSkip(t *testing.T) {
	config := &DegradationConfig{
		Enabled: true,
		Levels: []DegradationLevel{
			{
				Name:         "placeholder", // level 0 is normal mode
				TriggerCPU:   0.0,
				TriggerQueue: 0,
				Action: DegradationAction{
					AuditLevel: AuditLevelAll,
					SampleRate: 1.0,
				},
			},
			{
				Name:         "mild",
				TriggerCPU:   0.5,
				TriggerQueue: 500,
				Action: DegradationAction{
					AuditLevel: AuditLevelChangesOnly,
					SampleRate: 0.5,
				},
			},
			{
				Name:         "severe",
				TriggerCPU:   0.8,
				TriggerQueue: 800,
				Action: DegradationAction{
					AuditLevel: AuditLevelNone,
					SampleRate: 0.0,
				},
			},
		},
		RecoveryCooldown: 30 * time.Second,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)

	// 正常状态，不跳过
	queryEvent := &handler.Event{Operation: "query"}
	createEvent := &handler.Event{Operation: "create"}
	assert.False(t, controller.ShouldSkip(queryEvent))
	assert.False(t, controller.ShouldSkip(createEvent))

	// 设置降级级别到 level 1 (mild)
	controller.setLevel(1)

	// mild 级别：只跳过 query
	assert.True(t, controller.ShouldSkip(queryEvent))
	assert.False(t, controller.ShouldSkip(createEvent))

	// 设置降级级别到 level 2 (severe)
	controller.setLevel(2)

	// severe 级别：跳过所有事件
	assert.True(t, controller.ShouldSkip(queryEvent))
	assert.True(t, controller.ShouldSkip(createEvent))
}

func TestDegradationControllerSetLevel(t *testing.T) {
	config := &DegradationConfig{
		Enabled: true,
		Levels: []DegradationLevel{
			{
				Action: DegradationAction{SampleRate: 0.5},
			},
			{
				Action: DegradationAction{SampleRate: 0.1},
			},
		},
		RecoveryCooldown: 30 * time.Second,
	}
	sampler := NewRandomSampler(1.0)
	checker := &mockQueueChecker{depth: 0, capacity: 1000}

	controller := NewDegradationController(config, sampler, checker)

	// 设置降级级别
	controller.setLevel(1)
	assert.Equal(t, 1, controller.GetCurrentLevel())
	assert.InDelta(t, 0.1, sampler.GetEffectiveRate(), 0.01)

	// 恢复正常
	controller.setLevel(0)
	assert.Equal(t, 0, controller.GetCurrentLevel())
}
