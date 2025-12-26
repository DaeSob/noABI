package util

import (
	"testing"
)

func Test_MakeBridgeResponse(t *testing.T) {
	type TTemp struct {
		ABC string `json:"abc"`
		DEF string `json:"def"`
	}

	// res := MakeBridgeResponse(
	// 	"1",
	// 	TTemp{
	// 		"ABC",
	// 		"DEF",
	// 	},
	// 	"",
	// 	"",
	// 	false,
	// )

	// assert.Equal(
	// 	t,
	// 	res,
	// 	`{"id":"1","result":{"code":"","data":{"abc":"ABC","def":"DEF"},"error":""}}`,
	// )
}
