package format

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeStampUnitConstants(t *testing.T) {
	tests := []struct {
		name     string
		unit     TimeStampUnit
		expected string
	}{
		{"second", TimeStampUnitSecond, "second"},
		{"millisecond", TimeStampUnitMillisecond, "millisecond"},
		{"microsecond", TimeStampUnitMicrosecond, "microsecond"},
		{"nanosecond", TimeStampUnitNanosecond, "nanosecond"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.unit))
		})
	}
}

func TestTimeStampToTime(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		wantUnix  int64
		wantNanos int
		wantErr   bool
	}{
		{
			name:      "second",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			wantUnix:  baseTimestamp,
			wantNanos: 0,
		},
		{
			name:      "millisecond",
			timestamp: baseTimestamp * 1000,
			unit:      TimeStampUnitMillisecond,
			wantUnix:  baseTimestamp,
			wantNanos: 0,
		},
		{
			name:      "microsecond",
			timestamp: baseTimestamp * 1000 * 1000,
			unit:      TimeStampUnitMicrosecond,
			wantUnix:  baseTimestamp,
			wantNanos: 0,
		},
		{
			name:      "nanosecond",
			timestamp: baseTimestamp * 1000 * 1000 * 1000,
			unit:      TimeStampUnitNanosecond,
			wantUnix:  baseTimestamp,
			wantNanos: 0,
		},
		{
			name:      "nanosecond with fractional seconds",
			timestamp: baseTimestamp*1000*1000*1000 + 500*1000*1000, // +500ms in nanos
			unit:      TimeStampUnitNanosecond,
			wantUnix:  baseTimestamp,
			wantNanos: 500 * 1000 * 1000,
		},
		{
			name:      "zero timestamp",
			timestamp: 0,
			unit:      TimeStampUnitSecond,
			wantUnix:  0,
			wantNanos: 0,
		},
		{
			name:      "negative timestamp",
			timestamp: -1,
			unit:      TimeStampUnitSecond,
			wantUnix:  -1,
			wantNanos: 0,
		},
		{
			name:      "invalid unit",
			timestamp: baseTimestamp,
			unit:      TimeStampUnit("invalid"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timeStampToTime(tt.timestamp, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUnix, got.Unix())
			assert.Equal(t, tt.wantNanos, got.Nanosecond())
		})
	}
}

