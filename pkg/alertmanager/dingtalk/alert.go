package dingtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.piwriw.go-tools/pkg/alertmanager/common"

	"github.com/prometheus/alertmanager/template"
	"github.piwriw.go-tools/pkg/http"
)

// DingtalkAlert 钉钉告警策略
type DingtalkAlert struct {
	Client     *http.Client
	BeforeHook common.HookFunc
	AfterHook  common.HookFunc
	LimitFunc  common.LimitOption
}

var _ = common.AlertMethod(&DingtalkAlert{})

func NewDingtalkAlert() *DingtalkAlert {
	return &DingtalkAlert{Client: http.NewHTTPClient()}
}
func NewDingtalkAlertWithClient(client *http.Client) *DingtalkAlert {
	return &DingtalkAlert{Client: client}
}

func (d *DingtalkAlert) WithLimit(limit common.LimitOption) *DingtalkAlert {
	if limit == nil {
		return d
	}
	d.LimitFunc = limit
	return d
}

func (d *DingtalkAlert) WithBeforeHook(hook common.HookFunc) *DingtalkAlert {
	if hook == nil {
		return d
	}
	d.BeforeHook = hook
	return d
}

func (d *DingtalkAlert) WithAfterHook(hook common.HookFunc) *DingtalkAlert {
	if hook == nil {
		return d
	}
	d.AfterHook = hook
	return d
}

func (d *DingtalkAlert) Before(c common.AlertContext) error {
	if d.BeforeHook != nil {
		if err := d.BeforeHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (d *DingtalkAlert) After(c common.AlertContext) error {
	if d.AfterHook != nil {
		if err := d.AfterHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (d *DingtalkAlert) Limit(idx int, alerts template.Alerts) bool {
	if d.LimitFunc != nil {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		if d.LimitFunc.Allow() {
			return true
		}
		return false
	}
	return true
}

func (d *DingtalkAlert) Execute(c common.AlertContext, users []string) error {
	if err := d.Before(c); err != nil {
		return common.ErrorBeforeHook
	}
	dingtalkContext, ok := c.(*DingtalkAlertContext)
	if !ok {
		return errors.New("invalid context type for DingtalkAlert")
	}
	for idx, wa := range dingtalkContext.Message.Alerts {
		if d.Limit(idx, dingtalkContext.Message.Alerts) {
			return errors.New("alert is limited")
		}
		severity := wa.Labels["severity"]
		ss := strings.Split(severity, ":")
		status := ss[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(common.AlertFormat,
			dingtalkContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msgtype": "text",
			"text": map[string]string{
				"content": message + "\n",
			},
			"at": map[string]any{
				"atMobiles": users,
				"isAtAll":   false,
			},
		}
		req, _ := json.Marshal(data)
		body, err := d.Client.JSON().Post(dingtalkContext.Webhook, req)
		if err != nil {
			slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
			return err
		}
		dingTalkRes := &dingTalkRes{}
		if err = json.Unmarshal(body, &dingTalkRes); err != nil {
			slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
			return err
		}
		if dingTalkRes.Errcode != 0 {
			slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", dingTalkRes))
			return errors.New(dingTalkRes.Errmsg)
		}
	}
	if err := d.After(c); err != nil {
		return common.ErrorAfterHook
	}
	return nil
}

type dingTalkRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (d *DingtalkAlert) ExecuteTest(c common.AlertContext) error {
	dingtalkContext, ok := c.(*DingtalkAlertContext)
	if !ok {
		return errors.New("invalid context type for DingtalkAlert")
	}
	data := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("【%s】钉钉告警测试：收到本消息说明钉钉告警配置成功！", dingtalkContext.Title),
		},
	}

	slog.Debug("Dingtalk is ExecuteTest", slog.String("data", "告警测试邮件，请勿回复，谢谢"), slog.Any("Context", dingtalkContext))
	req, _ := json.Marshal(data)
	body, err := d.Client.JSON().Post(dingtalkContext.Webhook, req)
	if err != nil {
		slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
		return err
	}
	dingTalkRes := &dingTalkRes{}
	if err = json.Unmarshal(body, &dingTalkRes); err != nil {
		slog.Error("DingtalkAlert, Unmarshal response is failed", slog.Any("err", err))
		return err
	}
	if dingTalkRes.Errcode != 0 {
		slog.Error("DingtalkAlert, Send data to LarkAlert failed, ", slog.Any("err", dingTalkRes))
		return errors.New(dingTalkRes.Errmsg)
	}
	return nil
}
