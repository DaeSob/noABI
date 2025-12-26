package events

import (
	"cia/common/blockchain/abi"
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/rpcRequest/eth"
	"cia/common/blockchain/utils/encode"
	util "cia/common/utils"
	rpc "cia/common/utils/rpcRequest"

	"github.com/ethereum/go-ethereum/common"
)

func getTransferSig() string {
	hash := encode.EncodeEventSignature(
		"Transfer",
		"address",
		"address",
		"uint256",
	)

	return hash.String()
}

func _encodeAddress(_address string) string {
	res, _ := encode.EncodeParam(
		*abi.CreateAddressType("", false),
		common.HexToAddress(_address),
	)
	return util.BytesToHexString(res)
}

func GetTransferLog(
	_rpc string,
	_fromBlock uint64,
	_toBlock uint64,
	_tokenAddress string,
	_fromAddress string,
	_toAddress string,
	_id int64,
) *rpc.TRpcResponse {
	// ready eth struct
	eth := &eth.TEth{rpcRequest.TRPC{_rpc}}

	// ready topics
	topics := make([]interface{}, 3)

	// sig
	topics[0] = getTransferSig()

	// from address
	if _fromAddress != "" {
		topics[1] = _encodeAddress(_fromAddress)
	}

	// to address
	if _toAddress != "" {
		topics[2] = _encodeAddress(_toAddress)
	}

	// get log
	res := eth.GetLogs(_fromBlock, _toBlock, []string{_tokenAddress}, topics, _id)

	return res
}
