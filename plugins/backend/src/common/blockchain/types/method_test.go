package types

import (
	"cia/common/blockchain/abi"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_MethodType(t *testing.T) {
	mtd := *CreateMethod(
		"name",
		[]abi.TArgumentStr{},
		[]abi.TArgumentStr{
			*abi.CreateStringType("name", false),
		},
	)

	assert.Equal(t, mtd.Type, "function")
}

func Test_MethodSig(t *testing.T) {
	mtd := *CreateMethod(
		"name",
		[]abi.TArgumentStr{},
		[]abi.TArgumentStr{
			*abi.CreateStringType("name", false),
		},
	)

	assert.Equal(t, mtd.Sig(), "name()")
}

func Test_MethodToJsonString(t *testing.T) {
	mtd := *CreateMethod(
		"name",
		[]abi.TArgumentStr{},
		[]abi.TArgumentStr{
			*abi.CreateStringType("name", false),
		},
	)

	assert.Equal(
		t,
		mtd.ToJsonString(false),
		`{"inputs":[],"name":"name","outputs":[{"components":null,"indexed":false,"internalType":"string","name":"name","type":"string"}],"type":"function"}`,
	)
}
