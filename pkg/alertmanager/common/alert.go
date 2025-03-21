package common

// AlertMethod 告警策略接口
type AlertMethod interface {
	Execute(c AlertContext, users []string) error
	ExecuteTest(c AlertContext) error
	Before(c AlertContext) error
	After(c AlertContext) error
}
