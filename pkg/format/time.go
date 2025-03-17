package format

import (
	"fmt"
	"time"
)

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

func TimeParsedFormat(timeStr string) (*time.Time, error) {
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
