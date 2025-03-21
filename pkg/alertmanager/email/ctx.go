package email

import (
	"github.piwriw.go-tools/pkg/alertmanager/common"
)

// EmailAlertContext 邮件告警上下文
type EmailAlertContext struct {
	Message          common.CustomMsg
	Title            string
	SMTPSmarthost    string `json:"smtp_smarthost"`     // 邮件地址
	SMTPFrom         string `json:"smtp_from"`          // 邮件发送方
	SMTPAuthUsername string `json:"smtp_auth_username"` // 邮件用户
	SMTPAuthPassword string `json:"smtp_auth_password"` // 邮件密码
	SMTPRequireTLS   bool   `json:"smtp_require_tls"`   // 邮件是否使用tls
	SendTo           string `json:"send_to"`            // 发送给谁，测试使用字段
}

// 确保 EmailAlertContext 实现了 AlertContext 接口
var _ = common.AlertContext(&EmailAlertContext{})

func (e *EmailAlertContext) GetMethod() string {
	return common.MethodEmail
}
