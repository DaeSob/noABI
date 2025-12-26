package pack

import (
	"cia/common/utils"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

//////////////////////////////////////////////////////////////////
//	24.07.17
//	- use reflect
//		- can encounter **panic**
//		- ignore
//			- reflect.Func
//			- reflect.Chan
//	- ethereum@v1.13.14 unsupported abi type
//		- abi.HashTy
//		- abi.FixedPointTy
//
//	- todo-list
//		- zero safety
//		- optimizing
//		- abi.FunctionTy
//		- abi.AddressTy
//			- by bytes & *big.Int
//			- check size
//		- abi.FixedBytesTy
//			- check size
//		- abi.StringTy
//			-

func toArgs(abiArgs abi.Arguments, input any) ([]any, error) {
	values := make([]any, 0, len(abiArgs))
	for _, arg := range abiArgs {
		v, ok := toArg(arg.Name, arg.Type, input)
		if !ok {
			return nil, fmt.Errorf("failed toArgs: %s, %s", arg.Name, arg.Type.String())
		}
		values = append(values, v)
	}
	return values, nil
}

func toArg(name string, abiTyp abi.Type, input any) (any, bool) {
	rv := reflect.ValueOf(input)
	if name != "" {
		switch rv.Kind() {
		case reflect.Struct:
			rv = rv.FieldByName(abi.ToCamelCase(name))
		case reflect.Map:
			for _, key := range rv.MapKeys() {
				if key.String() == name {
					rv = rv.MapIndex(key)
					break
				}
			}
		case reflect.Func, reflect.Chan:
			// ignore
			return nil, true
		}
	}
	if rv.Kind() == reflect.Invalid {
		return nil, false
	}
	if (rv.Kind() == reflect.Pointer && rv.IsNil()) ||
		rv.Kind() == reflect.Func ||
		rv.Kind() == reflect.Chan {
		// ignore
		return nil, true
	}

	var (
		rawValue         = rv.Interface()
		strVal, isString = rawValue.(string)
		isHex            = isString && has0xPrefix(strVal)
		base             = utils.TernaryOp(isHex, 16, 10)
		typ              = abiTyp.T
		size             = abiTyp.Size
	)

	switch typ {
	case abi.IntTy:
		if size > 64 {
			return toBigInt(rawValue, base, true)
		}
		return toAbiInt(rawValue, base, size)

	case abi.UintTy:
		if size > 64 {
			return toBigInt(rawValue, base, false)
		}
		return toAbiUint(rawValue, base, size)

	case abi.BoolTy:
		if isString {
			return utils.TernaryOp(strVal == "true", true, false), strVal == "true" || strVal == "false"
		}
		b, ok := rawValue.(bool)
		return b, ok

	case abi.StringTy:
		if !isString || unsafe.Sizeof(strVal) > 32 {
			return nil, false
		}
		return strVal, isString

	case abi.AddressTy:
		return common.HexToAddress(strVal), true

	case abi.FixedBytesTy:
		if isHex {
			return common.Hex2BytesFixed(strVal[2:], size), true
		}
		bytes, ok := rawValue.([]byte)
		if !ok {
			return nil, false
		}

		blen := len(bytes)
		if blen == size {
			return bytes, true
		}
		if blen > size {
			return bytes[blen-size:], true
		}
		b := make([]byte, size)
		copy(b[size-blen:size], bytes)
		return b, true

	case abi.BytesTy:
		if isHex {
			return common.Hex2Bytes(strVal[2:]), true
		}
		bytes, ok := rawValue.([]byte)
		return bytes, ok

	case abi.FixedPointTy, abi.FunctionTy:
		return nil, false // not supported

	case abi.HashTy:
		if isHex {
			return common.HexToHash(strVal), true
		}
		bytes, ok := rawValue.([]byte)
		return common.BytesToHash(bytes), ok

	case abi.TupleTy:
		tuple, ok := toTuple(abiTyp, rawValue)
		if !ok {
			return nil, false
		}
		return tuple.Interface(), true

	case abi.SliceTy, abi.ArrayTy:
		elemTyp := abiTyp.Elem
		rv = reflect.ValueOf(rawValue)
		if !(rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
			return nil, false
		}

		slice := reflect.MakeSlice(toSliceTyp(*elemTyp), rv.Len(), rv.Len())
		if elemTyp.T == abi.TupleTy {
			for i := 0; i < rv.Len(); i++ {
				elem, ok := toTuple(*elemTyp, rv.Index(i).Interface())
				if !ok {
					return nil, false
				}
				slice.Index(i).Set(elem)
			}
		} else {
			for i := 0; i < rv.Len(); i++ {
				arg, ok := toArg("", *elemTyp, rv.Index(i).Interface())
				if !ok {
					return nil, false
				}
				slice.Index(i).Set(reflect.ValueOf(arg))
			}
		}
		return slice.Interface(), true

	default:
		return nil, false
	}
}

func toSliceTyp(abiTyp abi.Type) reflect.Type {
	if abiTyp.T == abi.TupleTy {
		tFields := make([]reflect.StructField, 0, len(abiTyp.TupleElems))
		for i, subTyp := range abiTyp.TupleElems {
			rawName := abiTyp.TupleRawNames[i]
			tFields = append(tFields, reflect.StructField{
				Name: abi.ToCamelCase(rawName),
				Type: subTyp.GetType(),
				Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", rawName)),
			})
		}
		tuple := reflect.New(reflect.StructOf(tFields)).Elem()
		return reflect.SliceOf(tuple.Type())
	} else {
		return reflect.SliceOf(abiTyp.GetType())
	}
}

func toTuple(abiTyp abi.Type, input any) (reflect.Value, bool) {
	tFields := make([]reflect.StructField, 0, len(abiTyp.TupleElems))
	values := make([]any, 0, len(abiTyp.TupleElems))
	for i, subTyp := range abiTyp.TupleElems {
		rawName := abiTyp.TupleRawNames[i]
		val, ok := toArg(rawName, *subTyp, input)
		if !ok {
			return reflect.Value{}, false
		}
		values = append(values, val)
		tFields = append(tFields, reflect.StructField{
			Name: abi.ToCamelCase(rawName),
			Type: reflect.TypeOf(val),
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", rawName)),
		})
	}

	tuple := reflect.New(reflect.StructOf(tFields)).Elem()
	for i, v := range values {
		tuple.Field(i).Set(reflect.ValueOf(v))
	}
	return tuple, true
}

func toBigInt(x any, base int, signed bool) (*big.Int, bool) {
	switch y := x.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,
		*int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		if signed {
			if z, ok := toInt64(y, 10, 64); ok {
				return new(big.Int).SetInt64(z), true
			}
		}
		if z, ok := toUint64(y, 10, 64); ok {
			return new(big.Int).SetUint64(z), true
		}
	case string:
		if base == 16 && has0xPrefix(y) {
			y = y[2:]
		}
		return new(big.Int).SetString(y, base)
	case *big.Int:
		return y, true
	case big.Int:
		return &y, true
	}
	return nil, false
}

