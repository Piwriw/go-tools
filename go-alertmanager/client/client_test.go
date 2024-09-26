package client

import (
	"testing"
)

func TestGetAlerts(t *testing.T) {
	name, err2 := ExtractServerName("http://notice:9003")
	if err2 != nil {
		t.Error(err2)
	}
	t.Log(name)
	alerts, err := NewClient("10.0.0.195:9003").GetActiveAlerts(false, false, 0, "instance=~\"instance-75|instance-56|instance-58\"", "severity=\"致命:#4A4A4A\"")
	if err != nil {
		t.Error(err)
	}
	for _, alert := range alerts {
		t.Log(*alert)
	}
	t.Log("All alerts len", len(alerts))

}
