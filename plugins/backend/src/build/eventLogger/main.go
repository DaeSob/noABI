package main

import (
	"cia/api/cmd"
	"cia/api/logHandler"
	"cia/api/preference"
	"cia/api/routers"

	eventActor "cia/pkg/eventLogger/eventActor"
	eventCollector "cia/pkg/eventLogger/eventCollector"
)

func main() {
	// get build mode
	build := cmd.GetBuildFromCommand()

	// log handler 초기화
	logHandler.Initialize(build)
	logHandler.Write("initialize", 0, "initialized log handler")

	logHandler.Write("initialize", 0, "initializing preference")
	err := preference.Initialize(build)
	if err != nil {
		logHandler.Write("initialize", 0, "panic", err.Error())
		return
	}

	logHandler.Write("initialize", 0, "start open P2E logger V3 logger")
	eventCollector.Initialize()
	eventCollector.Run()

	eventActor.Run()

	logHandler.Write("initialize", 0, "start open P2E logger V3 api")
	routers.Run(routers.EventLoggerTable())
}
