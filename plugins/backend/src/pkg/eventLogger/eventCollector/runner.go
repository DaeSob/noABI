package eventCollector

import (
	"cia/api/logHandler"
	"cia/api/preference"
	"cia/common/utils"
	"fmt"
)

// V2.0.0 By XeN
func Initialize() {
	chains := preference.GetEventLoggerChains()

	for _, chain := range chains {
		eventLoggerInfo := preference.GetEventLogger(chain)
		if eventLoggerInfo == nil {
			fmt.Printf("Warning: chain not found: %s\n", chain)
			continue
		}

		utils.Mkdir(eventLoggerInfo.Path)
		utils.Mkdir(eventLoggerInfo.Path + "/" + chain)
		// ready contract log path
		for name, _ := range eventLoggerInfo.Collectors {
			conPath := fmt.Sprintf("%v/%v/%v", eventLoggerInfo.Path, chain, name)
			utils.Mkdir(conPath)
		}
	}
}

func Run() error {
	worker()
	return nil
}

// InitializeChain initializes a specific chain (for dynamic addition)
func InitializeChain(_chain string) error {
	eventLoggerInfo := preference.GetEventLogger(_chain)
	if eventLoggerInfo == nil {
		return fmt.Errorf("chain not found: %s", _chain)
	}

	utils.Mkdir(eventLoggerInfo.Path)
	utils.Mkdir(eventLoggerInfo.Path + "/" + _chain)
	// ready contract log path
	for name, _ := range eventLoggerInfo.Collectors {
		conPath := fmt.Sprintf("%v/%v/%v", eventLoggerInfo.Path, _chain, name)
		utils.Mkdir(conPath)
	}
	return nil
}

// RunChain starts worker for a specific chain (for dynamic addition)
func RunChain(_chain string) error {
	go workerForChain(_chain)
	return nil
}

// StopChain stops all workers for a specific chain
func StopChain(_chain string) error {
	chainContextMutex.Lock()
	defer chainContextMutex.Unlock()

	if cancel, exists := chainCancels[_chain]; exists {
		cancel() // Cancel all goroutines for this chain
		delete(chainContexts, _chain)
		delete(chainCancels, _chain)
		logHandler.Write("trace", 0, "stopped all workers for chain", _chain)
		return nil
	}

	return fmt.Errorf("chain not found: %s", _chain)
}