package collection

import (
	"cia/api/db/event"
)

// collection
type TDocdbCollection string

const (
	COLLECTION_BRIDGE_ANCHOR TDocdbCollection = "bridge_anchors"
	COLLECTION_CHAIN         TDocdbCollection = "chain_metadata"
	COLLECTION_TOKEN         TDocdbCollection = "token_metadata"
	COLLECTION_BLOCK_TX      TDocdbCollection = "block_transactions"
	COLLECTION_LOG_EVENT     TDocdbCollection = "log_events"
	COLLECTION_NFT           TDocdbCollection = "nft_metadata"
)

func (c TDocdbCollection) String() string {
	return string(c)
}

// block_transaction schema
type TBlockTransactionSchema struct {
	BlockNumber      string       `json:"blockNumber"`
	ChainId          string       `json:"chainId"`
	From             string       `json:"from"`
	GasUsed          string       `json:"gasUsed"`
	Logs             []event.TLog `json:"logs"`
	Timestamp        string       `json:"timestamp"`
	TransactionIndex string       `json:"transactionIndex"`
	TransactionHash  string       `json:"transactionHash"`
	To               string       `json:"to"`
}

// log_event schema
type TLogEventSchema struct {
	BlockNumber      string                   `json:"blockNumber"`
	ChainId          string                   `json:"chainId"`
	LogIndex         string                   `json:"logIndex"`
	Address          string                   `json:"address"`
	Data             string                   `json:"data"`
	Topics           []string                 `json:"topics"`
	Event            TEvent                   `json:"event"`
	Transaction      *TBlockTransactionSchema `json:"transaction"`
	TransactionIndex string                   `json:"transactionIndex"`
}

type TEvent struct {
	EventName      string      `json:"eventName"`
	Parameters     interface{} `json:"parameters"`
	AdditionalData interface{} `json:"additionalData"`
}

type TLogEventQueryTransactionResult struct {
	From            string `json:"from"`
	GasUsed         string `json:"gasUsed"`
	Timestamp       string `json:"timestamp"`
	To              string `json:"to"`
	TransactionHash string `json:"transactionHash"`
}

type TLogEventQueryResult struct {
	Address     string                           `json:"address"`
	Data        string                           `json:"data"`
	Event       TEvent                           `json:"event"`
	Id          string                           `json:"_id"`
	Topics      []string                         `json:"topics"`
	Transaction *TLogEventQueryTransactionResult `json:"transaction"`
}
