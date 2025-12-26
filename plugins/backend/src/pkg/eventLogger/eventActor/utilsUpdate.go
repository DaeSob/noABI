package eventActor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	dbEvent "cia/api/db/event"
	"cia/api/logHandler"
	"cia/common/utils"
	jMapper "cia/common/utils/json/mapper"
	"cia/pkg/eventLogger/types"
)

func _insert(array []string, index int, element string) []string {
	result := append(array, element)
	copy(result[index+1:], result[index:])
	result[index] = element
	return result
}

// get first file
func getFilePaths(_dirPath string) (filePaths []string) {
	// file이 있는지 체크
	if isExistFile(_dirPath) {
		// get log file list
		var fileInfoList []utils.TFileInfo
		utils.SFolder(_dirPath, &fileInfoList)

		// process with first file
		for _, info := range fileInfoList {
			if strings.HasSuffix(info.Name, ".lock") {
				continue
			}
			if !strings.HasSuffix(info.Name, ".meta") && !strings.HasSuffix(info.Name, ".err") && !info.IsDir {
				sliceFileName := strings.Split(info.Name, ".")
				blockNumber, _ := utils.StringToI64(sliceFileName[0])
				if len(filePaths) == 0 {
					filePaths = append(filePaths, info.Name)
				} else {
					i := len(filePaths)
					for ; i > 0; i-- {
						preFileName := strings.Split(filePaths[i-1], ".")
						preBlockNumber, _ := utils.StringToI64(preFileName[0])
						if preBlockNumber < blockNumber {
							break
						}
					}
					filePaths = _insert(filePaths, i, info.Name)
				}
			}
		}
		return filePaths
	}
	return nil

}

func isExistFile(_filePath string) bool {
	// get log file list
	var fileInfoList []utils.TFileInfo
	utils.SFolder(_filePath, &fileInfoList)

	return len(fileInfoList) > 0
}

func errorHandler() {
	s := recover()
	if s != nil {
		err := fmt.Errorf("%v", s)
		logHandler.Write("trace", 0, "error", err.Error())
	}
}

func removeFile(_path string) {
	if _path != "" {
		logHandler.Write("trace", 0, "remove ", "[", _path, "]")
		utils.Remove(_path)
	}
}

func afterQuery(_commitedKey, _requestQuery string, _res *jMapper.TJsonMap) error {
	// convert to bytes from response
	b := bytes.NewBufferString(_res.PPrint())
	defer func() {
		b = nil
		_res = nil
	}()

	insertResult := &types.TDSLQueryResponse{}
	err := json.Unmarshal(b.Bytes(), insertResult)
	defer func() {
		insertResult = nil
	}()
	if err != nil {
		logHandler.Write("failed", 0, _commitedKey, "\n"+err.Error())
		return err
	}

	if insertResult.Result.Message != "" {
		logHandler.Write("failed", 0, _commitedKey, "\n"+insertResult.Result.Message)
		err = fmt.Errorf(insertResult.Result.Message)
		return err
	}

	/*
		if insertResult.Result.Message != "" && insertResult.Result.Message != "Error - undefined error : E11000 duplicate key error collection: log_events index: _id_" {
			logHandler.Write("failed", 0, _commitedKey, "\n"+insertResult.Result.Message)
			return false
		}
		if insertResult.Result.Message == "Error - undefined error : E11000 duplicate key error collection: log_events index: _id_" {
			logHandler.Write("failed", 0, _commitedKey, "\n"+insertResult.Result.Message)
			return false
		}
	*/

	return nil

}

func genAdditional(_taskEventLog *TTaskEventLog) map[string]interface{} {

	additional := map[string]interface{}{}

	switch _taskEventLog.Interface {
	case "P2PTrade":
		if _taskEventLog.Event.EventName == "eventRegisterTrade" {
			additional["status"] = string("0")
		}
	case "BridgeERC20", "BridgeERC721":
		// eventTransfer는 additional data에 claim 정보를 넣는다.(false)
		if _taskEventLog.Event.EventName == "eventTransfer" {
			// claim done
			additional["claimDone"] = false
			// erc type
			if _taskEventLog.Interface == "BridgeERC20" {
				additional["ercType"] = "ERC20"
			} else if _taskEventLog.Interface == "BridgeERC721" {
				additional["ercType"] = "ERC721"
			}
			// symbol
			additional["symbol"] = _taskEventLog.ExtraData.(string)
		}
	}
	return additional

}

func genLog(_taskEventLog *TTaskEventLog) (log dbEvent.TLog) {
	log.LogIndex = _taskEventLog.Event.LogIndex
	log.Address = _taskEventLog.Event.ContractAddress
	log.Data = _taskEventLog.Event.Data
	log.Topics = _taskEventLog.Event.Topics

	additional := genAdditional(_taskEventLog)

	//eventLogger.go Line 166 Issue 원복 code
	if _taskEventLog.Event.EventName == "TransferSingle" {
		num64, _ := strconv.ParseUint(_taskEventLog.Event.DecodeLog["value"].(string), 10, 64)
		_taskEventLog.Event.DecodeLog["value"] = num64
		num64, _ = strconv.ParseUint(_taskEventLog.Event.DecodeLog["id"].(string), 10, 64)
		_taskEventLog.Event.DecodeLog["id"] = num64
	}

	log.Event = dbEvent.TEvent{
		EventName:      _taskEventLog.Event.EventName,
		Parameters:     _taskEventLog.Event.DecodeLog,
		AdditionalData: additional,
	}
	return

}
