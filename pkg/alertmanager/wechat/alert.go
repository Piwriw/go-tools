package wechat

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

// WechatAlert 微信告警策略
type WechatAlert struct {
	Client     *http.Client
	BeforeHook common.HookFunc
	AfterHook  common.HookFunc
	LimitFunc  common.LimitOption
}

type wechatRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (w *WechatAlert) WithLimit(limit common.LimitOption) *WechatAlert {
	if limit == nil {
		return w
	}
	w.LimitFunc = limit
	return w
}

func (w *WechatAlert) WithBeforeHook(hook common.HookFunc) *WechatAlert {
	if hook == nil {
		return w
	}
	w.BeforeHook = hook
	return w
}

func (w *WechatAlert) WithAfterHook(hook common.HookFunc) *WechatAlert {
	if hook == nil {
		return w
	}
	w.AfterHook = hook
	return w
}

func (w *WechatAlert) Before(c common.AlertContext) error {
	if w.BeforeHook != nil {
		if err := w.BeforeHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (w *WechatAlert) After(c common.AlertContext) error {
	if w.AfterHook != nil {
		if err := w.AfterHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (w *WechatAlert) Limit(idx int, alerts template.Alerts) bool {
	if w.LimitFunc != nil {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		if w.LimitFunc.Allow() {
			return true
		}
		return false
	}
	return true
}

func (w *WechatAlert) Execute(c common.AlertContext, users []string) error {
	wechatContext, ok := c.(*WechatAlertContext)
	if !ok {
		return errors.New("invalid context type for WechatAlert")
	}
	for idx, wa := range wechatContext.Message.Alerts {
		if w.Limit(idx, wechatContext.Message.Alerts) {
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
			wechatContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msgtype": "text",
			"text": map[string]any{
				"content":               message,
				"mentioned_mobile_list": users,
			},
		}

		req, _ := json.Marshal(data)
		body, err := w.Client.JSON().Post(wechatContext.WechatWebhook, req)
		if err != nil {
			slog.Error("Wechat, Send data to Wechat failed, ", slog.Any("err", err))
			return err
		}
		wechatRes := &wechatRes{}
		if err = json.Unmarshal(body, &wechatRes); err != nil {
			slog.Error("WechatAlert, Unmarshal response is failed", slog.Any("err", err))
			return err
		}
		if wechatRes.Errcode != 0 {
			slog.Error("WechatAlert, Send data to WechatAlert failed, ", slog.Any("err", wechatRes))
			return errors.New(wechatRes.Errmsg)
		}
	}
	return nil
}
func (w *WechatAlert) ExecuteTest(c common.AlertContext) error {
	wechatContext, ok := c.(*WechatAlertContext)
	if !ok {
		return errors.New("invalid context type for WechatAlert")
	}
	data := map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content": fmt.Sprintf("【%s】微信告警测试：收到本消息说明微信告警配置成功！", wechatContext.Title),
		},
	}
	slog.Debug("Wechat is ExecuteTest", slog.String("data", "告警测试邮件，请勿回复，谢谢"), slog.Any("Context", wechatContext))
	req, _ := json.Marshal(data)
	body, err := w.Client.JSON().Post(wechatContext.WechatWebhook, req)
	if err != nil {
		slog.Error("Wechat, Send data to Wechat failed, ", slog.Any("err", err))
		return err
	}
	wechatRes := &wechatRes{}
	if err = json.Unmarshal(body, &wechatRes); err != nil {
		slog.Error("WechatAlert, Unmarshal response is failed", slog.Any("err", err))
		return err
	}
	if wechatRes.Errcode != 0 {
		slog.Error("WechatAlert, Send data to WechatAlert failed, ", slog.Any("err", wechatRes))
		return errors.New(wechatRes.Errmsg)
	}
	return nil
}
