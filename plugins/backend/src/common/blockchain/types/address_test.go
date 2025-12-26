package types

import (
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_HexToString(t *testing.T) {
	assert.Equal(
		t,
		HexToAddress("0x1").String(),
		"0x0000000000000000000000000000000000000001",
	)
	assert.Equal(t,
		HexToAddress("00000000000000000000000000000000000000001").String(),
		"0x0000000000000000000000000000000000000001",
	)
	assert.Equal(t,
		HexToAddress("0000000000000000000000000000000000000001").String(),
		"0x0000000000000000000000000000000000000001",
	)
}
