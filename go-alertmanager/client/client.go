package client

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/prometheus/alertmanager/api/v2/client/alertgroup"
	"github.com/prometheus/alertmanager/api/v2/client/receiver"
	"github.com/prometheus/alertmanager/api/v2/models"
)

var defaultTimeout = 10 * time.Second

type AlertManagerClient struct {
	client    *client.AlertmanagerAPI
	clientOpt *clientOpt
	formats   strfmt.Registry
	err       error
}
type AlertManagerClientOpt func(*AlertManagerClient)

func WithRegistry(registry strfmt.Registry) AlertManagerClientOpt {
	return func(c *AlertManagerClient) {
		c.formats = registry
	}
}

type clientOpt struct {
	silenced  bool
	inhibited bool
	filter    []string
	timeout   time.Duration
}

func ExtractServerName(fullURL string) (string, error) {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}
	return parsedURL.Host, nil
}

func NewClientWithPort(ip string, port int, opts ...AlertManagerClientOpt) *AlertManagerClient {
	address := fmt.Sprintf("%s:%d", ip, port)
	transport := httptransport.New(address, "/api/v2", []string{"http"})
	amClient := client.New(transport, nil)
	cli := &AlertManagerClient{
		client:    amClient,
		clientOpt: defaultClientOpt(),
		formats:   strfmt.Default,
	}

	// 解析可选参数
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func defaultClientOpt() *clientOpt {
	return &clientOpt{
		silenced:  false,
		inhibited: false,
		timeout:   defaultTimeout,
	}
}

// NewClient 创建一个新的 AlertManager 客户端
func NewClient(address string) *AlertManagerClient {
	transport := httptransport.New(address, "/api/v2", []string{"http"})
	amClient := client.New(transport, nil)
	return &AlertManagerClient{client: amClient, clientOpt: defaultClientOpt()}
}

func (am *AlertManagerClient) Silenced(silenced bool) *AlertManagerClient {
	if am.err != nil {
		return am
	}
	// 创建新的 clientOpt 副本，返回新的 AlertManagerClient
	newClientOpt := &clientOpt{
		silenced:  am.clientOpt.silenced,
		inhibited: am.clientOpt.inhibited,
		filter:    am.clientOpt.filter,
		timeout:   am.clientOpt.timeout,
	}
	newClientOpt.silenced = silenced
	return &AlertManagerClient{
		client:    am.client,
		clientOpt: newClientOpt,
		err:       am.err,
	}
}

func (am *AlertManagerClient) Inhibited(inhibited bool) *AlertManagerClient {
	if am.err != nil {
		return am
	}
	// 创建新的 clientOpt 副本，返回新的 AlertManagerClient
	newClientOpt := &clientOpt{
		silenced:  am.clientOpt.silenced,
		inhibited: am.clientOpt.inhibited,
		filter:    am.clientOpt.filter,
		timeout:   am.clientOpt.timeout,
	}
	newClientOpt.inhibited = inhibited
	return &AlertManagerClient{
		client:    am.client,
		clientOpt: newClientOpt,
		err:       am.err,
	}
}

func (am *AlertManagerClient) Timeout(timeout time.Duration) *AlertManagerClient {
	if am.err != nil {
		return am
	}
	// 创建新的 clientOpt 副本，返回新的 AlertManagerClient
	newClientOpt := &clientOpt{
		silenced:  am.clientOpt.silenced,
		inhibited: am.clientOpt.inhibited,
		filter:    am.clientOpt.filter,
		timeout:   am.clientOpt.timeout,
	}
	if timeout <= 0 {
		am.err = errors.New("timeout must be greater than zero")
	}
	return &AlertManagerClient{
		client:    am.client,
		clientOpt: newClientOpt,
		err:       am.err,
	}
}

func (am *AlertManagerClient) Filter(filter ...string) *AlertManagerClient {
	if am.err != nil {
		return am
	}
	// 创建新的 clientOpt 副本，返回新的 AlertManagerClient
	newClientOpt := &clientOpt{
		silenced:  am.clientOpt.silenced,
		inhibited: am.clientOpt.inhibited,
		filter:    am.clientOpt.filter,
		timeout:   am.clientOpt.timeout,
	}
	newClientOpt.filter = filter
	return &AlertManagerClient{
		client:    am.client,
		clientOpt: newClientOpt,
		err:       am.err,
	}
}

// GetAlerts 获取活跃告警，接受过滤器参数
func (am *AlertManagerClient) GetAlerts() ([]*models.GettableAlert, error) {
	if am.err != nil {
		return nil, am.err
	}
	var params = &alert.GetAlertsParams{
		Inhibited: &am.clientOpt.inhibited,
		Silenced:  &am.clientOpt.silenced,
	}
	if len(am.clientOpt.filter) > 0 {
		params.WithFilter(am.clientOpt.filter)
	}
	alertRes, err := am.client.Alert.GetAlerts(params)
	if err != nil {
		return nil, fmt.Errorf("error fetching alerts: %v", err)
	}
	if !alertRes.IsSuccess() {
		return nil, errors.New("error getting alerts")
	}
	return alertRes.GetPayload(), nil
}

// AddAlert 添加告警
func (am *AlertManagerClient) AddAlert(ale *models.PostableAlert) error {
	if am.err != nil {
		return am.err
	}
	if err := ale.Validate(am.formats); err != nil {
		am.err = errors.Join(am.err, err)
	}
	if am.err != nil {
		return am.err
	}
	alerts := make([]*models.PostableAlert, 0)
	alerts = append(alerts, ale)
	params := &alert.PostAlertsParams{
		Alerts: alerts,
	}

	params.WithTimeout(am.clientOpt.timeout)
	_, err := am.client.Alert.PostAlerts(params, nil)
	if err != nil {
		return fmt.Errorf("error posting alert: %v", err)
	}
	return nil
}

// AddAlerts 添加告警
func (am *AlertManagerClient) AddAlerts(alerts []*models.PostableAlert) error {
	if am.err != nil {
		return am.err
	}
	for _, alert := range alerts {
		if err := alert.Validate(am.formats); err != nil {
			am.err = errors.Join(am.err, err)
		}
	}
	if am.err != nil {
		return am.err
	}
	params := &alert.PostAlertsParams{
		Alerts: alerts,
	}
	params.WithTimeout(am.clientOpt.timeout)
	_, err := am.client.Alert.PostAlerts(params)
	if err != nil {
		return fmt.Errorf("error posting alerts: %v", err)
	}
	return nil
}

// GetReceivers 获取告警接收者
func (am *AlertManagerClient) GetReceivers() ([]*models.Receiver, error) {
	if am.err != nil {
		return nil, am.err
	}
	params := &receiver.GetReceiversParams{}
	params.WithTimeout(am.clientOpt.timeout)
	receiverRes, err := am.client.Receiver.GetReceivers(params)
	if err != nil {
		return nil, fmt.Errorf("error get alerts: %v", err)
	}
	if !receiverRes.IsSuccess() {
		return nil, errors.New("error getting alerts")
	}
	return receiverRes.GetPayload(), err
}

func (am *AlertManagerClient) GetAlertGroups() ([]*models.AlertGroup, error) {
	if am.err != nil {
		return nil, am.err
	}
	params := &alertgroup.GetAlertGroupsParams{
		Filter:    am.clientOpt.filter,
		Inhibited: &am.clientOpt.inhibited,
		Silenced:  &am.clientOpt.silenced,
	}
	params.WithTimeout(am.clientOpt.timeout)
	receiverRes, err := am.client.Alertgroup.GetAlertGroups(params)
	if err != nil {
		return nil, fmt.Errorf("error get alerts: %v", err)
	}
	if !receiverRes.IsSuccess() {
		return nil, errors.New("error getting alerts")
	}
	return receiverRes.GetPayload(), err
}

func (am *AlertManagerClient) ValidateRules() ([]*models.AlertGroup, error) {
	if am.err != nil {
		return nil, am.err
	}
	return nil, nil
}
