package alertmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.piwriw.go-tools/pkg/alertmanager/common"
	"github.piwriw.go-tools/pkg/alertmanager/lark"
)

func TestName(t *testing.T) {

	manager := NewAlertStrategyManager()

	manager.RegisterStrategy(common.MethodLark, lark.NewLarkAlert().WithBeforeHook(func(c common.AlertContext) error {
		fmt.Println("before hook")
		return nil
	}))

	// 创建 DingtalkAlertContext
	ctx := &lark.LarkAlertContext{
		Message: common.CustomMsg{
			Message: webhook.Message{
				Data: &template.Data{
					Receiver: "1",      // 接收者
					Status:   "firing", // 状态
					Alerts: template.Alerts{
						{
							Labels: map[string]string{"_VALUE_": "value",
								"alertname": "GaussDB 长时间运行SQL告警",
								"category":  "database.gaussdb",
								"datname":   "postgres",
								"groupId":   "0",
								"instance":  "instance-1",
								"ip":        "10.0.0.195",
								"job":       "interval1m",
								"name":      "gb",
								"severity":  "警告:#FFC92D",
								"usefor":    "resource.database",
								"usename":   "omm"}, // 示例标签
							Annotations: map[string]string{
								"description": "GaussDB数据库gb 出现长时间执行SQL，运行时间为7101分钟",
								"suggestion":  "确认长事务类型，遗漏提交操作时需尽快提交，以免引起其他的阻塞；未执行结束需排查SQL逻辑或日志参数配置，尽量减少长事务",
								"summary":     "GaussDB数据库gb 出现长时间执行SQL，运行时间为7101分钟",
							},
							StartsAt:     time.Now(),                    // 当前时间作为开始时间
							EndsAt:       time.Now().Add(1 * time.Hour), // 结束时间：1小时后
							GeneratorURL: "http://example.com",          // 示例 URL
						},
					},
					GroupLabels: nil,
					CommonLabels: map[string]string{
						"alertname": "GaussDB 长时间运行SQL告警",
					},
					CommonAnnotations: nil,
					ExternalURL:       "http://example.com/alerts", // 外部 URL
				},
				Version:         "1.0",    // 消息版本
				GroupKey:        "group1", // 组键
				TruncatedAlerts: 1,        // 截断的警报数量
			},
		},
		Title:   "DingTalk Alert", // 标题
		Webhook: "https:  ",       // Webhook URL
	}

	if err := manager.Execute(common.MethodLark, ctx, nil); err != nil {
		t.Error(err)
	}
}
