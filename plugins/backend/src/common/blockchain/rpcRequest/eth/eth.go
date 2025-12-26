package eth

import (
	"cia/common/blockchain/keystore/signer"
	"cia/common/blockchain/rpcRequest"
	rpc "cia/common/utils/rpcRequest"
	"reflect"

	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

type TEth struct {
	RPC rpcRequest.TRPC
}

type TCallData struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func (eth *TEth) Call(
	_from,
	_to,
	_data,
	_block string,
	_id int64,
) *rpc.TRpcResponse {
	const method = "eth_call"

	callData := &TCallData{_from, _to, _data}
	params := []interface{}{callData, _block}

	return eth.RPC.Request(method, params, _id)
}

type TTransaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	Nonce    string `json:"nonce"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
	Gas      string `json:"gas"`
}

type TSimpleTransaction struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func (eth *TEth) SimpleEstimateGas(_transaction TSimpleTransaction, _id int64) *rpc.TRpcResponse {
	const method = "eth_estimateGas"
	params := []interface{}{_transaction, "latest"}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) EstimateGas(_transaction TTransaction, _id int64) *rpc.TRpcResponse {
	const method = "eth_estimateGas"
	params := []interface{}{_transaction, "latest"}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetBalance(_address string, _id int64) *rpc.TRpcResponse {
	const method = "eth_getBalance"
	params := []interface{}{_address, "latest"}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetBlockByNumber(_number string, _isTxInstance bool, _id int64) *rpc.TRpcResponse {
	const method = "eth_getBlockByNumber"
	params := []interface{}{_number, _isTxInstance}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetBlockNumber(_id int64) *rpc.TRpcResponse {
	const method = "eth_blockNumber"
	params := []interface{}{}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetSyncingStatus(_id int64) *rpc.TRpcResponse {
	const method = "eth_syncing"
	params := []interface{}{}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetChainId(_id int64) *rpc.TRpcResponse {
	const method = "eth_chainId"
	params := []interface{}{}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetGasPrice(_id int64, _params []interface{}) *rpc.TRpcResponse {
	return eth.getGasPrice(_id, _params)
}

func (eth *TEth) getGasPrice(_id int64, _params []interface{}) *rpc.TRpcResponse {
	const method = "eth_gasPrice"
	if _params == nil {
		return eth.RPC.Request(method, []interface{}{}, _id)
	}
	return eth.RPC.Request(method, _params, _id)
}

func (eth *TEth) GetLogs(_fromBlock uint64, _toBlock uint64, _address []string, _topics []interface{}, _id int64) *rpc.TRpcResponse {
	const method = "eth_getLogs"

	// ready params
	type temp struct {
		FromBlock string        `json:"fromBlock"`
		ToBlock   string        `json:"toBlock"`
		Address   []string      `json:"address"`
		Topics    []interface{} `json:"topics"`
	}
	tmp := &temp{fmt.Sprintf("0x%x", _fromBlock), fmt.Sprintf("0x%x", _toBlock), _address, _topics}
	params := []interface{}{tmp}

	// release memory
	defer func() {
		tmp = nil
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetTransactionCount(_address string, _id int64) *rpc.TRpcResponse {
	const method = "eth_getTransactionCount"
	params := []interface{}{_address, "latest"}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetPendingNonce(_address string, _id int64) *rpc.TRpcResponse {
	const method = "eth_getTransactionCount"
	params := []interface{}{_address, "pending"}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

func (eth *TEth) GetTransactionReceipt(_txHash string, _id int64) *rpc.TRpcResponse {
	const method = "eth_getTransactionReceipt"
	params := []interface{}{_txHash}

	// release memory
	defer func() {
		params = nil
	}()

	return eth.RPC.Request(method, params, _id)
}

// sign이 완료된 tx를 send
func (eth *TEth) SendSignedRawTransaction(_signedTx interface{}, _id int64) *rpc.TRpcResponse {
	const method = "eth_sendRawTransaction"

	var params []interface{}
	// release memory
	defer func() {
		params = nil
	}()

	// hex string or types.Transaction만 지원
	switch reflect.ValueOf(_signedTx).Type() {
	case reflect.TypeOf(""): // hex tx string type
		params = []interface{}{_signedTx}
		return eth.RPC.Request(method, params, _id)
	case reflect.TypeOf((*types.Transaction)(nil)): // tx type
		rawTxHex := signer.SignedRawTx2HexString(_signedTx.(*types.Transaction))
		params = []interface{}{"0x" + rawTxHex}
		return eth.RPC.Request(method, params, _id)
	default:
		// no match
		println("no match")
	}

	// panic when no match
	panic(fmt.Errorf("Unsupported tx type."))
}

// sign 이 안된 tx를 send
func (eth *TEth) SendRawTransaction(_tx signer.TRawTransaction, _signer *signer.TSigner, _passPhrase string, _id int64) *rpc.TRpcResponse {
	// get signed tx
	signedTx, err := _signer.SignRawTx(_passPhrase, _tx)
	if err != nil {
		panic(err)
	}
	rawTxHex := signer.SignedRawTx2HexString(signedTx)

	// release memory
	defer func() {
		signedTx = nil
	}()

	return eth.SendSignedRawTransaction(rawTxHex, _id)
}
