const { _getParseInput } = require('./common/currentPath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');
const { RpcRequest } = require('../../lib/rpc/rpcRequest.js');

async function _commandEstimateGas(_inputTokens, _line) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  let optionValue = undefined;

  let from = parsedInput.opt['--from'];
  const to = parsedInput.opt['--to'];
  let gasPrice = parsedInput.opt['--gasPrice'];
  const callData = parsedInput.opt['--callData'];

  if (from === undefined) {
    from = user;
  }

  if (gasPrice === undefined) {
    gasPrice = '0x0';
  }

  const resGas = await RpcRequest(rpcUrl, 'eth_estimateGas', [{
    from: from,
    to: to,
    data: callData,
  }], 1);

  console.log(resGas);
}

module.exports._commandEstimateGas = _commandEstimateGas;