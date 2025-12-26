package main

import (
	"cia/api/cmd"
	"cia/api/logHandler"
	"cia/api/preference"
	"cia/api/routers"
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
	logHandler.Write("initialize", 0, "initialized preference")

	// api server 시작
	logHandler.Write("initialize", 0, "start management api server")
	routers.Run(routers.SignerTable())
}
