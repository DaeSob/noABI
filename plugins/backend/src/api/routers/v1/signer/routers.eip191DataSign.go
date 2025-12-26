package api_sgner

import (
	"cia/api/preference"
	apiUtil "cia/api/util"
	"cia/common/blockchain/keystore/signer"
	"cia/common/blockchain/signature"
	"cia/common/utils"
	jMapper "cia/common/utils/json/mapper"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

type TREQ_EIP191DataSign struct {
	RequestId     string      `json:"requestId"`
	ChainId       string      `json:"chainId"`
	SigDatas      []TSigData  `json:"sigDatas"`
	KeystorePairs interface{} `json:"signers,omitempty"`
}

type TSigData struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func EIP191DataSign(_c *gin.Context) {
	// body to TClaimCallDataParam
	bytes, _ := ioutil.ReadAll(_c.Request.Body)

	var param TREQ_EIP191DataSign
	jMapper.FromJson(bytes, &param)

	// error handling
	defer apiUtil.APIErrorHandler(_c, param.RequestId)

	var sigDatas []interface{}
	for _, v := range param.SigDatas {
		// 타입 끝이 []이면 배열 타입
		isArray := strings.HasSuffix(v.Type, "[]") || strings.LastIndex(v.Type, "[]") > 0

		if isArray {
			// 배열 값 처리
			values, err := convertValueToStringArray(v.Value, v.Name)
			if err != nil {
				panic(err.Error())
			}
			sigDatas = append(sigDatas, signature.ArrayData(v.Type, values))
		} else {
			// 단일 값 처리
			value, err := convertValueToString(v.Value, v.Name)
			if err != nil {
				panic(err.Error())
			}
			sigDatas = append(sigDatas, signature.Data(v.Type, value))
		}
	}

	// -----------------------------------------------------------
	// interface{} → []interface{} 변환 (KeystorePairs)
	// -----------------------------------------------------------
	// eip191은 jMapper를 사용하므로 별도 처리 필요
	var keystorePairs []interface{}
	{
		if param.KeystorePairs == nil {
			panic("KeystorePairs is required")
		}

		keyPairArray, ok := param.KeystorePairs.([]interface{})
		if !ok {
			// []interface{}가 아니면 jMapper로 변환 시도
			keyPairBytes, err := jMapper.ToJson(param.KeystorePairs)
			if err != nil {
				panic(fmt.Sprintf("failed to marshal KeystorePairs: %v", err))
			}
			if err := jMapper.FromJson(keyPairBytes, &keystorePairs); err != nil {
				panic(fmt.Sprintf("failed to unmarshal KeystorePairs: %v", err))
			}
		} else {
			keystorePairs = keyPairArray
		}

		if len(keystorePairs) == 0 {
			panic("KeystorePairs array is empty")
		}
	}

	//sig := rv2Utils.DataSignBy(param.ChainId, sigDatas, param.KeystorePairs)
	sigHashes, signatures := dataSignBy(param.ChainId, sigDatas, keystorePairs)

	// response
	apiUtil.Response(_c, http.StatusOK, apiUtil.MakeResponseString(param.RequestId, map[string]interface{}{
		"sigHashes":  sigHashes,
		"signatures": signatures,
	}, false))
}

func dataSign(
	_signer *signer.TSigner,
	_signerPhrase string,
	_sigData ...interface{},
) (sigHash string, signature string) {

	// SignMessage 내부에서 사용하는 해시 생성 로직과 동일하게 구현
	hashData := solsha3.SoliditySHA3(_sigData...)
	hashPrefix := solsha3.SoliditySHA3WithPrefix(
		solsha3.Bytes32("0x" + hex.EncodeToString(hashData)),
	)

	signedMsg, err := _signer.SignMessage(
		_signerPhrase,
		_sigData...,
	)
	if err != nil {
		panic(err)
	}

	sigHash = utils.BytesToHexString(hashPrefix)
	signature = utils.BytesToHexString(signedMsg)
	return

}

func dataSignBy(
	_chainId string,
	_sigDatas []interface{},
	_signers []interface{},
) (sigHashes []string, signatures []string) {

	for _, pair := range _signers {
		keystorePair := pair.(map[string]interface{})
		signer := signer.NewSigner(preference.GetKeystorePath(), utils.StringToBigInt(_chainId), common.HexToAddress(keystorePair["account"].(string)))
		sigHash, sig := dataSign(signer, keystorePair["phrase"].(string), _sigDatas...)
		sigHashes = append(sigHashes, sigHash)
		signatures = append(signatures, sig)
	}
	return

}
