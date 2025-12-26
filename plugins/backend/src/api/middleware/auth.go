package middleware

import (
	"cia/api/preference"
	apiUtil "cia/api/util"
	"cia/common/blockchain/keystore/signer"
	"cia/common/errors"
	"cia/common/utils"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthCheckDate(_c *gin.Context) {
	id := apiUtil.GetIdFromHeader(_c)

	// get date
	date := _c.GetHeader("Date")
	reqDate := utils.GMTStringToLocalTime(date).Unix()
	nowDate := utils.Now().Unix()

	// calc diff
	diffTime := nowDate - reqDate
	diffTime = int64(math.Abs(float64(diffTime)))

	// check date
	// 60초 이내의 request만 허용
	if diffTime > 60 || diffTime < 0 {
		apiUtil.Response(
			_c,
			http.StatusUnauthorized,
			apiUtil.MakeErrorResponseString(
				id,
				errors.TError{"401", "unauthorized"},
				false,
			),
		)

		_c.Abort()
		return
	}

	_c.Next()
}

func AuthCheckSignature(_c *gin.Context) {
	// set id
	id := apiUtil.GetIdFromHeader(_c)

	// get Authorization
	auth := _c.GetHeader("Authorization")
	arrAuth := strings.Split(auth, " ")
	// release memory
	defer func() {
		arrAuth = nil
	}()

	// check Authorization
	if len(arrAuth) != 2 || arrAuth[0] != "NOABI" {
		apiUtil.Response(
			_c,
			http.StatusUnauthorized,
			apiUtil.MakeErrorResponseString(
				id,
				errors.TError{"401", "unauthorized"},
				false,
			),
		)
		_c.Abort()
		return
	}

	// check signature
	sig := arrAuth[1]

	// Date & URI -> auth message
	date := _c.GetHeader("Date")
	uri := _c.Request.URL.Path
	authMsg := makeAuthMessage(uri, date)

	// recover signer
	recoverAddr := signer.RecoverDataSignerFromString(authMsg, sig)

	// check auth
	isAuth := preference.IsAuthSigner(recoverAddr)

	if isAuth {
		_c.Next()
	} else {
		apiUtil.Response(
			_c,
			http.StatusUnauthorized,
			apiUtil.MakeErrorResponseString(
				id,
				errors.TError{"401", "unauthorized"},
				false,
			),
		)
		_c.Abort()
	}
}

func makeAuthMessage(_uri string, _date string) string {
	return _uri + "\\r\\n" + _date + "\\r\\n"
}

func NoAuth(_c *gin.Context) {
	_c.Next()
}
