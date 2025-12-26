package types

import (
	"cia/common/blockchain/abi"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_GetType(t *testing.T) {
	eve := CreateEvent(
		"Transfer",
		[]abi.TArgumentStr{
			*abi.CreateAddressType("from", true),
			*abi.CreateAddressType("to", true),
			*abi.CreateUint256Type("value", false),
		},
	)

	assert.Equal(t, eve.Type, "event")
}

func Test_Sig(t *testing.T) {
	eve := CreateEvent(
		"Transfer",
		[]abi.TArgumentStr{
			*abi.CreateAddressType("from", true),
			*abi.CreateAddressType("to", true),
			*abi.CreateUint256Type("value", false),
		},
	)

	assert.Equal(t, eve.Sig(), "Transfer(address,address,uint256)")
}

func Test_EncodeSig(t *testing.T) {
	eve := CreateEvent(
		"Transfer",
		[]abi.TArgumentStr{
			*abi.CreateAddressType("from", true),
			*abi.CreateAddressType("to", true),
			*abi.CreateUint256Type("value", false),
		},
	)

	assert.Equal(t, eve.EncodeSig().String(), "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func Test_GetTopicOption(t *testing.T) {
	eve := CreateEvent(
		"Transfer",
		[]abi.TArgumentStr{
			*abi.CreateAddressType("from", true),
			*abi.CreateAddressType("to", true),
			*abi.CreateUint256Type("value", true),
		},
	)

	optsTopic := eve.GetTopicOption()
	assert.Equal(t, len(optsTopic), 4)
	assert.Equal(t, optsTopic[0].Name, "method")
	assert.Equal(t, optsTopic[0].Type, "function")
	assert.Equal(t, optsTopic[0].Indexed, true)

	assert.Equal(t, optsTopic[1].Name, "from")
	assert.Equal(t, optsTopic[1].Type, "address")
	assert.Equal(t, optsTopic[1].Indexed, true)

	assert.Equal(t, optsTopic[2].Name, "to")
	assert.Equal(t, optsTopic[2].Type, "address")
	assert.Equal(t, optsTopic[2].Indexed, true)

	assert.Equal(t, optsTopic[3].Name, "value")
	assert.Equal(t, optsTopic[3].Type, "uint256")
	assert.Equal(t, optsTopic[3].Indexed, true)
}

func Test_GetDecodedLog(t *testing.T) {

}

func Test_String(t *testing.T) {
	eve := CreateEvent(
		"Transfer",
		[]abi.TArgumentStr{
			*abi.CreateAddressType("from", true),
			*abi.CreateAddressType("to", true),
			*abi.CreateUint256Type("value", true),
		},
	)

	assert.Equal(
		t,
		eve.ToJsonString(false),
		`{"inputs":[{"components":null,"indexed":true,"internalType":"address","name":"from","type":"address"},{"components":null,"indexed":true,"internalType":"address","name":"to","type":"address"},{"components":null,"indexed":true,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"}`,
	)
}
