package eth

import (
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/types"
	"cia/common/blockchain/utils/encode"
	"cia/common/errors"
	"cia/common/utils"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func CallMethod(
	_mtd types.TMethod,
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) []string {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	data := encode.EncodeFunctionCallToHexString(_mtd, _args...)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	// decoding result
	hexResult := result.ResultToString()
	decodeResult, err := encode.DecodeParamsToString(
		_mtd.Outputs,
		utils.HexToBytes(hexResult),
	)
	if err != nil {
		panic(err)
	}

	return decodeResult
}

func CallMethodToHexString(
	_mtd types.TMethod,
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) string {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	data := encode.EncodeFunctionCallToHexString(_mtd, _args...)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	// decoding result
	return result.ResultToString()
}

type TMethod abi.Method

func (method TMethod) CallMethods(
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) []string {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}

	arguments, err := method.Inputs.Pack(_args...)
	if err != nil {
		panic(err)
	}

	bytes := append(method.ID, arguments...)
	data := utils.BytesToHexString(bytes)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	// decoding result
	hexResult := result.ResultToString()
	decodeResult, err := method.Outputs.Unpack(utils.HexToBytes(hexResult))
	if err != nil {
		panic(err)
	}

	res := []string{}
	for _, decode := range decodeResult {
		strDecode := fmt.Sprintf("%v", decode)
		res = append(res, strDecode)
	}

	return res
}

func (method TMethod) CallMethodsToHexString(
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) string {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	arguments, err := method.Inputs.Pack(_args...)
	if err != nil {
		panic(err)
	}

	bytes := append(method.ID, arguments...)
	data := utils.BytesToHexString(bytes)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	// decoding result
	return result.ResultToString()
}

func (method TMethod) CallMethodsToBytes(
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) []byte {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	arguments, err := method.Inputs.Pack(_args...)
	if err != nil {
		panic(err)
	}

	bytes := append(method.ID, arguments...)
	data := utils.BytesToHexString(bytes)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	// decoding result
	return result.ResultToBytes()
}

func (method TMethod) CallMethodsToInterface(
	_args []interface{},
	_rpc string,
	_from string,
	_to string,
	_id int64,
) interface{} {
	// ready call
	eth := TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	arguments, err := method.Inputs.Pack(_args...)
	if err != nil {
		panic(err)
	}

	bytes := append(method.ID, arguments...)
	data := utils.BytesToHexString(bytes)

	// do call
	result := eth.Call(_from, _to, data, "latest", _id)
	if result.Error != nil {
		panic(errors.ERROR_CALL(result.ErrorMessage()))
	}

	return method.Unpack(result.ResultToBytes())

}

// TODO
/*
func (method TMethod) Unpack(_in []byte) (res interface{}) {
	unpacked, err := method.Outputs.Unpack(_in)
	if err != nil {
		panic(err)
	}

	_byte, err := json.Marshal(unpacked[0])
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(_byte, &res)
	if err != nil {
		panic(err)
	}
	return
}
*/

func (method TMethod) Unpack(_in []byte) (res interface{}) {
	unpacked, err := method.Outputs.Unpack(_in)
	if err != nil {
		panic(err)
	}

	outputs := make(map[string]interface{})
	for i, unpacked := range unpacked {
		_byte, err := json.Marshal(unpacked)
		if err != nil {
			panic(err)
		}
		var unMarshal interface{}
		err = json.Unmarshal(_byte, &unMarshal)
		if err != nil {
			panic(err)
		}
		outputs[method.Outputs[i].Name] = unMarshal
	}
	res = outputs
	return
}
