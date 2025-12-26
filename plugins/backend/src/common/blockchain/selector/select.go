package selector

import (
	"cia/common/blockchain/types"
	jMapper "cia/common/utils/json/mapper"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	mABI         = map[string]types.TMethod{}
	mContractABI = make(map[string]abi.ABI)
	lockMutex    = sync.Mutex{}
)

func Select(_path, _contractName, _name string) types.TMethod {

	key := fmt.Sprintf("%v.%v", _contractName, _name)
	selectedAbi := mABI[key]
	if len(selectedAbi.Name) == 0 {
		var abi types.TMethod
		abiFile := fmt.Sprintf("%v/%v.%v.json", _path, _contractName, _name)
		byteData, _ := ioutil.ReadFile(abiFile)
		jMapper.FromJson(byteData, &abi)
		mABI[key] = abi
		selectedAbi = abi
	}
	if len(selectedAbi.Name) == 0 {
		panic("not found abi")
	}
	return selectedAbi

}

func SelectContractAbi(_path, _contractName string) (abi.ABI, error) {

	lockMutex.Lock()
	defer lockMutex.Unlock()
	selectedContract, exist := mContractABI[_contractName]
	if !exist {
		abiFile := fmt.Sprintf("%v/%v.json", _path, _contractName)
		byteData, err := ioutil.ReadFile(abiFile)
		if err != nil {
			return selectedContract, err
		}
		strSelectedAbi := string(byteData)
		parsed, err := abi.JSON(strings.NewReader(strSelectedAbi))
		if err != nil {
			return selectedContract, err
		}
		mContractABI[_contractName] = parsed
		return parsed, nil
	}
	return mContractABI[_contractName], nil

}

func SelectMethod(_path, _contractName, _name string) (method abi.Method, err error) {

	contractABI, err := SelectContractAbi(_path, _contractName)
	if err != nil {
		errMessage := fmt.Sprintf("not exist abi file(%v/%v.json)", _path, _contractName)
		panic(errMessage)
	}
	method, exist := contractABI.Methods[_name]
	if !exist {
		errMessage := fmt.Sprintf("not exist abi method(%v in %v/%v.json)", _name, _path, _contractName)
		panic(errMessage)
	}
	return

}
