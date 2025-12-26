package types

import (
	"cia/common/utils"
)

const (
	// length of the hash
	HashLength = 32
)

type THash [HashLength]byte

func BytesToHash(_byte []byte) THash {
	var h THash
	h.SetBytes(_byte)
	return h
}

func HexToHash(_str string) THash {
	return BytesToHash(utils.HexToBytes(_str))
}

func (h THash) Bytes() []byte {
	return h[:]
}

func (h THash) HexString() string {
	return utils.BytesToHexString(h.Bytes())
}

func (h THash) String() string {
	return h.HexString()
}

func (h *THash) SetBytes(_byte []byte) {
	if len(_byte) > len(h) {
		_byte = _byte[len(_byte)-HashLength:]
	}

	copy(h[HashLength-len(_byte):], _byte)
}
