package types_test

import (
	"cia/common/blockchain/abi"
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/rpcRequest/eth"
	"cia/common/blockchain/types"
	"cia/common/blockchain/utils/encode"
	"cia/common/utils"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	// testify assert
)

func Test_NewTx(t *testing.T) {
	address := types.HexToAddress("0x259cB235B89dEa84E19818C58364Ab4aA9F101fd")
	tokenAddr := types.HexToAddress("0x0ba5e6f582d25053af10b1d84747310595675a43")

	// get nonce
	rpc := "http://10.10.24.95:8551"
	eth := eth.TEth{rpcRequest.TRPC{rpc}}
	nonce := eth.GetTransactionCount(address.HexString(), int64(1))

	gasPrice := eth.GetGasPrice(int64(1))

	data := encode.EncodeFunctionCallToHexString(
		*types.CreateMethod(
			"transfer",
			[]abi.TArgumentStr{
				*abi.CreateAddressType("from", true),
				*abi.CreateAddressType("to", true),
				*abi.CreateUint256Type("value", true),
			},
			[]abi.TArgumentStr{},
		),
		common.HexToAddress(address.HexString()),
		common.HexToAddress(address.HexString()),
		big.NewInt(1),
	)

	tx := types.NewTx(
		nonce.ResultToUint64(),
		tokenAddr,
		big.NewInt(0),
		uint64(0),
		big.NewInt(gasPrice.ResultToInt64()),
		utils.HexToBytes(data),
	)

	txBytes, _ := tx.MarshalBinary()

	reTx := new(ethTypes.Transaction)
	rlp.DecodeBytes(txBytes, &reTx)

	signedTx := "0xf8cf80850ba43b740080940ba5e6f582d25053af10b1d84747310595675a4380b864beabacc8000000000000000000000000259cb235b89dea84e19818c58364ab4aa9f101fd000000000000000000000000259cb235b89dea84e19818c58364ab4aa9f101fd0000000000000000000000000000000000000000000000000000000000000001882f00213bf6dfd1cca00d6f1a83bcd2578844abcb03b40b6b5e3be4ef0fc4709a7999fc6f5f2c27a94ba0650176928eddff4dcae11f218bda518b67082f351f7f92d825e6fd0c5415525c"
	signTx := new(ethTypes.Transaction)
	rlp.DecodeBytes(utils.HexToBytes(signedTx), &signTx)

	assert.Equal(
		t,
		utils.InterfaceToJsonString(tx, true),
		utils.InterfaceToJsonString(reTx, true),
	)

	assert.Equal(
		t,
		reTx.Gas(),
		signTx.Gas(),
	)

	assert.Equal(
		t,
		reTx.To(),
		signTx.To(),
	)

	assert.Equal(
		t,
		reTx.Nonce(),
		signTx.Nonce(),
	)

	assert.Equal(
		t,
		signTx.ChainId(),
		big.NewInt(1693371730605631700),
	)
}
