package format

import (
	"errors"
	"time"
)

// TimeStampUnit 时间戳单位类型
type TimeStampUnit string

// 时间戳单位常量
const (
	TimeStampUnitSecond      TimeStampUnit = "second"
	TimeStampUnitMillisecond TimeStampUnit = "millisecond"
	TimeStampUnitMicrosecond TimeStampUnit = "microsecond"
	TimeStampUnitNanosecond  TimeStampUnit = "nanosecond"
)

// nowFunc 用于测试时替换当前时间
var nowFunc = time.Now

// errInvalidUnit 非法时间戳单位错误
var errInvalidUnit = errors.New("invalid timestamp unit")

// timeStampToTime 将时间戳转换为 time.Time（内部函数）。
// 注意：此函数名以小写开头，仅包内可用。
func timeStampToTime(timestamp int64, unit TimeStampUnit) (time.Time, error) {
	switch unit {
	case TimeStampUnitSecond:
		return time.Unix(timestamp, 0), nil
	case TimeStampUnitMillisecond:
		secs := timestamp / 1000
		nanos := (timestamp % 1000) * 1000 * 1000
		return time.Unix(secs, nanos), nil
	case TimeStampUnitMicrosecond:
		secs := timestamp / 1000 / 1000
		nanos := (timestamp % 1000) * 1000
		return time.Unix(secs, nanos), nil
	case TimeStampUnitNanosecond:
		secs := timestamp / 1000 / 1000 / 1000
		nanos := timestamp % (1000 * 1000 * 1000)
		return time.Unix(secs, nanos), nil
	default:
		return time.Time{}, errInvalidUnit
	}
}

// timeToTimeStamp 将 time.Time 转换为时间戳（内部函数）。
// 注意：此函数名以小写开头，仅包内可用。
func timeToTimeStamp(t time.Time, unit TimeStampUnit) (int64, error) {
	nanos := t.UnixNano()

	switch unit {
	case TimeStampUnitSecond:
		return nanos / 1000 / 1000 / 1000, nil
	case TimeStampUnitMillisecond:
		return nanos / 1000 / 1000, nil
	case TimeStampUnitMicrosecond:
		return nanos / 1000, nil
	case TimeStampUnitNanosecond:
		return nanos, nil
	default:
		return 0, errInvalidUnit
	}
}

// loadLocationOrDefault 加载时区，空字符串返回上海时区。
// 注意：此函数名以小写开头，仅包内可用。
func loadLocationOrDefault(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

// toShanghaiTime 将时间戳转换为上海时区的 time.Time。
// 注意：此函数名以小写开头，仅包内可用。
func toShanghaiTime(timestamp int64, unit TimeStampUnit) (time.Time, error) {
	t, err := timeStampToTime(timestamp, unit)
	if err != nil {
		return time.Time{}, err
	}
	loc, err := loadLocationOrDefault("")
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// ============================================================================
// Public Conversion APIs
// ============================================================================

// TimeStampToTime 将指定单位的时间戳转为 time.Time。
func TimeStampToTime(timestamp int64, unit TimeStampUnit) (time.Time, error) {
	return timeStampToTime(timestamp, unit)
}

// TimeToTimeStamp 将 time.Time 转为指定单位的时间戳。
func TimeToTimeStamp(t time.Time, unit TimeStampUnit) (int64, error) {
	return timeToTimeStamp(t, unit)
}

// ============================================================================
// Current Timestamp
// ============================================================================

// NowTimeStamp 返回当前时间的时间戳，按指定单位。
func NowTimeStamp(unit TimeStampUnit) (int64, error) {
	return timeToTimeStamp(nowFunc(), unit)
}

// NowTimeStampSecond 返回当前时间的秒级时间戳。
func NowTimeStampSecond() int64 {
	return nowFunc().Unix()
}

// NowTimeStampMillisecond 返回当前时间的毫秒级时间戳。
func NowTimeStampMillisecond() int64 {
	return nowFunc().UnixMilli()
}

// ============================================================================
// Timestamp → Formatted String
// ============================================================================

// TimeStampToDateTime 将时间戳转为 "2006-01-02 15:04:05" 格式字符串（上海时区）。
func TimeStampToDateTime(timestamp int64, unit TimeStampUnit) (string, error) {
	t, err := toShanghaiTime(timestamp, unit)
	if err != nil {
		return "", err
	}
	return t.Format(time.DateTime), nil
}

// TimeStampToFormat 将时间戳按指定 layout 格式化（上海时区）。
func TimeStampToFormat(timestamp int64, unit TimeStampUnit, layout string) (string, error) {
	t, err := toShanghaiTime(timestamp, unit)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// TimeStampToDateTimeIn 将时间戳转为 "2006-01-02 15:04:05" 格式字符串（指定时区）。
func TimeStampToDateTimeIn(timestamp int64, unit TimeStampUnit, timezone string) (string, error) {
	t, err := timeStampToTime(timestamp, unit)
	if err != nil {
		return "", err
	}
	loc, err := loadLocationOrDefault(timezone)
	if err != nil {
		return "", err
	}
	return t.In(loc).Format(time.DateTime), nil
}

// TimeStampToFormatIn 将时间戳按指定 layout 格式化（指定时区）。
func TimeStampToFormatIn(timestamp int64, unit TimeStampUnit, layout string, timezone string) (string, error) {
	t, err := timeStampToTime(timestamp, unit)
	if err != nil {
		return "", err
	}
	loc, err := loadLocationOrDefault(timezone)
	if err != nil {
		return "", err
	}
	return t.In(loc).Format(layout), nil
}

// ============================================================================
// Formatted String → Timestamp
// ============================================================================

// DateTimeToTimeStamp 将日期时间字符串解析后转为指定单位的时间戳。
// layout 为 Go 标准时间格式，timezone 为时区（空字符串默认上海时区）。
func DateTimeToTimeStamp(datetimeStr, layout string, unit TimeStampUnit, timezone string) (int64, error) {
	loc, err := loadLocationOrDefault(timezone)
	if err != nil {
		return 0, err
	}
	t, err := time.ParseInLocation(layout, datetimeStr, loc)
	if err != nil {
		return 0, err
	}
	return timeToTimeStamp(t, unit)
}

// ============================================================================
// Unit Detection
// ============================================================================

// DetectTimeStampUnit 根据时间戳数量级自动推断单位。
// 推断规则：
//   - < 1e10          → 秒（second）
//   - < 1e13          → 毫秒（millisecond）
//   - < 1e16          → 微秒（microsecond）
//   - 其他            → 纳秒（nanosecond）
func DetectTimeStampUnit(timestamp int64) TimeStampUnit {
	abs := timestamp
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs < 1e10:
		return TimeStampUnitSecond
	case abs < 1e13:
		return TimeStampUnitMillisecond
	case abs < 1e16:
		return TimeStampUnitMicrosecond
	default:
		return TimeStampUnitNanosecond
	}
}

// AutoTimeStampToTime 自动检测时间戳单位并转换为 time.Time（上海时区）。
func AutoTimeStampToTime(timestamp int64) (time.Time, error) {
	unit := DetectTimeStampUnit(timestamp)
	return toShanghaiTime(timestamp, unit)
}
