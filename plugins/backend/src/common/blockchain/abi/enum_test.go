package abi

import (
	"testing"
	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_TypeKind_String(t *testing.T) {
	assert.Equal(
		t,
		TYPE_ADDRESS.String(),
		"address",
	)

	assert.Equal(
		t,
		TYPE_STRING.String(),
		"string",
	)

	assert.Equal(
		t,
		TYPE_UINT8.String(),
		"uint8",
	)
	assert.Equal(
		t,
		TYPE_UINT16.String(),
		"uint16",
	)
	assert.Equal(
		t,
		TYPE_UINT32.String(),
		"uint32",
	)
	assert.Equal(
		t,
		TYPE_UINT64.String(),
		"uint64",
	)
	assert.Equal(
		t,
		TYPE_UINT128.String(),
		"uint128",
	)
	assert.Equal(
		t,
		TYPE_UINT256.String(),
		"uint256",
	)

	assert.Equal(
		t,
		TYPE_INT8.String(),
		"int8",
	)
	assert.Equal(
		t,
		TYPE_INT16.String(),
		"int16",
	)
	assert.Equal(
		t,
		TYPE_INT32.String(),
		"int32",
	)
	assert.Equal(
		t,
		TYPE_INT64.String(),
		"int64",
	)
	assert.Equal(
		t,
		TYPE_INT128.String(),
		"int128",
	)
	assert.Equal(
		t,
		TYPE_INT256.String(),
		"int256",
	)

	assert.Equal(
		t,
		TYPE_TUPLE.String(),
		"tuple",
	)

	assert.Equal(
		t,
		TYPE_ARRAY_ADDRESS.String(),
		"address[]",
	)
	assert.Equal(
		t,
		TYPE_ARRAY_BYTES.String(),
		"bytes[]",
	)
	assert.Equal(
		t,
		TYPE_ARRAY_UINT256.String(),
		"uint256[]",
	)
	assert.Equal(
		t,
		TYPE_ARRAY_STRING.String(),
		"string[]",
	)
}
