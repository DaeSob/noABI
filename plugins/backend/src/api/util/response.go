package util

import (
	"cia/api/logHandler"
	"cia/common/errors"
	"cia/common/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// //////////////////////////////////////////////////////////////////////////////
// api types
type TAPIResult struct {
	Data interface{} `json:"data"` // json format supports struct
}

type TAPIResultEx struct {
	Data    interface{} `json:"data,omitempty"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"error,omitempty"`
}

type TAPIResponse struct {
	Id     string     `json:"id"`
	Result TAPIResult `json:"result"`
}

type TAPIResponseEx struct {
	Id     string       `json:"id"`
	Result TAPIResultEx `json:"result"`
}

type TAPIErrorResponse struct {
	Id     string        `json:"id"`
	Result errors.TError `json:"result"`
}

type TAPIErrorExResponse struct {
	Id     string          `json:"id"`
	Result errors.TErrorEx `json:"result"`
}

func MakeResponseString(
	_id string,
	_data interface{},
	_pretty bool,
) string {
	res := TAPIResponse{
		_id,
		TAPIResult{_data},
	}

	return utils.InterfaceToJsonString(res, _pretty)
}

func MakeErrorResponseString(
	_id string,
	_error errors.TError,
	_pretty bool,
) string {
	res := TAPIErrorResponse{
		_id,
		_error,
	}

	return utils.InterfaceToJsonString(res, _pretty)
}

func MakeExErrorResponseString(
	_id string,
	_error errors.TError,
	_data interface{},
	_pretty bool,
) string {
	res := TAPIErrorExResponse{
		_id,
		errors.TErrorEx{
			Code:        _error.Code,
			CustomError: _error.CustomError,
			Data:        _data,
		},
	}
	return utils.InterfaceToJsonString(res, _pretty)
}

func Response(_c *gin.Context, _httpCode int, _result string) {
	// header
	SetDefaultHeader(_c)

	logHandler.Write("access", 0, _result)
	// response
	_c.String(_httpCode, _result)
}

func APIErrorHandler(_c *gin.Context, _id string) {
	s := recover()
	if s != nil {
		recovered, ok := s.(errors.TError)
		if ok {
			logHandler.Write("trace", 0, fmt.Sprintf("id:%s status:panic message:%s", _id, recovered.Error()))
			Response(
				_c,
				http.StatusOK,
				MakeErrorResponseString(_id, recovered, false),
			)
		} else {
			err := fmt.Errorf("%v", s)
			logHandler.Write("trace", 0, fmt.Sprintf("id:%s status:panic message:%s", _id, err.Error()))
			Response(
				_c,
				http.StatusOK,
				MakeErrorResponseString(_id, errors.TError{"", err.Error()}, false),
			)
		}
	}
}
