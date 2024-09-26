package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/prometheus/alertmanager/api/v2/models"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
)

var defaultTimeout = 10 * time.Second

type AlertManagerClient struct {
	client *client.AlertmanagerAPI
}

func ExtractServerName(fullURL string) (string, error) {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}
	return parsedURL.Host, nil
}

func NewClientWithPort(ip string, port int) *AlertManagerClient {
	address := fmt.Sprintf("%s:%d", ip, port)
	transport := httptransport.New(address, "/api/v2", []string{"http"})
	amClient := client.New(transport, nil)
	return &AlertManagerClient{client: amClient}
}

// NewClient 创建一个新的 AlertManager 客户端
func NewClient(address string) *AlertManagerClient {
	transport := httptransport.New(address, "/api/v2", []string{"http"})
	amClient := client.New(transport, nil)
	return &AlertManagerClient{client: amClient}
}

// GetActiveAlerts 获取活跃告警，接受过滤器参数
func (am *AlertManagerClient) GetActiveAlerts(silenced bool, inhibited bool, timeout time.Duration, filter ...string) ([]*models.GettableAlert, error) {
	var params = &alert.GetAlertsParams{
		Inhibited: &inhibited,
		Silenced:  &silenced,
	}
	if len(filter) > 0 {
		params.WithFilter(filter)
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	alerts, err := am.client.Alert.GetAlerts(params)
	if err != nil {
		return nil, errors.Errorf("error fetching alerts: %v", err)
	}
	return alerts.Payload, nil
}

func (am *AlertManagerClient) AddAlert(ale *models.PostableAlert, timeout time.Duration) error {
	alerts := make([]*models.PostableAlert, 0)
	alerts = append(alerts, ale)
	params := &alert.PostAlertsParams{
		Alerts: alerts,
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	params.WithTimeout(defaultTimeout)
	_, err := am.client.Alert.PostAlerts(params, nil)
	if err != nil {
		return errors.Errorf("error posting alert: %v", err)
	}
	return nil
}

func (am *AlertManagerClient) AddAlerts(alerts []*models.PostableAlert, timeout time.Duration) error {
	params := &alert.PostAlertsParams{
		Alerts: alerts,
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	params.WithTimeout(defaultTimeout)
	_, err := am.client.Alert.PostAlerts(params)
	if err != nil {
		return errors.Errorf("error posting alerts: %v", err)
	}
	return nil
}
