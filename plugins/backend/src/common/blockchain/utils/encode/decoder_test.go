package encode

import (
	"encoding/hex"
	"fmt"
	"testing"

	"cia/common/blockchain/abi"
	types "cia/common/blockchain/types"
	util "cia/common/utils"
	jMapper "cia/common/utils/json/mapper"

	// testify assert
	common "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func Test_DecodeParam(t *testing.T) {
	param := *abi.CreateStringType("", false)
	data, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096d65746120746573740000000000000000000000000000000000000000000000")

	decode, err := DecodeParam(param, data)
	if err != nil {
		t.Error(err)
	}

	result := fmt.Sprintf("%v", decode)

	assert.Equal(t, result, "meta test")
}

func Test_DecodeParamToString(t *testing.T) {
	param := *abi.CreateStringType("", false)
	data, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096d65746120746573740000000000000000000000000000000000000000000000")

	decode, err := DecodeParamToString(param, data)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, decode, "meta test")
}

func Test_DecodeParams(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateAddressType("jpwOwner", false),
	}
	data, _ := hex.DecodeString("000000000000000000000000b171fe0b0804651446a50344ae14e56596190bcf")

	decode, err := DecodeParams(params, data)
	if err != nil {
		t.Error(err)
	}

	addr1 := fmt.Sprintf("%v", decode[0])
	addr2 := fmt.Sprintf("%v", decode[1])

	assert.Equal(t, addr1, "0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd")
	assert.Equal(t, addr2, "0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E")
}

func Test_DecodeParams2(t *testing.T) {
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
	data, _ := hex.DecodeString("000000000000000000000000e9119ba33d4fff07cebf5a5f9f58a1ba14e127fd0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e")

	decode, err := DecodeParams(params, data)
	if err != nil {
		t.Error(err)
	}

	// [{0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd 0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E}]
	// fmt.Printf("%v \n", decode[0])

	type TTemp struct {
		FieldOne common.Address `json:"FieldOne"`
		FieldTwo common.Address `json:"FieldTwo"`
	}
	var tmp TTemp
	bytes, _ := jMapper.ToJson(decode[0])
	jMapper.FromJson(bytes, &tmp)

	assert.Equal(
		t,
		tmp.FieldOne.String(),
		"0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd",
	)
	assert.Equal(
		t,
		tmp.FieldTwo.String(),
		"0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E",
	)
}

func Test_DecodeParamsToString(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateAddressType("", false),
		*abi.CreateAddressType("", false),
	}
	data, _ := hex.DecodeString("000000000000000000000000e9119ba33d4fff07cebf5a5f9f58a1ba14e127fd0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e")

	decode, err := DecodeParamsToString(params, data)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, decode[0], "0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd")
	assert.Equal(t, decode[1], "0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E")
}

func Test_DecodeParamsToString2(t *testing.T) {
	params := []abi.TArgumentStr{
		*abi.CreateAddressType("", false),
		*abi.CreateAddressType("", false),
	}
	data, _ := hex.DecodeString("000000000000000000000000e9119ba33d4fff07cebf5a5f9f58a1ba14e127fd0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e")

	decode, err := DecodeParamsToString(params, data)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, decode[0], "0xe9119bA33d4FFf07CEbf5A5F9F58A1bA14e127Fd")
	assert.Equal(t, decode[1], "0x3cAe9255b7AD17Df118ed4F2338ff80C3ee3a15E")
}

func Test_DecodeLog(t *testing.T) {
	inputs := []abi.TArgumentStr{
		*abi.CreateDynamicType("method", "function", true),
		*abi.CreateAddressType("from", true),
		*abi.CreateAddressType("to", true),
		*abi.CreateUint256Type("value", false),
	}

	topics := []types.THash{
		types.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		types.HexToHash("0x000000000000000000000000afcfa70428f457bea047800bb1cea5b2dadca622"),
		types.HexToHash("0x00000000000000000000000099217904f418d1ede238c8a104d15603070cfa1e"),
	}

	logs := types.TLog{
		false,
		0x29,
		0x5,
		types.HexToHash("0xa1cf7d8cd03cd1e5210495393d119a5fa9039404ae9cee7d3bb0d63963df4dc0"),
		types.HexToHash("0x8a0ab9eaba5f7d720cbd8f8e8582d2ac1123833a30b7b660f662d7b5bf79e667"),
		0x562b940,
		types.HexToAddress("0xcee8faf64bb97a73bb51e115aa89c17ffa8dd167"),
		topics,
		util.HexToBytes("0x00000000000000000000000000000000000000000000000000000000000bb954"),
	}

	decodedLog := DecodeLog(inputs, logs)

	keys := make([]string, 0, len(decodedLog))
	for k := range decodedLog {
		keys = append(keys, k)
	}

	assert.Equal(t, len(keys), 4)

	/*
		// example of decode log
		res1 := fmt.Sprintf("%v", decodedLog[keys[0]])
		res2 := fmt.Sprintf("%v", decodedLog[keys[1]])
		res3 := fmt.Sprintf("%v", decodedLog[keys[2]])
		res4 := fmt.Sprintf("%v", decodedLog[keys[3]])

		println(res1)
		println(res2)
		println(res3)
		println(res4)

		bytes, _ := jMapper.ToJson(decodedLog)
		jMap, _ := jMapper.NewBytes(bytes)

		println(jMap.PPrint())
	*/
}

func Test_DecodeLog2(t *testing.T) {
	inputs := []abi.TArgumentStr{
		*abi.CreateDynamicType("method", "function", true),
	}

	topics := []types.THash{
		types.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		types.HexToHash("0x000000000000000000000000afcfa70428f457bea047800bb1cea5b2dadca622"),
		types.HexToHash("0x00000000000000000000000099217904f418d1ede238c8a104d15603070cfa1e"),
	}

	logs := types.TLog{
		false,
		0x29,
		0x5,
		types.HexToHash("0xa1cf7d8cd03cd1e5210495393d119a5fa9039404ae9cee7d3bb0d63963df4dc0"),
		types.HexToHash("0x8a0ab9eaba5f7d720cbd8f8e8582d2ac1123833a30b7b660f662d7b5bf79e667"),
		0x562b940,
		types.HexToAddress("0xcee8faf64bb97a73bb51e115aa89c17ffa8dd167"),
		topics,
		util.HexToBytes("0x00000000000000000000000000000000000000000000000000000000000bb954"),
	}

	decodedLog := DecodeLog(inputs, logs)

	keys := make([]string, 0, len(decodedLog))
	for k := range decodedLog {
		keys = append(keys, k)
	}

	res1 := fmt.Sprintf("%v", decodedLog[keys[0]])
	assert.Equal(t, len(keys), 1)
	assert.Equal(t, res1, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	/*
		// example of decode log
		res1 := fmt.Sprintf("%v", decodedLog[keys[0]])
		res2 := fmt.Sprintf("%v", decodedLog[keys[1]])
		res3 := fmt.Sprintf("%v", decodedLog[keys[2]])
		res4 := fmt.Sprintf("%v", decodedLog[keys[3]])

		println(res1)
		println(res2)
		println(res3)
		println(res4)

		bytes, _ := jMapper.ToJson(decodedLog)
		jMap, _ := jMapper.NewBytes(bytes)

		println(jMap.PPrint())
	*/
}
