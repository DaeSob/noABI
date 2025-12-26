package mapper

import "reflect"

type TJsonMap struct {
	m map[string]interface{}

	// for finder/remover
	splitKey []string
	cursor   int

	// for adder
	insertKey   string
	insertValue interface{}
}

/*** reflect type definition ***/
var IntType = reflect.TypeOf(1)
var Float64Type = reflect.TypeOf(float64(1))
var StringType = reflect.TypeOf(string(""))
var BoolType = reflect.TypeOf(false)
var SliceType = reflect.TypeOf([]interface{}(nil))
var SliceMapType = reflect.TypeOf([]map[string]interface{}(nil))
var JsonMapType = reflect.TypeOf((map[string]interface{})(nil))

/*** Token for key depth ***/
var SPLIT_TOKEN string = "."
