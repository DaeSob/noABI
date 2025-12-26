package types

import (
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_GetType(t *testing.T) {
	var v = 1

	res := GetType(v).String()

	assert.Equal(t, res, "int")
}

func Test_TypeToString(t *testing.T) {
	var v = 1

	res := TypeToString(v)

	assert.Equal(t, res, "int")
}
