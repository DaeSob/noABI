package utils

import (
	"testing"
	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_LocalTimeToGMTString(t *testing.T) {
	curTime := Now()

	println(curTime.Unix())
	println(LocalTimeToGMTString(curTime))
}

func Test_GMTStringToLocalTime(t *testing.T) {
	gmtString := "Wed, 28 Jun 2023 05:13:28 GMT"
	localTimeUnix := GMTStringToLocalTime(gmtString).Unix()

	assert.Equal(t, localTimeUnix, int64(1687929208))
}
