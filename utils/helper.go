package utils

import (
	"errors"
	"reflect"
)

func Call(m map[string]interface{}, name string, params ...interface{}) (rs []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	rs = f.Call(in)
	return
}
