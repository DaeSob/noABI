const fs = require('fs');
const promptSync = require('prompt-sync');
const prompt = promptSync();
const { OpenWalletFromPrivateKey, OpenWalletFromKeystoreV3 } = require('../../lib/import.js');
const { textColor } = require('../log/printScreen.js');
const { parseArgs } = require('../../lib/utils/args.js');
const {
  getSettingsFileName,
  setSettingsFileName,
  getObjSettings,
  setObjSettings,
  getCurrentSet,
  setCurrentSet,
  getRpcUrl,
  setRpcUrl,
  getBcc,
  setBcc,
  getGasBufferPercent,
  setGasBufferPercent,
  getCustomChain,
  setCustomChain,
  getDefaultCallGasLimit,
  setDefaultCallGasLimit,
  getKeyStore,
  setKeyStore,
  getPrivKey,
  setPrivKey,
  getUser,
  setUser,
  getWorkPath,
  setWorkPath,
  getLogs,
  setLogs
} = require('../state/globals.js');

async function _loadSettings() {
  const cliArgs = parseArgs(process.argv);

  let settingsFileName = getSettingsFileName();
  if (cliArgs.config !== undefined) {
    settingsFileName = cliArgs.config;
    setSettingsFileName(settingsFileName);
  }

  const rawSettingsFile = fs.readFileSync(settingsFileName, 'utf8');
  const objSettings = JSON.parse(rawSettingsFile);
  setObjSettings(objSettings);

  const currentSet = objSettings.currentSet;
  setCurrentSet(currentSet);
  setRpcUrl(objSettings[currentSet.chain].rpc);
  setBcc(objSettings[currentSet.chain].blockCreationCycle * 1000);
  setGasBufferPercent(objSettings[currentSet.chain].gasBufferPercent);
  setCustomChain(objSettings[currentSet.chain].custom);
  setKeyStore(objSettings[currentSet.wallet].keyStore);

  const customChain = getCustomChain();
  if (customChain !== undefined && customChain.gasLimit !== undefined) {
    setDefaultCallGasLimit(customChain.gasLimit);
  }

  let gasBufferPercent = getGasBufferPercent();
  if (gasBufferPercent === undefined) {
    // defaultGasBufferPercent는 global에 설정되어 있음
    setGasBufferPercent(global.defaultGasBufferPercent || 5);
  }

  const keyStore = getKeyStore();
  if (keyStore !== undefined) {
    const passPhrase = prompt(currentSet.wallet + ' Passphrase: ', { echo: '*' });
    const wallet = await OpenWalletFromKeystoreV3(keyStore, passPhrase);
    setPrivKey(wallet.getPrivateKeyString());
    setUser(wallet.getAddressString());
  } else {
    setPrivKey(objSettings[currentSet.wallet].privKey);
    const wallet = await OpenWalletFromPrivateKey(getPrivKey());
    setUser(wallet.getAddressString());
  }

  if (objSettings.workspace.root !== undefined && objSettings.workspace.root !== '') {
    setWorkPath(objSettings.workspace.root);
  } else {
    setWorkPath(process.cwd());
  }

  setLogs(objSettings.workspace.logs);
}

async function _setConfiguration() {
  const { _welcomeHead } = require('../../welcome.js');
  try {
    _welcomeHead();
    await _loadSettings();
  } catch (e) {
    console.log(textColor.error + e.message + textColor.white);
    process.exit();
  }
}

module.exports = {
  _loadSettings,
  _setConfiguration
};

