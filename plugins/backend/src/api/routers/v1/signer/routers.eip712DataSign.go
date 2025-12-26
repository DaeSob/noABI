package api_sgner

import (
	"cia/api/preference"
	apiUtil "cia/api/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cia/common/blockchain/keystore/signer"
	"cia/common/utils"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
)

type TREQ_EIP712DataSign struct {
	RequestId   string      `json:"requestId"`
	Doamin      interface{} `json:"domain,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	PrimaryType string      `json:"primaryType,omitempty"`
	KeyPair     interface{} `json:"KeyPair,omitempty"`
}

type TEIP712Domain struct {
	ChainId           string `json:"chainId"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	VerifyingContract string `json:"verifyingContract"`
}

// TSigData: value 하나로 단일/배열 모두 처리
type TEIP712Data struct {
	Type  string      `json:"type"`  // 예: "address", "uint256", "address[]"
	Name  string      `json:"name"`  // 필드명
	Value interface{} `json:"value"` // string 또는 []string
}

// convertInterfaceToArray: interface{}를 []interface{}로 변환하는 공통 함수
func convertInterfaceToArray(value interface{}, fieldName string) ([]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("%s is required", fieldName)
	}

	// []interface{}인지 확인
	keyPairArray, ok := value.([]interface{})
	if !ok {
		// []interface{}가 아니면 json.Marshal/Unmarshal로 변환 시도
		keyPairBytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %s: %v", fieldName, err)
		}
		var result []interface{}
		if err := json.Unmarshal(keyPairBytes, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %v", fieldName, err)
		}
		return result, nil
	}

	if len(keyPairArray) == 0 {
		return nil, fmt.Errorf("%s array is empty", fieldName)
	}

	return keyPairArray, nil
}

// convertValueToString: Value를 string으로 변환 (단일 값)
func convertValueToString(value interface{}, fieldName string) (string, error) {
	switch val := value.(type) {
	case string:
		return val, nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("field %s: single type requires string or number value", fieldName)
	}
}

// convertValueToStringArray: Value를 []string으로 변환 (배열 값)
func convertValueToStringArray(value interface{}, fieldName string) ([]string, error) {
	switch val := value.(type) {
	case []string:
		return val, nil
	case []interface{}:
		// JSON 배열은 기본적으로 []interface{}로 들어옴 → 변환 필요
		arr := make([]string, len(val))
		for i, iv := range val {
			strVal, ok := iv.(string)
			if !ok {
				return nil, fmt.Errorf("field %s: index %d value is not string", fieldName, i)
			}
			arr[i] = strVal
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("field %s: array type requires []string or []interface{} value", fieldName)
	}
}

// genTypedData: Types + Message 모두 자동 구성하는 함수
func genTypedData(_eip712Domain TEIP712Domain, _eip712DataType []TEIP712Data, _primaryType string) apitypes.TypedData {

	// ------------------------------------------
	// 1. 기본 Domain 타입 설정
	// ------------------------------------------
	types := apitypes.Types{
		"EIP712Domain": []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
	}

	// ------------------------------------------
	// 2. TEIP712SigData → apitypes.Type 변환 (Type, Name만 필요)
	// ------------------------------------------
	var convertedTypes []apitypes.Type

	for _, item := range _eip712DataType {
		convertedTypes = append(convertedTypes, apitypes.Type{
			Name: item.Name,
			Type: item.Type,
		})
	}

	// PrimaryType 타입 등록
	if _primaryType == "" {
		_primaryType = "EIP712Data" // 기본값
	}
	types[_primaryType] = convertedTypes

	// ------------------------------------------
	// 3. Message 자동 생성
	//    Type 이 "xxx[]" 이면 배열 처리
	//    그렇지 않으면 단일 값 처리
	// ------------------------------------------
	message := make(map[string]interface{})

	for _, item := range _eip712DataType {

		// 타입 끝이 []이면 배열 타입
		isArray := strings.HasSuffix(item.Type, "[]")

		if isArray {
			// 배열 값 처리
			arr, err := convertValueToStringArray(item.Value, item.Name)
			if err != nil {
				panic(err.Error())
			}
			message[item.Name] = arr
		} else {
			// 단일 값 처리
			value, err := convertValueToString(item.Value, item.Name)
			if err != nil {
				panic(err.Error())
			}
			message[item.Name] = value
		}
	}

	// ------------------------------------------
	// 4. TypedData 반환
	// ------------------------------------------
	return apitypes.TypedData{
		Types:       types,
		PrimaryType: _primaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              _eip712Domain.Name,
			Version:           _eip712Domain.Version,
			ChainId:           math.NewHexOrDecimal256(int64(utils.StringToUint64V2(_eip712Domain.ChainId))),
			VerifyingContract: _eip712Domain.VerifyingContract,
		},
		Message: message,
	}
}

func eip712DataSign(
	_eip712Data apitypes.TypedData,
	_signers []interface{},
) (sigHash []string, signature []string) {

	for _, pair := range _signers {
		keystorePair := pair.(map[string]interface{})
		signer := signer.NewSigner(preference.GetKeystorePath(), utils.StringToBigInt("0x1"), common.HexToAddress(keystorePair["account"].(string)))
		hash, sig, err := signer.EIP712DataSignature(keystorePair["phrase"].(string),
			_eip712Data,
		)
		if err != nil {
			panic(err)
		}
		sigHash = append(sigHash, utils.BytesToHexString(hash))
		signature = append(signature, utils.BytesToHexString(sig))
	}
	return

}

