package rpcRequest

import (
	"bytes"
	"encoding/json"

	"cia/common/errors"
	http "cia/common/utils/http/httpRequest"
)

type TPayload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
}

// url - rpc endpoint url
// method
// params
// id
func RpcRequest(_url string, _method string, _params []interface{}, _id int64) *TRpcResponse {
	// ready body
	data := &TPayload{
		Jsonrpc: "2.0",
		Method:  _method,
		Params:  _params,
		ID:      _id,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// do http request
	// timeout : 30 sec
	const timeout = 30 * 1000
	jsonResponse, err := http.PostRequestFromBytes(_url, payloadBytes, timeout)
	if err != nil {
		panic(errors.ERROR_RPC_REQUEST(err.Error()))
	}

	// convert to bytes from response
	bytes := bytes.NewBufferString(jsonResponse.PPrint())
	if err != nil {
		panic(err)
	}

	// result from bytes
	result := &TRpcResponse{}
	err = json.Unmarshal(bytes.Bytes(), result)
	if err != nil {
		panic(err)
	}

	data = nil
	payloadBytes = nil
	jsonResponse = nil
	bytes = nil

	// release memory
	defer func() {
		result = nil
	}()

	return result
}
