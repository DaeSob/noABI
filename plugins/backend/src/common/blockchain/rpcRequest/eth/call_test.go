package eth

import (
	"cia/common/blockchain/abi"
	"cia/common/blockchain/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CallMethod(t *testing.T) {

	rpc := "https://baas-rpc.luniverse.io:18545?lChainId=1666232515295616488"
	from := "0x052cf2e7809a82b9be6f50d0eb7d90ecd28e5932"
	to := "0x7896619f1d5c9ad9c1dc922dfc9b8597d9a5199e"
	id := int64(1)

	mtd := *types.CreateMethod(
		"totalSupply",
		[]abi.TArgumentStr{},
		[]abi.TArgumentStr{
			*abi.CreateUint256Type("totalSupply", false),
		},
	)

	totalSupply := CallMethod(
		mtd,
		[]interface{}{},
		rpc,
		from,
		to,
		id,
	)

	assert.Equal(t, totalSupply[0], "1400000000000000000000000000")
}
