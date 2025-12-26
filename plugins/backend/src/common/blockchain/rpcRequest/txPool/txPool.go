package txPool

import (
	"cia/common/blockchain/rpcRequest"
	rpc "cia/common/utils/rpcRequest"
)

type TTxPool struct {
	RPC rpcRequest.TRPC
}

func (txPool *TTxPool) Inspect(_id int64) *rpc.TRpcResponse {
	const method = "txpool_inspect"
	params := []interface{}{}

	return txPool.RPC.Request(method, params, _id)
}

func (txPool *TTxPool) Status(_id int64) *rpc.TRpcResponse {
	const method = "txpool_status"
	params := []interface{}{}

	return txPool.RPC.Request(method, params, _id)
}
