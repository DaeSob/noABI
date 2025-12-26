package abi

import (
	"encoding/json"

	jMapper "cia/common/utils/json/mapper"
)

type TArgumentStr struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	InternalType string         `json:"internalType"`
	Components   []TArgumentStr `json:"components,omitempty"`
	Indexed      bool           `json:"indexed"`
}

func (arg *TArgumentStr) Copy() *TArgumentStr {
	newLog := new(TArgumentStr)
	*newLog = *arg

	return newLog
}

func (arg TArgumentStr) Bytes() []byte {
	bytes, err := json.Marshal(arg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (arg TArgumentStr) String() string {
	jMap, err := jMapper.NewBytes(arg.Bytes())
	if err != nil {
		panic(err)
	}

	return jMap.Print()
}

func CreateArgument(
	_name,
	_type,
	_internalType string,
	_components []TArgumentStr,
	_indexed bool,
) *TArgumentStr {
	return &TArgumentStr{
		_name,
		_type,
		_internalType,
		_components,
		_indexed,
	}
}

func CreateDynamicType(_name, _type string, _indexed bool) *TArgumentStr {
	return CreateArgument(_name, _type, _type, nil, _indexed)
}

func CreateTuple(_name, _internalType string, _components []TArgumentStr) *TArgumentStr {
	return CreateArgument(
		_name,
		TYPE_TUPLE.String(),
		_internalType,
		_components,
		false,
	)
}

func CreateArrayTuple(_name, _internalType string, _components []TArgumentStr) *TArgumentStr {
	return CreateArgument(
		_name,
		TYPE_ARRAY_TUPLE.String(),
		_internalType,
		_components,
		false,
	)
}

// variable
func CreateAddressType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_ADDRESS.String(), _indexed)
}

func CreateStringType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_STRING.String(), _indexed)
}

func CreateBoolType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_BOOL.String(), _indexed)
}

func CreateBytesType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_BYTES.String(), _indexed)
}

// uint
func CreateUint8Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT8.String(), _indexed)
}

func CreateUint16Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT16.String(), _indexed)
}

func CreateUint32Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT32.String(), _indexed)
}

func CreateUint64Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT64.String(), _indexed)
}

func CreateUint128Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT128.String(), _indexed)
}

func CreateUint256Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_UINT256.String(), _indexed)
}

// int
func CreateInt8Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT8.String(), _indexed)
}

func CreateInt16Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT16.String(), _indexed)
}

func CreateInt32Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT32.String(), _indexed)
}

func CreateInt64Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT64.String(), _indexed)
}

func CreateInt128Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT128.String(), _indexed)
}

func CreateInt256Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_INT256.String(), _indexed)
}

// array
func CreateArrayAddressType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_ARRAY_ADDRESS.String(), _indexed)
}

func CreateArrayBytesType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_ARRAY_BYTES.String(), _indexed)
}

func CreateArrayUint256Type(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_ARRAY_UINT256.String(), _indexed)
}

func CreateArrayStringType(_name string, _indexed bool) *TArgumentStr {
	return CreateDynamicType(_name, TYPE_ARRAY_STRING.String(), _indexed)
}
