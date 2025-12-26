package unit

import (
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_ConvertFromWei(t *testing.T) {
	wei := ConvertFromWei("1", "ether")

	assert.Equal(t, wei, "0.000000000000000001")
}

func Test_ConvertToWei(t *testing.T) {
	wei := ConvertToWei("1", "ether")

	assert.Equal(t, wei, "1000000000000000000")
}
