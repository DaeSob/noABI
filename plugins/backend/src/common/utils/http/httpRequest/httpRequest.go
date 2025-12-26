package httpRequest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	jMapper "cia/common/utils/json/mapper"
)

var (
	defaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Millisecond * time.Duration(30000),
	}
)

// create request for get
func NewGetRequest(
	_strUri string,
) (*http.Request, error) {
	return http.NewRequest("GET", _strUri, nil)
}

// create request for post
func NewPostRequest(
	_strUri string,
	_body io.Reader,
) (*http.Request, error) {
	return http.NewRequest("POST", _strUri, _body)
}

/*
// prams
// i32Timeout - unit : millisecond
func NewHttpClient(_req *http.Request, _i32Timeout int) *http.Client {
	//필요시 헤더 추가 가능
	_req.Header.Add("Content-Type", "application/json")

	// set tls enable
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Client객체에서 Request 실행
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Millisecond * time.Duration(_i32Timeout),
	}
	return client
}
*/

func NewHttpClient(_req *http.Request, _i32Timeout int) *http.Client {
	//필요시 헤더 추가 가능
	_req.Header.Add("Content-Type", "application/json")

	return defaultClient
}

func PostRequestToBytes(_strUri string, _strBody string, _i32Timeout int) ([]byte, error) {
	req, err := NewPostRequest(_strUri, bytes.NewBufferString(_strBody))
	if err != nil {
		return nil, err
	}
	req.Close = true

	client := NewHttpClient(req, _i32Timeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 결과 출력
	bytes, e := ioutil.ReadAll(resp.Body)
	if e == nil && resp.StatusCode >= 400 {
		e = fmt.Errorf("%v %v", resp.StatusCode, string(bytes))
	}

	defer func() {
		resp.Body.Close()
		resp = nil
		client = nil
		bytes = nil
	}()

	return bytes, e
}

// prams
// i32Timeout - unit : millisecond
func PostRequestFromString(_strUri string, _strBody string, _i32Timeout int) (*jMapper.TJsonMap, error) {
	bytes, err := PostRequestToBytes(_strUri, _strBody, _i32Timeout)
	if err != nil {
		return nil, err
	}
	jsonMap, e := jMapper.NewBytes(bytes)

	bytes = nil
	return jsonMap, e
}

// prams
// i32Timeout - unit : millisecond
func PostRequestFromBytes(_strUri string, _byteBody []byte, _intTimeout int) (*jMapper.TJsonMap, error) {
	req, err := NewPostRequest(_strUri, bytes.NewBuffer(_byteBody))
	if err != nil {
		return nil, err
	}
	req.Close = true

	client := NewHttpClient(req, _intTimeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	jsonMap, e := jMapper.NewBytes(bytes)

	resp.Body.Close()
	resp = nil
	client = nil
	bytes = nil

	return jsonMap, e
}

// HTTP : 'GET'
// prams
// i32Timeout - unit : millisecond
func GetRequest(_strUri string, _i32Timeout int) (*jMapper.TJsonMap, error) {
	req, err := NewGetRequest(_strUri)
	if err != nil {
		return nil, err
	}
	req.Close = true

	client := NewHttpClient(req, _i32Timeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	jsonMap, e := jMapper.NewBytes(bytes)

	resp.Body.Close()
	resp = nil
	client = nil
	bytes = nil

	return jsonMap, e
}

func GetRequestToBytes(_strUri string, _i32Timeout int) ([]byte, error) {
	req, err := NewGetRequest(_strUri)
	if err != nil {
		return nil, err
	}
	req.Close = true

	client := NewHttpClient(req, _i32Timeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 결과 출력
	bytes, e := ioutil.ReadAll(resp.Body)

	defer func() {
		resp.Body.Close()
		resp = nil
		client = nil
		bytes = nil
	}()

	return bytes, e
}
