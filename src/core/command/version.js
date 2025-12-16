const { _getParseInput } = require('./common/currentPath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandVersion(_inputTokens) {
  printSuccess('0.6.0.0');
}

module.exports._commandVersion = _commandVersion;