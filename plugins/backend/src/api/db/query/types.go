package query

import "go.mongodb.org/mongo-driver/bson"

// query request
type TQueryOPRequest struct {
	OP        string      `json:"op"`
	Query     TDocdbQuery `json:"query"`
	Convert   bool        `json:"convert,omitempty"`
	Documents interface{} `json:"documents,omitempty"`
}

type TDocdbQuery struct {
	Collection               string            `json:"collection"`
	DocumentId               string            `json:"documentId,omitempty"`
	CaseInsensitivityId      *bool             `json:"caseInsensitivityId,omitempty"`
	CaseInsensitivityAddress *bool             `json:"caseInsensitivityAddress,omitempty"`
	Filter                   interface{}       `json:"filter,omitempty"`
	UpdateSet                interface{}       `json:"updateSet,omitempty"`
	Options                  TDocdbQueryOption `json:"options,omitempty"`
	Pipeline                 []interface{}     `json:"pipeline,omitempty"`
	IgnoreExt                *bool             `json:"ignoreExt,omitempty"`
}

type TDocdbQueryOption struct {
	Sort        interface{} `json:"sort,omitempty"`
	Limit       int64       `json:"limit,omitempty"`
	Skip        int64       `json:"skip,omitempty"`
	Upsert      bool        `json:"upsert,omitempty"`
	Project     bson.M      `json:"project,omitempty"`
	ProjectEasy []string    `json:"proj,omitempty"`
}

// query operator
type TRegexOperator struct {
	Regex string `json:"$regex"`
}
type TOROperator struct {
	OR []interface{} `json:"$or"`
}
type TANDOperator struct {
	AND []interface{} `json:"$and"`
}

// pipeline operator
type TMatchOperator struct {
	Match interface{} `json:"$match"`
}
type TAddFieldsOperator struct {
	AddFields interface{} `json:"$addFields"`
}
type TLookupOperator struct {
	Lookup interface{} `json:"$lookup"`
}
type TUnwindOperator struct {
	Unwind interface{} `json:"$unwind"`
}
type TSortOperator struct {
	Sort interface{} `json:"$sort"`
}
type TProjectOperator struct {
	Project interface{} `json:"$project"`
}

// query response
type TQueryResponse struct {
	Result TQueryResponseData `json:"result"`
}

type TQueryResponseData struct {
	Data interface{} `json:"data"`
}

// op
type TDocdbOp string

const (
	OP_AGGREGATE   TDocdbOp = "aggregate"
	OP_FIND        TDocdbOp = "find"
	OP_FIND_ONE    TDocdbOp = "findone"
	OP_INSERT      TDocdbOp = "insert"
	OP_UPSERT_ONE  TDocdbOp = "upsertone"
	OP_REPLACE_ONE TDocdbOp = "replaceone"
	OP_UPDATE_ONE  TDocdbOp = "updateone"
)

func (o TDocdbOp) String() string {
	return string(o)
}
