const fs = require('fs');
const { _getParseInput } = require('./common/currentPath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');
const { _extractStoreVarName, _saveToVarStore } = require('./common/keyValue.js');

async function _commandToUTC(_inputTokens, _line) {
  // => 연산자로 변수명 추출
  const { storeVarName, commandLine } = _extractStoreVarName(_line);
  let inputTokensForParsing = _inputTokens;

  if (commandLine && commandLine !== _line) {
    // => 이전 부분만 추출하여 토큰 재생성
    const { parseLine } = require('parse-line');
    inputTokensForParsing = parseLine(commandLine);
  }

  const parsedInput = _getParseInput(inputTokensForParsing, 1);
  let optionValue = undefined;

  let dateTime = new Date();

  if (parsedInput.data.length > 0) {
    if (isNaN(parsedInput.data[0])) {
      dateTime = new Date(parsedInput.data[0]);
    } else {
      dateTime = new Date(Number(parsedInput.data[0]) * 1000);
    }
  }

  let timestamp = Math.floor(dateTime.getTime() / 1000);
  let unixTime = timestamp;

  if (parsedInput.opt['-d']) {
    unixTime = (Math.floor(timestamp / 86400) * 86400);
  } else if (parsedInput.opt['-h']) {
    unixTime = (Math.floor(timestamp / 3600) * 3600);
  } else if (parsedInput.opt['-m']) {
    unixTime = (Math.floor(timestamp / 60) * 60);
  }

  optionValue = parsedInput.opt['--add'];
  if (optionValue !== undefined) {
    unixTime = unixTime + Number(optionValue);
  }

  let result = '';
  if (parsedInput.opt['-toString']) {
    dateTime = new Date(unixTime * 1000);
    result = dateTime.toISOString();
    console.log(textColor.success + result + textColor.white);
  } else {
    result = unixTime.toString();
    console.log(textColor.success + result + textColor.white);
  }

  // => 연산자가 있으면 store에 저장
  _saveToVarStore(storeVarName, result);
}

module.exports._commandToUTC = _commandToUTC;