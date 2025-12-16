const fs = require('fs');
const path = require('path');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');
const { _getParseInput } = require('./common/currentPath.js');
const { resolveWorkspacePath } = require('../utils/resolvePath.js');
const { _parseFunctionSignatureToAbi } = require('./abi.js');

const { _getCollectionName } = require('./common/keyValue.js');

async function _importAbi(_parsedInput) {
  const collectionName = _getCollectionName(_parsedInput);

  const importFile = resolveWorkspacePath(_parsedInput.data[0]);
  const rawDeployedCollection = fs.readFileSync(importFile, 'utf8');
  let importedContracts = JSON.parse(rawDeployedCollection);
  const contractName = _parsedInput.opt['--name'];

  if (Array.isArray(importedContracts)) {
    if (contractName === undefined) {
      throw new Error('require contract name(--name)');
    }
    importedContracts['contracts'] = {};
    importedContracts.contracts[contractName] = {
      abi: importedContracts
    };
  } else {
    const regParenthesis = /(?<=\{)(.*?)(?=\})/g;
    const selected = _parsedInput.opt['--select'];

    if (contractName === undefined) {
      throw new Error('require contract name(--name)');
    }
    if (selected === undefined) {
      throw new Error('require select abi in json file(--select {parent key}.{child key})');
    }
    const tempContractKey = selected.match(regParenthesis);
    tempContractKey.forEach(element => {
      importedContracts = importedContracts[element];
    });

    importedContracts['contracts'] = {};
    importedContracts.contracts[contractName] = {
      abi: importedContracts
    };
  }

  if (deployedContract[collectionName] === undefined
      || deployedContract[collectionName].contracts === undefined) {
    deployedContract[collectionName] = {
      contracts: {}
    };
  }

  deployedContract[collectionName] = {
    contracts: {
      ...deployedContract[collectionName].contracts,
      ...importedContracts.contracts,
    }
  };

  printSuccess('imported abi file');
}

/**
 * 함수 시그니처를 파싱하여 기존 ABI에 추가
 * @param {object} _parsedInput - 파싱된 입력
 * @param {string} _line - 원본 입력 라인
 */
async function _importAbiAdd(_parsedInput, _line) {
  const collectionName = _getCollectionName(_parsedInput);

  const contractName = _parsedInput.opt['--name'];
  if (contractName === undefined || contractName === '') {
    throw new Error('require contract name(--name)');
  }

  // 함수 시그니처 추출
  let signature = '';

  // _line에서 직접 추출 (따옴표 처리 개선)
  if (_line) {
    // --name 이후 부분 추출
    const nameIndex = _line.indexOf('--name');
    if (nameIndex > -1) {
      let afterName = _line.substring(nameIndex + 6).trim();
      // contractName 다음 부분 추출
      const parts = afterName.split(/\s+/);
      if (parts.length > 1) {
        // contractName 다음부터가 시그니처
        afterName = afterName.substring(parts[0].length).trim();

        // 따옴표로 시작하면 따옴표 쌍 찾기
        if (afterName.startsWith('"')) {
          const endQuote = afterName.indexOf('"', 1);
          if (endQuote > -1) {
            signature = afterName.substring(1, endQuote);
          } else {
            signature = afterName.substring(1);
          }
        } else if (afterName.startsWith("'")) {
          const endQuote = afterName.indexOf("'", 1);
          if (endQuote > -1) {
            signature = afterName.substring(1, endQuote);
          } else {
            signature = afterName.substring(1);
          }
        } else {
          // 따옴표가 없으면 나머지 전체
          signature = afterName;
        }
      }
    }
  }

  // _line에서 추출 실패 시 parsedInput.data 사용
  if (!signature && _parsedInput.data.length > 0) {
    signature = _parsedInput.data.join(' ').trim();
    // 따옴표 제거
    if ((signature.startsWith('"') && signature.endsWith('"')) ||
        (signature.startsWith("'") && signature.endsWith("'"))) {
      signature = signature.slice(1, -1);
    }
  }

  if (!signature || signature.length === 0) {
    throw new Error('Error: Function signature is required for -add option');
  }

  // 함수 시그니처를 ABI로 변환
  let newAbi;
  try {
    newAbi = _parseFunctionSignatureToAbi(signature);
  } catch (e) {
    throw new Error(e.message);
  }

  // collection과 contract 초기화
  if (deployedContract[collectionName] === undefined
      || deployedContract[collectionName].contracts === undefined) {
    deployedContract[collectionName] = {
      contracts: {}
    };
  }

  if (deployedContract[collectionName].contracts[contractName] === undefined) {
    deployedContract[collectionName].contracts[contractName] = {
      abi: []
    };
  }

  // 기존 ABI 배열 가져오기
  const existingAbi = deployedContract[collectionName].contracts[contractName].abi;
  if (!Array.isArray(existingAbi)) {
    throw new Error(`Error: Invalid ABI structure for contract ${contractName}`);
  }

  // 중복 체크: 함수 이름과 입력 파라미터 타입이 같으면 중복
  const isDuplicate = existingAbi.some(existing => {
    if (existing.type !== 'function' || newAbi.type !== 'function') {
      return false;
    }
    if (existing.name !== newAbi.name) {
      return false;
    }
    if (existing.inputs.length !== newAbi.inputs.length) {
      return false;
    }
    // 모든 입력 파라미터 타입이 같은지 확인
    for (let i = 0; i < existing.inputs.length; i++) {
      if (existing.inputs[i].type !== newAbi.inputs[i].type) {
        return false;
      }
    }
    return true;
  });

  if (isDuplicate) {
    throw new Error(`Error: Function "${newAbi.name}" with the same signature already exists in contract "${contractName}"`);
  }

  // ABI 추가
  existingAbi.push(newAbi);

  printSuccess(`Added function "${newAbi.name}" to contract "${contractName}" in collection "${collectionName}"`);
}

async function _commandImport(_inputTokens, _line) {
  const parsedInput = _getParseInput(_inputTokens, 1);

  try {
    if (parsedInput.opt['-abi']) {
      if (parsedInput.opt['-add']) {
        await _importAbiAdd(parsedInput, _line);
      } else {
        await _importAbi(parsedInput);
      }
    } else if (parsedInput.opt['-memory']) {
      const rawStore = fs.readFileSync(parsedInput.data[0], 'utf8');
      varStore = JSON.parse(rawStore);
      printSuccess('imported store file');
    }
  } catch (e) {
    throw e;
  }
}

module.exports._commandImport = _commandImport;