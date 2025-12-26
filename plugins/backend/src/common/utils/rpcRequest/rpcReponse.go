package rpcRequest

import (
	"cia/common/utils"
	jMapper "cia/common/utils/json/mapper"
	"math/big"
	"strconv"
)

type TRpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// define rpc response structure
type TRpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (res *TRpcResponse) ToString(_pretty bool) string {
	bytesResult, err := jMapper.ToJson(res)
	if err != nil {
		return ""
	}

	jMap, err := jMapper.NewBytes(bytesResult)
	if err != nil {
		return ""
	}

	if _pretty == true {
		return jMap.PPrint()
	} else {
		return jMap.Print()
	}
}

// result to string
func (res *TRpcResponse) ResultToString() string {
	if res.Result != nil {
		return res.Result.(string)
	} else {
		return ""
	}
}

// result to bytes
func (res *TRpcResponse) ResultToBytes() []byte {
	if res.Result != nil {
		return utils.HexToBytes(res.Result.(string))
	} else {
		return ([]byte{})
	}
}

// result to json string
func (res *TRpcResponse) ResultToJsonString(pretty bool) string {
	if res.Result != nil {
		return utils.InterfaceToJsonString(res.Result, pretty)
	} else {
		return ""
	}
}

// result to uint64
func (res *TRpcResponse) ResultToUint64() uint64 {
	strResult := res.ResultToString()
	resultUint64, err := strconv.ParseUint(strResult, 0, 64)

	if err != nil {
		panic(err)
	}

	return resultUint64
}

// result to Big Int
func (res *TRpcResponse) ResultToBigInt() *big.Int {
	return utils.StringToBigInt(res.ResultToString())
}

// result to int64
func (res *TRpcResponse) ResultToInt64() int64 {
	strResult := res.ResultToString()
	resultInt64, err := strconv.ParseInt(strResult, 0, 64)

	if err != nil {
		panic(err)
	}

	return resultInt64
}

// result to json bytes
func (res *TRpcResponse) ResultToJson() *jMapper.TJsonMap {
	bytesError, err := jMapper.ToJson(res.Result)
	if err != nil {
		return nil
	}

	jMap, err := jMapper.NewBytes(bytesError)
	if err != nil {
		return nil
	}

	return jMap
}

// error to string
func (res *TRpcResponse) ErrorToString(pretty bool) string {
	jMap := res.ErrorToJson()

	if pretty == true {
		return jMap.PPrint()
	} else {
		return jMap.Print()
	}
}

// error to json bytes
func (res *TRpcResponse) ErrorToJson() *jMapper.TJsonMap {
	bytesError, err := jMapper.ToJson(res.Error)
	if err != nil {
		return nil
	}

	jMap, err := jMapper.NewBytes(bytesError)
	if err != nil {
		return nil
	}

	return jMap
}

// get error message
func (res *TRpcResponse) ErrorMessage() string {
	return res.ErrorToJson().Find("message").(string)
}
