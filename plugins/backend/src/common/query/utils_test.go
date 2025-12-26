package query

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type UserOperation struct {
	Sender               string  `json:"sender"`
	Nonce                *string `json:"nonce"`
	InitCode             string  `json:"initCode"`
	CallData             string  `json:"callData"`
	CallGasLimit         string  `json:"callGasLimit"`
	VerificationGasLimit string  `json:"verificationGasLimit"`
	PreVerificationGas   string  `json:"preVerificationGas"`
	MaxFeePerGas         string  `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string  `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string  `json:"paymasterAndData"`
	Signature            string  `json:"signature"`
	Foobar               map[string]string
	Barfoo               map[string]map[string]string
}

type QueuedUserOperation struct {
	UUID          string        `json:"uuid"`
	RequestId     *string       `json:"requestId"`
	EstimateGas   string        `json:"estimategas"`
	UserOperation UserOperation `json:"userOperation"`
}

func s2p(s string) *string { return &s }

func TestStructOrMapToBasicUpdateSet(t *testing.T) {
	foo := QueuedUserOperation{
		UUID:        "uuid",
		RequestId:   nil,
		EstimateGas: "",
		UserOperation: UserOperation{
			Sender: "sender",
			Nonce:  s2p("nonce"),
			Foobar: map[string]string{"foo": "bar"},
			Barfoo: map[string]map[string]string{"barfoo": {"bar": "foo"}},
		},
	}
	bar := make(bson.M)
	StructOrMapToBasicUpdateSet(&bar, foo, "slots")
	foobar, _ := json.MarshalIndent(bar, "", "    ")
	t.Log("\n\n\n", string(foobar), "\n\n\n")
	t.Fatal("SUCCESS")
}
