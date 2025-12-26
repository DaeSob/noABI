package eth

import (
	"cia/common/blockchain/keystore/signer"
	"cia/common/blockchain/rpcRequest"
	"cia/common/errors"
	"cia/common/utils"
	request "cia/common/utils/rpcRequest"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func SendTransaction(
	_rpc string,
	_to string,
	_data []byte,
	_value string,
	_gasLimit uint64,
	_signer signer.TSigner,
	_phrase string,
) *request.TRpcResponse {
	// set eth
	ethcall := &TEth{RPC: rpcRequest.TRPC{URL: _rpc}}

	// gas price
	resGasPrice := ethcall.GetGasPrice(int64(0), []interface{}{"latest"})

	return sendTransaction(
		_rpc,
		_to,
		_data,
		_value,
		_gasLimit,
		resGasPrice.ResultToString(),
		_signer,
		_phrase,
	)
}

func SendTxLuniverse(
	_rpc string,
	_to string,
	_data []byte,
	_value string,
	_gasLimit uint64,
	_signer signer.TSigner,
	_phrase string,
) *request.TRpcResponse {
	// set eth
	ethcall := &TEth{RPC: rpcRequest.TRPC{URL: _rpc}}

	// gas price
	resGasPrice := ethcall.GetGasPrice(int64(0), []interface{}{})

	return sendTransaction(
		_rpc,
		_to,
		_data,
		_value,
		_gasLimit,
		resGasPrice.ResultToString(),
		_signer,
		_phrase,
	)
}

func sendTransaction(
	_rpc string,
	_to string,
	_data []byte,
	_value string,
	_gasLimit uint64,
	_gasPrice string,
	_signer signer.TSigner,
	_phrase string,
) *request.TRpcResponse {
	// set eth
	ethcall := &TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	id := int64(1)

	// nonce
	resNonce := ethcall.GetTransactionCount(_signer.Address.String(), id)
	nonce := resNonce.ResultToUint64()

	// value to big int
	value, _ := big.NewInt(0).SetString(_value, 16)

	// make tx
	tx := signer.TRawTransaction{
		To:       common.HexToAddress(_to),
		Value:    value,
		Data:     _data,
		Nonce:    nonce,
		GasPrice: utils.StringToBigInt(_gasPrice),
		GasLimit: _gasLimit,
	}

	estimateGas := ethcall.SimpleEstimateGas(
		TSimpleTransaction{
			From: _signer.Address.String(),
			To:   _to,
			Data: utils.BytesToHexString(_data),
		},
		id,
	)

	if estimateGas.Error != nil {
		panic(errors.ERROR_ESTIMATE_GAS(estimateGas.ErrorMessage()))
	}

	// sign tx
	signedTx, err := _signer.SignRawTx(
		_phrase,
		tx,
	)
	if err != nil {
		panic(err)
	}

	// signed tx to raw signed tx
	rawSignedTx := "0x" + signer.SignedRawTx2HexString(signedTx)

	// send tx
	result := ethcall.SendSignedRawTransaction(rawSignedTx, id)

	return result
}
