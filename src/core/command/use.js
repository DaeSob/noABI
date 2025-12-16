const promptSync = require('prompt-sync');
const prompt = promptSync();
const { OpenWalletFromPrivateKey,
  OpenWalletFromKeystoreV3
} = require('../../lib/import.js');
const { _getParseInput } = require('./common/currentPath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _changeSettings(_newCurrnetSet) {
  if (objSettings[_newCurrnetSet.chain] === undefined
    && objSettings[_newCurrnetSet.wallet] === undefined) {
    throw new Error('not exist chain & wallet key from settings.json');
  }
  currentSet = _newCurrnetSet;
  objSettings.currentSet = _newCurrnetSet;

  rpcUrl = objSettings[currentSet.chain].rpc;
  blockCreationCycle = objSettings[currentSet.chain].blockCreationCycle * 1000;
  gasBufferPercent = objSettings[currentSet.chain].gasBufferPercent;
  customChain = objSettings[currentSet.chain].custom;
  keyStore = objSettings[currentSet.wallet].keyStore;

  if (customChain !== undefined && customChain.gasLimit !== undefined) {
    defaultCallGasLimit = customChain.gasLimit;
  }

  if (gasBufferPercent === undefined) {
    gasBufferPercent = defaultGasBufferPercent;
  }

  if (keyStore !== undefined) {
    const passPhrase = prompt(currentSet.wallet + ' Passphrase: ', { echo: '*' });
    const wallet = await OpenWalletFromKeystoreV3(keyStore, passPhrase);
    privKey = wallet.getPrivateKeyString();
    user = wallet.getAddressString();
  } else {
    privKey = objSettings[currentSet.wallet].privKey;
    const wallet = await OpenWalletFromPrivateKey(privKey);
    user = wallet.getAddressString();
  }
}

async function _commandUse(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  let optionValue = undefined;
  const newSettings = { ...currentSet };

  optionValue = parsedInput.opt['--chainSet'];
  if (optionValue !== undefined) {
    newSettings.chain = optionValue;
  }

  optionValue = parsedInput.opt['--walletSet'];
  if (optionValue !== undefined) {
    newSettings.wallet = optionValue;
  }

  await _changeSettings(newSettings);

  console.log(textColor.gray + 'rpc  ' + textColor.blue + rpcUrl + textColor.white);
  console.log(textColor.gray + 'user ' + textColor.blue + user + textColor.white);
}

module.exports._commandUse = _commandUse;