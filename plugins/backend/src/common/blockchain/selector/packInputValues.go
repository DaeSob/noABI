package selector

import (
	typeABI "cia/common/blockchain/abi"
	"cia/common/blockchain/types"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func PackInputValues(_abi types.TMethod, _values []interface{}) []interface{} {

	var inputValues []interface{}
	for i := 0; i < len(_abi.Inputs); i++ {
		idx := strings.LastIndex(_abi.Inputs[i].Type, "[]")
		if idx == -1 {
			inputValues = append(inputValues, typeABI.TypeToArg(typeABI.StringToType(_abi.Inputs[i].Type), _values[i].(string)))
		} else {
			inputValues = append(inputValues, typeABI.TypeToArr(typeABI.StringToType(_abi.Inputs[i].Type), _values[i].([]string)))
		}
	}
	return inputValues

}

func bindInputs(_method abi.Method, _values []interface{}) []interface{} {

	var inputValues []interface{}
	for i := 0; i < len(_method.Inputs); i++ {
		paramType := _method.Inputs[i].Type.String()
		idx := strings.LastIndex(paramType, "[]")
		if idx == -1 {
			inputValues = append(inputValues, typeABI.TypeToArg(typeABI.StringToType(paramType), _values[i].(string)))
		} else {
			inputValues = append(inputValues, typeABI.TypeToArr(typeABI.StringToType(paramType), _values[i].([]string)))
		}
	}
	return inputValues

}

func Bind(_method abi.Method, _values []interface{}) ([]byte, error) {

	boundInput := bindInputs(_method, _values)
	arguments, err := _method.Inputs.Pack(boundInput...)
	if err != nil {
		return nil, err
	}
	return append(_method.ID, arguments...), nil

}

func PackParams(_path, _contractName, _name string, _values []interface{}) ([]byte, error) {

	method, err := SelectMethod(_path, _contractName, _name)
	if err != nil {
		return nil, err
	}
	return Bind(method, _values)

}

func Pack(_path, _contractName, _name string, _args ...interface{}) ([]byte, error) {

	method, err := SelectMethod(_path, _contractName, _name)
	if err != nil {
		return nil, err
	}
	arguments, err := method.Inputs.Pack(_args...)
	if err != nil {
		return nil, err
	}
	return append(method.ID, arguments...), nil

}
