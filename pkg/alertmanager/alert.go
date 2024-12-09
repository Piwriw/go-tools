package alertmanager

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"gopkg.in/gomail.v2"
	"yunqutech.com/hawkeye/athena/pkg/http"
)

// AlertMethod 告警策略接口
type AlertMethod interface {
	Execute(c AlertContext, users []string) error
	ExecuteTest(c AlertContext) error
}

// AlertContext 告警上下文接口，不同策略有不同的上下文
type AlertContext interface {
	GetMethod() string
}

type QCloudAlertContext struct {
	API    string `json:"api"`
	AppID  int    `json:"app_id"`
	AppKey string `json:"app_key"`
	Sign   string `json:"sign"`
	TplID  int    `json:"tpl_id"`
	Mobile int    `json:"mobile"`
}

func (q *QCloudAlertContext) GetMethod() string {
	return MethodQCloud
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

func (l *LarkAlertContext) GetMethod() string {
	return MethodLark
}

// EmailAlert 邮件告警策略
type EmailAlert struct {
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
		for _, wa := range emailContext.Message.Alerts {
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
	client *http.Client
}

func (w *WechatAlert) Execute(c AlertContext, users []string) error {
	wechatContext, ok := c.(*WechatAlertContext)
	if !ok {
		return errors.New("invalid context type for WechatAlert")
	}
	for _, wa := range wechatContext.Message.Alerts {
		severity := wa.Labels["severity"]
		severities := strings.Split(severity, ":")
		status := severities[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}

		message := fmt.Sprintf(AgilexFormat,
			wechatContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msgtype": "text",
			"text": map[string]any{
				"content":               message,
				"mentioned_mobile_list": users,
			},
		}

		req, _ := json.Marshal(data)
		_, err := w.client.JSON().Post(wechatContext.WechatWebhook, req)
		if err != nil {
			slog.Error("Wechat, Send data to Wechat failed, ", slog.Any("err", err))
			return err
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
	_, err := w.client.JSON().Post(wechatContext.WechatWebhook, req)
	if err != nil {
		slog.Error("Wechat, Send data to Wechat failed, ", slog.Any("err", err))
		return err
	}
	return nil
}

// DingtalkAlert 钉钉告警策略
type DingtalkAlert struct {
	client *http.Client
}

func (d *DingtalkAlert) Execute(c AlertContext, users []string) error {
	dingtalkContext, ok := c.(*DingtalkAlertContext)
	if !ok {
		return errors.New("invalid context type for DingtalkAlert")
	}
	for _, wa := range dingtalkContext.Message.Alerts {
		severity := wa.Labels["severity"]
		ss := strings.Split(severity, ":")
		status := ss[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(AgilexFormat,
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
		_, err := d.client.JSON().Post(dingtalkContext.Webhook, req)
		if err != nil {
			slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
			return err
		}
	}
	return nil
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
	_, err := d.client.JSON().Post(dingtalkContext.Webhook, req)
	if err != nil {
		slog.Error("Dingtalk, Send data to Dingtalk failed, ", slog.Any("err", err))
		return err
	}
	return nil
}

type QCloudAlert struct {
	client *http.Client
}

// Execute 没有找到相关代码
func (q *QCloudAlert) Execute(c AlertContext, users []string) error {
	return nil
}

func (q *QCloudAlert) ExecuteTest(c AlertContext) error {
	qcContext, ok := c.(*QCloudAlertContext)
	if !ok {
		return errors.New("invalid context type for QCloudAlert")
	}
	// 执行脚本
	ts := time.Now().Unix()
	s := fmt.Sprintf("appkey=%s&random=7226249334&time=%d&mobile=%d", qcContext.AppKey, ts, qcContext.Mobile)
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return err
	}
	sig := fmt.Sprintf("%x", h.Sum(nil))
	data := map[string]any{
		"ext":    "",
		"extend": "",
		"params": []string{},
		"sig":    sig,
		"sign":   qcContext.Sign,
		"tel": map[string]string{
			"mobile":     strconv.Itoa(qcContext.Mobile),
			"nationcode": "86",
		},
		"time":   ts,
		"tpl_id": qcContext.TplID,
	}
	req, _ := json.Marshal(data)
	url := fmt.Sprintf("%s?sdkappid=%d&random=7226249334", qcContext.API, qcContext.AppID)
	_, err = q.client.JSON().Post(url, req)
	if err != nil {
		slog.Error("QCloud, Send data to QCloud failed, ", slog.Any("err", err))
		return err
	}
	return nil
}

// ScriptAlert 脚本告警策略
type ScriptAlert struct {
	client *http.Client
}

func (s *ScriptAlert) Execute(c AlertContext, users []string) error {
	scriptContext, ok := c.(*ScriptAlertContext)
	if !ok {
		return errors.New("invalid context type for ScriptAlert")
	}
	// 执行脚本
	var errArr error
	for _, touser := range users {
		for _, alert := range scriptContext.Message.Alerts {
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
	client *http.Client
}

const larkTemplate = `<at user_id="%s">Tom</at>`

func getLarkContent(userIDS []string, context string) string {
	res := ""
	for _, userID := range userIDS {
		msg := fmt.Sprintf(larkTemplate, userID)
		res += msg
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

	for _, wa := range larkContext.Message.Alerts {
		severity := wa.Labels["severity"]
		severities := strings.Split(severity, ":")
		status := severities[0]
		typ := strings.Split(wa.Labels["category"], ".")[1]
		if wa.Status == "resolved" {
			status = "恢复"
		}
		message := fmt.Sprintf(AgilexFormat,
			larkContext.Title, typ, wa.Labels["name"], wa.Labels["ip"], wa.Labels["alertname"], status, wa.Annotations["description"], wa.StartsAt.Format("2006-01-02 15:04:05"))
		data := map[string]any{
			"msg_type": "text",
			"content": map[string]any{
				"text": getLarkContent(users, message),
			},
		}
		req, _ := json.Marshal(data)
		_, err := l.client.JSON().Post(larkContext.Webhook, req)
		if err != nil {
			slog.Error("Lark, Send data to Lark failed, ", slog.Any("err", err))
			return err
		}
	}
	return nil
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
	_, err := l.client.JSON().Post(larkAlertContext.Webhook, req)
	if err != nil {
		slog.Error("LarkAlert, Send data to LarkAlert failed, ", slog.Any("err", err))
		return err
	}
	return nil
}

type AlertStrategyManager struct {
	strategies map[string]AlertMethod
}

func NewAlertStrategyManager() *AlertStrategyManager {
	manager := &AlertStrategyManager{
		strategies: make(map[string]AlertMethod),
	}
	client := http.NewHTTPClient()
	// 注册所有告警策略
	manager.strategies[MethodEmail] = &EmailAlert{}
	manager.strategies[MethodWechat] = &WechatAlert{
		client: client,
	}
	manager.strategies[MethodDingtalk] = &DingtalkAlert{
		client: client,
	}
	manager.strategies[MethodScript] = &ScriptAlert{client: client}
	manager.strategies[MethodLark] = &LarkAlert{client: client}
	manager.strategies[MethodQCloud] = &QCloudAlert{client: client}

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
