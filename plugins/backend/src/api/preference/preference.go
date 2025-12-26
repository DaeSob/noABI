package preference

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"cia/api/logHandler"
	YamlLoader "cia/common/yaml"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - server
type TServer struct {
	Host       string
	Debug      bool
	Supervisor string
	PathPrefix string
	Shutdown   bool
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - authorization
type TAuthorization struct {
	Enable bool
	Signer string
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// section - auth address list
type TAuthAddress struct {
	address string
	enable  bool
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// preference struct
type TPreference struct {
	strModulePath  string
	iModulePathLen int
	strConfigPath  string

	//yaml
	mapYaml map[string]interface{}

	//lock
	lockMutex sync.Mutex

	// attributes
	// common
	tServer        TServer
	tAuthorization TAuthorization

	timeZone int64

	// RPC
	mRPC map[string]TChainInfo

	// Keystore
	keystorePath string

	// Event Log
	tEventLogger TEventLogger
}

// singleton instance for preference
var lockMutex sync.Mutex
var instance *TPreference

func GetInstance() *TPreference {
	lockMutex.Lock()
	defer lockMutex.Unlock()
	if instance == nil {
		instance = &TPreference{}
	}
	return instance
}

// function for initialize for preference
func Initialize(_build string) (err error) {
	defer func() {
		s := recover()
		if s != nil {
			err = fmt.Errorf("%v", s)
		}
	}()

	Loadenv()
	Loadyaml(_build)

	SetPreference()

	return
}

func (inst *TPreference) _find(_strKey string) interface{} {
	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()
	return inst.mapYaml[_strKey]
}

func Loadenv() {
	inst := GetInstance()

	inst.strModulePath, _ = os.Getwd()
	inst.iModulePathLen = len(inst.strModulePath)
}

func Loadyaml(_build string) {
	inst := GetInstance()

	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()

	inst.mapYaml = nil
	inst.mapYaml = make(map[string]interface{})

	for {
		path := "./config/" + _build + "/preference.yaml"
		e := YamlLoader.LoadFromFile(path, &inst.mapYaml)
		if e != nil {
			logHandler.Write("initialize", 0, e.Error())
			panic(e)
		}
		inst.strConfigPath = filepath.Dir(path)
		break
	}
}

func SetPreference() {
	// common
	onSetCommonPreference()

	onSetTimeZone()

	// set Chain RPC
	onSetRPC()

	// keystore info
	onSetKeystore()

	// Logger
	onSetEventLogger()

}
