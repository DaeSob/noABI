package encode

import (
	"fmt"

	"cia/common/blockchain/abi"
	types "cia/common/blockchain/types"
)

func DecodeParams(_params []abi.TArgumentStr, _data []byte) ([]interface{}, error) {
	// 'unpack'은 non indexed argument만 decoding 해줌 (go-ethereum)
	// 따라서 _params의 모든 indexed 값을 false로 치환하여 unpack을 진행한다.
	var params []abi.TArgumentStr
	for _, param := range _params {
		// turn false
		tmp := param.Copy()
		tmp.Indexed = false

		params = append(params, *tmp)
	}

	args, err := getArguments(params)
	if err != nil {
		return nil, err
	}

	return args.Unpack(_data)
}

func DecodeParamsToString(_params []abi.TArgumentStr, _data []byte) ([]string, error) {
	decodes, err := DecodeParams(_params, _data)
	if err != nil {
		return []string{""}, err
	}

	result := []string{}
	for _, decode := range decodes {
		strDecode := fmt.Sprintf("%v", decode)
		result = append(result, strDecode)
	}

	return result, nil
}

func DecodeParam(_param abi.TArgumentStr, _data []byte) (interface{}, error) {
	decoded, err := DecodeParams([]abi.TArgumentStr{_param}, _data)
	if err != nil {
		return nil, err
	}

	return decoded[0], nil
}

func DecodeParamToString(_param abi.TArgumentStr, _data []byte) (string, error) {
	decode, err := DecodeParam(_param, _data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", decode), nil
}

func DecodeLog(_abiInputs []abi.TArgumentStr, _log types.TLog) map[string]interface{} {
	// sort indexed or non indexed
	var indexed, nonIndexed []abi.TArgumentStr

	for _, arg := range _abiInputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}

	// decode indexed
	var decodedIndexed []interface{}
	for i, arg := range indexed {
		if i == 0 {
			decodedIndexed = append(decodedIndexed, _log.Topics[0].String())
		} else {
			decoded, err := DecodeParam(arg, _log.Topics[i].Bytes())
			if err != nil {
				panic(err)
			}

			decodedIndexed = append(decodedIndexed, decoded)
		}
	}

	// decode non indexed
	decodedNonIndexed, err := DecodeParams(nonIndexed, _log.Data)
	if err != nil {
		panic(err)
	}

	// make result
	result := map[string]interface{}{}
	for _, arg := range _abiInputs {
		if arg.Indexed {
			result[arg.Name] = decodedIndexed[0]
			decodedIndexed = decodedIndexed[1:]
		} else {
			result[arg.Name] = decodedNonIndexed[0]
			decodedNonIndexed = decodedNonIndexed[1:]
		}
	}

	return result
}
