package query

import (
	"bytes"
	"cia/api/preference"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var defaultClient = &http.Client{Timeout: time.Second * 3}

type QueryPayload struct {
	Op        string      `json:"op"` // operation code
	Query     DocdbQuery  `json:"query"`
	Convert   bool        `json:"convert"`
	Documents interface{} `json:"documents,omitempty"`
}

type DocdbQuery struct {
	IgnoreExtensions bool `json:"ignoreExt"`

	Collection               string   `json:"collection"`
	DocumentId               string   `json:"documentId,omitempty"` // if not nil => will overwrite to Filter[_id]
	CaseInsensitivityId      *bool    `json:"caseInsensitivityId,omitempty"`
	CaseInsensitivityAddress *bool    `json:"caseInsensitivityAddress,omitempty"`
	Filter                   bson.M   `json:"filter,omitempty"`
	UpdateSet                bson.M   `json:"updateSet,omitempty"`
	Pipeline                 []bson.M `json:"pipeline,omitempty"`
	Options                  struct {
		Sort        *bson.D  `json:"sort,omitempty"`
		Limit       *int64   `json:"limit,omitempty"`
		Skip        *int64   `json:"skip,omitempty"`
		Upsert      bool     `json:"upsert,omitempty"`
		Project     bson.M   `json:"project,omitempty"`
		ProjectEasy []string `json:"proj,omitempty"`
	} `json:"options,omitempty"`
}

type responseQuery[T any] struct {
	Result struct {
		Data    T      `json:"data"`
		Message string `json:"message"`
	} `json:"result"`
}

func ExecQuery[T map[string]interface{} | []map[string]interface{} | string /*write result*/ | any](query QueryPayload, zeroSafety ...bool) (data T, err error) {
	url := preference.GetApiServerUrl("dslQuery")
	if url == "" {
		return data, errors.New("empty uri, check preference")
	}
	return ExecQueryWithUrl[T](url, query, zeroSafety...)
}

func ExecQueryWithUrl[T map[string]interface{} | []map[string]interface{} | string /*write result*/ | any](url string, query QueryPayload, zeroSafety ...bool) (data T, err error) {
	reqBody, err := json.Marshal(query)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := defaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var parsed responseQuery[T]
	if err = json.Unmarshal(resBody, &parsed); err != nil {
		return
	}

	switch msg := parsed.Result.Message; {
	case msg == "":
		// skip
	case zeroSafety != nil && zeroSafety[0] && strings.Contains(msg, notFoundDocumentSubstr):
		// skip
	default:
		err = errors.New(msg)
		return
	}
	return parsed.Result.Data, nil
}
