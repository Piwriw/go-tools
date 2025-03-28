package client

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/models"
)

func TestAddAlerts(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	var (
		alerts = []*models.PostableAlert{
			{
				Annotations: models.LabelSet{
					"name": "test",
				},
				EndsAt:   strfmt.DateTime{},
				StartsAt: strfmt.DateTime{},
				Alert:    models.Alert{},
			},
		}
	)
	client.AddAlerts()
}

func TestGetAlertGroups(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	alertGroups, err := client.Silenced(false).Inhibited(false).GetAlertGroups()
	if err != nil {
		t.Error(err)
	}
	for _, group := range alertGroups {
		t.Logf("%+v\n", group)
	}
}

func TestGetReceivers(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	receivers, err := client.GetReceivers()
	if err != nil {
		t.Fatal(err)
	}
	for _, receiver := range receivers {
		t.Log(*receiver.Name)
	}
}

func TestGetAlert(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	alerts, err := client.Silenced(false).
		Filter("severity=\"致命:#4A4A4A\"\n").Inhibited(false).GetAlerts()
	if err != nil {
		t.Error(err)
	}
	for _, alert := range alerts {
		t.Log(alert)
	}
}

func TestGetAlerts(t *testing.T) {
	// name, err2 := ExtractServerName("http://notice:9003")
	// if err2 != nil {
	//	t.Error(err2)
	// }
	// t.Log(name)
	// false, false, 0, , "severity=\"致命:#4A4A4A\""
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
