package preference

import (
	"strings"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - server
func GetHost() string {
	inst := GetInstance()
	return inst.tServer.Host
}

func IsDebugMode() bool {
	inst := GetInstance()

	return inst.tServer.Debug
}

func PathPrefix() string {
	inst := GetInstance()

	return inst.tServer.PathPrefix
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - authorization
func IsEnableAuth() bool {
	inst := GetInstance()

	return inst.tAuthorization.Enable
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - auth address list
func IsAuthSigner(_address string) bool {
	inst := GetInstance()

	if inst.tAuthorization.Signer == strings.ToLower(_address) {
		return true
	}
	return false
}
