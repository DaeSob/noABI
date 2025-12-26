package preference

import (
	"cia/api/logHandler"
	"strings"
)

func onSetCommonPreference() {
	inst := GetInstance()

	// lock
	inst.lockMutex.Lock()
	defer inst.lockMutex.Unlock()

	logHandler.Write("initialize", 0, "loading [common]")

	// Host
	_findToString("server.host", &inst.tServer.Host)

	// Debug mode
	_findToBool("server.debug", &inst.tServer.Debug)

	// Supervisor
	_findToString("server.supervisor", &inst.tServer.Supervisor)
	if len(inst.tServer.Supervisor) > 0 {
		inst.tServer.Supervisor = strings.ToLower(inst.tServer.Supervisor)
	}

	// path prefix
	_findToString("server.prefix", &inst.tServer.PathPrefix)

	// authorization
	_findToBool("authorization.enable", &inst.tAuthorization.Enable)
	_findToString("authorization.signer", &inst.tAuthorization.Signer)
	if len(inst.tAuthorization.Signer) > 0 {
		inst.tAuthorization.Signer = strings.ToLower(inst.tAuthorization.Signer)
	}

	logHandler.Write("initialize", 0, "loaded [common]")

}
