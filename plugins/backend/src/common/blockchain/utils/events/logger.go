package events

import (
	"cia/common/blockchain/abi"
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/rpcRequest/eth"
	"cia/common/blockchain/types"
	"cia/common/blockchain/utils/encode"
	"cia/common/blockchain/utils/events/abiMapper"
	rpc "cia/common/utils/rpcRequest"
	utilTypes "cia/common/utils/types"

	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// get log ex
// return : TRpcResponse -> []TLog
func GetLogEx(
	_rpc string,
	_fromBlock uint64,
	_toBlock uint64,
	_contractAddr []string,
	_topics []interface{},
	_id int64,
) []types.TLog {
	// get log
	eth := &eth.TEth{RPC: rpcRequest.TRPC{URL: _rpc}}
	res := eth.GetLogs(_fromBlock, _toBlock, _contractAddr, _topics, _id)
	// result to logs
	if res.Error != nil {
		rpcError := rpc.TRpcError{}
		mapstructure.Decode(res.Error, &rpcError)
		panicMessage := fmt.Sprintf("request GetLogs from %v to %v(%v)", _fromBlock, _toBlock, rpcError.Message)
		panic(panicMessage)
	} else {
		ref := utilTypes.GetValue(res.Result)
		var logs []types.TLog
		defer func() { logs = nil }()
		if ref.Kind() == reflect.Slice {
			for i := 0; i < ref.Len(); i++ {
				var log types.TLog
				log.MapToLog(ref.Index(i).Interface().(map[string]interface{}))
				logs = append(logs, log)
			}
		}
		return logs
	}
}

func GetDecodedLog(
	_log types.TLog,
	_abiMapper abiMapper.TABIMappers,
) map[string]interface{} {
	// get sig mapper from contract address
	address := _log.Address.HexString()
	sigMapper := _abiMapper[address]

	// get signature from log
	optsTopic := *abi.CreateDynamicType(
		"method", "function", true,
	)
	decode := encode.DecodeLog([]abi.TArgumentStr{optsTopic}, _log)
	keys := make([]string, 0, len(decode))
	for k := range decode {
		keys = append(keys, k)
	}

	// key가 없으면 중단
	if len(keys) != 1 {
		return nil
	}

	sig := fmt.Sprintf("%v", decode[keys[0]])

	// get abi from mapper
	event := sigMapper[sig]

	// release memory
	defer func() {
		sigMapper = nil
		decode = nil
		keys = nil
	}()

	// decode
	return encode.DecodeLog(event.GetTopicOption(), _log)
}
