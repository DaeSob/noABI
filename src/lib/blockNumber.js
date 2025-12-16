import { RpcRequest } from "../rpcRequest.js";

async function _blockNumber(rpc, id) {
  try {
    const result = await RpcRequest(rpc, "eth_blockNumber", [], id);
    const blockNumber = parseInt(result.result, 16);

    return blockNumber;
  } catch (e) {
    console.log(e);
  }
}

export const getLatestBlockNumber = _blockNumber;
