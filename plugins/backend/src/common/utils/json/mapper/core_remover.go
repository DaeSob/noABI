package mapper

import (
	"fmt"
	"reflect"
	"strconv"
)

// remover's action is almost same as finder.
// 	 : remover will remove element/value instead of return

// remover_search_root : remover recursive function start point
func (j *TJsonMap) remover_search_root() {
	var currentKey string
	currentKey = j.splitKey[j.cursor]

	// goto: switch type for next key
SWITCH_TYPE:
	switch reflect.TypeOf(j.m[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			delete(j.m, currentKey)
			return
		}
		j.remover_search_map_r(j.m[currentKey].(map[string]interface{}))
		return

	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			delete(j.m, currentKey)
			return
		}
		// slice must set return value
		j.m[currentKey] = j.remover_search_slice_r(j.m[currentKey].([]interface{}))
		return

	case Float64Type, StringType, BoolType, IntType:
		delete(j.m, currentKey)
		return

	default:
		// set next key
		if j.cursor < len(j.splitKey)-1 {
			j.cursor = j.cursor + 1
			currentKey = currentKey + "." + j.splitKey[j.cursor]
			goto SWITCH_TYPE
		}
	}
}

// remover_search_slice_r : remover recursive function for slice
//   - slice must have slice return value for set new slice
func (j *TJsonMap) remover_search_slice_r(_sub []interface{}) []interface{} {
	// slice key is always integer
	var currentKey int
	var err error

	currentKey, err = strconv.Atoi(j.splitKey[j.cursor])
	if nil != err {
		fmt.Println("### ERROR ### :", err)
		return _sub
	}

	if currentKey >= len(_sub) || currentKey < 0 {
		fmt.Println("### ERROR ### : index is out of range")
		return _sub
	}

	switch reflect.TypeOf(_sub[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			_sub = j.removeSliceElement(_sub, currentKey)
			return _sub
		}
		j.remover_search_map_r(_sub[currentKey].(map[string]interface{}))
		return _sub

		// --- must not hit this section. json doesn't allow netsted array.
	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			_sub = j.removeSliceElement(_sub, currentKey)
			return _sub
		}
		// slice must set return value
		_sub[currentKey] = j.remover_search_slice_r(_sub[currentKey].([]interface{}))
		return _sub
		// --- end of section

	case Float64Type, StringType, BoolType, IntType:
		_sub = j.removeSliceElement(_sub, currentKey)
		return _sub

	default:
		// it must be error, must not hit
		fmt.Println("### LOG ### default: undefined error")
	}
	return _sub
}

// remover_search_map_r : remover recursive function for map
func (j *TJsonMap) remover_search_map_r(_sub map[string]interface{}) {
	var currentKey string
	currentKey = j.splitKey[j.cursor]

	// goto: switch type for next key
SWITCH_TYPE:
	switch reflect.TypeOf(_sub[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			delete(_sub, currentKey)
			return
		}
		j.remover_search_map_r(_sub[currentKey].(map[string]interface{}))
		return

	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			delete(_sub, currentKey)
			return
		}
		_sub[currentKey] = j.remover_search_slice_r(_sub[currentKey].([]interface{}))
		return

	case Float64Type, StringType, BoolType, IntType:
		delete(_sub, currentKey)
		return

	default:
		// set next key
		if j.cursor < len(j.splitKey)-1 {
			j.cursor = j.cursor + 1
			currentKey = currentKey + "." + j.splitKey[j.cursor]
			goto SWITCH_TYPE
		}
	}
}

// delete slice element
func (j *TJsonMap) removeSliceElement(s []interface{}, index int) []interface{} {
	copy(s[index:], s[index+1:])
	s[len(s)-1] = nil
	s = s[:len(s)-1]
	return s
}
