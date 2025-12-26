package event

type TEventLogs struct {
	ChainId      string        `json:"chainId"`
	BlockNumber  string        `json:"blockNumber"`
	Timestamp    string        `json:"timestamp"`
	Transactions []interface{} `json:"transactions"`
}

type TTransaction struct {
	TransactionIndex string `json:"transactionIndex"`
	TransactionHash  string `json:"transactionHash"`
	GasUsed          string `json:"gasUsed"`
	From             string `json:"from"`
	To               string `json:"to"`
	Logs             []TLog `json:"logs"`
}

type TLog struct {
	LogIndex string   `json:"logIndex"`
	Address  string   `json:"address"`
	Data     string   `json:"data"`
	Topics   []string `json:"topics"`
	Event    TEvent   `json:"event"`
}

type TEvent struct {
	EventName      string      `json:"eventName"`
	Parameters     interface{} `json:"parameters"`
	AdditionalData interface{} `json:"additionalData,omitempty"`
}
