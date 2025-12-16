const { _getParseInput } = require('./common/currentPath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandSet(_inputTokens, _line) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  let optionValue = undefined;

  optionValue = parsedInput.opt['--var'];
  if (optionValue !== undefined) {
    const inputJson = _line.substr(_line.indexOf('--var') + ('--var').length);
    const varInputs = JSON.parse(inputJson);
    varStore = {
      ...varStore,
      ...varInputs,
    };
    printSuccess('set varialbe success');
  }

  optionValue = parsedInput.opt['--gasLimit'];
  if (optionValue !== undefined) {
    defaultCallGasLimit = optionValue;
  }

  optionValue = parsedInput.opt['--output'];
  if (optionValue !== undefined) {
    if (optionValue === 'minimal' || optionValue === 'full') {
      outputFormat = optionValue;
      printSuccess('set output format: ' + optionValue);
    } else {
      printError('Error: invalid output format. Use \'minimal\' or \'full\'');
    }
  }
}

module.exports._commandSet = _commandSet;