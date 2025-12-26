package routers

import (
	apiSigner "cia/api/routers/v1/signer"
	apiEventLogger "cia/api/routers/v1/eventLogger"
)

func SignerTable() TRouterTable {
	rt := TRouterTable{}

	//Data Sign
	rt.AddPostRouter("/v1/sign/eip191", apiSigner.EIP191DataSign, false)
	rt.AddPostRouter("/v1/sign/eip712", apiSigner.EIP712DataSign, false)
	rt.AddPostRouter("/v1/recover/signer", apiSigner.RecoverDataSigner, false)

	// ping
	rt.AddRouter(ROUTER_PING)

	return rt
}

func EventLoggerTable() TRouterTable {
	rt := TRouterTable{}

	// Event Logger APIs
	rt.AddGetRouter("/v1/event/chains", apiEventLogger.GetChains, false)
	rt.AddGetRouter("/v1/event/chains/:alias", apiEventLogger.GetChain, false)
	rt.AddGetRouter("/v1/event/remove/:alias", apiEventLogger.RemoveChain, false)
	rt.AddPostRouter("/v1/event/chains/:alias", apiEventLogger.UpdateOrAddChain, false)

	// ping
	rt.AddRouter(ROUTER_PING)

	return rt
}