func TestTimeToTimeStamp(t *testing.T) {
	// 2024-01-01 00:00:00 UTC + 500ms
	tm := time.Unix(1704067200, 500*1000*1000)

	tests := []struct {
		name    string
		time    time.Time
		unit    TimeStampUnit
		want    int64
		wantErr bool
	}{
		{
			name: "second",
			time: tm,
			unit: TimeStampUnitSecond,
			want: 1704067200,
		},
		{
			name: "millisecond",
			time: tm,
			unit: TimeStampUnitMillisecond,
			want: 1704067200*1000 + 500,
		},
		{
			name: "microsecond",
			time: tm,
			unit: TimeStampUnitMicrosecond,
			want: 1704067200*1000*1000 + 500*1000,
		},
		{
			name: "nanosecond",
			time: tm,
			unit: TimeStampUnitNanosecond,
			want: 1704067200*1000*1000*1000 + 500*1000*1000,
		},
		{
			name:    "invalid unit",
			time:    tm,
			unit:    TimeStampUnit("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timeToTimeStamp(tt.time, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadLocationOrDefault(t *testing.T) {
	tests := []struct {
		name    string
		tz      string
		wantErr bool
	}{
		{"empty returns Shanghai", "", false},
		{"Asia/Shanghai", "Asia/Shanghai", false},
		{"UTC", "UTC", false},
		{"America/New_York", "America/New_York", false},
		{"invalid timezone", "Invalid/Zone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := loadLocationOrDefault(tt.tz)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, loc)
			if tt.tz == "" {
				assert.Equal(t, "Asia/Shanghai", loc.String())
			} else {
				assert.Equal(t, tt.tz, loc.String())
			}
		})
	}
}

func TestToShanghaiTime(t *testing.T) {
	// 2024-01-01 00:00:00 UTC = 2024-01-01 08:00:00 Shanghai
	const baseTimestamp int64 = 1704067200

	got, err := toShanghaiTime(baseTimestamp, TimeStampUnitSecond)
	assert.NoError(t, err)
	assert.Equal(t, "Asia/Shanghai", got.Location().String())
	assert.Equal(t, 2024, got.Year())
	assert.Equal(t, time.January, got.Month())
	assert.Equal(t, 1, got.Day())
	assert.Equal(t, 8, got.Hour())
}

// ============================================================================
// Public Conversion APIs
// ============================================================================

func TestTimeStampToTime_Public(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		wantUnix  int64
		wantErr   bool
	}{
		{
			name:      "second",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			wantUnix:  baseTimestamp,
		},
		{
			name:      "millisecond",
			timestamp: baseTimestamp * 1000,
			unit:      TimeStampUnitMillisecond,
			wantUnix:  baseTimestamp,
		},
		{
			name:      "microsecond",
			timestamp: baseTimestamp * 1000 * 1000,
			unit:      TimeStampUnitMicrosecond,
			wantUnix:  baseTimestamp,
		},
		{
			name:      "nanosecond",
			timestamp: baseTimestamp * 1000 * 1000 * 1000,
			unit:      TimeStampUnitNanosecond,
			wantUnix:  baseTimestamp,
		},
		{
			name:      "invalid unit",
			timestamp: baseTimestamp,
			unit:      TimeStampUnit("invalid"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStampToTime(tt.timestamp, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUnix, got.Unix())
		})
	}
}

func TestTimeToTimeStamp_Public(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	tm := time.Unix(1704067200, 0)

	tests := []struct {
		name    string
		time    time.Time
		unit    TimeStampUnit
		want    int64
		wantErr bool
	}{
		{
			name: "second",
			time: tm,
			unit: TimeStampUnitSecond,
			want: 1704067200,
		},
		{
			name: "millisecond",
			time: tm,
			unit: TimeStampUnitMillisecond,
			want: 1704067200 * 1000,
		},
		{
			name: "microsecond",
			time: tm,
			unit: TimeStampUnitMicrosecond,
			want: 1704067200 * 1000 * 1000,
		},
		{
			name: "nanosecond",
			time: tm,
			unit: TimeStampUnitNanosecond,
			want: 1704067200 * 1000 * 1000 * 1000,
		},
		{
			name:    "invalid unit",
			time:    tm,
			unit:    TimeStampUnit("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeToTimeStamp(tt.time, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ============================================================================
// Formatting APIs
// ============================================================================

func TestTimeStampToDateTime(t *testing.T) {
	// 2024-01-01 08:00:00 Shanghai = 2024-01-01 00:00:00 UTC
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		want      string
		wantErr   bool
	}{
		{
			name:      "second",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			want:      "2024-01-01 08:00:00",
		},
		{
			name:      "millisecond",
			timestamp: baseTimestamp * 1000,
			unit:      TimeStampUnitMillisecond,
			want:      "2024-01-01 08:00:00",
		},
		{
			name:      "invalid unit",
			timestamp: baseTimestamp,
			unit:      TimeStampUnit("invalid"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStampToDateTime(tt.timestamp, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStampToFormat(t *testing.T) {
	// 2024-01-01 08:00:00 Shanghai
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		layout    string
		want      string
		wantErr   bool
	}{
		{
			name:      "date only",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			layout:    "2006-01-02",
			want:      "2024-01-01",
		},
		{
			name:      "RFC3339",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			layout:    time.RFC3339,
			want:      "2024-01-01T08:00:00+08:00",
		},
		{
			name:      "invalid unit",
			timestamp: baseTimestamp,
			unit:      TimeStampUnit("invalid"),
			layout:    time.RFC3339,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStampToFormat(tt.timestamp, tt.unit, tt.layout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStampToDateTimeIn(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		timezone  string
		want      string
		wantErr   bool
	}{
		{
			name:      "UTC",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			timezone:  "UTC",
			want:      "2024-01-01 00:00:00",
		},
		{
			name:      "America/New_York",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			timezone:  "America/New_York",
			want:      "2023-12-31 19:00:00",
		},
		{
			name:      "invalid timezone",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			timezone:  "Invalid/Zone",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStampToDateTimeIn(tt.timestamp, tt.unit, tt.timezone)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeStampToFormatIn(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	const baseTimestamp int64 = 1704067200

	tests := []struct {
		name      string
		timestamp int64
		unit      TimeStampUnit
		layout    string
		timezone  string
		want      string
		wantErr   bool
	}{
		{
			name:      "UTC RFC3339",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			layout:    time.RFC3339,
			timezone:  "UTC",
			want:      "2024-01-01T00:00:00Z",
		},
		{
			name:      "invalid timezone",
			timestamp: baseTimestamp,
			unit:      TimeStampUnitSecond,
			layout:    time.RFC3339,
			timezone:  "Invalid/Zone",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStampToFormatIn(tt.timestamp, tt.unit, tt.layout, tt.timezone)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
