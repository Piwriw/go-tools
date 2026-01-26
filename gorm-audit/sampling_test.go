package audit

import (
	"context"
	"testing"

	"github.com/piwriw/gorm/gorm-audit/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewRandomSampler(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	assert.NotNil(t, sampler)
	assert.Equal(t, 0.5, sampler.GetEffectiveRate())
	assert.Equal(t, "random", sampler.String())
}

func TestRandomSamplerShouldSample(t *testing.T) {
	tests := []struct {
		name     string
		rate     float64
		minCount int
		maxCount int
	}{
		{"always sample", 1.0, 1000, 1000},
		{"never sample", 0.0, 0, 0},
		{"50% sample", 0.5, 400, 600},
		{"10% sample", 0.1, 50, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := NewRandomSampler(tt.rate)
			event := &handler.Event{}

			count := 0
			trials := 1000
			for i := 0; i < trials; i++ {
				if sampler.ShouldSample(context.Background(), event) {
					count++
				}
			}

			assert.GreaterOrEqual(t, count, tt.minCount)
			assert.LessOrEqual(t, count, tt.maxCount)
		})
	}
}

func TestRandomSamplerUpdateRate(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	assert.Equal(t, 0.5, sampler.GetEffectiveRate())

	sampler.UpdateRate(0.1)
	assert.Equal(t, 0.1, sampler.GetEffectiveRate())
}

func TestRandomSamplerThreadSafe(t *testing.T) {
	sampler := NewRandomSampler(0.5)
	event := &handler.Event{}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				sampler.ShouldSample(context.Background(), event)
				sampler.UpdateRate(0.3)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证没有 panic 且采样率在合理范围内
	rate := sampler.GetEffectiveRate()
	assert.GreaterOrEqual(t, rate, 0.0)
	assert.LessOrEqual(t, rate, 1.0)
}
