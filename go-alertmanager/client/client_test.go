package client

import (
	"testing"
	"time"
)

func TestGetAlerts(t *testing.T) {
	//name, err2 := ExtractServerName("http://notice:9003")
	//if err2 != nil {
	//	t.Error(err2)
	//}
	//t.Log(name)
	//false, false, 0, , "severity=\"致命:#4A4A4A\""
	client := NewClient("10.0.0.195:9003")
	_, err := client.Silenced(false).
		Filter("instance=~\"instance-75|instance-56|instance-58\"").Inhibited(false).GetAlerts()
	if err != nil {
		t.Error(err)
	}

	go func(managerClient *AlertManagerClient) {
		for i := 0; i < 10; i++ {
			alerts, err := managerClient.Inhibited(false).GetAlerts()
			if err != nil {
				t.Error(err)
			}

			t.Log("All alerts len", len(alerts))
		}
	}(client)

	for i := 0; i < 10; i++ {
		alerts, err := client.Filter("instance=\"instance-46\"").GetAlerts()
		if err != nil {
			t.Error(err)
		}

		t.Log("All alerts len", len(alerts))

	}
	time.Sleep(1 * time.Minute)

}
