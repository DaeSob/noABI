package json

import (
	"encoding/json"
	"fmt"
	"reflect"

	jmapper "cia/common/utils/json/mapper"
)

func ConvertToJSON(_in interface{}) ([]jmapper.TJsonMap, error) {

	var sliceJson []jmapper.TJsonMap

	defer func() {
		sliceJson = nil
	}()

	switch reflect.TypeOf(_in).Kind() {
	case reflect.Map:
		b1, _ := json.Marshal(_in)
		j, e := jmapper.NewBytes(b1)
		if nil != e {
			return sliceJson, e
		}
		sliceJson = append(sliceJson, *j)
		break
	case reflect.Slice:
		v := reflect.ValueOf(_in)
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i).Interface()
			b1, _ := json.Marshal(value)
			j, e := jmapper.NewBytes(b1)
			if nil != e {
				return nil, e
			}
			sliceJson = append(sliceJson, *j)
		}
		break
	case reflect.String:
		j, e := jmapper.NewBytes([]byte(fmt.Sprintf("%v", _in)))
		if nil != e {
		}
		sliceJson = append(sliceJson, *j)
		break
	}

	return sliceJson, nil

}

func PPrint(_in []byte) string {
	j, _ := jmapper.NewBytes(_in)
	return j.PPrint()
}
