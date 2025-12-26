package mapper

import (
	"fmt"
	"reflect"
	"strconv"
)

// finder_search_root : finder recursive function start point
func (j *TJsonMap) finder_search_root() interface{} {
	var currentKey string
	currentKey = j.splitKey[j.cursor]

	// goto: switch type for next key
SWITCH_TYPE:
	switch reflect.TypeOf(j.m[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return j.m[currentKey]
		}
		return j.finder_search_map_r(j.m[currentKey].(map[string]interface{}))

	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return j.m[currentKey]
		}
		return j.finder_search_slice_r(j.m[currentKey].([]interface{}))

	case Float64Type, StringType, BoolType, IntType:
		return j.m[currentKey]

	default:
		// set next key
		if j.cursor < len(j.splitKey)-1 {
			j.cursor = j.cursor + 1
			currentKey = currentKey + "." + j.splitKey[j.cursor]
			goto SWITCH_TYPE
		}
	}
	return nil
}

// finder_search_slice_r : finder recursive function for slice
func (j *TJsonMap) finder_search_slice_r(_sub []interface{}) interface{} {
	// slice key is always integer
	var currentKey int
	var err error

	currentKey, err = strconv.Atoi(j.splitKey[j.cursor])
	if nil != err {
		fmt.Println("### ERROR ### :", err)
		return nil
	}

	if currentKey >= len(_sub) || currentKey < 0 {
		fmt.Println("### ERROR ### : index is out of range")
		return nil
	}

	switch reflect.TypeOf(_sub[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return _sub[currentKey]
		}
		return j.finder_search_map_r(_sub[currentKey].(map[string]interface{}))

		// --- must not hit this section. json doesn't allow netsted array.
	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return _sub[currentKey]
		}
		return j.finder_search_slice_r(_sub[currentKey].([]interface{}))

	case Float64Type, StringType, BoolType, IntType:
		return _sub[currentKey]
		// --- end of section

	default:
		// it must be error, must not hit
		fmt.Println("### LOG ### default: undefined error")
		return nil
	}
}

// finder_search_map_r : finder recursive function for map
func (j *TJsonMap) finder_search_map_r(_sub map[string]interface{}) interface{} {
	var currentKey string
	currentKey = j.splitKey[j.cursor]

	// goto: switch type for next key
SWITCH_TYPE:
	switch reflect.TypeOf(_sub[currentKey]) {
	case JsonMapType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return _sub[currentKey]
		}
		return j.finder_search_map_r(_sub[currentKey].(map[string]interface{}))

	case SliceType:
		j.cursor = j.cursor + 1
		if j.cursor >= len(j.splitKey) {
			return _sub[currentKey]
		}
		return j.finder_search_slice_r(_sub[currentKey].([]interface{}))

	case Float64Type, StringType, BoolType, IntType:
		return _sub[currentKey]

	default:
		// set next key
		if j.cursor < len(j.splitKey)-1 {
			j.cursor = j.cursor + 1
			currentKey = currentKey + "." + j.splitKey[j.cursor]
			goto SWITCH_TYPE
		}
	}
	return nil
}
