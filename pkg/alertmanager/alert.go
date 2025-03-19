package alertmanager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/alertmanager/template"
	"gopkg.in/gomail.v2"

	"github.piwriw.go-tools/pkg/http"
)

// AlertMethod 告警策略接口
type AlertMethod interface {
	Execute(c AlertContext, users []string) error
	ExecuteTest(c AlertContext) error
}

// EmailAlert 邮件告警策略
type EmailAlert struct {
	limitFunc LimitOption
}

func (e *EmailAlert) Execute(c AlertContext, users []string) error {
	emailContext, ok := c.(*EmailAlertContext)
	if !ok {
		return errors.New("invalid context type for EmailAlert")
	}
	var errArr error
	// 执行邮件发送操作
	for _, user := range users {
		var message string
		for idx, wa := range emailContext.Message.Alerts {
			if e.limitFunc != nil && shouldLimit(e.limitFunc.Allow, idx, emailContext.Message.Alerts) {
				return errors.New("alert is limited")
			}
			severity := wa.Labels["severity"]
			ss := strings.Split(severity, ":")
			alertStatus := AlertStatusFiring
			alertLevel := ss[0]
			if wa.Status == "resolved" {
				alertStatus = AlertStatusResolved
			}
			message += fmt.Sprintf(emailTableTemplate,
				wa.Labels["alertname"], wa.Annotations["summary"], wa.Labels["name"], wa.Labels["ip"], wa.StartsAt.Format("2006-01-02 15:04:05"), alertLevel, alertStatus, wa.Annotations["description"], wa.Annotations["suggestion"])
		}
		htmlEmail := emailTemplateHeader + message + "</div>"
		m := gomail.NewMessage()
		m.SetHeader("From", emailContext.SMTPFrom)
		m.SetHeader("To", user)
		m.SetHeader("Subject", fmt.Sprintf("【%s】告警邮件", emailContext.Title))
		m.SetBody("text/html", htmlEmail)

		d, err := checkAndNewDialer(emailContext.SMTPSmarthost, emailContext.SMTPAuthUsername, emailContext.SMTPAuthPassword)
		if err != nil {
			slog.Error("checkAndNewDialer is failed", slog.Any("err", err))
			errArr = errors.Join(errArr, err)
			continue
		}

		if emailContext.SMTPRequireTLS {
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
		if err := d.DialAndSend(m); err != nil {
			slog.Error("Send Email is failed ", slog.Any("err", err))
			errArr = errors.Join(errArr, err)
			continue
		}
	}
	if errArr != nil {
		return errArr
	}
	return nil
}

// checkAndNewDialer TODO 这是一个原来的func 但是好像需要优化
func checkAndNewDialer(host string, username string, password string) (*gomail.Dialer, error) {
	var d *gomail.Dialer

	hp := strings.Split(host, ":")
	switch len(hp) {
	case 1:
		d = gomail.NewDialer(hp[0], SMTPDefaultPort, username, password)
		return d, nil
	case 2:
		port, err := strconv.Atoi(hp[1])
		if err != nil {
			err := errors.New(fmt.Sprintf("invalid SMTP `smtp_smarthost` port:%s", hp[1]))
			slog.Error("invalid mail format error", slog.Any("SMTPSmarthost", host), slog.Any("err", err))
			return nil, err
		}
		d = gomail.NewDialer(hp[0], port, username, password)
		return d, err
	default:
		err := errors.New("invalid SMTP `smtp_smarthost` format")
		slog.Error("invalid mail format error", slog.Any("SMTPSmarthost", host), slog.Any("err", err))
		return nil, err
	}
}

func (e *EmailAlert) ExecuteTest(c AlertContext) error {
	emailContext, ok := c.(*EmailAlertContext)
	if !ok {
		return errors.New("invalid context type for WechatAlert")
	}
	slog.Debug("Email is ExecuteTest", slog.String("data", "告警测试邮件，请勿回复，谢谢"), slog.Any("Context", emailContext))

	m := gomail.NewMessage()
	m.SetHeader("From", emailContext.SMTPFrom)
	m.SetHeader("To", emailContext.SendTo)
	m.SetHeader("Subject", fmt.Sprintf("【%s】告警测试邮件：收到本邮件表示邮件告警配置成功！", emailContext.Title))
	m.SetBody("text/html", "告警测试邮件，请勿回复，谢谢。")

	d, err := checkAndNewDialer(emailContext.SMTPSmarthost, emailContext.SMTPAuthUsername, emailContext.SMTPAuthPassword)
	if err != nil {
		slog.Error("checkAndNewDialer is failed", slog.Any("error", err))
		return err
	}
	if emailContext.SMTPRequireTLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// WechatAlert 微信告警策略
type WechatAlert struct {
	client    *http.Client
	limitFunc LimitOption
}

type wechatRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (w *WechatAlert) Execute(c AlertContext, users []string) error {
	wechatContext, ok := c.(*WechatAlertContext)
	if !ok {
		return errors.New("invalid context type for WechatAlert")
	}
	for idx, wa := range wechatContext.Message.Alerts {
		if w.limitFunc != nil && shouldLimit(w.limitFunc.Allow, idx, wechatContext.Message.Alerts) {
			return errors.New("alert is limited")
		}
		severity := wa.Labels["severity"]
		severities := strings.Split(severity, ":")
		status := severities[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}

		message := fmt.Sprintf(AlertFormat,
			wechatContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msgtype": "text",
			"text": map[string]any{
				"content":               message,
				"mentioned_mobile_list": users,
			},
		}

		req, _ := json.Marshal(data)
		body, err := w.client.JSON().Post(wechatContext.WechatWebhook, req)
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
func (w *WechatAlert) ExecuteTest(c AlertContext) error {
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
	body, err := w.client.JSON().Post(wechatContext.WechatWebhook, req)
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

// DingtalkAlert 钉钉告警策略
type DingtalkAlert struct {
	client    *http.Client
	hook      *HookHandler
	limitFunc LimitOption
}

func NewDingtalkAlert() *DingtalkAlert {
	return &DingtalkAlert{}
}

func (d *DingtalkAlert) WithLimit(limit LimitOption) *DingtalkAlert {
	d.limitFunc = limit
	return d
}

func (d *DingtalkAlert) WithHook(hook *HookHandler) *DingtalkAlert {
	d.hook = hook
	return d
}

func shouldLimit(limitFunc func() bool, idx int, alerts template.Alerts) bool {
	if limitFunc != nil && !limitFunc() {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		return true
	}
	return false
}

func (d *DingtalkAlert) Execute(c AlertContext, users []string) error {
	if callBeforeHook(d.hook, c) != nil {
		return ErrorBeforeHook
	}
	dingtalkContext, ok := c.(*DingtalkAlertContext)
	if !ok {
		return errors.New("invalid context type for DingtalkAlert")
	}
	for idx, wa := range dingtalkContext.Message.Alerts {
		if d.limitFunc != nil && shouldLimit(d.limitFunc.Allow, idx, dingtalkContext.Message.Alerts) {
			return errors.New("alert is limited")
		}
		severity := wa.Labels["severity"]
		ss := strings.Split(severity, ":")
		status := ss[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(AlertFormat,
			dingtalkContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msgtype": "text",
			"text": map[string]string{
				"content": message + "\n",
			},
			"at": map[string]any{
				"atMobiles": users,
				"isAtAll":   false,
			},
		}
		req, _ := json.Marshal(data)
		body, err := d.client.JSON().Post(dingtalkContext.Webhook, req)
		if err != nil {
			slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
			return err
		}
		dingTalkRes := &dingTalkRes{}
		if err = json.Unmarshal(body, &dingTalkRes); err != nil {
			slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
			return err
		}
		if dingTalkRes.Errcode != 0 {
			slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", dingTalkRes))
			return errors.New(dingTalkRes.Errmsg)
		}
	}
	if callAfterHook(d.hook, c) != nil {
		return ErrorAfterHook
	}
	return nil
}

type dingTalkRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (d *DingtalkAlert) ExecuteTest(c AlertContext) error {
	dingtalkContext, ok := c.(*DingtalkAlertContext)
	if !ok {
		return errors.New("invalid context type for DingtalkAlert")
	}
	data := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("【%s】钉钉告警测试：收到本消息说明钉钉告警配置成功！", dingtalkContext.Title),
		},
	}

	slog.Debug("Dingtalk is ExecuteTest", slog.String("data", "告警测试邮件，请勿回复，谢谢"), slog.Any("Context", dingtalkContext))
	req, _ := json.Marshal(data)
	body, err := d.client.JSON().Post(dingtalkContext.Webhook, req)
	if err != nil {
		slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
		return err
	}
	dingTalkRes := &dingTalkRes{}
	if err = json.Unmarshal(body, &dingTalkRes); err != nil {
		slog.Error("DingtalkAlert, Unmarshal response is failed", slog.Any("err", err))
		return err
	}
	if dingTalkRes.Errcode != 0 {
		slog.Error("DingtalkAlert, Send data to LarkAlert failed, ", slog.Any("err", dingTalkRes))
		return errors.New(dingTalkRes.Errmsg)
	}
	return nil
}

// ScriptAlert 脚本告警策略
type ScriptAlert struct {
	client    *http.Client
	limitFunc LimitOption
	hook      *HookHandler
}

func (s *ScriptAlert) Execute(c AlertContext, users []string) error {
	if callBeforeHook(s.hook, c) != nil {
		return ErrorBeforeHook
	}
	scriptContext, ok := c.(*ScriptAlertContext)
	if !ok {
		return errors.New("invalid context type for ScriptAlert")
	}
	// 执行脚本
	var errArr error
	for _, touser := range users {
		for idx, alert := range scriptContext.Message.Alerts {
			if s.limitFunc != nil && shouldLimit(s.limitFunc.Allow, idx, scriptContext.Message.Alerts) {
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
	if callAfterHook(s.hook, c) != nil {
		return ErrorAfterHook
	}
	return nil
}

func (s *ScriptAlert) ExecuteTest(c AlertContext) error {
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

// LarkAlert  Lark告警策略
type LarkAlert struct {
	client    *http.Client
	limitFunc LimitOption
}

const larkTemplate = `<at user_id="%s">Tom</at>`

func getLarkContent(userIDS []string, context string) string {
	res := ""
	for _, userID := range userIDS {
		msg := fmt.Sprintf(larkTemplate, userID)
		res += msg
	}
	if res == "" {
		return context
	}
	return res + "\n" + context
}

func (l *LarkAlert) Execute(c AlertContext, users []string) error {
	larkContext, ok := c.(*LarkAlertContext)
	if !ok {
		return errors.New("invalid context type for LarkAlert")
	}
	// 执行 Lark 发送操作
	slog.Debug("Lark is ExecuteTest", slog.Any("sendToUser", users), slog.Any("Context", larkContext))

	for idx, wa := range larkContext.Message.Alerts {
		if l.limitFunc != nil && shouldLimit(l.limitFunc.Allow, idx, larkContext.Message.Alerts) {
			return errors.New("alert is limited")
		}
		severity := wa.Labels["severity"]
		severities := strings.Split(severity, ":")
		status := severities[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(AlertFormat,
			larkContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msg_type": "text",
			"content": map[string]any{
				"text": getLarkContent(users, message),
			},
		}
		req, _ := json.Marshal(data)
		body, err := l.client.JSON().Post(larkContext.Webhook, req)
		if err != nil {
			slog.Error("Lark, Send data to Lark failed, ", slog.Any("err", err))
			return err
		}
		larkRes := &larkRes{}
		if err = json.Unmarshal(body, &larkRes); err != nil {
			slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
			return err
		}
		if larkRes.Code != 0 {
			slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", larkRes))
			return errors.New(larkRes.Msg)
		}
	}
	return nil
}

type larkRes struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func (l *LarkAlert) ExecuteTest(c AlertContext) error {
	larkAlertContext, ok := c.(*LarkAlertContext)
	if !ok {
		return errors.New("invalid context type for LarkAlert")
	}
	data := map[string]any{
		"msg_type": "text",
		"content": map[string]any{
			"text": fmt.Sprintf("【%s】飞书告警测试：收到本消息说明飞书告警配置成功！", larkAlertContext.Title),
		},
	}

	slog.Debug("LarkAlert is ExecuteTest", slog.String("data", "飞书告警测试，请勿回复，谢谢"), slog.Any("Context", larkAlertContext))
	req, _ := json.Marshal(data)
	body, err := l.client.JSON().Post(larkAlertContext.Webhook, req)
	if err != nil {
		slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", err))
		return err
	}
	larkRes := &larkRes{}
	if err = json.Unmarshal(body, &larkRes); err != nil {
		slog.Error("LarkAlert, Unmarshal response is failed", slog.Any("err", err))
		return err
	}
	if larkRes.Code != 0 {
		slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", larkRes))
		return errors.New(larkRes.Msg)
	}
	return nil
}

type AlertStrategyManager struct {
	strategies map[string]AlertMethod
	limitFunc  LimitOption
	hookFunc   LimitOption
}

// registerStrategies 统一注册所有告警策略
func (manager *AlertStrategyManager) registerStrategies(client *http.Client) {
	manager.strategies[MethodEmail] = &EmailAlert{}
	manager.strategies[MethodWechat] = &WechatAlert{client: client}
	manager.strategies[MethodDingtalk] = &DingtalkAlert{client: client}
	manager.strategies[MethodScript] = &ScriptAlert{client: client}
	manager.strategies[MethodLark] = &LarkAlert{client: client}
}

// RegisterStrategy 统一注册所有告警策略
func (manager *AlertStrategyManager) RegisterStrategy(methodName string, alertMethod AlertMethod) error {
	_, ok := manager.strategies[methodName]
	// Same Method Name also can register
	manager.strategies[methodName] = alertMethod
	if ok {
		slog.Warn("registerStrategy is same name", slog.Any("methodName", methodName))
		return ErrorSameMethod
	}
	return nil
}

func NewAlertStrategyManager(options ...Option) *AlertStrategyManager {
	manager := &AlertStrategyManager{
		strategies: make(map[string]AlertMethod),
	}

	client := http.NewHTTPClient()
	// 注册所有告警策略
	manager.registerStrategies(client)
	for _, option := range options {
		option(manager)
	}
	if manager.limitFunc != nil {
		manager.strategies[MethodEmail] = &EmailAlert{
			limitFunc: manager.limitFunc,
		}
		manager.strategies[MethodWechat] = &WechatAlert{
			client:    client,
			limitFunc: manager.limitFunc,
		}
		manager.strategies[MethodDingtalk] = &DingtalkAlert{
			client:    client,
			limitFunc: manager.limitFunc,
		}
		manager.strategies[MethodScript] = &ScriptAlert{client: client,
			limitFunc: manager.limitFunc}
		manager.strategies[MethodLark] = &LarkAlert{client: client, limitFunc: manager.limitFunc}
	}
	return manager
}

// Execute 根据告警方法执行相应的策略
func (manager *AlertStrategyManager) Execute(method string, c AlertContext, users []string) error {
	strategy, exists := manager.strategies[method]
	if !exists {
		return errors.New(fmt.Sprintf("method %s not found", method))
	}

	// 执行策略
	return strategy.Execute(c, users)
}

func (manager *AlertStrategyManager) ExecuteTest(method string, c AlertContext) error {
	strategy, exists := manager.strategies[method]
	if !exists {
		return errors.New(fmt.Sprintf("method %s not found", method))
	}

	// 执行策略
	return strategy.ExecuteTest(c)
}
