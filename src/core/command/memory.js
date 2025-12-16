const fs = require('fs');
const path = require('path');
const { _getCurrentPath, _getParseInput } = require('./common/currentPath.js');
const { _getNestedValue } = require('./common/keyValue.js');
const { textColor, printBlue, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandMemory(_inputTokens) {
  const inputLine = Array.isArray(_inputTokens) ? _inputTokens.join(' ') : _inputTokens;

  const parsedInput = _getParseInput(_inputTokens, 1);
  const optionValue = parsedInput.opt['-clear'];

  if (optionValue !== undefined) {
    varStore = {};
    printSuccess('memory clean');
    return;
  }

  if (parsedInput.data.length === 0) {
    printDefault(varStore);
  } else {
    const filtering = {};
    parsedInput.data.forEach(element => {
      const value = _getNestedValue(element, varStore);
      if (value !== undefined) {
        filtering[element] = value;
      } else {
        filtering[element] = undefined;
      }
    });

    // 출력: object인 경우 JSON.stringify로 전체 내용 표시
    const output = {};
    Object.keys(filtering).forEach(key => {
      const value = filtering[key];
      if (value !== undefined) {
        // object이고 null이 아니면 JSON.stringify로 전체 내용 표시
        if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
          output[key] = value; // 객체 전체를 저장
        } else {
          output[key] = value;
        }
      } else {
        output[key] = undefined;
      }
    });

    // 출력 포맷팅: object는 JSON.stringify로 예쁘게 출력
    if (Object.keys(output).length === 1) {
      // 단일 키인 경우 값만 출력
      const key = Object.keys(output)[0];
      const value = output[key];
      if (value === undefined) {
        console.log(textColor.error + `Error: "${key}" is undefined` + textColor.white);
      } else if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
        console.log(JSON.stringify(value, null, 2));
      } else {
        console.log(value);
      }
    } else {
      // 여러 키인 경우 객체로 출력
      console.log(JSON.stringify(output, null, 2));
    }
  }
}

module.exports._commandMemory = _commandMemory;