func EIP712DataSign(_c *gin.Context) {
	// body 읽기 (Gin의 GetRawData 사용)
	bytes, err := _c.GetRawData()
	if err != nil {
		panic(fmt.Sprintf("failed to read request body: %v", err))
	}

	// 빈 바이트 체크
	if len(bytes) == 0 {
		panic("request body is empty")
	}

	var param TREQ_EIP712DataSign

	// JSON 파싱 시도
	if err := json.Unmarshal(bytes, &param); err != nil {
		// JSON이 문자열로 이중 인코딩된 경우 처리
		var jsonStr string
		if err2 := json.Unmarshal(bytes, &jsonStr); err2 == nil {
			// 문자열로 파싱 성공하면 다시 JSON으로 파싱 시도
			if err3 := json.Unmarshal([]byte(jsonStr), &param); err3 == nil {
				// 성공적으로 파싱됨
			} else {
				// 여전히 실패하면 원래 에러와 함께 상세 정보 제공
				errorMsg := fmt.Sprintf("failed to parse request body (even after string decode): %v", err3)
				if len(jsonStr) > 200 {
					errorMsg += fmt.Sprintf(" (first 200 chars: %s...)", jsonStr[:200])
				} else {
					errorMsg += fmt.Sprintf(" (body: %s)", jsonStr)
				}
				panic(errorMsg)
			}
		} else {
			// 원래 에러와 함께 상세 정보 제공
			errorMsg := fmt.Sprintf("failed to parse request body: %v", err)
			if len(bytes) > 200 {
				errorMsg += fmt.Sprintf(" (first 200 chars: %s...)", string(bytes[:200]))
			} else {
				errorMsg += fmt.Sprintf(" (body: %s)", string(bytes))
			}
			panic(errorMsg)
		}
	}

	// error handling
	defer apiUtil.APIErrorHandler(_c, param.RequestId)

	// -----------------------------------------------------------
	// 1. interface{} → TEIP712Domain 변환
	// -----------------------------------------------------------
	var domain TEIP712Domain
	{
		if param.Doamin == nil {
			panic("domain is required")
		}
		domainBytes, err := json.Marshal(param.Doamin)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal domain: %v", err))
		}
		if err := json.Unmarshal(domainBytes, &domain); err != nil {
			panic(fmt.Sprintf("failed to unmarshal domain: %v", err))
		}
	}

	// -----------------------------------------------------------
	// 2. interface{} → []TEIP712Data 변환
	// -----------------------------------------------------------
	var data []TEIP712Data
	{
		if param.Data == nil {
			panic("data is required")
		}

		// param.Data가 []interface{}인지 확인
		dataArray, ok := param.Data.([]interface{})
		if !ok {
			// []interface{}가 아니면 json.Marshal/Unmarshal로 변환 시도
			dataBytes, err := json.Marshal(param.Data)
			if err != nil {
				panic(fmt.Sprintf("failed to marshal data: %v", err))
			}
			if err := json.Unmarshal(dataBytes, &data); err != nil {
				panic(fmt.Sprintf("failed to unmarshal data: %v", err))
			}
		} else {
			// []interface{}를 []TEIP712Data로 직접 변환
			data = make([]TEIP712Data, len(dataArray))
			for i, item := range dataArray {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					// map이 아니면 json.Marshal/Unmarshal로 변환
					itemBytes, err := json.Marshal(item)
					if err != nil {
						panic(fmt.Sprintf("failed to marshal data[%d]: %v", i, err))
					}
					if err := json.Unmarshal(itemBytes, &data[i]); err != nil {
						panic(fmt.Sprintf("failed to unmarshal data[%d]: %v", i, err))
					}
				} else {
					// map에서 직접 변환
					if typeVal, ok := itemMap["type"].(string); ok {
						data[i].Type = typeVal
					}
					if nameVal, ok := itemMap["name"].(string); ok {
						data[i].Name = nameVal
					}
					if valueVal, ok := itemMap["value"]; ok {
						data[i].Value = valueVal
					}
				}
			}
		}

		if len(data) == 0 {
			panic("data array is empty")
		}
	}

	// -----------------------------------------------------------
	// 3. genTypedData(domain, data, primaryType) 호출
	// -----------------------------------------------------------
	primaryType := param.PrimaryType
	if primaryType == "" {
		primaryType = "EIP712Data" // 기본값
	}
	typedData := genTypedData(domain, data, primaryType)

	// -----------------------------------------------------------
	// 4. interface{} → []interface{} 변환 (KeyPair)
	// -----------------------------------------------------------
	keyPairs, err := convertInterfaceToArray(param.KeyPair, "KeyPair")
	if err != nil {
		panic(err.Error())
	}

	// -----------------------------------------------------------
	// 5. 실제 서명 실행
	//    KeyPair = { address, passphrase }
	// -----------------------------------------------------------
	sigHashes, signatures := eip712DataSign(typedData, keyPairs)

	// response
	apiUtil.Response(_c, http.StatusOK, apiUtil.MakeResponseString(param.RequestId, map[string]interface{}{
		"sigHashes":  sigHashes,
		"signatures": signatures,
	}, false))

}
