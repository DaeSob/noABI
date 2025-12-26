package query

import (
	"cia/api/db/query"
	"cia/api/preference"
	"cia/common/utils"
	"cia/common/utils/http/httpRequest"
	jMapper "cia/common/utils/json/mapper"
)

type TSplitOperator struct {
	Split []string `json:"$split"`
}
type TArrayElemAt struct {
	ArrayElemAt []interface{} `json:"$arrayElemAt"`
}
type TConcat struct {
	Concat []interface{} `json:"$concat"`
}
type TBlockTx struct {
	BlockTx interface{} `json:"block_tx"`
}

// add fields
func getAddFields() query.TAddFieldsOperator {
	return query.TAddFieldsOperator{
		TBlockTx{
			TConcat{
				[]interface{}{
					TArrayElemAt{[]interface{}{TSplitOperator{[]string{`$_id`, "#"}}, 0}},
					"#",
					TArrayElemAt{[]interface{}{TSplitOperator{[]string{`$_id`, "#"}}, 1}},
					"#",
					TArrayElemAt{[]interface{}{TSplitOperator{[]string{`$_id`, "#"}}, 2}},
				},
			},
		},
	}
}

// lookup
func getLookup() query.TLookupOperator {
	return query.TLookupOperator{
		struct {
			LocalField   string `json:"localField"`
			ForeignField string `json:"foreignField"`
			From         string `json:"from"`
			As           string `json:"as"`
		}{
			"block_tx",
			"_id",
			"block_transactions",
			"transaction",
		},
	}
}

// unwind
func getUnwind() query.TUnwindOperator {
	return query.TUnwindOperator{
		"$transaction",
	}
}

// sort
// recently -> 내림차순
// 1 (for ascending) or -1 (for descending)
func getSort(_opt int) query.TSortOperator {
	return query.TSortOperator{
		struct {
			Timestamp int `json:"transaction.timestamp"`
		}{
			_opt,
		},
	}
}

func getProject() query.TProjectOperator {
	return query.TProjectOperator{
		struct {
			Id        int `json:"_id"`
			Address   int `json:"address"`
			Data      int `json:"data"`
			Topics    int `json:"topics"`
			Event     int `json:"event"`
			Timestamp int `json:"transaction.timestamp"`
			TxHash    int `json:"transaction.transactionHash"`
			GasUsed   int `json:"transaction.gasUsed"`
			From      int `json:"transaction.from"`
			To        int `json:"transaction.to"`
		}{
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
		},
	}
}

func makeEventLogQuery(_id, _eventName string) string {
	// match id & address & event name
	match := query.TMatchOperator{
		struct {
			Id        query.TRegexOperator `json:"_id"`
			EventName string               `json:"event.eventName"`
		}{
			Id:        query.TRegexOperator{"^" + _id},
			EventName: _eventName,
		},
	}

	// add fields
	addFields := getAddFields()

	// lookup
	lookup := getLookup()

	// unwind
	unwind := getUnwind()

	// sort
	// recently -> 내림차순
	// 1 (for ascending) or -1 (for descending)
	sort := getSort(-1)

	// project
	project := getProject()

	return utils.InterfaceToJsonString(
		query.MakeLogEventAggregateQuery(
			[]interface{}{
				match,
				addFields,
				lookup,
				unwind,
				sort,
				project,
			},
		),
		false,
	)
}

type TFindQueryData struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type TFindQueryResponse struct {
	Result TFindQueryData `json:"result"`
}

func findQueryOnDB(_requestData string) TFindQueryResponse {
	// request find query
	apiUrl := preference.GetApiServerUrl("dslQuery")
	bytes, err := httpRequest.PostRequestToBytes(
		apiUrl,
		_requestData,
		5*1000,
	)
	if err != nil {
		panic(err.Error())
	}

	// response
	var res TFindQueryResponse
	jMapper.FromJson(bytes, &res)

	return res
}

func QueryEventLog(_id, _eventName string) TFindQueryResponse {
	query := makeEventLogQuery(_id, _eventName)

	return findQueryOnDB(query)
}
