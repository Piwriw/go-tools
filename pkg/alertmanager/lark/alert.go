package lark

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

// LarkAlert  Lark告警策略
type LarkAlert struct {
	Client     *http.Client
	BeforeHook common.HookFunc
	AfterHook  common.HookFunc
	LimitFunc  common.LimitOption
}

var _ = common.AlertMethod(&LarkAlert{})

func NewLarkAlert() *LarkAlert {
	return &LarkAlert{Client: http.NewHTTPClient()}
}

func NewLarkAlertWithClient(client *http.Client) *LarkAlert {
	return &LarkAlert{Client: client}
}

func (l *LarkAlert) WithLimit(limit common.LimitOption) *LarkAlert {
	if limit == nil {
		return l
	}
	l.LimitFunc = limit
	return l
}

func (l *LarkAlert) WithBeforeHook(hook common.HookFunc) *LarkAlert {
	if hook == nil {
		return l
	}
	l.BeforeHook = hook
	return l
}

func (l *LarkAlert) WithAfterHook(hook common.HookFunc) *LarkAlert {
	if hook == nil {
		return l
	}
	l.AfterHook = hook
	return l
}

func (l *LarkAlert) Before(c common.AlertContext) error {
	if l.BeforeHook != nil {
		if err := l.BeforeHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (l *LarkAlert) After(c common.AlertContext) error {
	if l.AfterHook != nil {
		if err := l.AfterHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (l *LarkAlert) Limit(idx int, alerts template.Alerts) bool {
	if l.LimitFunc != nil {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		if !l.LimitFunc.Allow() {
			return true
		}
		return false
	}
	return false
}

func (l *LarkAlert) Execute(c common.AlertContext, users []string) error {
	if err := l.Before(c); err != nil {
		return common.ErrorBeforeHook
	}
	larkContext, ok := c.(*LarkAlertContext)
	if !ok {
		return errors.New("invalid context type for LarkAlert")
	}
	// 执行 Lark 发送操作
	slog.Debug("Lark is ExecuteTest", slog.Any("sendToUser", users), slog.Any("Context", larkContext))

	for idx, wa := range larkContext.Message.Alerts {
		if l.Limit(idx, larkContext.Message.Alerts) {
			return errors.New("alert is limited")
		}
		severity := wa.Labels["severity"]
		severities := strings.Split(severity, ":")
		status := severities[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(common.AlertFormat,
			larkContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msg_type": "text",
			"content": map[string]any{
				"text": getLarkContent(users, message),
			},
		}
		req, _ := json.Marshal(data)
		body, err := l.Client.JSON().Post(larkContext.Webhook, req)
		if err != nil {
			slog.Error("Lark, Send data to Lark failed, ", slog.Any("err", err))
			return err
		}
		larkRes := &larkRes{}
		if err = json.Unmarshal(body, &larkRes); err != nil {
			slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
			return err
		}
		if larkRes.Code != 0 {
			slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", larkRes))
			return errors.New(larkRes.Msg)
		}
	}
	if err := l.After(c); err != nil {
		return common.ErrorAfterHook
	}
	return nil
}

type larkRes struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func (l *LarkAlert) ExecuteTest(c common.AlertContext) error {
	larkAlertContext, ok := c.(*LarkAlertContext)
	if !ok {
		return errors.New("invalid context type for LarkAlert")
	}
	data := map[string]any{
		"msg_type": "text",
		"content": map[string]any{
			"text": fmt.Sprintf("【%s】飞书告警测试：收到本消息说明飞书告警配置成功！", larkAlertContext.Title),
		},
	}

	slog.Debug("LarkAlert is ExecuteTest", slog.String("data", "飞书告警测试，请勿回复，谢谢"), slog.Any("Context", larkAlertContext))
	req, _ := json.Marshal(data)
	body, err := l.Client.JSON().Post(larkAlertContext.Webhook, req)
	if err != nil {
		slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", err))
		return err
	}
	larkRes := &larkRes{}
	if err = json.Unmarshal(body, &larkRes); err != nil {
		slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
		return err
	}
	if larkRes.Code != 0 {
		slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", larkRes))
		return errors.New(larkRes.Msg)
	}
	return nil
}
