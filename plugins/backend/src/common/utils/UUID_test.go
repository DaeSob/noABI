package utils

import (
	"testing"
	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_CreateUUID(t *testing.T) {
	uuid1 := CreateUUID()
	uuid2 := CreateUUID()

	assert.NotEqual(
		t,
		uuid1,
		uuid2,
	)
}
