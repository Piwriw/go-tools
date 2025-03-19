package alertmanager

import (
	"encoding/json"

	"github.com/prometheus/alertmanager/notify/webhook"
)

// AlertContext 告警上下文接口，不同策略有不同的上下文
type AlertContext interface {
	GetMethod() string
}

// EmailAlertContext 邮件告警上下文
type EmailAlertContext struct {
	Message          webhook.Message
	Title            string
	SMTPSmarthost    string `json:"smtp_smarthost"`     // 邮件地址
	SMTPFrom         string `json:"smtp_from"`          // 邮件发送方
	SMTPAuthUsername string `json:"smtp_auth_username"` // 邮件用户
	SMTPAuthPassword string `json:"smtp_auth_password"` // 邮件密码
	SMTPRequireTLS   bool   `json:"smtp_require_tls"`   // 邮件是否使用tls
	SendTo           string `json:"send_to"`            // 发送给谁，测试使用字段
}

// 确保 EmailAlertContext 实现了 AlertContext 接口
var _ = AlertContext(&EmailAlertContext{})

func (e *EmailAlertContext) GetMethod() string {
	return MethodEmail
}

// WechatAlertContext 微信告警上下文
type WechatAlertContext struct {
	Message webhook.Message
	Title   string `json:"title"`
	// todo 删除wechat,只保留Webhook 这是个历史一致问题
	WechatWebhook string `json:"wechat_webhook"`
}

func (w *WechatAlertContext) GetMethod() string {
	return MethodWechat
}

// DingtalkAlertContext 钉钉告警上下文
type DingtalkAlertContext struct {
	Message webhook.Message
	Title   string `json:"title"`
	Webhook string `json:"webhook"` // 钉钉配置
}

var _ = AlertContext(&DingtalkAlertContext{})

func (d *DingtalkAlertContext) GetMethod() string {
	return MethodDingtalk
}

// ScriptAlertContext 脚本告警上下文
type ScriptAlertContext struct {
	Message        webhook.Message
	Title          string `json:"title"`
	ScriptReceiver string `json:"script_receiver"`
	AlertID        string `json:"alert_id"`
	Script         string `json:"script"`
	Args           Arg    `json:"args"`
}

var _ = AlertContext(&ScriptAlertContext{})

func (s *ScriptAlertContext) GetAlertString(fingerprint string) string {
	for _, alert := range s.Message.Alerts {
		if fingerprint == alert.Fingerprint {
			marshal, err := json.Marshal(alert)
			if err != nil {
				return ""
			}
			return string(marshal)
		}
	}
	return ""
}

type Arg struct {
	ConfigFile string `json:"configFile"`
}

func (a *Arg) String() string {
	marshal, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func (s *ScriptAlertContext) GetMethod() string {
	return MethodScript
}

// LarkAlertContext Lark告警上下文
type LarkAlertContext struct {
	Message webhook.Message
	Title   string `json:"title"`
	Webhook string `json:"webhook"`
}

var _ = AlertContext(&LarkAlertContext{})

func (l *LarkAlertContext) GetMethod() string {
	return MethodLark
}
