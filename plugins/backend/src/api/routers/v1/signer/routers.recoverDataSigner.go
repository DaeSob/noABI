package api_sgner

import (
	apiUtil "cia/api/util"
	"fmt"
	"io/ioutil"

	"cia/common/blockchain/keystore/signer"
	"cia/common/utils"
	jMapper "cia/common/utils/json/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TREQ_RecoverDataParam struct {
	RequestId string `json:"requestId"`
	SigHash   string `json:"sigHash"`
	Signature string `json:"signature"`
}

func RecoverDataSigner(_c *gin.Context) {
	// body to TUserOperationParam
	bytes, _ := ioutil.ReadAll(_c.Request.Body)

	var param TREQ_RecoverDataParam
	jMapper.FromJson(bytes, &param)

	// error handling
	defer apiUtil.APIErrorHandler(_c, param.RequestId)

	// gen user operation
	address, err := signer.RecoverDataSigner(utils.HexToBytes(param.SigHash), utils.HexToBytes(param.Signature))

	if err != nil {
		panic(fmt.Errorf("%v", err.Error()))
	}

	res := fmt.Sprintf("%v", address)

	// response
	apiUtil.Response(_c, http.StatusOK, apiUtil.MakeResponseString(param.RequestId, res, false))
}
