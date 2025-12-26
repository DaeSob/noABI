package types

import (
	"cia/common/utils"
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_HexToTopics(t *testing.T) {
	strTopics := []string{
		"0xc36e2e152ff8ae843ed0c8b86c2e40c378cd3c154151fe44c54430f7517037d6",
	}

	topics := HexToTopics(strTopics)
	assert.Equal(t, len(topics), 1)
	assert.Equal(t, topics[0].HexString(), "0xc36e2e152ff8ae843ed0c8b86c2e40c378cd3c154151fe44c54430f7517037d6")
}

func Test_CreateLog(t *testing.T) {
	address := "0xcfa002f78ed8f008ca2a5c08b7ce611dde3f3f88"
	blockHash := "0xb9b15aa8945272d8d836eec41afdfc4272d4c43502cf5a9c3d13295dd86ea0c6"
	blockNumber := "0x6562990"
	data := "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000002d3e297470b3c17fd9c2500625c53fca04bf24df0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e0000000000000000000000002db72e0565cd7607b77203301aafdaf0518b4e420000000000000000000000007b4cc5d0af2c5a3f96384324e0ad91d3e42c84bd000000000000000000000000e57c7ea7e43467b54900a4e37e31d7e0cd490382"
	logIndex := "0xa"
	removed := "false"
	topics := []string{"0xc36e2e152ff8ae843ed0c8b86c2e40c378cd3c154151fe44c54430f7517037d6"}
	transactionHash := "0xf257619cac21d33989844bc366b31cf94f7c92912ba6a81b4fa438d3e85adf6e"
	transactionIndex := "0x0"

	log := CreateLog(
		removed,
		logIndex,
		transactionIndex,
		transactionHash,
		blockHash,
		blockNumber,
		address,
		topics,
		data,
	)

	assert.Equal(
		t,
		log.Address.HexString(),
		"0xcfa002f78ed8f008ca2a5c08b7ce611dde3f3f88",
	)
	assert.Equal(
		t,
		log.BlockHash.HexString(),
		"0xb9b15aa8945272d8d836eec41afdfc4272d4c43502cf5a9c3d13295dd86ea0c6",
	)
	assert.Equal(
		t,
		log.BlockNumber,
		utils.HexToUint64("0x6562990"),
	)
	assert.Equal(
		t,
		utils.BytesToHexString(log.Data),
		"0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000002d3e297470b3c17fd9c2500625c53fca04bf24df0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e0000000000000000000000002db72e0565cd7607b77203301aafdaf0518b4e420000000000000000000000007b4cc5d0af2c5a3f96384324e0ad91d3e42c84bd000000000000000000000000e57c7ea7e43467b54900a4e37e31d7e0cd490382",
	)
	assert.Equal(
		t,
		log.LogIndex,
		utils.HexToUint64("0xa"),
	)
	assert.Equal(
		t,
		log.Removed,
		false,
	)
	assert.Equal(
		t,
		len(log.Topics),
		1,
	)
	assert.Equal(
		t,
		log.TransactionHash.HexString(),
		"0xf257619cac21d33989844bc366b31cf94f7c92912ba6a81b4fa438d3e85adf6e",
	)
	assert.Equal(
		t,
		log.TransactionIndex,
		utils.HexToUint64("0x0"),
	)

}

func Test_GetEventSignature(t *testing.T) {
	address := "0xcfa002f78ed8f008ca2a5c08b7ce611dde3f3f88"
	blockHash := "0xb9b15aa8945272d8d836eec41afdfc4272d4c43502cf5a9c3d13295dd86ea0c6"
	blockNumber := "0x6562990"
	data := "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000002d3e297470b3c17fd9c2500625c53fca04bf24df0000000000000000000000003cae9255b7ad17df118ed4f2338ff80c3ee3a15e0000000000000000000000002db72e0565cd7607b77203301aafdaf0518b4e420000000000000000000000007b4cc5d0af2c5a3f96384324e0ad91d3e42c84bd000000000000000000000000e57c7ea7e43467b54900a4e37e31d7e0cd490382"
	logIndex := "0xa"
	removed := "false"
	topics := []string{"0xc36e2e152ff8ae843ed0c8b86c2e40c378cd3c154151fe44c54430f7517037d6"}
	transactionHash := "0xf257619cac21d33989844bc366b31cf94f7c92912ba6a81b4fa438d3e85adf6e"
	transactionIndex := "0x0"

	log := CreateLog(
		removed,
		logIndex,
		transactionIndex,
		transactionHash,
		blockHash,
		blockNumber,
		address,
		topics,
		data,
	)

	sig := log.GetEventSignature().HexString()

	assert.Equal(t, sig, "0xc36e2e152ff8ae843ed0c8b86c2e40c378cd3c154151fe44c54430f7517037d6")
}

func Test_MapToLog(t *testing.T) {
	mp := map[string]interface{}{}
	mp["address"] = "0xbaffb941d4acb8c6068f8bcaf6810e0c3df9a345"
	mp["blockHash"] = "0xb9b15aa8945272d8d836eec41afdfc4272d4c43502cf5a9c3d13295dd86ea0c6"
	mp["blockNumber"] = "0x6562990"
	mp["data"] = "0x0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000e7f0b8a2df657903910dbfae0f8ae85138b5cee6000000000000000000000000000000000000000000000d8ca71c301620870000"
	mp["logIndex"] = "0x3"
	mp["removed"] = false
	mp["topics"] = []interface{}{
		"0x9d5cac4d5d1c90463091cb0a5d1f027a1a297e640f21d2587915c31b3eef1ff2",
	}
	mp["transactionHash"] = "0xf257619cac21d33989844bc366b31cf94f7c92912ba6a81b4fa438d3e85adf6e"
	mp["transactionIndex"] = "0x0"

	log := &TLog{}
	log.MapToLog(mp)

	assert.Equal(t, log.ToJsonString(false), `{"address":"0xbaffb941d4acb8c6068f8bcaf6810e0c3df9a345","blockHash":"0xb9b15aa8945272d8d836eec41afdfc4272d4c43502cf5a9c3d13295dd86ea0c6","blockNumber":"0x6562990","data":"0x0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000e7f0b8a2df657903910dbfae0f8ae85138b5cee6000000000000000000000000000000000000000000000d8ca71c301620870000","logIndex":"0x3","removed":"false","topics":["0x9d5cac4d5d1c90463091cb0a5d1f027a1a297e640f21d2587915c31b3eef1ff2"],"transactionHash":"0xf257619cac21d33989844bc366b31cf94f7c92912ba6a81b4fa438d3e85adf6e","transactionIndex":"0x0"}`)
}
