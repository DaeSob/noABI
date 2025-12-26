package preference

//--------------------------------------------------------------
//V2.0.0 By XeN

func GetKeystorePath() string {
	inst := GetInstance()

	return inst.keystorePath
}
