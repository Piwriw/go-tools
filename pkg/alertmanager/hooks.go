package alertmanager

type HookFunc func(c AlertContext) error

type HookHandler struct {
	Before HookFunc
	After  HookFunc
}

func callBeforeHook(hook *HookHandler, c AlertContext) error {
	if hook == nil {
		return nil
	}
	return hook.Before(c)
}

func callAfterHook(hook *HookHandler, c AlertContext) error {
	if hook == nil {
		return nil
	}
	return hook.After(c)
}
