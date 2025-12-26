package query

import (
	"bytes"
	"cia/common/utils/http/httpRequest"
	jMapper "cia/common/utils/json/mapper"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

func BenchmarkJson(b *testing.B) {
	url := "http://43.202.6.216:5000/dc/v1/documents/queries/op"
	query := QueryPayload{Op: "find", Query: DocdbQuery{Collection: "token_metadata"}}
	client := &http.Client{}

	for i := 0; i < b.N; i++ {
		reqBody, err := json.Marshal(query)
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		res.Body.Close()
		parsed := struct {
			Result struct {
				Data    interface{} `json:"data"`
				Message string      `json:"message"`
			} `json:"result"`
		}{}
		if err = json.Unmarshal(resBody, &parsed); err != nil {
			panic(err)
		}

		switch msg := parsed.Result.Message; {
		case msg == "":
			// skip
		default:
			panic(errors.New(msg))
		}

		foo := parsed.Result.Data
		_ = foo
	}

}

func BenchmarkJmap(b *testing.B) {
	url := "http://43.202.6.216:5000/dc/v1/documents/queries/op"
	query := QueryPayload{Op: "find", Query: DocdbQuery{Collection: "token_metadata"}}

	for i := 0; i < b.N; i++ {
		jBytes, err := jMapper.ToJson(query)
		if err != nil {
			panic(err)
		}
		jMap, err := jMapper.NewBytes(jBytes)
		if err != nil {
			panic(err)
		}

		bytes, err := httpRequest.PostRequestToBytes(url, jMap.Print(), 5*1000)
		if err != nil {
			panic(err)
		}
		res := struct {
			Result struct {
				Data    interface{} `json:"data"`
				Message string      `json:"message"`
			} `json:"result"`
		}{}
		if err = json.Unmarshal(bytes, &res); err != nil {
			panic(err)
		}

		switch msg := res.Result.Message; {
		case msg == "":
			// skip
		default:
			panic(errors.New(msg))
		}
		foo := res.Result.Data
		_ = foo
	}
}
