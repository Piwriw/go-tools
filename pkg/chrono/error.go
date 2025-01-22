package chrono

import "errors"

var (
	ErrInvalidTask = errors.New("invalid task")
	ErrTaskFuncNil = errors.New("task function cannot be nil")
)
