package preference

import (
	"cia/api/logHandler"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Constants for event logger modes and fields
const (
	ModeLive = "live"
	ModeDev  = "dev"

	FieldDB     = "db"
	FieldActors = "actors"

	DefaultVersion = "V2.0.0"
	DefaultMode    = ModeLive
)

// Error messages
const (
	ErrChainNotFound              = "chain not found: %s"
	ErrChainAlreadyExists         = "chain already exists: %s"
	ErrEventLoggerSectionNotFound = "event-logger section not found in configuration. Please add 'event-logger' section with 'version' field in yaml file"
	ErrInvalidEventLoggerFormat   = "invalid event-logger section format"
	ErrVersionNotFound            = "version field not found in event-logger section. Please add 'version' field in event-logger section"
)

type TCollection struct {
	Address   string      `json:"address"`
	Enable    bool        `json:"enable"`
	Interface string      `json:"interface"`
	ExtraData interface{} `json:"extraData"`
}

type TCollector struct {
	StartBlockNumber uint64
	Collections      []string
}

type TActor struct {
	Collectors []string
}

type TEvetLoggerInfo struct {
	Mode             string
	ChainId          string
	RPC              string
	DB               string
	Retries          int
	Path             string
	CommitBlockCount uint8
	Period           int64
	LogRange         uint64
	InputLogging     bool
	BlockLogging     bool
	Collections      map[string]TCollection
	Collectors       map[string]TCollector
	Actors           map[string]TActor
}

type TEventLogger struct {
	Version     string
	Chains      []string
	Interfaces  map[string][]string
	EventLogger map[string]*TEvetLoggerInfo
}

func onSetEventLogger() {
	inst := GetInstance()

	inst.tEventLogger.EventLogger = make(map[string]*TEvetLoggerInfo)
	inst.tEventLogger.Chains = make([]string, 0) // Initialize empty slice
	inst.tEventLogger.Interfaces = make(map[string][]string)
	inst.tEventLogger.Version = DefaultVersion // Default version

	eventLogger := inst.mapYaml["event-logger"]

	if eventLogger != nil {
		logHandler.Write("initialize", 0, "loading [event-logger]")

		logger := eventLogger.(map[interface{}]interface{})
		if logger["version"] != nil {
			inst.tEventLogger.Version = logger["version"].(string)
		}

		if logger["chains"] != nil {
			chains := logger["chains"].([]interface{})
			mapstructure.Decode(chains, &inst.tEventLogger.Chains)
			for _, ch := range chains {
				chain := logger[ch].(map[interface{}]interface{})
				eventLogger, err := decodeChainDataFromYaml(chain)
				if err != nil {
					logHandler.Write("initialize", 0, "failed to decode chain", ch.(string), "error:", err.Error())
					continue
				}

				logHandler.Write("initialize", 0, "chain id:", eventLogger.ChainId+",", "rpc:", eventLogger.RPC+",", "log db:", eventLogger.DB)
				inst.tEventLogger.EventLogger[ch.(string)] = eventLogger
			}
		}

		logHandler.Write("initialize", 0, "loaded [event-logger]")
	}

	solInterfaces := inst.mapYaml["interfaces"]
	if solInterfaces != nil {
		logHandler.Write("initialize", 0, "loading [interfaces]")
		solInf := solInterfaces.(map[interface{}]interface{})
		for name, v := range solInf {
			var infs []string
			mapstructure.Decode(v, &infs)
			inst.tEventLogger.Interfaces[name.(string)] = infs
			logHandler.Write("initialize", 0, "interface name:", name.(string))
		}
		logHandler.Write("initialize", 0, "loaded [interfaces]")
	}

}

func GetEventLoggerChains() []string {
	inst := GetInstance()

	return inst.tEventLogger.Chains
}

func GetEventLoggerEventInterfaces() *map[string][]string {
	inst := GetInstance()

	return &inst.tEventLogger.Interfaces
}

func GetEventLogger(_chain string) *TEvetLoggerInfo {
	inst := GetInstance()

	return inst.tEventLogger.EventLogger[_chain]
}

func GetEventLoggerFilePath(_chain string) string {
	inst := GetInstance()

	return inst.tEventLogger.EventLogger[_chain].Path
}

// decodeChainData decodes chain data from map to TEvetLoggerInfo
func decodeChainData(_chainData map[string]interface{}) (*TEvetLoggerInfo, error) {
	eventLogger := new(TEvetLoggerInfo)
	eventLogger.Collections = make(map[string]TCollection)
	eventLogger.Collectors = make(map[string]TCollector)
	eventLogger.Actors = make(map[string]TActor)

	// Decode mode
	if mode, ok := _chainData["mode"].(string); ok {
		eventLogger.Mode = mode
	} else {
		eventLogger.Mode = DefaultMode
	}

	// Decode chainId
	if chainId, ok := _chainData["chainId"].(string); ok {
		eventLogger.ChainId = chainId
	}

	// Decode RPC
	if rpc, ok := _chainData["rpc"].(string); ok {
		eventLogger.RPC = rpc
	}

	// Decode DB (optional, may not exist in dev mode)
	if db, ok := _chainData["db"].(string); ok {
		eventLogger.DB = db
	}

	// Decode path
	if path, ok := _chainData["path"].(string); ok {
		eventLogger.Path = path
	}

	// Decode commitBlockCount
	if _chainData["commitBlockCount"] != nil {
		if err := mapstructure.Decode(_chainData["commitBlockCount"], &eventLogger.CommitBlockCount); err != nil {
			return nil, fmt.Errorf("failed to decode commitBlockCount: %w", err)
		}
	}

	// Decode period
	if _chainData["period"] != nil {
		if err := mapstructure.Decode(_chainData["period"], &eventLogger.Period); err != nil {
			return nil, fmt.Errorf("failed to decode period: %w", err)
		}
	}

	// Decode logRange
	if _chainData["logRange"] != nil {
		if err := mapstructure.Decode(_chainData["logRange"], &eventLogger.LogRange); err != nil {
			return nil, fmt.Errorf("failed to decode logRange: %w", err)
		}
	}

	// Decode inputLogging
	if inputLogging, ok := _chainData["inputLogging"].(bool); ok {
		eventLogger.InputLogging = inputLogging
	} else {
		eventLogger.InputLogging = false
	}

	// Decode blockLogging
	if blockLogging, ok := _chainData["blockLogging"].(bool); ok {
		eventLogger.BlockLogging = blockLogging
	} else {
		eventLogger.BlockLogging = false
	}

	// Decode collections
	if _chainData["collections"] != nil {
		collections, ok := _chainData["collections"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid collections format")
		}
		for name, v := range collections {
			var info TCollection
			if err := mapstructure.Decode(v, &info); err != nil {
				return nil, fmt.Errorf("failed to decode collection %s: %w", name, err)
			}
			info.Address = strings.ToLower(info.Address)
			eventLogger.Collections[name] = info
		}
	}

	// Decode collectors
	if _chainData["collectors"] != nil {
		collectors, ok := _chainData["collectors"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid collectors format")
		}
		for name, v := range collectors {
			var info TCollector
			if err := mapstructure.Decode(v, &info); err != nil {
				return nil, fmt.Errorf("failed to decode collector %s: %w", name, err)
			}
			eventLogger.Collectors[name] = info
		}
	}

	// Decode actors (optional)
	if _chainData["actors"] != nil {
		actors, ok := _chainData["actors"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid actors format")
		}
		for name, v := range actors {
			var info TActor
			if err := mapstructure.Decode(v, &info); err != nil {
				return nil, fmt.Errorf("failed to decode actor %s: %w", name, err)
			}
			eventLogger.Actors[name] = info
		}
	}

	return eventLogger, nil
}

// decodeChainDataFromYaml decodes chain data from YAML format (map[interface{}]interface{})
func decodeChainDataFromYaml(_chain map[interface{}]interface{}) (*TEvetLoggerInfo, error) {
	// Convert to map[string]interface{} for decodeChainData
	chainData := make(map[string]interface{})

	if _chain["mode"] != nil {
		chainData["mode"] = _chain["mode"].(string)
	} else {
		chainData["mode"] = DefaultMode
	}
	if _chain["chainId"] != nil {
		chainData["chainId"] = _chain["chainId"].(string)
	}
	if _chain["rpc"] != nil {
		chainData["rpc"] = _chain["rpc"].(string)
	}
	if _chain["db"] != nil {
		chainData["db"] = _chain["db"].(string)
	}
	if _chain["path"] != nil {
		chainData["path"] = _chain["path"].(string)
	}
	if _chain["commitBlockCount"] != nil {
		chainData["commitBlockCount"] = _chain["commitBlockCount"]
	}
	if _chain["period"] != nil {
		chainData["period"] = _chain["period"]
	}
	if _chain["logRange"] != nil {
		chainData["logRange"] = _chain["logRange"]
	}
	if _chain["inputLogging"] != nil {
		chainData["inputLogging"] = _chain["inputLogging"].(bool)
	}
	if _chain["blockLogging"] != nil {
		chainData["blockLogging"] = _chain["blockLogging"].(bool)
	}
	// Convert collections from map[interface{}]interface{} to map[string]interface{}
	if _chain["collections"] != nil {
		collectionsYaml, ok := _chain["collections"].(map[interface{}]interface{})
		if ok {
			collections := make(map[string]interface{})
			for k, v := range collectionsYaml {
				collections[k.(string)] = v
			}
			chainData["collections"] = collections
		}
	}

	// Convert collectors from map[interface{}]interface{} to map[string]interface{}
	if _chain["collectors"] != nil {
		collectorsYaml, ok := _chain["collectors"].(map[interface{}]interface{})
		if ok {
			collectors := make(map[string]interface{})
			for k, v := range collectorsYaml {
				collectors[k.(string)] = v
			}
			chainData["collectors"] = collectors
		}
	}

	// Convert actors from map[interface{}]interface{} to map[string]interface{}
	if _chain["actors"] != nil {
		actorsYaml, ok := _chain["actors"].(map[interface{}]interface{})
		if ok {
			actors := make(map[string]interface{})
			for k, v := range actorsYaml {
				actors[k.(string)] = v
			}
			chainData["actors"] = actors
		}
	}

	return decodeChainData(chainData)
}

// updateChainFields updates specific fields in eventLogger from updateData
func updateChainFields(_eventLogger *TEvetLoggerInfo, _updateData map[string]interface{}) error {
	if _updateData["rpc"] != nil {
		if rpc, ok := _updateData["rpc"].(string); ok {
			_eventLogger.RPC = rpc
		}
	}
	if _updateData["path"] != nil {
		if path, ok := _updateData["path"].(string); ok {
			_eventLogger.Path = path
		}
	}
	if _updateData["commitBlockCount"] != nil {
		if err := mapstructure.Decode(_updateData["commitBlockCount"], &_eventLogger.CommitBlockCount); err != nil {
			return fmt.Errorf("failed to decode commitBlockCount: %w", err)
		}
	}
	if _updateData["period"] != nil {
		if err := mapstructure.Decode(_updateData["period"], &_eventLogger.Period); err != nil {
			return fmt.Errorf("failed to decode period: %w", err)
		}
	}
	if _updateData["logRange"] != nil {
		if err := mapstructure.Decode(_updateData["logRange"], &_eventLogger.LogRange); err != nil {
			return fmt.Errorf("failed to decode logRange: %w", err)
		}
	}
	if _updateData["inputLogging"] != nil {
		if inputLogging, ok := _updateData["inputLogging"].(bool); ok {
			_eventLogger.InputLogging = inputLogging
		}
	}
	if _updateData["blockLogging"] != nil {
		if blockLogging, ok := _updateData["blockLogging"].(bool); ok {
			_eventLogger.BlockLogging = blockLogging
		}
	}
	return nil
}

// updateCollections updates collections in eventLogger
func updateCollections(_eventLogger *TEvetLoggerInfo, _collections map[string]interface{}) error {
	_eventLogger.Collections = make(map[string]TCollection)
	for name, v := range _collections {
		var info TCollection
		if err := mapstructure.Decode(v, &info); err != nil {
			return fmt.Errorf("failed to decode collection %s: %w", name, err)
		}
		info.Address = strings.ToLower(info.Address)
		_eventLogger.Collections[name] = info
	}
	return nil
}

// updateCollectors updates collectors in eventLogger
func updateCollectors(_eventLogger *TEvetLoggerInfo, _collectors map[string]interface{}) error {
	_eventLogger.Collectors = make(map[string]TCollector)
	for name, v := range _collectors {
		var info TCollector
		if err := mapstructure.Decode(v, &info); err != nil {
			return fmt.Errorf("failed to decode collector %s: %w", name, err)
		}
		_eventLogger.Collectors[name] = info
	}
	return nil
}

// UpdateEventLoggerChain updates an existing chain in memory
func UpdateEventLoggerChain(_alias string, _updateData map[string]interface{}) error {
	inst := GetInstance()
	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()

	eventLogger := inst.tEventLogger.EventLogger[_alias]
	if eventLogger == nil {
		return fmt.Errorf(ErrChainNotFound, _alias)
	}

	// Update mapYaml
	eventLoggerYaml := inst.mapYaml["event-logger"].(map[interface{}]interface{})

	// Get chainYaml - handle both map[string]interface{} and map[interface{}]interface{}
	chainYamlRaw := eventLoggerYaml[_alias]
	if chainYamlRaw == nil {
		return fmt.Errorf("chain not found in mapYaml: %s", _alias)
	}

	// Convert to map[interface{}]interface{} if needed
	var chainYaml map[interface{}]interface{}
	if chainYamlStr, ok := chainYamlRaw.(map[string]interface{}); ok {
		// Convert map[string]interface{} to map[interface{}]interface{}
		chainYaml = make(map[interface{}]interface{})
		for k, v := range chainYamlStr {
			chainYaml[k] = v
		}
		// Update the original mapYaml with converted type
		eventLoggerYaml[_alias] = chainYaml
	} else if chainYamlInt, ok := chainYamlRaw.(map[interface{}]interface{}); ok {
		chainYaml = chainYamlInt
	} else {
		return fmt.Errorf("invalid chain data type in mapYaml: %s", _alias)
	}

	// Update fields from updateData
	for key, value := range _updateData {
		if key == FieldDB || key == FieldActors {
			continue // Skip db and actors in dev mode
		}
		chainYaml[key] = value
	}

	// Update memory structure
	if err := updateChainFields(eventLogger, _updateData); err != nil {
		return err
	}

	// Update collections
	if _updateData["collections"] != nil {
		collections, ok := _updateData["collections"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid collections format")
		}
		if err := updateCollections(eventLogger, collections); err != nil {
			return err
		}
		chainYaml["collections"] = _updateData["collections"]
	}

	// Update collectors
	if _updateData["collectors"] != nil {
		collectors, ok := _updateData["collectors"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid collectors format")
		}
		if err := updateCollectors(eventLogger, collectors); err != nil {
			return err
		}
		chainYaml["collectors"] = _updateData["collectors"]
	}

	return nil
}

// AddEventLoggerChain adds a new chain in memory
func AddEventLoggerChain(_chainId string, _chainData map[string]interface{}) error {
	inst := GetInstance()
	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()

	// Check if event-logger section exists in mapYaml
	if inst.mapYaml["event-logger"] == nil {
		return fmt.Errorf(ErrEventLoggerSectionNotFound)
	}

	eventLoggerYaml, ok := inst.mapYaml["event-logger"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(ErrInvalidEventLoggerFormat)
	}

	// Check if version exists
	if eventLoggerYaml["version"] == nil {
		return fmt.Errorf(ErrVersionNotFound)
	}

	// Check if chain already exists
	if inst.tEventLogger.EventLogger[_chainId] != nil {
		return fmt.Errorf(ErrChainAlreadyExists, _chainId)
	}

	// Ensure mode is dev
	_chainData["mode"] = ModeDev

	// Decode chain data
	eventLogger, err := decodeChainData(_chainData)
	if err != nil {
		return fmt.Errorf("failed to decode chain data: %w", err)
	}

	// Add to memory
	inst.tEventLogger.EventLogger[_chainId] = eventLogger
	inst.tEventLogger.Chains = append(inst.tEventLogger.Chains, _chainId)

	// Add to mapYaml (eventLoggerYaml is already checked above)
	eventLoggerYaml[_chainId] = _chainData

	// Update chains list in mapYaml
	if eventLoggerYaml["chains"] == nil {
		eventLoggerYaml["chains"] = []interface{}{}
	}
	chains := eventLoggerYaml["chains"].([]interface{})
	eventLoggerYaml["chains"] = append(chains, _chainId)

	return nil
}

// DeleteEventLoggerChain deletes a chain from memory
func DeleteEventLoggerChain(_chainId string) error {
	inst := GetInstance()
	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()

	// Check if chain exists
	if inst.tEventLogger.EventLogger[_chainId] == nil {
		return fmt.Errorf(ErrChainNotFound, _chainId)
	}

	// Delete from memory
	delete(inst.tEventLogger.EventLogger, _chainId)

	// Remove from chains list
	newChains := []string{}
	for _, ch := range inst.tEventLogger.Chains {
		if ch != _chainId {
			newChains = append(newChains, ch)
		}
	}
	inst.tEventLogger.Chains = newChains

	// Delete from mapYaml
	if inst.mapYaml["event-logger"] != nil {
		eventLoggerYaml := inst.mapYaml["event-logger"].(map[interface{}]interface{})
		delete(eventLoggerYaml, _chainId)

		// Update chains list in mapYaml
		if eventLoggerYaml["chains"] != nil {
			chains := eventLoggerYaml["chains"].([]interface{})
			newChainsYaml := []interface{}{}
			for _, ch := range chains {
				if ch.(string) != _chainId {
					newChainsYaml = append(newChainsYaml, ch)
				}
			}
			eventLoggerYaml["chains"] = newChainsYaml
		}
	}

	return nil
}
