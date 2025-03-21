package script

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/alertmanager/template"
	"github.piwriw.go-tools/pkg/alertmanager/common"
	"github.piwriw.go-tools/pkg/http"
)

// ScriptAlert 脚本告警策略
type ScriptAlert struct {
	Client     *http.Client
	LimitFunc  common.LimitOption
	BeforeHook common.HookFunc
	AfterHook  common.HookFunc
}

var _ = common.AlertMethod(&ScriptAlert{})

func (s *ScriptAlert) Before(c common.AlertContext) error {
	if s.BeforeHook != nil {
		if err := s.BeforeHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScriptAlert) After(c common.AlertContext) error {
	if s.AfterHook != nil {
		if err := s.AfterHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScriptAlert) Limit(idx int, alerts template.Alerts) bool {
	if s.LimitFunc != nil {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		if s.LimitFunc.Allow() {
			return true
		}
		return false
	}
	return true
}

func (s *ScriptAlert) Execute(c common.AlertContext, users []string) error {
	if s.Before(c) != nil {
		return common.ErrorBeforeHook
	}
	scriptContext, ok := c.(*ScriptAlertContext)
	if !ok {
		return errors.New("invalid context type for ScriptAlert")
	}
	// 执行脚本
	var errArr error
	for _, touser := range users {
		for idx, alert := range scriptContext.Message.Alerts {
			if s.Limit(idx, scriptContext.Message.Alerts) {
				return errors.New("alert is limited")
			}
			severity := alert.Labels["severity"]
			ss := strings.Split(severity, ":")
			if ss[0] == "严重" || ss[0] == "致命" {
				status := "告警"
				if alert.Status == "resolved" {
					status = "恢复"
				}
				message := fmt.Sprintf("[%s] [%s]%v", scriptContext.Title, status, alert.Annotations["description"])
				// 构建 args 参数数组
				args := []string{
					"--user", touser,
					"--msg", message,
					"--alert", scriptContext.GetAlertString(alert.Fingerprint),
					"--args", scriptContext.Args.String(),
				}

				// 使用变参形式传递参数
				cmd := exec.Command(scriptContext.Script, args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					slog.Error("Script(), exec cmd is failed", slog.Any("err", err))
					errArr = errors.Join(errArr, err)
				}
			}
		}
	}
	if s.After(c) != nil {
		return common.ErrorAfterHook
	}
	return nil
}

func (s *ScriptAlert) ExecuteTest(c common.AlertContext) error {
	scriptContext, ok := c.(*ScriptAlertContext)
	if !ok {
		return errors.New("invalid context type for ScriptAlert")
	}
	message := fmt.Sprintf("[%s] [%s]%v", scriptContext.Title, "测试", "发送成功")
	cmd := exec.Command(scriptContext.Script, scriptContext.ScriptReceiver, message)
	slog.Debug("ScriptAlert INFO", slog.Any("message", message))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		slog.Error("Script(), exec cmd is failed", slog.Any("err", err))
		return err
	}
	return nil
}
