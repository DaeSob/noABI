package encode

import (
	"testing"

	"cia/common/blockchain/abi"
	types "cia/common/blockchain/types"
	util "cia/common/utils"

	"github.com/ethereum/go-ethereum/common"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_EncodeFunctionCallToByte(t *testing.T) {
	// [{
	//		"inputs":[],
	//		"name":"name",
	//		"outputs":[{"name":"","type":"string"}],
	//		"type":"function"
	// }]
	input := []abi.TArgumentStr{}
	out := []abi.TArgumentStr{*abi.CreateStringType("", false)}

	testAbiItem := *types.CreateMethod("name", input, out)
	result, err := EncodeFunctionCallToByte(testAbiItem)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, util.BytesToHexString(result), "0x06fdde03")
}

func Test_EncodeFunctionCallHexString(t *testing.T) {
	// [{
	//		"inputs":[{"name":"","type":"address"}],
	//		"name":"isBar",
	//		"outputs":[{"name":"","type":"bool"}],
	//		"type":"function"
	// }]
	input := []abi.TArgumentStr{*abi.CreateAddressType("", false)}
	out := []abi.TArgumentStr{*abi.CreateBoolType("", false)}
	testAbiItem := *types.CreateMethod("isBar", input, out)

	result := EncodeFunctionCallToHexString(testAbiItem, common.HexToAddress("01"))

	assert.Equal(t, result, "0x1f2c40920000000000000000000000000000000000000000000000000000000000000001")
}

func Test_EncodeEvnetSignature(t *testing.T) {
	// transfer(address,address,uint256)
	const funcName = "transfer"
	const arg1 = "address"
	const arg2 = "address"
	const arg3 = "uint256"

	hash := EncodeEventSignature(funcName, arg1, arg2, arg3)

	const exp = "0xbeabacc8ffedac16e9a60acdb2ca743d80c2ebb44977a93fa8e483c74d2b35a8"
	assert.Equal(t, hash.String(), exp)
}

func Test_EncodeParams(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateAddressType("", false),
		*abi.CreateAddressType("", false),
	}

	encodeParams, err := EncodeParams(
		params,
		common.HexToAddress("0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd"),
		common.HexToAddress("0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E"),
	)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(encodeParams), 64)
	assert.Equal(
		t,
		util.BytesToHexString(encodeParams),
		"0x000000000000000000000000e9119ba33d4fff07cebf5a5f9f58a1ba14e127fd0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e",
	)
}

func Test_EncodeParams2(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateArrayAddressType("", false),
		*abi.CreateArrayBytesType("", false),
	}

	encodeParams, err := EncodeParams(
		params,
		[]common.Address{
			common.HexToAddress("0x96b6b4f90bbf2468d41f8030f377e253ee2da546"),
		},
		[][]byte{
			util.HexToBytes("0x095ea7b30000000000000000000000003fab080ce2d4b7b31508b368df15380e4f6f5724000000000000000000000000000000000000000000084595161401484a000000"),
		},
	)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(encodeParams), 320)
	assert.Equal(
		t,
		util.BytesToHexString(encodeParams),
		"0x00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000000000000096b6b4f90bbf2468d41f8030f377e253ee2da546000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000003fab080ce2d4b7b31508b368df15380e4f6f5724000000000000000000000000000000000000000000084595161401484a00000000000000000000000000000000000000000000000000000000000000",
	)
}

func Test_EncodeParamsTuple(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateTuple(
			"JPA",
			"TJPA",
			[]abi.TArgumentStr{
				*abi.CreateAddressType("FieldOne", false),
				*abi.CreateAddressType("FieldTwo", false),
			},
		),
	}

	encodeParams, err := EncodeParams(
		params,
		struct {
			FieldOne common.Address
			FieldTwo common.Address
		}{
			common.HexToAddress("0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd"),
			common.HexToAddress("0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E"),
		},
	)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(encodeParams), 64)
	assert.Equal(
		t,
		util.BytesToHexString(encodeParams),
		"0x000000000000000000000000e9119ba33d4fff07cebf5a5f9f58a1ba14e127fd0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e",
	)
}

func Test_EncodeParam(t *testing.T) {
	param := *abi.CreateStringType("", false)

	encodeParam, err := EncodeParam(param, "meta test")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(encodeParam), 96)
	assert.Equal(
		t,
		util.BytesToHexString(encodeParam),
		"0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096d65746120746573740000000000000000000000000000000000000000000000",
	)
}
