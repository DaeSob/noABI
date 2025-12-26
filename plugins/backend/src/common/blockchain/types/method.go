package types

import (
	"fmt"
	"strings"

	"cia/common/blockchain/abi"
	"cia/common/utils"
)

// `CreateMethod` 함수로만 객체 생성
type TMethod struct {
	Name            string             `json:"name"`
	Type            string             `json:"type"`
	Inputs          []abi.TArgumentStr `json:"inputs"`
	Outputs         []abi.TArgumentStr `json:"outputs"`
	StateMutability string             `json:"stateMutability,omitempty"`
}

// `CreateMethod` 함수로만 객체 생성
func CreateMethod(
	_name string,
	_inputs []abi.TArgumentStr,
	_outputs []abi.TArgumentStr,
) *TMethod {
	var mth TMethod
	mth.Name = _name
	mth.Type = "function"
	mth.Inputs = _inputs
	mth.Outputs = _outputs

	return &mth
}

func CreateMethodEx(
	_type string,
	_name string,
	_stateMutability string,
	_inputs []abi.TArgumentStr,
	_outputs []abi.TArgumentStr,
) *TMethod {
	var mth TMethod
	mth.Name = _name
	mth.Type = _type
	mth.StateMutability = _stateMutability
	mth.Inputs = _inputs
	mth.Outputs = _outputs

	return &mth
}

func (m TMethod) getInputTypes() []string {
	var types []string
	for _, input := range m.Inputs {
		types = append(types, input.Type)
	}

	return types
}

func (m TMethod) Sig() string {
	types := m.getInputTypes()

	return fmt.Sprintf("%v(%v)", m.Name, strings.Join(types, ","))
}

func (m TMethod) ToJsonString(_pretty bool) string {
	return utils.InterfaceToJsonString(m, _pretty)
}
