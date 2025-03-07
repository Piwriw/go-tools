package format

import (
	"testing"
)

// 测试有效的时区
func TestTimeZoneTime_ValidTimeZones(t *testing.T) {
	time, err := TimeZoneTime("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}

func TestCSTTime(t *testing.T) {
	time, err := CSTTime()
	if err != nil {
		t.Error(err)
	}
	t.Log(time)
}

func TestTimeParsedFormat(t *testing.T) {
	str := "10:10:10"
	time, err := TimeParsedFormat(str)
	if err != nil {
		t.Error(err)
	}
	t.Log(time)
}
