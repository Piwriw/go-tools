package gocron

import (
	"fmt"
	"time"
)

const (
	DayTimeType   string = "%d %d * * * "
	WeekTimeType  string = "%d %d * * %d "
	MonthTimeType string = "%d %d * %d * "
)

type TimeType string

// DayTimeToCron 将 time.Time 转换为 Cron 表达式
func DayTimeToCron(t time.Time) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(DayTimeType, minute, hour)
}

func WeekTimeToCron(t time.Time, week time.Weekday) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(WeekTimeType, minute, hour, week)
}

func MonthTimeToCron(t time.Time, month time.Month) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(MonthTimeType, minute, hour, month)
}
