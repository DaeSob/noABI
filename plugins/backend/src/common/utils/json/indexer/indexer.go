package indexer

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TJSONIndexer struct {
	//Private
	tempJsonData map[string]interface{}

	//Public
	RawData      string
	SliceKeyword []map[string]interface{}
}

func (j *TJSONIndexer) _mapParse(_key string, _in interface{}) {
	iter := reflect.ValueOf(_in).MapRange()
	for iter.Next() {
		key := _key + "." + iter.Key().String()
		value := iter.Value().Interface()
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			j._mapParse(key, value)
		case reflect.Slice:
			j._sliceParse(key, value)
		default:
			j.tempJsonData[key] = value
		}
	}
}

func (j *TJSONIndexer) _sliceParse(_key string, _in interface{}) {
	j.tempJsonData[_key] = _in
	v := reflect.ValueOf(_in)
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		key := _key + "[" + fmt.Sprint(i) + "]"
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			j._mapParse(key, value)
		case reflect.Slice:
			j._sliceParse(key, value)
		default:
			j.tempJsonData[key] = value
		}
	}
}

func (j *TJSONIndexer) ConvertStringToJSON(_data string) error {
	var mapParse map[string]interface{}
	j.tempJsonData = make(map[string]interface{})

	defer func() {
		mapParse = nil
		j.tempJsonData = nil
	}()

	err := json.Unmarshal([]byte(_data), &mapParse)
	if err != nil {
		return err
	}

	iter := reflect.ValueOf(mapParse).MapRange()
	for iter.Next() {
		key := iter.Key().String()
		value := iter.Value().Interface()
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			j._mapParse(key, value)
		case reflect.Slice:
			j._sliceParse(key, value)
		default:
			j.tempJsonData[key] = value
		}
	}
	j.SliceKeyword = append(j.SliceKeyword, j.tempJsonData)
	return nil
}

func (j *TJSONIndexer) ConvertMapToJSON(_in interface{}) error {
	j.tempJsonData = make(map[string]interface{})
	defer func() {
		j.tempJsonData = nil
	}()

	iter := reflect.ValueOf(_in).MapRange()
	for iter.Next() {
		key := iter.Key().String()
		value := iter.Value().Interface()
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			j._mapParse(key, value)
		case reflect.Slice:
			j._sliceParse(key, value)
		default:
			j.tempJsonData[key] = value
		}
	}
	j.SliceKeyword = append(j.SliceKeyword, j.tempJsonData)
	return nil
}

func (j *TJSONIndexer) ConvertSliceToJSON(_in interface{}) error {
	v := reflect.ValueOf(_in)
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			{
				j.ConvertMapToJSON(value)
			}
		case reflect.String:
			{
				j.ConvertStringToJSON(fmt.Sprintf("%v", value))
			}
		default:
		}
	}
	return nil
}

func (j *TJSONIndexer) ConvertToJSON(_in interface{}) error {
	switch reflect.TypeOf(_in).Kind() {
	case reflect.Map:
		j.ConvertMapToJSON(_in)
	case reflect.Slice:
		j.ConvertSliceToJSON(_in)
	case reflect.String:
		j.ConvertStringToJSON(fmt.Sprintf("%v", _in))
	}
	return nil
}

func (j *TJSONIndexer) Find(idx int64, _keyword string) interface{} {
	return j.SliceKeyword[idx][_keyword]
}
