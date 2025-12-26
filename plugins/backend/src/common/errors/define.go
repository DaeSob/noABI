package errors

////////////////////////////////////////////////////////////////////////////////
// common
// define code & error (code have to pair with error)
const (
	// test error
	TEST_CODE  = "-1"
	TEST_ERROR = "Test Message"

	// common error
)

// define custom type error
func ERROR_TEST() TError {
	return TError{
		TEST_CODE,
		TEST_ERROR,
	}
}

////////////////////////////////////////////////////////////////////////////////
// common
// define code (message is passed as error message of the function)
const (
	// rpc request (http request)
	RPC_REQUEST_CODE = "10001"
)

// define custom type error
// rpc request (http request)
func ERROR_RPC_REQUEST(_error string) TError {
	return TError{RPC_REQUEST_CODE, _error}
}

////////////////////////////////////////////////////////////////////////////////
// bridge api (11001 ~ 12000)
// define code & message
const (
	// unsupported chain id
	UNSUPPORTED_CHAIN_ID_CODE  = "11001"
	UNSUPPORTED_CHAIN_ID_ERROR = "unsupported chain id"

	// invalid claim id
	INVALID_CLAIM_ID_CODE  = "11002"
	INVALID_CLAIM_ID_ERROR = "invalid claim id"
)

// define custom type error
// unsupported chain id
func ERROR_UNSUPPORTED_CHAIN_ID() TError {
	return TError{UNSUPPORTED_CHAIN_ID_CODE, UNSUPPORTED_CHAIN_ID_ERROR}
}

// invalid claim id
func ERROR_INVALID_CLAIM_ID() TError {
	return TError{INVALID_CLAIM_ID_CODE, INVALID_CLAIM_ID_ERROR}
}

////////////////////////////////////////////////////////////////////////////////
// blockchain
// define code (message is passed as error message of the chain)
const (
	// rpc request : eth
	CALL_CODE                    = "20001" // eth_call
	ESTIMATE_GAS_CODE            = "20002" // eth_estimateGas
	GET_BALANCE_CODE             = "20003" // eth_getBalance
	GET_BLOCK_BY_NUMBER_CODE     = "20004" // eth_getBlockByNumber
	BLOCK_NUMBER_CODE            = "20005" // eth_blockNumber
	CHAIN_ID_CODE                = "20006" // eth_chainId
	GAS_PRICE_CODE               = "20007" // eth_gasPrice
	GET_LOGS_CODE                = "20008" // eth_getLogs
	GET_TRANSACTION_COUNT_CODE   = "20009" // eth_getTransactionCount
	GET_TRANSACTION_RECEIPT_CODE = "20010" // eth_getTransactionReceipt
	SEND_RAW_TRANSACTION_CODE    = "20011" // eth_sendRawTransaction

	// rpc request : tx pool
	INSPECT_CODE = "20012" // txpool_inspect
	STATUS_CODE  = "20013" // txpool_status
)

// define custom type error
// rpc request : eth
func ERROR_CALL(_error string) TError {
	return TError{CALL_CODE, _error}
}
func ERROR_ESTIMATE_GAS(_error string) TError {
	return TError{ESTIMATE_GAS_CODE, _error}
}
func ERROR_GET_BALANCE(_error string) TError {
	return TError{GET_BALANCE_CODE, _error}
}
func ERROR_GET_BLOCK_BY_NUMBER(_error string) TError {
	return TError{GET_BLOCK_BY_NUMBER_CODE, _error}
}
func ERROR_BLOCK_NUMBER(_error string) TError {
	return TError{BLOCK_NUMBER_CODE, _error}
}
func ERROR_CHAIN_ID(_error string) TError {
	return TError{CHAIN_ID_CODE, _error}
}
func ERROR_GAS_PRICE(_error string) TError {
	return TError{GAS_PRICE_CODE, _error}
}
func ERROR_GET_LOGS(_error string) TError {
	return TError{GET_LOGS_CODE, _error}
}
func ERROR_GET_TRANSACTION_COUNT(_error string) TError {
	return TError{GET_TRANSACTION_COUNT_CODE, _error}
}
func ERROR_GET_TRANSACTION_RECEIPT(_error string) TError {
	return TError{GET_TRANSACTION_RECEIPT_CODE, _error}
}
func ERROR_SEND_RAW_TRANSACTION(_error string) TError {
	return TError{SEND_RAW_TRANSACTION_CODE, _error}
}

// rpc request : tx pool
func ERROR_INSPECT(_error string) TError {
	return TError{INSPECT_CODE, _error}
}
func ERROR_STATUS(_error string) TError {
	return TError{STATUS_CODE, _error}
}

////////////////////////////////////////////////////////////////////////////////
// hexlant
// define code (message is passed as error message of the hexlant api)
const (
	CREATE_WALLET_CODE          = "30001" // createWallet
	REQUEST_DATA_SIGN_CODE      = "30002" // requestDataSign
	REQUEST_DATA_SIGN_INFO_CODE = "30003" // requestDataSignInfo
	REQUEST_TX_SIGN_CODE        = "30004" // requestTxSign
	REQUEST_TX_SIGN_INFO_CODE   = "30005" // requestTxSignInfo
)

func ERROR_CREATE_WALLET(_error string) TError {
	return TError{CREATE_WALLET_CODE, _error}
}
func ERROR_REQUEST_DATA_SIGN(_error string) TError {
	return TError{REQUEST_DATA_SIGN_CODE, _error}
}
func ERROR_REQUEST_DATA_SIGN_INFO(_error string) TError {
	return TError{REQUEST_DATA_SIGN_INFO_CODE, _error}
}
func ERROR_REQUEST_TX_SIGN(_error string) TError {
	return TError{REQUEST_TX_SIGN_CODE, _error}
}
func ERROR_REQUEST_TX_SIGN_INFO(_error string) TError {
	return TError{REQUEST_TX_SIGN_INFO_CODE, _error}
}
