package call

import (
	rpcRequest "cia/common/blockchain/rpcRequest"
	txPool "cia/common/blockchain/rpcRequest/txPool"
	rpc "cia/common/utils/rpcRequest"
	types "cia/common/utils/types"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_MethodCallByName(t *testing.T) {
	// ready tx pool
	const url = "http://10.10.24.95:8551"
	txpool := &txPool.TTxPool{rpcRequest.TRPC{url}}

	// ready params
	id := (int64)(1)

	res := MethodCallByName(txpool, "Inspect", id)

	assert.Equal(t, len(res), 1)
	assert.Equal(t, types.TypeToString(res[0]), "*rpcRequest.TRpcResponse")
	assert.Equal(t, res[0].(*rpc.TRpcResponse).ResultToJsonString(false), `{"pending":{},"queued":{}}`)
}
