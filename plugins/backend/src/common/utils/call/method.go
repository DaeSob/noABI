package call

import (
	"reflect"
)

// method call by name
// _struct : methodë¥?ê°€ì§?struct
// _methodName : ?¤í–‰?œí‚¬ method ?´ë¦„
// _params : method ??parameter
func MethodCallByName(_struct interface{}, _methodName string, _params ...interface{}) []interface{} {
	// ready params
	params := make([]reflect.Value, len(_params))
	for i := range _params {
		params[i] = reflect.ValueOf(_params[i])
	}

	// call
	ref := reflect.ValueOf(_struct).MethodByName(_methodName)
	returns := ref.Call(params)

	// ready result
	results := make([]interface{}, len(returns))
	for i := range returns {
		results[i] = returns[i].Interface()
	}

	return results
}
