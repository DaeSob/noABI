package types

import (
	"reflect"
)

func GetValue(_variable interface{}) reflect.Value {
	return reflect.ValueOf(_variable)
}

func GetType(_variable interface{}) reflect.Type {
	return GetValue(_variable).Type()
}

func TypeToString(_variable interface{}) string {
	return GetType(_variable).String()
}
