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

func TimeParsedFormat(timeStr string) (*time.Time, error) {
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
