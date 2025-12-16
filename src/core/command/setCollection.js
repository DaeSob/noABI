const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandSetCollection(_inputTokens) {
  if (_inputTokens.length > 1) {
    if (currentCollection !== _inputTokens[1]) {
      currentCollection = _inputTokens[1];
      currentPath = _inputTokens[1];
    }
  } else {
    currentCollection = '';
    currentPath = '';
  }
  printSuccess('set collection: ' + currentCollection);
}

module.exports._commandSetCollection = _commandSetCollection;  