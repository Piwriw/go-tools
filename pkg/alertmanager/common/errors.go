package common

import "errors"

var ErrorBeforeHook = errors.New("[alert]:before hook is failed")
var ErrorAfterHook = errors.New("[alert]:after hook is failed")
var ErrorSameMethod = errors.New("[alert]:same alert method name")
var ErrorNilLimitFunc = errors.New("[alert]:limitFunc is nil")
