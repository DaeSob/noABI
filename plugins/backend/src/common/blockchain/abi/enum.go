package abi

import (
	"cia/common/utils"
	"math/big"

	common "github.com/ethereum/go-ethereum/common"
)

type TTypeKind int

const (
	// unsupported
	TYPE_UNSUPPORTED TTypeKind = iota
	// variable
	TYPE_ADDRESS
	TYPE_STRING
	TYPE_BOOL
	TYPE_BYTES
	// uint
	TYPE_UINT8
	TYPE_UINT16
	TYPE_UINT32
	TYPE_UINT64
	TYPE_UINT128
	TYPE_UINT256
	// int
	TYPE_INT8
	TYPE_INT16
	TYPE_INT32
	TYPE_INT64
	TYPE_INT128
	TYPE_INT256
	// struct
	TYPE_TUPLE

	// array
	TYPE_ARRAY_ADDRESS
	TYPE_ARRAY_BYTES
	TYPE_ARRAY_UINT256
	TYPE_ARRAY_STRING
	TYPE_ARRAY_TUPLE
)

func (tk TTypeKind) String() string {
	names := [...]string{
		"unsupported type",
		// variable
		"address",
		"string",
		"bool",
		"bytes",
		// uint
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uint128",
		"uint256",
		// int
		"int8",
		"int16",
		"int32",
		"int64",
		"int128",
		"int256",
		// struct
		"tuple",
		// array
		"address[]",
		"bytes[]",
		"uint256[]",
		"string[]",
		"tuple[]",
	}

	return names[tk]
}

func StringToType(_type string) TTypeKind {
	switch _type {
	case "address":
		return TYPE_ADDRESS
	case "string":
		return TYPE_STRING
	case "bool":
		return TYPE_BOOL
	case "uint8":
		return TYPE_UINT8
	case "uint16":
		return TYPE_UINT16
	case "uint32":
		return TYPE_UINT32
	case "uint64":
		return TYPE_UINT64
	case "uint128":
		return TYPE_UINT128
	case "uint256":
		return TYPE_UINT256
	case "int8":
		return TYPE_INT8
	case "int16":
		return TYPE_INT16
	case "int32":
		return TYPE_INT32
	case "int64":
		return TYPE_INT64
	case "int128":
		return TYPE_INT128
	case "int256":
		return TYPE_INT256
	case "address[]":
		return TYPE_ARRAY_ADDRESS
	case "bytes[]":
		return TYPE_ARRAY_BYTES
	case "uint256[]":
		return TYPE_ARRAY_UINT256
	case "string[]":
		return TYPE_ARRAY_STRING
	default:
		return TYPE_UNSUPPORTED
	}
}

// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/type_test.go
// -> TestTypeCheck() 함수 참고
func TypeToArg(_type TTypeKind, _value string) interface{} {
	switch _type.String() {
	case "address":
		return common.HexToAddress(_value)
	case "string":
		return _value
	case "bool":
		return utils.StringToBool(_value)
	case "uint8":
		u := utils.StringToUint64(_value, 10)
		return uint8(u)
	case "uint16":
		u := utils.StringToUint64(_value, 10)
		return uint16(u)
	case "uint32":
		u := utils.StringToUint64(_value, 10)
		return uint32(u)
	case "uint64":
		return utils.StringToUint64(_value, 10)
	case "uint128", "uint256":
		u, _ := new(big.Int).SetString(_value, 10)
		return u
	case "int8":
		i, _ := utils.StringToI64(_value)
		return int8(i)
	case "int16":
		i, _ := utils.StringToI64(_value)
		return int16(i)
	case "int32":
		i, _ := utils.StringToI64(_value)
		return int32(i)
	case "int64":
		i, _ := utils.StringToI64(_value)
		return int64(i)
	case "int128", "int256":
		i, _ := new(big.Int).SetString(_value, 10)
		return i
	default:
		return _value
	}
}

func TypeToArr(_type TTypeKind, _values []string) interface{} {
	switch _type.String() {
	case "address[]":
		return _convertStrArrToAddrArr(_values)
	case "bytes[]":
		return _convertStrArrToBytesArr(_values)
	case "uint256[]":
		return _convertStrArrToUint256Arr(_values)
	case "string[]":
		return _values
	default:
		return _values
	}
}

// utils
func _convertStrArrToAddrArr(
	_addressList []string,
) []common.Address {
	var res []common.Address
	for _, addr := range _addressList {
		res = append(res, TypeToArg(TYPE_ADDRESS, addr).(common.Address))
	}

	return res
}

func _convertStrArrToBytesArr(
	_bytesList []string,
) [][]byte {
	var res [][]byte
	for _, bytes := range _bytesList {
		res = append(res, utils.HexToBytes(bytes))
	}

	return res
}

func _convertStrArrToUint256Arr(
	_uint256List []string,
) []*big.Int {
	var res []*big.Int
	for _, uint256 := range _uint256List {
		res = append(res, TypeToArg(TYPE_UINT256, uint256).(*big.Int))
	}

	return res
}
