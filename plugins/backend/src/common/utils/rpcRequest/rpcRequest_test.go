package rpcRequest

import (
	"fmt"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_RpcRequest(t *testing.T) {
	// ready rpc request
	const rpcUrl = "http://10.10.24.95:8551"
	const method = "eth_getBlockByNumber"
	params := []interface{}{fmt.Sprintf("0x%x", -1), true}
	const id = 1

	// do rpc request
	res := RpcRequest(rpcUrl, method, params, id)

	assert.Equal(t, res.Result, nil, "result should be nil")
	assert.NotNil(t, res.Error)
}
