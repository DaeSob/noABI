package types

type TDSLQueryResponse struct {
	Result struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message"`
	} `json:"result"`
}
