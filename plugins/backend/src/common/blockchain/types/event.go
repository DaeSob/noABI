package types

import (
	"fmt"
	"strings"

	"cia/common/blockchain/abi"
	"cia/common/utils"

	crypto "github.com/ethereum/go-ethereum/crypto"
)

// `CreateEvent` 함수로만 객체를 생성할 것!
type TEvent struct {
	Type   string             `json:"type"`
	Name   string             `json:"name"`
	Inputs []abi.TArgumentStr `json:"inputs"`
}

// `CreateEvent` 함수로만 객체를 생성할 것!
func CreateEvent(_name string, _inputs []abi.TArgumentStr) *TEvent {
	var eve TEvent
	eve.Type = "event"
	eve.Name = _name
	eve.Inputs = _inputs
	return &eve
}

func (e TEvent) Sig() string {
	var types []string
	for _, input := range e.Inputs {
		types = append(types, input.Type)
	}

	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ","))
}

func (e TEvent) EncodeSig() THash {
	return BytesToHash(crypto.Keccak256([]byte(e.Sig())))
}

func (e TEvent) GetTopicOption() []abi.TArgumentStr {
	var opts []abi.TArgumentStr

	// method
	opts = append(opts, *abi.CreateDynamicType("method", "function", true))

	// inputs
	opts = append(opts, e.Inputs...)

	return opts
}

func (e TEvent) ToJsonString(_pretty bool) string {
	return utils.InterfaceToJsonString(e, _pretty)
}
