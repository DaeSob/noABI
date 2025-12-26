package eventLog

import (
	"cia/api/logHandler"
	"cia/common/blockchain/types"
	"cia/common/blockchain/utils/events"
	"cia/common/blockchain/utils/events/abiMapper"
	"cia/common/utils"
	"fmt"
	"math/big"
	"os"
)

type TBlockLog struct {
	BlockNumber string     `json:"blockNumber"`
	BlockHash   string     `json:"blockHash"`
	EventLogs   TEventLogs `json:"eventLogs"`
}

func GenerateBlockLog(
	_blockNumber,
	_blockHash string,
	_eventLogs TEventLogs) *TBlockLog {
	return &TBlockLog{
		_blockNumber,
		_blockHash,
		_eventLogs,
	}
}

func (b *TBlockLog) ToJsonString(_pretty bool) string {
	return utils.InterfaceToJsonString(b, _pretty)
}

func (e *TBlockLog) SaveToFile(_logPath string) {
	// file name : ${blocknumber}.json
	lockFilename := e.BlockNumber + ".json.lock"
	fileName := e.BlockNumber + ".json"

	// file path : _logPath + filename
	lockFilePath := fmt.Sprintf("%v/%v", _logPath, lockFilename)
	filePath := fmt.Sprintf("%v/%v", _logPath, fileName)

	// to string
	logText := e.ToJsonString(true)

	// save to file
	_writeFile(lockFilePath, logText)

	//unlock
	os.Rename(lockFilePath, filePath)

	// trace log
	logHandler.Write("trace", 0, "save to log: ", e.ToJsonString(false))
}

type TEventLog struct {
	BlockNumber      string                 `json:"blockNumber"`
	BlockHash        string                 `json:"blockHash"`
	ContractAddress  string                 `json:"contractAddress"`
	ContractName     string                 `json:"contractName"`
	DecodeLog        map[string]interface{} `json:"decodeLog"`
	EventName        string                 `json:"eventName"`
	TransactionHash  string                 `json:"transactionHash"`
	TransactionIndex string                 `json:"transactionIndex"`
	LogIndex         string                 `json:"logIndex"`
	Topics           []string               `json:"topics"`
	Data             string                 `json:"data"`
}

type TEventLogs []TEventLog

func (e *TEventLogs) GetBlockLog() []TBlockLog {
	// key : block number
	blockLog := make(map[string]TBlockLog)

	// separate event log with block number
	for _, log := range *e {
		block := blockLog[log.BlockNumber]
		if block.BlockNumber == "" {
			tmp := *GenerateBlockLog(
				log.BlockNumber,
				log.BlockHash,
				TEventLogs{log},
			)
			blockLog[log.BlockNumber] = tmp
		} else {
			block.EventLogs = append(block.EventLogs, log)
			blockLog[log.BlockNumber] = block
		}
	}

	// make result
	result := []TBlockLog{}
	for key := range blockLog {
		result = append(result, blockLog[key])
	}

	// release memory
	defer func() {
		blockLog = nil
		result = nil
	}()

	return result
}

func GenerateEventLog(
	_blockNumber,
	_blockHash,
	_addr,
	_name,
	_event string,
	_decodeLog map[string]interface{},
	_hash string,
	_txIdx string,
	_logIdx string,
	_topics []string,
	_data string,
) *TEventLog {
	return &TEventLog{
		_blockNumber,
		_blockHash,
		_addr,
		_name,
		_decodeLog,
		_event,
		_hash,
		_txIdx,
		_logIdx,
		_topics,
		_data,
	}
}

func (e *TEventLog) ToJsonString(_pretty bool) string {
	return utils.InterfaceToJsonString(e, _pretty)
}

