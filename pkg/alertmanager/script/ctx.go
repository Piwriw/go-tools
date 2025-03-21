package script

import (
	"encoding/json"

	"github.piwriw.go-tools/pkg/alertmanager/common"
)

// ScriptAlertContext 脚本告警上下文
type ScriptAlertContext struct {
	Message        common.CustomMsg
	Title          string `json:"title"`
	ScriptReceiver string `json:"script_receiver"`
	AlertID        string `json:"alert_id"`
	Script         string `json:"script"`
	Args           Arg    `json:"args"`
}

var _ = common.AlertContext(&ScriptAlertContext{})

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
	return common.MethodScript
}
