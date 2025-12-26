// V2.0.0 By XeN
package preference

import (
	"cia/api/logHandler"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// V2::Default Chain

type TChainInfo struct {
	ChainId    string `json:"chainId"`
	ChainName  string `json:"chainName"`
	RPCUrl     string `json:"rpcUrl"`
	RPCOpt     string `json:"rpcOpt"`
	GasLimit   uint64 `json:"gasLimit"`
	BlockCycle uint64 `json:"blockCycle"` //block 생성 주기
	Supported  string `json:"supported"`
	Migration  bool   `json:"migration"`
	MSupported map[string]bool
}

func onSetRPC() {
	inst := GetInstance()

	inst.mRPC = make(map[string]TChainInfo)

	rpcs := inst.mapYaml["rpc"]
	if rpcs != nil {
		logHandler.Write("initialize", 0, "loading [rpc]")
		infos := rpcs.(map[interface{}]interface{})
		for key, value := range infos {
			var rpcInfo TChainInfo
			mapstructure.Decode(value, &rpcInfo)
			inst.mRPC[key.(string)] = TChainInfo{
				ChainId:    key.(string),
				ChainName:  rpcInfo.ChainName,
				RPCUrl:     rpcInfo.RPCUrl,
				RPCOpt:     rpcInfo.RPCOpt,
				GasLimit:   rpcInfo.GasLimit,
				BlockCycle: rpcInfo.BlockCycle,
				Migration:  rpcInfo.Migration,
				MSupported: make(map[string]bool),
			}
			supported := strings.Split(rpcInfo.Supported, ",")
			for _, val := range supported {
				if len(val) > 0 {
					inst.mRPC[key.(string)].MSupported[val] = true
				}
			}
			logHandler.Write("initialize", 0, "name:", rpcInfo.ChainName+",", "chainId:", key.(string)+",", "rpc:", rpcInfo.RPCUrl+",", "supported:", rpcInfo.Supported)
		}
		logHandler.Write("initialize", 0, "loaded [rpc]")
	}
}

func GetChainInfo(_chainId string) TChainInfo {
	inst := GetInstance()

	return inst.mRPC[_chainId]
}

func CheckVPASupportedChainId(_chainId string) {
	inst := GetInstance()

	if len(_chainId) == 0 {
		panic("required chainId")
	}
	chainInfo, exist := inst.mRPC[_chainId]
	if !exist {
		panic("not support chain id")
	}
	if !chainInfo.MSupported["VPA"] {
		panic("not support VPA for this chain id")
	}
}
