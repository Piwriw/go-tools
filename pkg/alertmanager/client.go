package alertmanager

import (
	"errors"
	"fmt"
	"log/slog"

	"github.piwriw.go-tools/pkg/alertmanager/common"
	"github.piwriw.go-tools/pkg/alertmanager/dingtalk"
	"github.piwriw.go-tools/pkg/alertmanager/email"
	"github.piwriw.go-tools/pkg/alertmanager/lark"
	"github.piwriw.go-tools/pkg/alertmanager/script"
	"github.piwriw.go-tools/pkg/alertmanager/wechat"
	"github.piwriw.go-tools/pkg/http"
)

type Option func(*AlertStrategyManager)

type AlertStrategyManager struct {
	strategies map[string]common.AlertMethod
	LimitFunc  common.LimitOption
	beforeHook common.HookFunc
	afterHook  common.HookFunc
}

// WithGlobalHook 设置全局 Hook
func WithGlobalHook(beforeHook, afterHook common.HookFunc) Option {
	return func(manager *AlertStrategyManager) {
		manager.beforeHook = beforeHook
		manager.afterHook = afterHook
	}
}

// WithGlobalLimit 设置全局限流
func WithGlobalLimit(limitFunc common.LimitOption) Option {
	return func(manager *AlertStrategyManager) {
		manager.LimitFunc = limitFunc
	}
}

// registerStrategies 统一注册所有告警策略
func (manager *AlertStrategyManager) registerStrategies(client *http.Client) {
	manager.strategies[common.MethodEmail] = &email.EmailAlert{}
	manager.strategies[common.MethodWechat] = &wechat.WechatAlert{Client: client}
	manager.strategies[common.MethodDingtalk] = &dingtalk.DingtalkAlert{Client: client}
	manager.strategies[common.MethodScript] = &script.ScriptAlert{Client: client}
	manager.strategies[common.MethodLark] = &lark.LarkAlert{Client: client}
}

// RegisterStrategy 统一注册所有告警策略
func (manager *AlertStrategyManager) RegisterStrategy(methodName string, alertMethod common.AlertMethod) error {
	_, ok := manager.strategies[methodName]
	// Same Method Name also can register
	manager.strategies[methodName] = alertMethod
	if ok {
		slog.Warn("registerStrategy is same name", slog.Any("methodName", methodName))
		return common.ErrorSameMethod
	}
	return nil
}

func (manager *AlertStrategyManager) EnableLimit() bool {
	if manager.LimitFunc == nil {
		return false
	}
	return true
}

// NewAlertStrategyManager 创建告警策略管理器
// 注册钩子，需要重新自己注册
func NewAlertStrategyManager(options ...Option) *AlertStrategyManager {
	manager := &AlertStrategyManager{
		strategies: make(map[string]common.AlertMethod),
	}

	client := http.NewHTTPClient()
	// 注册所有告警策略
	manager.registerStrategies(client)
	for _, option := range options {
		option(manager)
	}
	return manager
}

func (manager *AlertStrategyManager) before(c common.AlertContext) error {
	if manager.beforeHook == nil {
		return nil
	}
	return manager.beforeHook(c)
}

func (manager *AlertStrategyManager) after(c common.AlertContext) error {
	if manager.afterHook == nil {
		return nil
	}
	return manager.afterHook(c)
}

// Execute 根据告警方法执行相应的策略
func (manager *AlertStrategyManager) Execute(method string, c common.AlertContext, users []string) error {
	strategy, exists := manager.strategies[method]
	if !exists {
		return errors.New(fmt.Sprintf("method %s not found", method))
	}
	// 执行前钩子
	if err := manager.before(c); err != nil {
		return err
	}

	// 执行策略
	if err := strategy.Execute(c, users); err != nil {
		return err
	}

	if err := manager.after(c); err != nil {
		return err
	}

	return nil
}

func (manager *AlertStrategyManager) ExecuteTest(method string, c common.AlertContext) error {
	strategy, exists := manager.strategies[method]
	if !exists {
		return errors.New(fmt.Sprintf("method %s not found", method))
	}

	// 执行策略
	return strategy.ExecuteTest(c)
}

// ExecuteTestWithHook 根据告警方法执行相应的策略,测试可能会不想要测试钩子
func (manager *AlertStrategyManager) ExecuteTestWithHook(method string, c common.AlertContext) error {
	strategy, exists := manager.strategies[method]
	if !exists {
		return errors.New(fmt.Sprintf("method %s not found", method))
	}
	// 执行前钩子
	if err := manager.before(c); err != nil {
		return err
	}
	// 执行策略
	if err := strategy.ExecuteTest(c); err != nil {
		return err
	}
	if err := manager.after(c); err != nil {
		return err
	}
	return nil
}
