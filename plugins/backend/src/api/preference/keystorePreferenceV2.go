// V2.0.0 By XeN
package preference

import "cia/api/logHandler"

func onSetKeystore() {
	inst := GetInstance()

	logHandler.Write("initialize", 0, "loading [keystore]")

	_findToString("keystore.path", &inst.keystorePath)
	logHandler.Write("initialize", 0, "keystore path:", inst.keystorePath)

	logHandler.Write("initialize", 0, "loaded [keystore]")
}
