package pack

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func PackAny(abi abi.ABI, name string, input any) (packed []byte, err error) {
	method, exist := abi.Methods[name]
	if !exist {
		return nil, errors.New("not found method")
	}

	args, err := toArgs(method.Inputs, input)
	if err != nil {
		return nil, err
	}
	return abi.Pack(method.Name, args...)
}

func PackAny_safe(abi abi.ABI, name string, input any) (packed []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
			err = fmt.Errorf("%v", r)
		}
	}()
	return PackAny(abi, name, input)
}

func PackAny2Hex(abi abi.ABI, name string, input any, safe ...bool) (hex string, err error) {
	var packed []byte
	if safe != nil && safe[0] {
		packed, err = PackAny_safe(abi, name, input)
	} else {
		packed, err = PackAny(abi, name, input)
	}
	if err != nil {
		return "", err
	}
	return hexutil.Encode(packed), nil
}

// /////////////////////////////////////////////
//	for overloaded method
//	find method by Id(method-sig)

func PackAnyBySig(abi abi.ABI, sig []byte, input any) (packed []byte, err error) {
	method, err := abi.MethodById(sig)
	if err != nil {
		return nil, err
	}
	args, err := toArgs(method.Inputs, input)
	if err != nil {
		return nil, err
	}
	return abi.Pack(method.Name, args...)
}

func PackAnyBySig_safe(abi abi.ABI, sig []byte, input any) (packed []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
			err = fmt.Errorf("%v", r)
		}
	}()
	return PackAnyBySig(abi, sig, input)
}

func PackAny2HexBySig(abi abi.ABI, sig []byte, input any, safe ...bool) (hex string, err error) {
	var packed []byte
	if safe != nil && safe[0] {
		packed, err = PackAnyBySig_safe(abi, sig, input)
	} else {
		packed, err = PackAnyBySig(abi, sig, input)
	}
	if err != nil {
		return "", err
	}
	return hexutil.Encode(packed), nil
}
