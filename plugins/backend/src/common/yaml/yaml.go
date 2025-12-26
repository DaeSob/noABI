package yaml

import (
	"cia/common/utils/http/httpRequest"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func LoadFromFile(_strFilePath string, _mapUnmarshal *map[string]interface{}) error {
	byteData, err := ioutil.ReadFile(_strFilePath)
	if err != nil {
		return err
	}

	err = loadImportFromBytes(byteData, _mapUnmarshal)
	if err != nil {
		return err
	}
	return loadFromBytes(byteData, _mapUnmarshal)
}

func LoadFromHttp(_strUrl string, _mapUnmarshal *map[string]interface{}) error {
	byteData, err := httpRequest.GetRequestToBytes(
		_strUrl,
		5*1000,
	)
	if err != nil {
		return err
	}
	err = loadImportFromBytes(byteData, _mapUnmarshal)
	if err != nil {
		return err
	}
	return loadFromBytes(byteData, _mapUnmarshal)
}

func loadImportFromBytes(_byteData []byte, _mapUnmarshal *map[string]interface{}) error {

	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal(_byteData, &m)
	if err != nil {
		return err
	}

	var importFiles []string

	for k := range m {
		strNewKey := fmt.Sprintf("%v", k)
		if strNewKey != "import" {
			continue
		}

		mapstructure.Decode(m[k], &importFiles)
	}

	type TImporting struct {
		Proto string
		Uri   string
	}

	var importing []TImporting
	for _, v := range importFiles {
		idx := strings.Index(v, "://")
		var proto string
		var uri string
		if idx == -1 {
			proto = "file://"
			uri = strings.Trim(v, " ")
		} else {
			proto = v[0 : idx+3]
			uri = v[idx+3:]
			proto = strings.Trim(proto, " ")
		}
		importing = append(importing, TImporting{
			Proto: proto,
			Uri:   uri,
		})
	}

	for _, v := range importing {
		switch v.Proto {
		case "file://":
			{
				err = LoadFromFile(v.Uri, _mapUnmarshal)
				if err != nil {
					return err
				}
			}
		case "http://", "https://":
			{
				err = LoadFromHttp(v.Proto+v.Uri, _mapUnmarshal)
				if err != nil {
					return err
				}
			}
		}
	}

	m = nil
	return err

}

func loadFromBytes(_byteData []byte, _mapUnmarshal *map[string]interface{}) error {

	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal(_byteData, &m)
	if err != nil {
		return err
	}

	for k := range m {
		strNewKey := fmt.Sprintf("%v", k)
		if strNewKey == "import" {
			continue
		}
		(*_mapUnmarshal)[strNewKey] = m[k]
		err = _marshaling(strNewKey, m[k], _mapUnmarshal)
		if err != nil {
			break
		}
	}

	m = nil
	return err
}

func _marshaling(_strKey string, _m interface{}, _mapUnmarshal *map[string]interface{}) error {

	byteData, err := yaml.Marshal(_m)
	if err != nil {
		return err
	}

	unMar := make(map[interface{}]interface{})
	err = yaml.Unmarshal(byteData, &unMar)

	if err != nil {
		return err
	}

	for k := range unMar {
		strCurrentKey := fmt.Sprintf("%v", k)
		strNewKey := _strKey + "." + strCurrentKey
		(*_mapUnmarshal)[strNewKey] = unMar[k]
		_marshaling(strNewKey, unMar[k], _mapUnmarshal)
	}
	unMar = nil
	byteData = nil

	return err
}
