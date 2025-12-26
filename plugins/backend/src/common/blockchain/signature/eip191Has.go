package signature

import (
	"encoding/hex"

	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

func GenStringHashData(_str string) []byte {
	var data []interface{}
	data = append(data, Data("string", _str))

	hashData := solsha3.SoliditySHA3(data...)
	return solsha3.SoliditySHA3WithPrefix(
		solsha3.Bytes32("0x" + hex.EncodeToString(hashData)),
	)
}
