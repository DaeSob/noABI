package eventActor

import (
	dbEvent "cia/api/db/event"
	"cia/api/logHandler"
	"cia/api/preference"
	"cia/common/blockchain/eventLog"
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/rpcRequest/eth"
	"cia/common/utils"
	jMapper "cia/common/utils/json/mapper"
	rpc "cia/common/utils/rpcRequest"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
)

type TTaskEventLog struct {
	CTX         context.Context
	ChainId     string
	Eth         eth.TEth
	Url         string
	Retries     int
	BlockNumber string
	BlockHash   string
	Timestamp   string
	Interface   string
	ExtraData   interface{}
	Event       *eventLog.TEventLog
	TemporaryTx map[string]interface{}
}

func readEventLogFile(_chain, _logPath string, _eventLoggerInfo *preference.TEvetLoggerInfo) {
	// error handler
	defer errorHandler()

	// get log file path
	filePaths := getFilePaths(_logPath)

	defer func() {
		filePaths = nil
		s := recover()
		if s != nil {
			err := fmt.Errorf("%v", s)
			logHandler.Write("trace", 0, "panic", err.Error())
		}
	}()

	for _, fileName := range filePaths {
		// remove file when exit function
		filePath := fmt.Sprintf("%v/%v", _logPath, fileName)
		if filePath != "" {
			// read file
			bytesData, err := ioutil.ReadFile(filePath)
			defer func() {
				bytesData = nil
			}()
			if err != nil {
				errFilePath := fmt.Sprintf("%v/failed/%v", _logPath, fileName)
				utils.Rename(filePath, errFilePath)
				logHandler.Write("failed", 0, "read file", "\nfile backup path: ", errFilePath)
				continue
			}

			// get block log from file data
			blockLogs := eventLog.TBlockLog{}
			err = jMapper.FromJson(bytesData, &blockLogs)
			if err != nil {
				errFilePath := fmt.Sprintf("%v/failed/%v", _logPath, fileName)
				utils.Rename(filePath, errFilePath)
				logHandler.Write("failed", 0, "convert to json", "\nfile backup path: ", errFilePath)
				continue
			}

			logSwitch(_chain, _eventLoggerInfo, &blockLogs)
			removeFile(filePath)

		}
	}
}

func logSwitch(_chain string, _eventLoggerInfo *preference.TEvetLoggerInfo, _logs *eventLog.TBlockLog) {

	// get block info
	eth := eth.TEth{RPC: rpcRequest.TRPC{URL: _eventLoggerInfo.RPC}}
	block, _ := utils.StringToI64(_logs.BlockNumber)
	hexBlock := fmt.Sprintf("0x%x", block)

	res := eth.GetBlockByNumber(hexBlock, false, int64(1))
	if res.Error != nil {
		rpcError := rpc.TRpcError{}
		mapstructure.Decode(res.Error, &rpcError)
		panicMessage := fmt.Sprintf("request GetBlockByNumber blockNumber %v(%v)", hexBlock, rpcError.Message)
		panic(panicMessage)
	}
	if res.Result == nil {
		panicMessage := fmt.Sprintf("can't look up %v block", hexBlock)
		panic(panicMessage)
	}
	blockInfo := res.Result.(map[string]interface{})

	taskEventLog := TTaskEventLog{
		CTX:         context.Background(),
		ChainId:     _eventLoggerInfo.ChainId,
		Eth:         eth,
		Url:         _eventLoggerInfo.DB,
		Retries:     _eventLoggerInfo.Retries,
		BlockNumber: _logs.BlockNumber,
		BlockHash:   _logs.BlockHash,
		Timestamp:   blockInfo["timestamp"].(string),
		TemporaryTx: map[string]interface{}{},
	}

	for _, event := range _logs.EventLogs {

		_, ok := taskEventLog.TemporaryTx[event.TransactionHash]
		if !ok {
			res := eth.GetTransactionReceipt(event.TransactionHash, 1)
			if res.Error != nil {
				rpcError := rpc.TRpcError{}
				mapstructure.Decode(res.Error, &rpcError)
				panicMessage := fmt.Sprintf("request GetTransactionReceipt txHash %v(%v)", event.TransactionHash, rpcError.Message)
				panic(panicMessage)
			}
			txReceipt := res.Result.(map[string]interface{})

			var gasUsed, from, to string
			if txReceipt["gasUsed"] != nil {
				gasUsed = txReceipt["gasUsed"].(string)
			}
			if txReceipt["from"] != nil {
				from = txReceipt["from"].(string)
			}
			if txReceipt["to"] != nil {
				to = txReceipt["to"].(string)
			}

			tx := dbEvent.TTransaction{
				TransactionIndex: event.TransactionIndex,
				TransactionHash:  event.TransactionHash,
				GasUsed:          gasUsed,
				From:             from,
				To:               to,
			}
			taskEventLog.TemporaryTx[event.TransactionHash] = tx
		}

		taskEventLog.Interface = _eventLoggerInfo.Collections[event.ContractName].Interface
		taskEventLog.ExtraData = _eventLoggerInfo.Collections[event.ContractName].ExtraData

		taskEventLog.Event = &event
		switch taskEventLog.Interface {
		default:
			onDefaultEventLog(&taskEventLog)
		}
	}
}
