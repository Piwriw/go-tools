package common

// AlertContext 告警上下文接口，不同策略有不同的上下文
type AlertContext interface {
	GetMethod() string
}
