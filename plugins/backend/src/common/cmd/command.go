package cmd

import (
	"cia/common/errors"
	"flag"
)

type TArgType string

const (
	STRING TArgType = "string"
	INT             = "int"
	BOOL            = "bool"
)

func GetArg(
	_type TArgType,
	_name string,
	_defaultValue interface{},
	_desc string,
) interface{} {
	switch _type {
	case STRING:
		value := _defaultValue.(string)
		arg := flag.String(_name, value, _desc)
		flag.Parse()
		return arg
	case INT:
		value := _defaultValue.(int)
		arg := flag.Int(_name, value, _desc)
		flag.Parse()
		return arg
	case BOOL:
		value := _defaultValue.(bool)
		arg := flag.Bool(_name, value, _desc)
		flag.Parse()
		return arg
	default:
		panic(errors.TError{"", "unsupported argument type"})
	}
}
