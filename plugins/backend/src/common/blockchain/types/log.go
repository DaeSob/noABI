package types

import (
	util "cia/common/utils"
	"fmt"
	"strconv"
)

type TTopics []THash

func HexToTopics(_str []string) TTopics {
	var topics TTopics
	for _, topic := range _str {
		topics = append(topics, HexToHash(topic))
	}

	return topics
}

type TLog struct {
	Removed          bool     `json:"removed"`
	LogIndex         uint64   `json:"logIndex"`
	TransactionIndex uint64   `json:"transactionIndex"`
	TransactionHash  THash    `json:"transactionHash"`
	BlockHash        THash    `json:"blockHash"`
	BlockNumber      uint64   `json:"blockNumber"`
	Address          TAddress `json:"address"`
	Topics           TTopics  `json:"topics"`
	Data             []byte   `json:"data"`
}

func CreateLog(
	_removed,
	_logIndex,
	_transactionIndex,
	_transactionHash,
	_blockHash,
	_blockNumber,
	_address string,
	_topics []string,
	_data string,
) TLog {
	var log TLog
	log.Removed = util.StringToBool(_removed)
	log.LogIndex = util.HexToUint64(_logIndex)
	log.TransactionIndex = util.HexToUint64(_transactionIndex)
	log.TransactionHash = HexToHash(_transactionHash)
	log.BlockHash = HexToHash(_blockHash)
	log.BlockNumber = util.HexToUint64(_blockNumber)
	log.Address = HexToAddress(_address)
	log.Topics = HexToTopics(_topics)
	log.Data = util.HexToBytes(_data)

	return log
}

// topics to string array
func (log *TLog) TopicsToStringArray() (result []string) {
	for _, topic := range log.Topics {
		result = append(result, topic.String())
	}

	return
}

// get event signature -> first element of topics
func (log *TLog) GetEventSignature() THash {
	return log.Topics[0]
}

func (log *TLog) MapToLog(_map map[string]interface{}) {
	// removed
	log.Removed = _map["removed"].(bool)

	// blockHash
	log.BlockHash = HexToHash(_map["blockHash"].(string))

	// txHash
	log.TransactionHash = HexToHash(_map["transactionHash"].(string))

	// topics
	for _, topic := range _map["topics"].([]interface{}) {
		log.Topics = append(log.Topics, HexToHash(topic.(string)))
	}

	// address
	log.Address = HexToAddress(_map["address"].(string))

	// log index
	log.LogIndex = util.StringToUint64(_map["logIndex"].(string)[2:], 16)

	// tx index
	log.TransactionIndex = util.StringToUint64(_map["transactionIndex"].(string)[2:], 16)

	// block number
	log.BlockNumber = util.StringToUint64(_map["blockNumber"].(string)[2:], 16)

	// data
	log.Data = util.HexToBytes(_map["data"].(string))
}

func (log *TLog) ToJsonString(_pretty bool) string {
	type strLog struct {
		Removed          string   `json:"removed"`
		LogIndex         string   `json:"logIndex"`
		TransactionIndex string   `json:"transactionIndex"`
		TransactionHash  string   `json:"transactionHash"`
		BlockHash        string   `json:"blockHash"`
		BlockNumber      string   `json:"blockNumber"`
		Address          string   `json:"address"`
		Topics           []string `json:"topics"`
		Data             string   `json:"data"`
	}

	tlog := strLog{
		fmt.Sprintf("%v", log.Removed),
		"0x" + strconv.FormatUint(log.LogIndex, 16),
		"0x" + strconv.FormatUint(log.TransactionIndex, 16),
		log.TransactionHash.String(),
		log.BlockHash.String(),
		"0x" + strconv.FormatUint(log.BlockNumber, 16),
		log.Address.String(),
		[]string{},
		util.BytesToHexString(log.Data),
	}

	for _, topic := range log.Topics {
		tlog.Topics = append(tlog.Topics, topic.String())
	}

	return util.InterfaceToJsonString(tlog, _pretty)
}
