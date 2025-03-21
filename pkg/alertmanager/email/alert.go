package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.piwriw.go-tools/pkg/alertmanager/common"

	"github.com/prometheus/alertmanager/template"

	"gopkg.in/gomail.v2"
)

// EmailAlert 邮件告警策略
type EmailAlert struct {
	LimitFunc  common.LimitOption
	BeforeHook common.HookFunc
	AfterHook  common.HookFunc
}

var _ = common.AlertMethod(&EmailAlert{})

func (e *EmailAlert) WithLimit(limit common.LimitOption) *EmailAlert {
	if limit == nil {
		return e
	}
	e.LimitFunc = limit
	return e
}

func (e *EmailAlert) WithBeforeHook(hook common.HookFunc) *EmailAlert {
	if hook == nil {
		return e
	}
	e.BeforeHook = hook
	return e
}

func (e *EmailAlert) WithAfterHook(hook common.HookFunc) *EmailAlert {
	if hook == nil {
		return e
	}
	e.AfterHook = hook
	return e
}

func (e *EmailAlert) Before(c common.AlertContext) error {
	if e.BeforeHook != nil {
		if err := e.BeforeHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (e *EmailAlert) After(c common.AlertContext) error {
	if e.AfterHook != nil {
		if err := e.AfterHook(c); err != nil {
			return err
		}
	}
	return nil
}

func (e *EmailAlert) Limit(idx int, alerts template.Alerts) bool {
	if e.LimitFunc != nil {
		slog.Warn("Current Alert is Limited",
			slog.Int("remaining_alert_count", len(alerts[idx:])),
			slog.Any("remaining_alerts", alerts[idx:]),
		)
		if e.LimitFunc.Allow() {
			return true
		}
		return false
	}
	return true
}

func (e *EmailAlert) Execute(c common.AlertContext, users []string) error {
	emailContext, ok := c.(*EmailAlertContext)
	if !ok {
		return errors.New("invalid context type for EmailAlert")
	}
	var errArr error
	// 执行邮件发送操作
	for _, user := range users {
		var message string
		for idx, wa := range emailContext.Message.Alerts {
			if e.Limit(idx, emailContext.Message.Alerts) {
				return errors.New("alert is limited")
			}
			severity := wa.Labels["severity"]
			ss := strings.Split(severity, ":")
			alertStatus := common.AlertStatusFiring
			alertLevel := ss[0]
			if wa.Status == "resolved" {
				alertStatus = common.AlertStatusResolved
			}
			message += fmt.Sprintf(EmailTableTemplate,
				wa.Labels["alertname"], wa.Annotations["summary"], wa.Labels["name"], wa.Labels["ip"], wa.StartsAt.Format("2006-01-02 15:04:05"), alertLevel, alertStatus, wa.Annotations["description"], wa.Annotations["suggestion"])
		}
		htmlEmail := EmailTableTemplate + message + "</div>"
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

func (e *EmailAlert) ExecuteTest(c common.AlertContext) error {
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
