package compare

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsNumericInt tests IsNumeric for all integer types
func TestIsNumericInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOk  bool
		wantVal int64
	}{
		{"valid positive", "123", true, 123},
		{"valid negative", "-456", true, -456},
		{"zero", "0", true, 0},
		{"with plus sign", "+789", true, 789},
		{"invalid alphabet", "abc", false, 0},
		{"invalid float", "123.45", false, 0},
		{"empty", "", false, 0},
		{"overflow int64", "9223372036854775808", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("int", func(t *testing.T) {
				ok, val := IsNumeric[int](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.Equal(t, int(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("int8", func(t *testing.T) {
				ok, val := IsNumeric[int8](tt.input)
				if tt.wantOk && tt.wantVal >= math.MinInt8 && tt.wantVal <= math.MaxInt8 {
					require.True(t, ok)
					assert.Equal(t, int8(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("int16", func(t *testing.T) {
				ok, val := IsNumeric[int16](tt.input)
				if tt.wantOk && tt.wantVal >= math.MinInt16 && tt.wantVal <= math.MaxInt16 {
					require.True(t, ok)
					assert.Equal(t, int16(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("int32", func(t *testing.T) {
				ok, val := IsNumeric[int32](tt.input)
				if tt.wantOk && tt.wantVal >= math.MinInt32 && tt.wantVal <= math.MaxInt32 {
					require.True(t, ok)
					assert.Equal(t, int32(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("int64", func(t *testing.T) {
				ok, val := IsNumeric[int64](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.Equal(t, tt.wantVal, int64(val))
				} else {
					assert.False(t, ok)
				}
			})
		})
	}
}

// TestIsNumericUint tests IsNumeric for all unsigned integer types
func TestIsNumericUint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOk  bool
		wantVal uint64
	}{
		{"valid positive", "123", true, 123},
		{"zero", "0", true, 0},
		{"invalid negative", "-456", false, 0},
		{"invalid alphabet", "abc", false, 0},
		{"invalid float", "123.45", false, 0},
		{"empty", "", false, 0},
		{"overflow uint64", "18446744073709551616", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("uint", func(t *testing.T) {
				ok, val := IsNumeric[uint](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.Equal(t, uint64(tt.wantVal), uint64(val))
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("uint8", func(t *testing.T) {
				ok, val := IsNumeric[uint8](tt.input)
				if tt.wantOk && tt.wantVal <= math.MaxUint8 {
					require.True(t, ok)
					assert.Equal(t, uint8(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("uint16", func(t *testing.T) {
				ok, val := IsNumeric[uint16](tt.input)
				if tt.wantOk && tt.wantVal <= math.MaxUint16 {
					require.True(t, ok)
					assert.Equal(t, uint16(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("uint32", func(t *testing.T) {
				ok, val := IsNumeric[uint32](tt.input)
				if tt.wantOk && tt.wantVal <= math.MaxUint32 {
					require.True(t, ok)
					assert.Equal(t, uint32(tt.wantVal), val)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("uint64", func(t *testing.T) {
				ok, val := IsNumeric[uint64](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.Equal(t, tt.wantVal, uint64(val))
				} else {
					assert.False(t, ok)
				}
			})
		})
	}
}

// TestIsNumericFloat tests IsNumeric for float types
func TestIsNumericFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOk  bool
		wantVal float64
	}{
		{"valid integer", "123", true, 123.0},
		{"valid negative", "-456.78", true, -456.78},
		{"valid float", "123.45", true, 123.45},
		{"zero", "0", true, 0.0},
		{"with plus sign", "+789.5", true, 789.5},
		{"scientific notation", "1e5", true, 1e5},
		{"dot prefix", ".5", true, 0.5},
		{"dot suffix", "5.", true, 5.0},
		{"invalid alphabet", "abc", false, 0},
		{"empty", "", false, 0},
		{"multiple dots", "1.2.3", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("float32", func(t *testing.T) {
				ok, val := IsNumeric[float32](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.InDelta(t, tt.wantVal, float64(val), 1e-3)
				} else {
					assert.False(t, ok)
				}
			})
			t.Run("float64", func(t *testing.T) {
				ok, val := IsNumeric[float64](tt.input)
				if tt.wantOk {
					require.True(t, ok)
					assert.InDelta(t, tt.wantVal, val, 1e-6)
				} else {
					assert.False(t, ok)
				}
			})
		})
	}
}

// TestIsNumericEdgeCases 测试边界情况
func TestIsNumericEdgeCases(t *testing.T) {
	t.Run("max values", func(t *testing.T) {
		ok, _ := IsNumeric[int](strconv.Itoa(math.MaxInt))
		assert.True(t, ok)
		ok, _ = IsNumeric[int64](strconv.FormatInt(math.MaxInt64, 10))
		assert.True(t, ok)
		ok, _ = IsNumeric[uint](strconv.FormatUint(math.MaxUint, 10))
		assert.True(t, ok)
		ok, _ = IsNumeric[uint64](strconv.FormatUint(math.MaxUint64, 10))
		assert.True(t, ok)
	})

	t.Run("min values", func(t *testing.T) {
		ok, _ := IsNumeric[int](strconv.Itoa(math.MinInt))
		assert.True(t, ok)
		ok, _ = IsNumeric[int64](strconv.FormatInt(math.MinInt64, 10))
		assert.True(t, ok)
	})
}
