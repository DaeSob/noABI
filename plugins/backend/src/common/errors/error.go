package errors

import "cia/common/utils"

type TError struct {
	Code        string `json:"code,omitempty"`
	CustomError string `json:"error,omitempty"`
}

type TErrorEx struct {
	Code        string      `json:"code,omitempty"`
	CustomError string      `json:"error,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

func (e TError) String() string {
	return utils.InterfaceToJsonString(e, true)
}

func (e TError) Error() string {
	return e.CustomError
}
