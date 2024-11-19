package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
	"github.piwriw.alertmanager/webhook/utils"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogLevel      = "info"
	defaultLogFilePath   = "sms.log"
	defaultLogMaxSize    = 100
	defaultLogMaxBackups = 5
	defaultLogMaxAge     = 30
	defaultLogCompress   = true
)

var CFG = new(AppConfig)

type AppConfig struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	SMS     SMS    `json:"SMS"`
	Log     Logger `json:"Log"`
}

type Logger struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func DefaultLogger() *Logger {
	return &Logger{
		Level:      defaultLogLevel,
		FilePath:   defaultLogFilePath,
		MaxSize:    defaultLogMaxSize,
		MaxBackups: defaultLogMaxBackups,
		MaxAge:     defaultLogMaxAge,
		Compress:   defaultLogCompress,
	}
}

type SMS struct {
	Address           map[string]string `json:"address"`
	SysCode           string            `json:"sysCode"`
	SysName           string            `json:"sysName"`
	EventType         string            `json:"eventType"`
	EventLevel        []string          `json:"eventLevel"`
	EventRecoveryType string            `json:"eventRecoveryType"`
}

func (s *SMS) getAddress(eventType SMSEventType) string {
	switch eventType {
	case AddType:
		return s.getAddAddress()
	case ResolvedType:
		return s.getRecoverAddress()
	default:
		return ""
	}
}
func (s *SMS) getRecoverAddress() string {
	if s.Address == nil {
		return ""
	}
	address, ok := s.Address["recover"]
	if !ok {
		return ""
	}
	return address
}

func (s *SMS) getAddAddress() string {
	if s.Address == nil {
		return ""
	}
	address, ok := s.Address["add"]
	if !ok {
		return ""
	}
	return address
}

type BaseSMSReq struct {
}

func (s *BaseSMSReq) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *BaseSMSReq) DoHttpRequest(url string) error {
	marshal, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return utils.DoHTTPRequest(marshal, url, utils.DefaultHTTPClientConfig())
}

type SMSRecover struct {
	BaseSMSReq
	Address           string    `json:"-"`
	EventID           string    `json:"eventid" validate:"required"`
	EventRecoveryType string    `json:"eventRecoveryType" `
	EventRecoveryTime time.Time `json:"eventRecoveryTime" validate:"required"`
}

func NewSMSRecover(eventID string, eventRecoveryType string, eventRecoveryTime time.Time) *SMSRecover {
	return &SMSRecover{
		EventID:           eventID,
		EventRecoveryType: eventRecoveryType,
		EventRecoveryTime: eventRecoveryTime,
	}
}

func (s *SMSRecover) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *SMSRecover) DoHttpRequest(url string) error {
	marshal, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return utils.DoHTTPRequest(marshal, url, utils.DefaultHTTPClientConfig())
}

type SMSReq struct {
	BaseSMSReq
	Address    string `json:"-"`
	SysCode    string `json:"sysCode"`
	SysName    string `json:"sysName"`
	EventType  string `json:"eventType"`
	EventLevel string `json:"eventLevel"`
	// 系统负责人
	ProjectManager string `json:"projectManager"`
	// DB-AGILEX+年月日时分秒+序列号
	EventID        string    `json:"eventid" validate:"required"`
	EventSummary   string    `json:"eventSummary" validate:"required"`
	EventDetail    string    `json:"eventDetail"`
	EventDatetime  time.Time `json:"eventDatetime" validate:"required"`
	EventIP        string    `json:"eventIp"`
	AlertType      string    `json:"alertType"`
	AlertCellphone string    `json:"alertCellphone"`
	AlertEhr       string    `json:"alertEhr"`
}

