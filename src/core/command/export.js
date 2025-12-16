const fs = require('fs');
const path = require('path');
const { _getParseInput } = require('./common/currentPath.js');
const { _getCollectionName } = require('./common/keyValue.js');
const { resolveWorkspacePath } = require('../utils/resolvePath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _exportAbi(_parsedInput) {
  let collectionName = _parsedInput.opt['--collection'];
  if (collectionName === undefined) {
    collectionName = currentCollection;
    if (collectionName === undefined) {
      throw new Error('require collection name(--collection)');
    }
  }

  const contractName = _parsedInput.opt['--name'];
  if (contractName === undefined) {
    throw new Error('require contract name(--name)');
  }

  const exportFile = resolveWorkspacePath(_parsedInput.data[0]);
  const mkPath = path.dirname(exportFile);
  fs.mkdirSync(mkPath, { recursive: true });

  try {
    fs.writeFileSync(
      exportFile,
      JSON.stringify(deployedContract[collectionName]['contracts'][contractName]['abi'], null, 2)
    );
    printSuccess('abi exported to ' + _parsedInput.data[0]);
  } catch (error) {
    printError(error);
  }
}

async function _commandExport(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  try {
    if (parsedInput.opt['-abi']) {
      await _exportAbi(parsedInput);
    } else if (parsedInput.opt['-memory']) {
      const exportFile = resolveWorkspacePath(parsedInput.data[0]);
      const mkPath = path.dirname(exportFile);
      fs.mkdirSync(mkPath, { recursive: true });
      try {
        fs.writeFileSync(
          exportFile,
          JSON.stringify(varStore, null, 2)
        );
        printSuccess('store exported to ' + parsedInput.data[0]);
      } catch (error) {
        printError(error);
      }
    }
  } catch (e) {
    throw e;
  }
}

module.exports._commandExport = _commandExport;