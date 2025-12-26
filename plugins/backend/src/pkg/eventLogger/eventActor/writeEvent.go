package eventActor

import (
	"cia/api/preference"
	"fmt"
)

func onWriteEvent(_chain string, _collectors []string, _eventLoggerInfo *preference.TEvetLoggerInfo) {
	for _, name := range _collectors {
		logPath := fmt.Sprintf("%v/%v/%v", _eventLoggerInfo.Path, _chain, name)
		readEventLogFile(_chain, logPath, _eventLoggerInfo)
	}
}
