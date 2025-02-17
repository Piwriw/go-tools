package progressbar

import (
	"fmt"
	"reflect"
)

func callFunc(fn any, params ...any) error {
	if fn == nil {
		return nil
	}
	f := reflect.ValueOf(fn)
	if f.IsZero() {
		return nil
	}
	if len(params) != f.Type().NumIn() {
		return fmt.Errorf("expected function with %d parameters, got one with %d", f.Type().NumIn(), len(params))
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	returnValues := f.Call(in)
	for _, val := range returnValues {
		i := val.Interface()
		if err, ok := i.(error); ok {
			return err
		}
	}
	return nil
}
