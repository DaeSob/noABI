const fs = require('fs');
const { _getParseInput } = require('./common/currentPath.js');
const { _parseKeyValue, _extractStoreVarName, _saveToVarStore } = require('./common/keyValue.js');
const { RpcRequest } = require('../../lib/rpc/rpcRequest.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandGetBalance(_inputTokens, _line) {
  // => 연산자로 변수명 추출
  const { storeVarName, commandLine } = _extractStoreVarName(_line);
  let addressLine = commandLine;

  // _line에서 직접 주소 추출 (parseLine이 $ 문자를 제거하는 문제 방지)
  let address = '';
  if (addressLine) {
    const match = addressLine.match(/getBalance\s+(.+)/);
    if (match && match[1]) {
      address = match[1].trim();
    }
  }
  
  if (!address || address === '') {
    throw new Error("Error: address is required");
  }
  
  // 따옴표 제거 (따옴표로 감싸진 경우)
  const trimmedAddress = address.trim();
  if ((trimmedAddress.startsWith('"') && trimmedAddress.endsWith('"')) ||
      (trimmedAddress.startsWith("'") && trimmedAddress.endsWith("'"))) {
    address = trimmedAddress.slice(1, -1);
  } else {
    address = trimmedAddress;
  }
  
  // 빈 문자열 체크
  if (address === '' || address === '""' || address === "''") {
    throw new Error("Error: address cannot be empty");
  }
  
  try {
    // _parseKeyValue는 예약 변수(USER 등)와 store 변수를 모두 처리
    address = _parseKeyValue(address);
    
    // 매핑 후에도 빈 문자열인지 체크
    if (address === '' || address === null || address === undefined) {
      throw new Error("Error: address cannot be empty");
    }
    
    // 값이 객체인 경우 문자열로 변환
    if (typeof address === 'object') {
      address = JSON.stringify(address);
    } else {
      address = String(address);
    }
    
    // 최종 변환 후 빈 문자열 체크
    const finalAddress = address.trim();
    if (finalAddress === '' || finalAddress === '""' || finalAddress === "''") {
      throw new Error("Error: address cannot be empty");
    }
    
    address = finalAddress;
  } catch (e) {
    // 변수가 없으면 에러 메시지
    throw new Error(e.toString());
  }
  
  const res = await RpcRequest(rpcUrl, 'eth_getBalance', [address, 'latest'], 1);
  if (res.error !== undefined) {
    throw Error(res.error.message);
  }
  // uint256 값을 문자열로 표시 (BigInt 사용)
  const balance = BigInt(res.result);
  const balanceString = balance.toString();

  // => 연산자가 있으면 store에 저장
  _saveToVarStore(storeVarName, balanceString);

  // 출력은 현재와 동일하게 유지
  console.log(textColor.success + balanceString + textColor.white);
}

module.exports._commandGetBalance = _commandGetBalance;