func NewSMSReq(sms SMS, eventID, projectManager, eventSummary, eventDetail, eventIP, eventLevel, alertType, alertCellphone, alertEhr string, eventDateTime time.Time) (*SMSReq, error) {
	if eventID == "" {
		return nil, errors.New("eventID is empty")
	}
	if eventIP == "" {
		return nil, errors.New("eventIP is empty")
	}
	if projectManager == "" {
		return nil, errors.New("projectManager is empty")
	}
	s := &SMSReq{
		SysCode:        sms.SysCode,
		SysName:        sms.SysName,
		EventType:      sms.EventType,
		EventLevel:     eventLevel,
		ProjectManager: projectManager,
		EventSummary:   eventSummary,
		EventDetail:    eventDetail,
		EventDatetime:  eventDateTime,
		EventIP:        eventIP,
		AlertType:      alertType,
		AlertCellphone: alertCellphone,
		AlertEhr:       alertEhr,
	}
	s.generaEventID(eventID)
	return s, nil
}

func (s *SMSReq) DoHttpRequest(url string) error {
	marshal, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return utils.DoHTTPRequest(marshal, url, utils.DefaultHTTPClientConfig())
}

func (s *SMSReq) Validate() error {
	if s.SysCode == "" {
		return errors.New("SysCode is empty")
	}
	if s.SysName == "" {
		return errors.New("SysName is empty")
	}
	if s.ProjectManager == "" {
		return errors.New("ProjectManager is empty")
	}
	validate := validator.New()
	return validate.Struct(s)
}

func (s *SMSReq) generaEventID(eventID string) string {
	s.EventID = fmt.Sprintf("%s%s%s", s.SysCode, time.Now().Format("20060102150405"), eventID)
	return s.EventID
}

func InitConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("viper.ReadInConfig failed", slog.Any("err", err))
		return err
	}
	if err := viper.Unmarshal(CFG); err != nil {
		slog.Error("viper.Unmarshal failed", slog.Any("err", err))
		return err
	}
	return nil
}

