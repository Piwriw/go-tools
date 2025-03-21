package dingtalk

import (
	"github.piwriw.go-tools/pkg/alertmanager/common"
)

// DingtalkAlertContext 钉钉告警上下文
type DingtalkAlertContext struct {
	Message common.CustomMsg
	Title   string `json:"title"`
	Webhook string `json:"webhook"` // 钉钉配置
}

var _ = common.AlertContext(&DingtalkAlertContext{})

func (d *DingtalkAlertContext) GetMethod() string {
	return common.MethodDingtalk
}