func toInt64(x any, base, bitSize int) (int64, bool) {
	switch y := x.(type) {
	case string:
		if base == 16 && has0xPrefix(y) {
			y = y[2:]
		}
		signed, err := strconv.ParseInt(y, base, bitSize)
		return signed, err == nil
	case *big.Int:
		if y == nil {
			return 0, true
		}
		return y.Int64(), true
	case big.Int:
		return y.Int64(), true
	case uint:
		return int64(y), true
	case uint8:
		return int64(y), true
	case uint16:
		return int64(y), true
	case uint32:
		return int64(y), true
	case uint64:
		return int64(y), true
	case int:
		return int64(y), true
	case int8:
		return int64(y), true
	case int16:
		return int64(y), true
	case int32:
		return int64(y), true
	case int64:
		return y, true
	case float32:
		return int64(y), true
	case float64:
		return int64(y), true
	default:
		rv := reflect.ValueOf(y)
		if rv.Kind() == reflect.Pointer {
			return toInt64(reflect.Indirect(rv).Interface(), base, bitSize)
		}
		return 0, false
	}
}

func toUint64(x any, base, bitSize int) (uint64, bool) {
	switch y := x.(type) {
	case string:
		if base == 16 && has0xPrefix(y) {
			y = y[2:]
		}
		unsigned, err := strconv.ParseUint(y, base, bitSize)
		return unsigned, err == nil
	case *big.Int:
		if y == nil {
			return 0, true
		}
		return y.Uint64(), y.Sign() >= 0
	case uint:
		return uint64(y), true
	case uint8:
		return uint64(y), true
	case uint16:
		return uint64(y), true
	case uint32:
		return uint64(y), true
	case uint64:
		return y, true
	case int:
		return uint64(y), y >= 0
	case int8:
		return uint64(y), y >= 0
	case int16:
		return uint64(y), y >= 0
	case int32:
		return uint64(y), y >= 0
	case int64:
		return uint64(y), y >= 0
	case float32:
		return uint64(y), y >= 0
	case float64:
		return uint64(y), y >= 0
	default:
		rv := reflect.ValueOf(y)
		if rv.Kind() == reflect.Pointer {
			return toUint64(reflect.Indirect(rv).Interface(), base, bitSize)
		}
		return 0, false
	}
}

func toAbiInt(x any, base, bitSize int) (any, bool) {
	if bitSize > 64 || bitSize%8 != 0 {
		return nil, false
	}
	i64, ok := toInt64(x, base, bitSize)
	if !ok {
		return nil, false
	}

	max := int64(1<<(bitSize-1) - 1)
	min := int64(-1 << (bitSize - 1))
	if i64 > max || i64 < min {
		return nil, false
	}

	switch bitSize {
	case 8:
		return int8(i64), true
	case 16:
		return int16(i64), true
	case 32:
		return int32(i64), true
	case 64:
		return i64, true
	default:
		return toBigInt(i64, base, true)
	}
}

func toAbiUint(x any, base, bitSize int) (any, bool) {
	if bitSize > 64 || bitSize%8 != 0 {
		return nil, false
	}
	ui64, ok := toUint64(x, base, bitSize)
	if !ok {
		return nil, false
	}

	max := uint64(1<<bitSize - 1)
	if ui64 > max {
		return nil, false
	}

	switch bitSize {
	case 8:
		return uint8(ui64), true
	case 16:
		return uint16(ui64), true
	case 32:
		return uint32(ui64), true
	case 64:
		return ui64, true
	default:
		return toBigInt(ui64, base, false)
	}
}
