package lark

import (
	"github.piwriw.go-tools/pkg/alertmanager/common"
)

// LarkAlertContext Lark告警上下文
type LarkAlertContext struct {
	Message common.CustomMsg
	Title   string `json:"title"`
	Webhook string `json:"webhook"`
}

var _ = common.AlertContext(&LarkAlertContext{})

func (l *LarkAlertContext) GetMethod() string {
	return common.MethodLark
}
