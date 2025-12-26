package abiMapper

import (
	types "cia/common/blockchain/types"
	"strings"
)

// key : encoded signature
type TSigMapper map[string]types.TEvent

// key : contract address
type TABIMappers map[string]TSigMapper

// key는 모두 lower case를 적용하도록 한다.
func (abiMapper *TABIMappers) SetSigMapper(_key string, _sig TSigMapper) {
	key := strings.ToLower(_key)
	(*abiMapper)[key] = _sig
}
