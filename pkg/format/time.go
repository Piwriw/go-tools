package format

import "time"

// ConvertToWeekdays 将[]int转换为[]time.Weekday
func ConvertToWeekdays(intArray []int) []time.Weekday {
	weekdays := make([]time.Weekday, len(intArray))
	for i, v := range intArray {
		weekdays[i] = time.Weekday(v)
	}
	return weekdays
}
