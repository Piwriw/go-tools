package format

import (
	"fmt"
	"time"
)

// ParseCSTTime 解析一个中国标准时间（CST）格式的时间字符串为 time.Time
func ParseCSTTime(layout, cstTime string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai") // 指定 CST（中国标准时间）
	if err != nil {
		return time.Time{}, err
	}

	// 用指定时区解析时间字符串
	t, err := time.ParseInLocation(layout, cstTime, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// TimeStampDateTime 将时间戳转换为时间格式
func TimeStampDateTime(timestamp int64) (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	t := time.Unix(timestamp, 0).In(loc)
	return t.Format(time.DateTime), nil
}

// TimeStampFormat 将时间戳转换为指定格式的时间字符串
func TimeStampFormat(timestamp int64, layout string) (string, error) {
	t := time.Unix(timestamp, 0)
	return t.Format(layout), nil
}

// CSTTime 返回当前北京时间（CST）
func CSTTime() (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Now(), err
	}
	return time.Now().In(loc), nil
}

// TimeZoneTime 返回指定时区的当前时间
func TimeZoneTime(timezone string) (time.Time, error) {
	// 默认时区
	locName := "Asia/Shanghai"
	if timezone != "" {
		locName = timezone
	}

	// 加载时区
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location %s: %v", locName, err)
	}

	// 返回当前时间（带时区信息）
	return time.Now().In(loc), nil
}

// ConvertToWeekdays 将[]int转换为[]time.Weekday
func ConvertToWeekdays(intArray []int) []time.Weekday {
	weekdays := make([]time.Weekday, len(intArray))
	for i, v := range intArray {
		weekdays[i] = time.Weekday(v)
	}
	return weekdays
}

// IsSameDay 判断两个时间是否为同一天（按各自本地时区）
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsToday 判断时间是否为今天（按该时间自身时区）
func IsToday(t time.Time) bool {
	now := time.Now().In(t.Location())
	return IsSameDay(t, now)
}

// IsWeekend 判断时间是否为周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// StartOfDay 返回当天开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回当天结束时间
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - time.Nanosecond)
}

// StartOfMonth 返回当月开始时间
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 返回当月结束时间
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// Between 判断目标时间是否在开始和结束时间之间（闭区间）
func Between(start, end, target time.Time) bool {
	return !target.Before(start) && !target.After(end)
}

// StartOfWeek 返回当周开始时间（周一为一周起点）
func StartOfWeek(t time.Time) time.Time {
	startOfDay := StartOfDay(t)
	weekday := int(startOfDay.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return startOfDay.AddDate(0, 0, -(weekday - 1))
}

// EndOfWeek 返回当周结束时间（周日为一周终点）
func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// DaysBetween 返回两个时间之间的自然日差
func DaysBetween(start, end time.Time) int {
	startDate := StartOfDay(start)
	endDate := StartOfDay(end)
	return int(endDate.Sub(startDate) / (24 * time.Hour))
}

// StartOfYear 返回当年开始时间
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 返回当年结束时间
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// IsZeroTime 判断是否为零值时间
func IsZeroTime(t time.Time) bool {
	return t.IsZero()
}

func TimeParsedFormat(timeStr string) (*time.Time, error) {
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
