package wechat

import (
	"github.piwriw.go-tools/pkg/alertmanager/common"
)

// WechatAlertContext 微信告警上下文
type WechatAlertContext struct {
	Message common.CustomMsg
	Title   string `json:"title"`
	// todo 删除wechat,只保留Webhook 这是个历史一致问题
	WechatWebhook string `json:"wechat_webhook"`
}

func (w *WechatAlertContext) GetMethod() string {
	return common.MethodWechat
}
