package types

import "cia/common/utils"

const (
	// length of the address
	AddressLength = 20
)

type TAddress [AddressLength]byte

func HexToAddress(_str string) TAddress {
	return BytesToAddress(utils.HexToBytes(_str))
}

func BytesToAddress(_byte []byte) TAddress {
	var a TAddress
	a.SetBytes(_byte)
	return a
}

func (a TAddress) Bytes() []byte {
	return a[:]
}

func (a TAddress) HexString() string {
	return utils.BytesToHexString(a.Bytes())
}

func (a TAddress) String() string {
	return a.HexString()
}

func (a *TAddress) SetBytes(_byte []byte) {
	if len(_byte) > len(a) {
		_byte = _byte[len(_byte)-AddressLength:]
	}

	copy(a[AddressLength-len(_byte):], _byte)
}
