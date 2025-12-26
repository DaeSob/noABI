package eventActor

import (
	dbEvent "cia/api/db/event"
	"cia/api/logHandler"
	"cia/common/utils"
	"fmt"
)

func onDefaultEventLog(_taskEventLog *TTaskEventLog) {

	tx := _taskEventLog.TemporaryTx[_taskEventLog.Event.TransactionHash].(dbEvent.TTransaction)
	// make tx
	newLog := genLog(_taskEventLog)

	committedKey := fmt.Sprintf("%v-%v-%v-%v", _taskEventLog.ChainId, _taskEventLog.Event.BlockNumber, _taskEventLog.Event.TransactionIndex, _taskEventLog.Event.LogIndex)

	logData := struct {
		ChainId         string      `json:"chainId"`
		TxHash          string      `json:"txHash"`
		ContractAddress string      `json:"contractAddress"`
		ContractName    string      `json:"contractName"`
		EventName       string      `json:"eventName"`
		EventInputs     interface{} `json:"eventInputs"`
	}{
		ChainId:         _taskEventLog.ChainId,
		TxHash:          tx.TransactionHash,
		ContractAddress: _taskEventLog.Event.ContractAddress,
		ContractName:    _taskEventLog.Event.ContractName,
		EventName:       newLog.Event.EventName,
		EventInputs:     newLog.Event.Parameters,
	}

	logHandler.Write("trace", 0, "\n", "commitedKey:", committedKey, "\n", "logData:", utils.InterfaceToJsonString(logData, true))

}
