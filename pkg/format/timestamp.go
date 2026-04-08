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
