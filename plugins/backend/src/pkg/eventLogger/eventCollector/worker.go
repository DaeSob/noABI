package eventCollector

import (
	"cia/api/logHandler"
	"cia/api/preference"
	"cia/common/blockchain/eventLog"
	"cia/common/blockchain/rpcRequest"
	"cia/common/blockchain/rpcRequest/eth"
	"cia/common/blockchain/utils/events/abiMapper"
	"cia/common/utils"
	rpc "cia/common/utils/rpcRequest"
	"context"
	"fmt"
	"math"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

// metaFileLocks manages mutexes for event.meta files
// Key format: "chain:collector"
var metaFileLocks = make(map[string]*sync.Mutex)
var metaFileLocksMutex sync.Mutex

// chainContexts manages contexts for each chain to allow cancellation
var chainContexts = make(map[string]context.Context)
var chainCancels = make(map[string]context.CancelFunc)
var chainContextMutex sync.Mutex

// getMetaFileLock returns a mutex for a specific chain:collector combination
func getMetaFileLock(_chain string, _collectorName string) *sync.Mutex {
	key := fmt.Sprintf("%s:%s", _chain, _collectorName)

	metaFileLocksMutex.Lock()
	defer metaFileLocksMutex.Unlock()

	if lock, exists := metaFileLocks[key]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	metaFileLocks[key] = lock
	return lock
}

type TNodeSyncingStatus struct {
	StartingBlock uint64 `json:"startingBlock"`
	CurrentBlock  uint64 `json:"currentBlock"`
	HighestBlock  uint64 `json:"highestBlock"`
}

func worker() (err error) {
	chains := preference.GetEventLoggerChains()
	for _, chain := range chains {
		go workerForChain(chain)
	}
	return nil
}

func workerForChain(_chain string) {
	// Create context for this chain
	chainContextMutex.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	chainContexts[_chain] = ctx
	chainCancels[_chain] = cancel
	chainContextMutex.Unlock()

	eventLoggerInfo := preference.GetEventLogger(_chain)

	collectors := eventLoggerInfo.Collectors

	// get event mapper
	mapper := getEventMapper(eventLoggerInfo)

	for collectorName, v := range collectors {
		var collector preference.TCollector
		mapstructure.Decode(v, &collector)
		var addresses []string
		activateCollections := map[string]string{}
		for _, name := range collector.Collections {
			if eventLoggerInfo.Collections[name].Enable {
				addresses = append(addresses, eventLoggerInfo.Collections[name].Address)
				activateCollections[eventLoggerInfo.Collections[name].Address] = name
			}
		}

		if len(addresses) > 0 {
			go workerForContract(
				ctx,
				_chain,
				activateCollections,
				collectorName,
				eventLoggerInfo,
				mapper,
			)

		}
	}
}

func workerForContract(
	_ctx context.Context,
	_chain string,
	_collections map[string]string,
	_collectorName string,
	_eventLoggerInfo *preference.TEvetLoggerInfo,
	_mapper abiMapper.TABIMappers,
) {
	// latest block number
	// init latestBlockNumber from event.meta file and current startBlockNumber
	var latestBlockNumber uint64

	// Get initial values from preference and event.meta
	currentEventLoggerInfo := preference.GetEventLogger(_chain)
	var currentStartBlockNumber uint64 = 0

	if currentEventLoggerInfo != nil {
		if collector, exists := currentEventLoggerInfo.Collectors[_collectorName]; exists {
			currentStartBlockNumber = collector.StartBlockNumber
		}
	}

	blockFromMeta := GetEventMetaInfo(_chain, _collectorName)
	latestBlockNumber = uint64(math.Max(float64(blockFromMeta), float64(currentStartBlockNumber)))

	for {
		// Check if context is cancelled (chain was deleted)
		select {
		case <-_ctx.Done():
			logHandler.Write("trace", 0, "stopping worker for chain", _chain, "collector", _collectorName, "reason: chain deleted")
			return
		default:
		}

		// Read latest startBlockNumber from preference on each cycle
		// This allows latestBlockNumber to be updated when startBlockNumber changes
		currentEventLoggerInfo := preference.GetEventLogger(_chain)
		var currentStartBlockNumber uint64 = 0

		if currentEventLoggerInfo != nil {
			if collector, exists := currentEventLoggerInfo.Collectors[_collectorName]; exists {
				currentStartBlockNumber = collector.StartBlockNumber
			}
		}

		// Read event.meta file on each cycle to get latest value
		blockFromMeta = GetEventMetaInfo(_chain, _collectorName)

		// Use the maximum of event.meta and current startBlockNumber
		// This ensures we always use the latest configuration
		newLatestBlockNumber := uint64(math.Max(float64(blockFromMeta), float64(currentStartBlockNumber)))

		// Update latestBlockNumber if it changed (especially when startBlockNumber is reduced)
		if newLatestBlockNumber != latestBlockNumber {
			latestBlockNumber = newLatestBlockNumber
			logHandler.Write("trace", 0, "updated latestBlockNumber for chain", _chain, "collector", _collectorName, "to", latestBlockNumber, "startBlockNumber", currentStartBlockNumber)
		}

		// init eth for request
		getLogAndSaveToFile(
			_chain,
			_collections,
			_collectorName,
			&latestBlockNumber,
			_eventLoggerInfo,
			_mapper,
		)

		// apply period with context cancellation check
		select {
		case <-_ctx.Done():
			logHandler.Write("trace", 0, "stopping worker for chain", _chain, "collector", _collectorName, "reason: chain deleted")
			return
		case <-time.After(time.Second * time.Duration(_eventLoggerInfo.Period)):
			// Continue to next cycle
		}
	}
}

func getLogAndSaveToFile(
	_chain string,
	_collections map[string]string,
	_collectorName string,
	_latestBlockNumber *uint64,
	_eventLoggerInfo *preference.TEvetLoggerInfo,
	_mapper abiMapper.TABIMappers,
) {
	defer func() {
		s := recover()
		if s != nil {
			err := fmt.Errorf("%v", s)
			logHandler.Write("trace", 0, "panic", err.Error())
			/*
				node가 syncing 진행 중일 때 발생 하는거 같다
					- One of the blocks specified in filter (fromBlock, toBlock or blockHash) cannot be found.
				RPC Server의 LP에서 응답을 하는거 같다
					- no backend is currently healthy to serve traffic
			*/
		}
	}()

	eth := eth.TEth{RPC: rpcRequest.TRPC{URL: _eventLoggerInfo.RPC}}
	logPath := fmt.Sprintf("%v/%v", _eventLoggerInfo.Path, _chain)
	id := int64(1)
	var toBlockNumber uint64
	var syncing bool

	/*
		node의 동기화 상태를 조회 한다
		- 동기화 중일 땐 JSON 값을 수신
		- 동기화가 완료 되었을 땐 false 값을 수신 한다
	*/
	resSyncingStatus := eth.GetSyncingStatus(id)
	if resSyncingStatus.Error != nil {
		rpcError := rpc.TRpcError{}
		mapstructure.Decode(resSyncingStatus.Error, &rpcError)
		panicMessage := fmt.Sprintf("request GetSyncingStatus", rpcError.Message)
		panic(panicMessage)
	}

	if reflect.TypeOf(resSyncingStatus.Result).Kind() == reflect.Bool { // 동기화 완료 상태
		res := eth.GetBlockNumber(id)
		if res.Error != nil {
			rpcError := rpc.TRpcError{}
			mapstructure.Decode(res.Error, &rpcError)
			panicMessage := fmt.Sprintf("request GetBlockNumber(%v)", rpcError.Message)
			panic(panicMessage)
		}
		hexBlockNumber := res.ResultToString()
		toBlockNumber = (utils.HexToUint64(hexBlockNumber) - uint64(_eventLoggerInfo.CommitBlockCount))
	} else {
		var nodeSyncingStatus TNodeSyncingStatus
		mapstructure.Decode(resSyncingStatus.Result, nodeSyncingStatus)
		if nodeSyncingStatus.CurrentBlock > 0 {
			toBlockNumber = nodeSyncingStatus.CurrentBlock - 1
		} else {
			toBlockNumber = 0
		}
		syncing = true
		logHandler.Write("trace", 0, "node syncing", resSyncingStatus.Result)
	}

	if *_latestBlockNumber >= toBlockNumber {
		logHandler.Write("block", 0, _chain, _collectorName, "syncing:", syncing, "ignore block from:", (*_latestBlockNumber + 1), "to:", toBlockNumber)
		return
	}

	// set block range
	fromBlock := *_latestBlockNumber + 1 // [이전 블록 + 1]부터 log 수집
	maxBlock := *_latestBlockNumber + _eventLoggerInfo.LogRange
	toBlock := uint64(math.Min(float64(maxBlock), float64(toBlockNumber)))

	// get event logs
	logs := eventLog.GetContractEventLogsEx(
		_eventLoggerInfo.RPC,
		uint64(fromBlock),
		uint64(toBlock),
		_collections,
		id,
		_mapper,
		_eventLoggerInfo.InputLogging,
	)

	// write log file
	blockLogs := logs.GetBlockLog()
	savePath := fmt.Sprintf("%v/%v", logPath, _collectorName)
	for _, blockLog := range blockLogs {
		if _eventLoggerInfo.InputLogging {
			logHandler.Write("logCollector", 0, utils.InterfaceToJsonString(blockLog, true))
		}
		blockLog.SaveToFile(savePath)
	}

	// update latest number
	*_latestBlockNumber = uint64(toBlock)
	UpdateEventMetaInfo(_chain, _collectorName, *_latestBlockNumber)
	if _eventLoggerInfo.BlockLogging {
		logHandler.Write("block", 0, _chain, _collectorName, "subscribe block from:", fromBlock, "to:", *_latestBlockNumber)
	}

	// release memory
	blockLogs = nil
	logs = nil
}

// write file for save log file
func _getEventMetaFilePath(_chain string, _contractName string) string {
	logPath := preference.GetEventLoggerFilePath(_chain)
	filename := "event.meta"

	return fmt.Sprintf("%v/%v/%v/%v", logPath, _chain, _contractName, filename)
}

// rewrite on meta file with lock
func _writeToMetaFile(_filePath string, _strText string, _lock *sync.Mutex) {
	_lock.Lock()
	defer _lock.Unlock()

	f, e := os.OpenFile(_filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	f.WriteString(_strText)
}

// ReadFromMetaFile : read from meta file with lock
func _readFromMetaFile(_metaFilePath string, _lock *sync.Mutex) (string, error) {
	_lock.Lock()
	defer _lock.Unlock()

	// Check if file exists
	if _, err := os.Stat(_metaFilePath); os.IsNotExist(err) {
		return "", nil // File doesn't exist, return empty string
	}

	f, e := os.OpenFile(_metaFilePath, os.O_RDONLY, 0644)
	if e != nil {
		fmt.Println(e)
		return "", e
	}
	defer f.Close()

	fileSize := utils.GetFileSize(_metaFilePath)
	if fileSize == 0 {
		return "", nil
	}

	buffer := make([]byte, fileSize)
	_, e = f.Read(buffer)
	if e != nil {
		return "", e
	}

	// release memory
	defer func() {
		buffer = nil
	}()

	return string(buffer), nil
}

func GetEventMetaInfo(_chain string, _contractName string) uint64 {
	metaPath := _getEventMetaFilePath(_chain, _contractName)
	lock := getMetaFileLock(_chain, _contractName)
	info, err := _readFromMetaFile(metaPath, lock)
	if err != nil || info == "" {
		return 0
	}

	return utils.StringToUint64(info, 10)
}

func UpdateEventMetaInfo(_chain string, _contractName string, _latestBlockNumber uint64) {
	metaPath := _getEventMetaFilePath(_chain, _contractName)
	lock := getMetaFileLock(_chain, _contractName)
	_writeToMetaFile(
		metaPath,
		fmt.Sprintf("%v", _latestBlockNumber),
		lock,
	)
}
