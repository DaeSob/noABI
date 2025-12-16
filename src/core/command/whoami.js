const { OpenWalletFromPrivateKey } = require('../../lib/import.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandWhoAmI(_inputTokens) {
  if (privKey === undefined || privKey.length !== privKey.length) {
    throw new String('no account');
  }
  const wallet = await OpenWalletFromPrivateKey(privKey);
  printSuccess(wallet.getAddressString());
}

module.exports._commandWhoAmI = _commandWhoAmI;  