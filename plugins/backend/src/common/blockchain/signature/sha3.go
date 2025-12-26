package signature

import (
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

func Data(_type string, _value string) interface{} {
	return makeData(_type, _value)
}

func ArrayData(_type string, _values []string) interface{} {
	return makeArrayData(_type, _values)
}

func Bytes(_value string) []byte {
	return makeData("bytes", _value).([]byte)
}

// make data
// solsha3
// address
// bool
// string
// bytes4 & bytes32 & bytes
// uint8 & uint256
func makeData(_type string, _value string) interface{} {
	var data interface{}
	switch _type {
	case "address":
		data = solsha3.Address(_value)
	case "bool":
		data = solsha3.Bool(_value)
	case "string":
		data = solsha3.String(_value)
	case "bytes4":
		data = solsha3.Bytes4(_value)
	case "bytes32":
		data = solsha3.Bytes32(_value)
	case "bytes":
		data = solsha3.Bytes32(_value)
	case "uint8":
		data = solsha3.Uint8(_value)
	case "uint256":
		data = solsha3.Uint256(_value)
	default:
		panic("not supported type")
	}

	return data
}

func makeArrayData(_type string, _values []string) interface{} {
	var data interface{}
	switch _type {
	case "address[]":
		data = solsha3.AddressArray(_values)
	case "uint256[]":
		data = solsha3.Uint256Array(_values)
	default:
		panic("not supported type")
	}

	return data
}
