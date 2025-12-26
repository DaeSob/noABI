package types

import (
	"fmt"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_BytesToHash(t *testing.T) {
	bytes := []byte{5}
	hash := BytesToHash(bytes)

	var exp THash
	exp[31] = 5

	assert.Equal(t, hash, exp)
}

func Test_HexToHash(t *testing.T) {
	const hexString = "0xb26f2b342aab24bcf63ea218c6a9274d30ab9a15a218c6a9274d30ab9a151000"
	hash := HexToHash(hexString)

	assert.Equal(t, hash.String(), hexString)
}

func Test_Hash_Format(t *testing.T) {
	var hash THash
	hash.SetBytes([]byte{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
		0x10, 0x00,
	})

	assert.Equal(
		t,
		fmt.Sprintln(hash),
		"0xb26f2b342aab24bcf63ea218c6a9274d30ab9a15a218c6a9274d30ab9a151000\n",
	)
	assert.Equal(
		t,
		fmt.Sprint(hash),
		"0xb26f2b342aab24bcf63ea218c6a9274d30ab9a15a218c6a9274d30ab9a151000",
	)
}
