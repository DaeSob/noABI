const path = require('path');
const fs = require('fs');
const { OpenWalletFromPrivateKey } = require('../../lib/import.js');
const { _getParseInput } = require('./common/currentPath.js');
const { _commandSetCollection } = require('./setCollection.js');
const { resolveWorkspacePath } = require('../utils/resolvePath.js');
const { _deploy } = require('../deploy/deploy.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandDeploy(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  const composeFile = resolveWorkspacePath(parsedInput.data[0]);
  const optionalDeploy = parsedInput.data.slice(1);

  const rawInputCompose = fs.readFileSync(composeFile, 'utf8');
  const objComposeFile = JSON.parse(rawInputCompose);

  if (objComposeFile.collection === undefined) {
    throw new String('not found a collection in compose-file');
  }

  _commandSetCollection(['use', objComposeFile.collection]);

  objComposeFile['chain'] = {
    rpc: rpcUrl
  };
  objComposeFile['wallet'] = {
    privKey: privKey
  };

  const wallet = await OpenWalletFromPrivateKey(privKey);
  console.log(textColor.warning + 'Compose file:', composeFile);
  console.log(textColor.warning + 'URL RPC:', rpcUrl);
  console.log(textColor.warning + 'Account:', wallet.getAddressString());

  await _deploy(objComposeFile, optionalDeploy);
}

module.exports._commandDeploy = _commandDeploy;