type Arg struct {
	ConfigFile     string `json:"configFile"`
	ProjectManager string `json:"projectManager"`
	AlertType      string `json:"alertType"`
	AlertCellphone string `json:"alertCellphone"`
	AlertEhr       string `json:"alertEhr"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

func (a *Alert) getLevel() string {
	severity, ok := a.Labels["severity"]
	if !ok {
		return ""
	}
	if strings.Contains(severity, "致命") {
		return "04"
	} else if strings.Contains(severity, "严重") {
		return "03"
	}
	return ""
}

// getID 根据Fingerprint和satrtat 生成唯一ID
func (a *Alert) getID() string {
	return a.Fingerprint + a.StartsAt.Format("20060102150405")
}

func (a *Alert) getDetail() string {
	if a.Annotations == nil {
		return ""
	}
	return a.Annotations["summary"]
}

func (a *Alert) getSummary() string {
	if a.Annotations == nil {
		return ""
	}
	return a.Annotations["summary"]
}

func (a *Alert) getInstanceIP() string {
	if a.Labels == nil {
		return ""
	}
	return a.Labels["ip"]
}

// 定义命令行标志
var user string
var msg string
var argStr string
var alertStr string

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo})))
}

func InitLogger(filePath string, logLevel string, maxSize int, maxBackups int, maxAge int, isCompress bool) *lumberjack.Logger {
	// 设置日志文件轮转
	logFile := &lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件的最大尺寸（MB）
		MaxBackups: maxBackups, // 保留的旧日志文件数
		MaxAge:     maxAge,     // 日志文件保存的天数
		Compress:   isCompress, // 是否压缩/归档旧日志文件
	}

	// 设置控制台和文件输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var level slog.Level
	switch logLevel {
	case slog.LevelInfo.String(), strings.ToLower(slog.LevelInfo.String()):
		level = slog.LevelInfo
	case slog.LevelDebug.String(), strings.ToLower(slog.LevelDebug.String()):
		level = slog.LevelDebug
	case slog.LevelWarn.String(), strings.ToLower(slog.LevelWarn.String()):
		level = slog.LevelWarn
	case slog.LevelError.String(), strings.ToLower(slog.LevelError.String()):
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 创建自定义的 `slog.Handler`，将日志输出到 multiWriter
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{AddSource: true, Level: level})

	// 创建 logger 实例
	logger := slog.New(handler)
	slog.SetDefault(logger)
	// 在程序结束时关闭日志文件
	return logFile
}

func main() {
	flag.StringVar(&user, "user", "", "User name")
	flag.StringVar(&msg, "msg", "", "Message")
	flag.StringVar(&argStr, "args", "", "JSON string for additional parameters")
	flag.StringVar(&alertStr, "alert", "", "JSON string for additional parameters")
	// 解析命令行参数
	flag.Parse()

	// TODO：在此之前的Logger不会被写入文件，因为配置文件依赖于文件路径，写入默认配置
	logger := InitLogger(DefaultLogger().FilePath, DefaultLogger().Level, DefaultLogger().MaxSize, DefaultLogger().MaxBackups, DefaultLogger().MaxAge, DefaultLogger().Compress)
	defer logger.Close()
	slog.Info("cmd argStr info...", slog.String("user", user), slog.String("msg", msg), slog.Any("args", argStr), slog.Any("alert", alertStr))
	arg := &Arg{}
	if err := json.Unmarshal([]byte(argStr), arg); err != nil {
		slog.Error("Args json.Unmarshal is  failed", slog.Any("err", err))
		return
	}

	alert := &Alert{}
	if err := json.Unmarshal([]byte(alertStr), alert); err != nil {
		slog.Error("alert json.Unmarshal is  failed", slog.Any("err", err))
		return
	}

	if err := InitConfig(arg.ConfigFile); err != nil {
		slog.Error("InitConfig failed", slog.Any("err", err))
		return
	}
	logger = InitLogger(CFG.Log.FilePath, CFG.Log.Level, CFG.Log.MaxSize, CFG.Log.MaxBackups, CFG.Log.MaxAge, CFG.Log.Compress)
	defer logger.Close()
	msg = replaceTitle(msg, CFG.Name)
	eventType := StringToSMSEventType(alert.Status)
	smsService, err := NewSMSServiceInterface(eventType, &CFG.SMS, alert, arg)
	if err != nil {
		slog.Error("NewSMSServiceInterface is failed", slog.Any("err", err))
		return
	}

	if err = smsService.Validate(); err != nil {
		slog.Error("Validate is failed", slog.Any("err", err))
		return
	}

	if err = smsService.DoHttpRequest(CFG.SMS.getAddress(eventType)); err != nil {
		slog.Error("DoHttpRequest is failed", slog.Any("err", err))
		return
	}
}

func replaceTitle(msg string, app string) string {
	return strings.ReplaceAll(msg, "云趣科技", app)
}

type SMSSendInterface interface {
	DoHttpRequest(url string) error
	Validate() error
}

// NewSMSServiceInterface 创建并返回 SMSSendInterface 类型的实例
func NewSMSServiceInterface(smsType SMSEventType, sms *SMS, alert *Alert, arg *Arg) (SMSSendInterface, error) {
	switch smsType {
	case AddType:
		req, err := NewSMSReq(*sms, alert.getID(), arg.ProjectManager, alert.getSummary(), alert.getDetail(),
			alert.getInstanceIP(), alert.getLevel(), arg.AlertType, arg.AlertCellphone, arg.AlertEhr, alert.StartsAt)
		if err != nil {
			return nil, err
		}
		return req, nil
	case ResolvedType:
		req := NewSMSRecover(alert.getID(), sms.EventRecoveryType, alert.EndsAt)
		return req, nil
	default:
		return nil, errors.New("unknown smsType")
	}
}

type SMSEventType string

const (
	AddType      SMSEventType = "add"
	ResolvedType SMSEventType = "resolved"
)

func StringToSMSEventType(str string) SMSEventType {
	switch str {
	case "firing":
		return AddType
	case "resolved":
		return ResolvedType
	default:
		return ""
	}
}
