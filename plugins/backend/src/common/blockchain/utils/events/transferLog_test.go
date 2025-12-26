package events

import (
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_GetTransferLog(t *testing.T) {
	url := "https://public-en-cypress.klaytn.net"
	fromBlock := 90356032
	toBlock := 90356032
	contractAddr := "0xcee8faf64bb97a73bb51e115aa89c17ffa8dd167"
	fromAddr := ""
	toAddr := "0x99217904f418d1ede238c8a104d15603070cfa1e"
	id := 1

	res := GetTransferLog(
		url,
		uint64(fromBlock),
		uint64(toBlock),
		contractAddr,
		fromAddr,
		toAddr,
		int64(id),
	)

	assert.Equal(t, res.ToString(false), `{"error":null,"id":1,"jsonrpc":"2.0","result":[{"address":"0xcee8faf64bb97a73bb51e115aa89c17ffa8dd167","blockHash":"0x8a0ab9eaba5f7d720cbd8f8e8582d2ac1123833a30b7b660f662d7b5bf79e667","blockNumber":"0x562b940","data":"0x00000000000000000000000000000000000000000000000000000000000bb954","logIndex":"0x29","removed":false,"topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x000000000000000000000000afcfa70428f457bea047800bb1cea5b2dadca622","0x00000000000000000000000099217904f418d1ede238c8a104d15603070cfa1e"],"transactionHash":"0xa1cf7d8cd03cd1e5210495393d119a5fa9039404ae9cee7d3bb0d63963df4dc0","transactionIndex":"0x5"}]}`)
}
