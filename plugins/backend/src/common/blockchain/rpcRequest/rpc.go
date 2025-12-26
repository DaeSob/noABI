package rpcRequest

import (
	rpc "cia/common/utils/rpcRequest"
)

type TRPC struct {
	URL string // rpc endpoint url
}

func (trpc *TRPC) Request(_method string, _params []interface{}, _id int64) *rpc.TRpcResponse {
	result := rpc.RpcRequest(trpc.URL, _method, _params, _id)

	// release memory
	defer func() {
		result = nil
	}()

	return result
}
