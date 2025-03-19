package alertmanager

import (
	"github.com/prometheus/alertmanager/notify/webhook"
	"testing"
)

func TestName(t *testing.T) {

	manager := NewAlertStrategyManager()
	ctx := &LarkAlertContext{
		Message: webhook.Message{},
		Title:   "",
		Webhook: "",
	}
	manager.Execute(MethodDingtalk, ctx, nil)
}
