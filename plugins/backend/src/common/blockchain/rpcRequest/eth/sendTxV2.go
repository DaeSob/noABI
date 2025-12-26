package eth

import (
	"cia/common/blockchain/keystore/signer"
	"cia/common/errors"
	"cia/common/utils"
	request "cia/common/utils/rpcRequest"

	"github.com/ethereum/go-ethereum/common"
)

func SendRawTx(_ethCall *TEth,
	_signer *signer.TSigner,
	_phrase, _to, _nonce, _value, _gasPrice, _gasLimit string,
	_callData []byte) *request.TRpcResponse {

	nonce := utils.StringToUint64V2(_nonce)
	value := utils.StringToBigInt(_value)
	gasPrice := utils.StringToBigInt(_gasPrice)
	gasLimit := utils.StringToUint64V2(_gasLimit)

	// make tx
	tx := signer.TRawTransaction{
		From:     _signer.Address,
		To:       common.HexToAddress(_to),
		Value:    value,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		Data:     _callData,
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
	return _ethCall.SendSignedRawTransaction(rawSignedTx, int64(0))

}

type TSendTxParams struct {
	To          string        `json:"to"`
	Nonce       string        `json:"nonce"`
	Value       string        `json:"value,omitempty"`
	GasPrice    string        `json:"gasPrice"`
	GasLimit    string        `json:"gasLimit"`
	CallData    []byte        `json:"callData"`
	GasPriceTip uint64        `json:"gasPriceTip"`
	GasLimitTip uint64        `json:"gasLimitTip"`
	Opts        []interface{} `json:"options"`
}

func SendTx(
	_ethCall *TEth,
	_signer *signer.TSigner, _phrase string,
	_requestId string,
	_params TSendTxParams,
) *request.TRpcResponse {

	id := int64(0)
	if len(_requestId) > 0 {
		id, _ = utils.StringToI64(_requestId)
	}

	// gas
	if len(_params.GasPrice) == 0 {
		_params.GasPrice = _ethCall.GetGasPrice(id, _params.Opts).ResultToString()
		if _params.GasPriceTip > 0 {
			gasPrice := utils.StringToUint64V2(_params.GasPrice)
			gasPrice = gasPrice + (gasPrice * _params.GasPriceTip / 10000)
			_params.GasPrice = utils.Uint64ToString(gasPrice)
		}
	}

	// nonce
	if len(_params.Nonce) == 0 {
		_params.Nonce = _ethCall.GetTransactionCount(_signer.Address.String(), id).ResultToString()
	}

	// estimate gas
	if len(_params.GasLimit) == 0 {
		resEstimateGas := _ethCall.SimpleEstimateGas(
			TSimpleTransaction{
				From: _signer.Address.String(),
				To:   _params.To,
				Data: utils.BytesToHexString(_params.CallData),
			},
			id,
		)
		if resEstimateGas.Error != nil {
			panic(errors.ERROR_ESTIMATE_GAS(resEstimateGas.ErrorMessage()))
		}
		//측정된 Gas Limit에 10% 더 허용
		tempGasLimit := utils.StringToUint64V2(resEstimateGas.ResultToString())
		tempGasLimit = tempGasLimit + (tempGasLimit * _params.GasLimitTip / 10000)
		_params.GasLimit = utils.Uint64ToString(tempGasLimit)
	}

	return SendRawTx(
		_ethCall,
		_signer, _phrase,
		_params.To,
		_params.Nonce,
		_params.Value,
		_params.GasPrice,
		_params.GasLimit,
		_params.CallData,
	)

}
