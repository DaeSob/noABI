package eventActor

import (
	"cia/api/preference"
	"time"
)

func worker() (err error) {

	chains := preference.GetEventLoggerChains()
	for _, chain := range chains {
		eventLoggerInfo := preference.GetEventLogger(chain)
		for name, _ := range eventLoggerInfo.Actors {
			go workerForChain(chain, name)
		}
	}
	return nil
}

func workerForChain(_chain string, _actorName string) {

	eventLoggerInfo := preference.GetEventLogger(_chain)
	actor := eventLoggerInfo.Actors[_actorName]

	activate := false
	for _, name := range actor.Collectors {
		collector := eventLoggerInfo.Collectors[name]
		for _, contractName := range collector.Collections {
			if eventLoggerInfo.Collections[contractName].Enable {
				activate = true
				break
			}
		}
		if activate {
			break
		}
	}

	if activate {
		for {
			onWriteEvent(_chain, actor.Collectors, eventLoggerInfo) //eventLoggerInfo.DB, eventLoggerInfo.Path, actor.Collectors)
			// apply period
			time.Sleep(time.Second * time.Duration(eventLoggerInfo.Period))
		}
	}
}
