package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type TTransaction ethTypes.Transaction

func NewTx(
	_nonce uint64,
	_to TAddress,
	_value *big.Int,
	_gasLimit uint64,
	_gasPrice *big.Int,
	_data []byte,
) *ethTypes.Transaction {
	return ethTypes.NewTransaction(
		_nonce,
		common.HexToAddress(_to.HexString()),
		_value,
		_gasLimit,
		_gasPrice,
		_data,
	)
}
