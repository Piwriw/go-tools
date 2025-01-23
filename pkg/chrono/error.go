package chrono

import "errors"

var (
	ErrInvalidJob  = errors.New("invalid job")
	ErrTaskFuncNil = errors.New("task function cannot be nil")
)
