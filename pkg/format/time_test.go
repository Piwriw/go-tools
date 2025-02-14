package format

import "testing"

func TestName(t *testing.T) {
	str := "10:10:10"
	time, err := parseTimeOfDay(str)
	if err != nil {
		t.Error(err)
	}
	t.Log(time)
}
