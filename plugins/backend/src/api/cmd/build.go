package cmd

import (
	commonCmd "cia/common/cmd"
)

type TBuildMode string

const (
	BUILD_DEV    TBuildMode = "dev"
	BUILD_LIVE              = "live"
	BUILD_LIVE_A            = "live-a"
	BUILD_LIVE_B            = "live-b"
	BUILD_LIVEQA            = "live-qa"
	BUILD_LOCAL             = "local"
	BUILD_QA                = "qa"
	BUILD_WSL               = "wsl"
)

var currentBuildMode string

// get build mode from command arguments
// [dev|live|live-a|live-b|live-qa|local|qa]
// dev :
// live : aws live
// live-a, live-b : aws live 이중화
// live-qa : aws live qa(main chain api)
// local : developer local pc
// qa : aws qa
func GetBuildFromCommand() string {
	if currentBuildMode != "" {
		return currentBuildMode
	}

	arg := commonCmd.GetArg(
		commonCmd.STRING, // arguments type
		"build",          // arguments name
		"local",          // default value
		"build mode[dev|live|live-a|live-b|live-qa|local|qa|wsl]", // desc
	)

	currentBuildMode = *(arg.(*string))
	switch TBuildMode(currentBuildMode) {
	case BUILD_DEV:
	case BUILD_LIVE:
	case BUILD_LIVE_A:
	case BUILD_LIVE_B:
	case BUILD_LIVEQA:
	case BUILD_LOCAL:
	case BUILD_QA:
	case BUILD_WSL:
	default:
		currentBuildMode = BUILD_LOCAL
		return BUILD_LOCAL
	}

	return currentBuildMode
}
