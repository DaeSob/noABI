package encode

import (
	"encoding/json"
	"fmt"
	"strings"

	"cia/common/blockchain/abi"
	types "cia/common/blockchain/types"
	utils "cia/common/utils"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	crypto "github.com/ethereum/go-ethereum/crypto"

	jMapper "cia/common/utils/json/mapper"
)

func EncodeFunctionCallToByte(_functionAbi types.TMethod, _args ...interface{}) ([]byte, error) {
	// abi item to json
	bytes, err := json.Marshal(_functionAbi)
	if err != nil {
		return nil, err
	}

	jMap, err := jMapper.NewBytes(bytes)
	if err != nil {
		return nil, err
	}

	strObjAbiItem := "[" + jMap.Print() + "]"
	funcAbi, err := ethAbi.JSON(strings.NewReader(strObjAbiItem))
	if err != nil {
		return nil, err
	}

	// encoding
	return funcAbi.Pack(_functionAbi.Name, _args...)
}

func EncodeFunctionCallToHexString(_functionAbi types.TMethod, _args ...interface{}) string {
	bytes, err := EncodeFunctionCallToByte(_functionAbi, _args...)
	if err != nil {
		panic(err)
	}

	return utils.BytesToHexString(bytes)
}

func EncodeEventSignature(_funcName string, _args ...string) types.THash {
	sig := fmt.Sprintf("%v(%v)", _funcName, strings.Join(_args, ","))
	return types.BytesToHash(crypto.Keccak256([]byte(sig)))
}

func getArguments(_params []abi.TArgumentStr) (ethAbi.Arguments, error) {
	args := ethAbi.Arguments{}
	for _, param := range _params {
		// element to json
		arg := ethAbi.Argument{}
		err := arg.UnmarshalJSON(param.Bytes())
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return args, nil
}

func EncodeParams(_params []abi.TArgumentStr, _args ...interface{}) ([]byte, error) {
	args, err := getArguments(_params)
	if err != nil {
		return nil, err
	}

	return args.Pack(_args...)
}

func EncodeParam(_param abi.TArgumentStr, arg interface{}) ([]byte, error) {
	return EncodeParams([]abi.TArgumentStr{_param}, arg)
}