func GetEventLogs(
	_rpc string,
	_fromBlock uint64,
	_toBlock uint64,
	_contractAddr []string,
	_contractName map[string]string,
	_id int64,
	_abiMapper abiMapper.TABIMappers,
	_inputLogging bool,
) TEventLogs {
	// get logs
	logs := events.GetLogEx(
		_rpc,
		_fromBlock,
		_toBlock,
		_contractAddr,
		[]interface{}{},
		_id,
	)

	type TInputLog struct {
		Status        string      `json:"status,omitempty"`
		InputLogCount int         `json:"inputLogCount,omitempty"`
		FromBlock     uint64      `json:"fromBlockNumber,omitempty"`
		ToBlock       uint64      `json:"toBlockNumber,omitempty"`
		RawData       interface{} `json:"rawData,omitempty"`
		EventName     string      `json:"eventName,omitempty"`
	}

	if _inputLogging {
		logHandler.Write("logCollector", 0, utils.InterfaceToJsonString(TInputLog{
			Status:        "input log",
			InputLogCount: len(logs),
			FromBlock:     _fromBlock,
			ToBlock:       _toBlock,
		}, false))
	}

	var eventLogs []TEventLog
	for _, log := range logs {
		// get decode log
		decodeLog := events.GetDecodedLog(log, _abiMapper)

		// big int to string
		// -> type 변환을 거치면서 float 오차 발생으로 문제가 발생하는 것을 방지
		for key, value := range decodeLog {
			bigIntValue, ok := value.(*big.Int)
			if ok {
				decodeLog[key] = bigIntValue.Text(10)
			}
		}

		eventName := _getEventName(log, _abiMapper)

		if len(decodeLog) == 0 {
			if _inputLogging {
				logHandler.Write("logCollector", 0, utils.InterfaceToJsonString(TInputLog{
					Status:    "failed",
					EventName: eventName,
					RawData:   decodeLog,
				}, false))
			}
		} else if eventName == "" {
			if _inputLogging {
				logHandler.Write("logCollector", 0, utils.InterfaceToJsonString(TInputLog{
					Status:  "not registered event",
					RawData: decodeLog,
				}, false))
			}
		} else {
			// generate event log
			eventLog := GenerateEventLog(
				fmt.Sprint(log.BlockNumber),
				log.BlockHash.HexString(),
				log.Address.String(),
				_contractName[log.Address.String()],
				eventName,
				decodeLog,
				log.TransactionHash.HexString(),
				utils.Uint64ToString(log.TransactionIndex),
				utils.Uint64ToString(log.LogIndex),
				log.TopicsToStringArray(),
				utils.BytesToHexString(log.Data),
			)
			// add to result
			eventLogs = append(eventLogs, *eventLog)
		}
	}

	// release memory
	defer func() {
		logs = nil
		eventLogs = nil
	}()

	return eventLogs
}

func _getEventName(_log types.TLog, _abiMapper abiMapper.TABIMappers) string {
	// get sig mapper
	contractAddr := _log.Address.HexString()
	sigMapper := _abiMapper[contractAddr]

	// get event from sig mapper
	hexSig := _log.GetEventSignature().HexString()
	event := sigMapper[hexSig]

	// release memory
	defer func() {
		sigMapper = nil
	}()

	// get event name
	return event.Name
}

func GetContractEventLogs(
	_rpc string,
	_fromBlock uint64,
	_toBlock uint64,
	_contractAddr string,
	_contractName string,
	_id int64,
	_abiMapper abiMapper.TABIMappers,
	_inputLogging bool,
) TEventLogs {
	return GetEventLogs(
		_rpc,
		_fromBlock,
		_toBlock,
		[]string{_contractAddr},
		map[string]string{_contractAddr: _contractName},
		_id,
		_abiMapper,
		_inputLogging,
	)
}

func GetContractEventLogsEx(
	_rpc string,
	_fromBlock uint64,
	_toBlock uint64,
	_collections map[string]string,
	_id int64,
	_abiMapper abiMapper.TABIMappers,
	_inputLogging bool,
) TEventLogs {

	var addresses []string
	for address, _ := range _collections {
		addresses = append(addresses, address)
	}
	return GetEventLogs(
		_rpc,
		_fromBlock,
		_toBlock,
		addresses,
		_collections,
		_id,
		_abiMapper,
		_inputLogging,
	)
}

// write file for save log file
//var LockMutex sync.Mutex

func _writeFile(_filePath string, _strText string) {
	// rewrite
	//LockMutex.Lock()
	//defer LockMutex.Unlock()
	f, e := os.OpenFile(_filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	f.WriteString(_strText)
}
