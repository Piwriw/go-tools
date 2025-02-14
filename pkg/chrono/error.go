package chrono

import "errors"

var (
	ErrInvalidJob      = errors.New("invalid job")
	ErrTaskFuncNil     = errors.New("task function cannot be nil")
	ErrTaskTimeout     = errors.New("task is timeout")
	ErrValidateTimeout = errors.New("task timeout must be greater than 0")
	ErrTaskFailed      = errors.New("task is failed")
	ErrFoundAlias      = errors.New("can not found  by alias")
)
