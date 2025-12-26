package preference

import (
	"cia/common/utils"
	"fmt"
	"math/big"
	"strconv"
)

func Find(_strKey string) interface{} {
	inst := GetInstance()
	return inst._find(_strKey)
}

func _findToInt(key string, out *int) (err error) {
	inst := GetInstance()
	tempVal := fmt.Sprint(inst.mapYaml[key])
	if tempVal == "<nil>" {
		*out = int(0)
		return nil
	}
	tmp, err := strconv.ParseInt(tempVal, 10, 64)
	*out = int(tmp)
	return
}

func _findToString(key string, out *string) (err error) {
	inst := GetInstance()
	*out = fmt.Sprint(inst.mapYaml[key])
	if *out == "<nil>" {
		*out = ""
	}
	return
}

func _findToInt64(key string, out *int64) (err error) {
	inst := GetInstance()
	tempVal := fmt.Sprint(inst.mapYaml[key])
	if tempVal == "<nil>" {
		*out = int64(0)
		return nil
	}
	*out, err = strconv.ParseInt(tempVal, 10, 64)
	return
}

func _findToUint64(key string, out *uint64) (err error) {
	inst := GetInstance()
	tempVal := fmt.Sprint(inst.mapYaml[key])
	if tempVal == "<nil>" {
		*out = uint64(0)
		return nil
	}
	*out, err = strconv.ParseUint(tempVal, 10, 64)
	return
}

func _findToBigint(key string, out **big.Int) (err error) {
	inst := GetInstance()
	tempVal := fmt.Sprint(inst.mapYaml[key])
	if tempVal == "<nil>" {
		*out = big.NewInt(0)
		return nil
	}
	i64, err := strconv.ParseInt(tempVal, 10, 64)
	*out = big.NewInt(i64)
	return
}

func _findToBool(_key string, out *bool) (err error) {
	inst := GetInstance()
	flag := inst.mapYaml[_key]

	if flag == nil {
		*out = false
		return nil
	}

	*out = flag.(bool)
	return
}

func _findToStringArray(_key string, out *[]string) (err error) {
	inst := GetInstance()
	arrValue, ok := inst.mapYaml[_key].([]interface{})
	if ok {
		*out = utils.InterfaceArrayToStringArray(arrValue)
	}

	return
}
