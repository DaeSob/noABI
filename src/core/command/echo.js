const { printError, printGray, textColor } = require('../log/printScreen.js');
const { _parseKeyValue, _getNestedValue } = require('./common/keyValue.js');

/**
 * echo 명령어의 변수를 매핑하여 결과 문자열 반환
 * @param {Array<string>} _inputTokens - 입력 토큰 배열
 * @param {string} _line - 원본 라인
 * @returns {string} 매핑된 결과 문자열
 */
function _mapEchoVariables(_inputTokens, _line) {
  let echoArg = '';
  
  // echo 명령어가 포함된 경우
  if (_line.indexOf('echo') >= 0) {
    if (_inputTokens.length < 2) {
      return '';
    }
    // 첫 번째 인자부터 모두 합치기 (공백 포함)
    echoArg = _line.substring(_line.indexOf('echo') + 4).trim();
    
    // 따옴표로 감싸진 경우 따옴표 제거
    if ((echoArg.startsWith('"') && echoArg.endsWith('"')) ||
        (echoArg.startsWith("'") && echoArg.endsWith("'"))) {
      echoArg = echoArg.slice(1, -1);
    }
  } else {
    // echo 명령어 없이 입력된 경우 (예: ${pciAddress})
    echoArg = _line.trim();
    
    // 세미콜론 제거
    if (echoArg.endsWith(';')) {
      echoArg = echoArg.slice(0, -1).trim();
    }
  }
  
  // ${변수명} 형식 처리
  const regVarName = /\$\{([^}]+)\}/g;
  let match;
  let lastIndex = 0;
  let output = '';
  
  while ((match = regVarName.exec(echoArg)) !== null) {
    // 변수 앞의 텍스트 추가
    output += echoArg.substring(lastIndex, match.index);
    
    // 변수 값 가져오기
    try {
      // ${변수명} 형식에서 변수명 추출
      const varName = match[1];
      
      // 먼저 예약 변수인지 확인 (예: USER)
      const reservedVars = ['USER'];
      const isReservedVar = reservedVars.some(rv => varName === rv || varName.startsWith(rv + '.'));
      
      let varValue;
      if (isReservedVar) {
        // 예약 변수는 _parseKeyValue 사용 (기존 로직 유지)
        varValue = _parseKeyValue(match[0]);
      } else {
        // store 변수는 _getNestedValue로 subkey 지원
        varValue = _getNestedValue(varName, global.varStore);
        if (varValue === undefined) {
          throw new Error(`Variable "${varName}" is undefined`);
        }
      }
      
      // 변수 값을 문자열로 변환
      const varValueStr = typeof varValue === 'object' ? JSON.stringify(varValue) : String(varValue);
      output += varValueStr;
    } catch (e) {
      // 변수가 없으면 원본 유지
      output += match[0];
    }
    
    lastIndex = match.index + match[0].length;
  }
  
  // 마지막 텍스트 추가
  output += echoArg.substring(lastIndex);
  
  return output;
}

/**
 * echo 명령어 처리 (프롬프트에만 표시하므로 출력하지 않음)
 * @param {Array<string>} _inputTokens - 입력 토큰 배열
 * @param {string} _line - 원본 라인
 */
async function _commandEcho(_inputTokens, _line) {
  // echo 명령어는 프롬프트에만 표시하므로 여기서는 아무것도 하지 않음
  // 실제 매핑은 프롬프트 표시 시점에 수행됨
}

module.exports._commandEcho = _commandEcho;
module.exports._mapEchoVariables = _mapEchoVariables;

