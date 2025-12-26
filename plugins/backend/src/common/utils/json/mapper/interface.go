package mapper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// NewBytes : new jmap from json bytes
func NewBytes(_b []byte) (*TJsonMap, error) {
	j := &TJsonMap{
		m: make(map[string]interface{}),
	}
	err := FromJson(_b, &j.m)
	return j, err
}

// PPrint : pretty print
func (j *TJsonMap) PPrint() string {
	b, _ := ToJson(j.m)
	str, _ := prettyPrint(b)
	return str
}

// Print : print
func (j *TJsonMap) Print() string {
	b, _ := ToJson(j.m)
	return string(b)
}

// toMap
func (j *TJsonMap) ToMap() interface{} {
	return j.m
}

// Find : find value from json key
func (j *TJsonMap) Find(_k string) interface{} {
	if _k == "" {
		return j.m
	}

	j.splitKey = strings.Split(_k, SPLIT_TOKEN)
	j.cursor = 0

	v := j.finder_search_root()

	return v
}

// Remove : remove value from key. prevent remove root.
func (j *TJsonMap) Remove(_k string) {
	if _k == "" {
		fmt.Println("### ERROR ### cannot delete root")
		return
	}

	j.splitKey = strings.Split(_k, SPLIT_TOKEN)
	j.cursor = 0

	j.remover_search_root()
}

// Insert : insert/update, when insert root, set [base == ""]
func (j *TJsonMap) Insert(_base, _k string, _v interface{}) {
	j.splitKey = strings.Split(_base, SPLIT_TOKEN)
	j.splitKey = append(j.splitKey, _k)
	j.cursor = 0

	j.insertKey = _k
	// type cast : []map[string]interface{} -> []interface{}
	vTmp := make([]interface{}, 0)
	switch reflect.TypeOf(_v) {
	case SliceMapType:
		for i := range _v.([]map[string]interface{}) {
			vTmp = append(vTmp, _v.([]map[string]interface{})[i])
		}
		j.insertValue = vTmp
		break

	default:
		j.insertValue = _v
	}

	j.adder_search_root()
}

// ToJson : object(struct) to json bytes
func ToJson(_o interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(_o)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// FromJson : json bytes to object(struct)
func FromJson(_byte []byte, _o interface{}) error {
	err := json.Unmarshal(_byte, _o)
	return err
}
