const { _getParseInput }                                               = require("./common/currentPath.js");
const { EncodeFunctionSignature }                                      = require("../../lib/utils/encoder.js");
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

/**
 * 함수 시그니처를 파싱하여 ABI 객체 생성
 * @param {string} signature - 함수 시그니처 (예: "function balanceOf( address owner ) public view returns( uint256 amount )")
 * @returns {object} ABI 객체
 */
function _parseFunctionSignatureToAbi(signature) {
  // 공백 정규화
  const normalized = signature.trim().replace(/\s+/g, ' ');
  
  // function 키워드 확인
  if (!normalized.toLowerCase().startsWith('function ')) {
    throw new Error('Error: Function signature must start with "function"');
  }
  
  // 괄호 균형 검증
  let parenCount = 0;
  let inString = false;
  let stringChar = '';
  
  for (let i = 0; i < normalized.length; i++) {
    const char = normalized[i];
    
    // 문자열 내부 체크
    if ((char === '"' || char === "'") && (i === 0 || normalized[i-1] !== '\\')) {
      if (!inString) {
        inString = true;
        stringChar = char;
      } else if (char === stringChar) {
        inString = false;
        stringChar = '';
      }
      continue;
    }
    
    if (!inString) {
      if (char === '(') {
        parenCount++;
      } else if (char === ')') {
        parenCount--;
        if (parenCount < 0) {
          throw new Error('Error: Unmatched closing parenthesis in function signature');
        }
      }
    }
  }
  
  if (parenCount !== 0) {
    throw new Error('Error: Unmatched parentheses in function signature');
  }
  
  // 함수 이름 추출: function name(
  const nameMatch = normalized.match(/function\s+(\w+)\s*\(/i);
  if (!nameMatch) {
    throw new Error('Error: Could not parse function name - missing opening parenthesis after function name');
  }
  const functionName = nameMatch[1];
  
  // 입력 파라미터 추출: ( ... )
  // 첫 번째 괄호 쌍 찾기 (함수 이름 뒤)
  let inputStart = normalized.indexOf('(', nameMatch.index + nameMatch[0].length - 1);
  if (inputStart === -1) {
    throw new Error('Error: Missing opening parenthesis for function parameters');
  }
  
  let inputEnd = -1;
  parenCount = 1;
  for (let i = inputStart + 1; i < normalized.length; i++) {
    const char = normalized[i];
    if (char === '(') {
      parenCount++;
    } else if (char === ')') {
      parenCount--;
      if (parenCount === 0) {
        inputEnd = i;
        break;
      }
    }
  }
  
  if (inputEnd === -1) {
    throw new Error('Error: Missing closing parenthesis for function parameters');
  }
  
  const inputParamsStr = normalized.substring(inputStart + 1, inputEnd).trim();
  
  // 입력 파라미터 파싱
  const inputs = [];
  if (inputParamsStr.length > 0) {
    // 쉼표로 분리 (하지만 튜플 내부의 쉼표는 제외해야 함 - 간단한 버전으로 구현)
    const params = inputParamsStr.split(',').map(p => p.trim()).filter(p => p.length > 0);
    params.forEach(param => {
      // "type name" 형식 파싱
      const parts = param.trim().split(/\s+/);
      if (parts.length >= 2) {
        // 마지막이 이름, 나머지가 타입
        const name = parts[parts.length - 1];
        const type = parts.slice(0, -1).join(' ');
        inputs.push({ name, type });
      } else if (parts.length === 1) {
        // 이름이 없는 경우
        inputs.push({ name: '', type: parts[0] });
      }
    });
  }
  
  // stateMutability 추출 (public, view, pure, payable, nonpayable)
  let stateMutability = 'nonpayable';
  const lowerSignature = normalized.toLowerCase();
  if (lowerSignature.includes(' view ')) {
    stateMutability = 'view';
  } else if (lowerSignature.includes(' pure ')) {
    stateMutability = 'pure';
  } else if (lowerSignature.includes(' payable ')) {
    stateMutability = 'payable';
  }
  
  // 출력 파라미터 추출: returns( ... )
  const outputs = [];
  const returnsMatch = normalized.match(/returns\s*\(/i);
  if (returnsMatch) {
    const returnsStart = normalized.indexOf('(', returnsMatch.index);
    if (returnsStart === -1) {
      throw new Error('Error: Missing opening parenthesis for returns clause');
    }
    
    let returnsEnd = -1;
    let returnsParenCount = 1;
    for (let i = returnsStart + 1; i < normalized.length; i++) {
      const char = normalized[i];
      if (char === '(') {
        returnsParenCount++;
      } else if (char === ')') {
        returnsParenCount--;
        if (returnsParenCount === 0) {
          returnsEnd = i;
          break;
        }
      }
    }
    
    if (returnsEnd === -1) {
      throw new Error('Error: Missing closing parenthesis for returns clause');
    }
    
    const returnsStr = normalized.substring(returnsStart + 1, returnsEnd).trim();
    if (returnsStr.length > 0) {
      const returnParams = returnsStr.split(',').map(p => p.trim()).filter(p => p.length > 0);
      returnParams.forEach(param => {
        const parts = param.trim().split(/\s+/);
        if (parts.length >= 2) {
          const name = parts[parts.length - 1];
          const type = parts.slice(0, -1).join(' ');
          outputs.push({ name, type });
        } else if (parts.length === 1) {
          outputs.push({ name: '', type: parts[0] });
        }
      });
    }
  }
  
  // ABI 객체 생성
  const abi = {
    type: 'function',
    name: functionName,
    inputs: inputs,
    outputs: outputs,
    stateMutability: stateMutability
  };
  
  return abi;
}

async function _commandAbi(_inputTokens, _line) {

  let parsedInput = _getParseInput( _inputTokens, 1 );
  let optionValue = undefined;

  // -create 옵션 처리
  if (parsedInput.opt['-create'] !== undefined) {
    // _line에서 직접 시그니처 추출 (따옴표 처리 개선)
    let signature = '';
    if (_line) {
      // -create 이후 부분 추출
      const createIndex = _line.indexOf('-create');
      if (createIndex > -1) {
        let afterCreate = _line.substring(createIndex + 7).trim();
        // 따옴표로 시작하면 따옴표 쌍 찾기
        if (afterCreate.startsWith('"')) {
          const endQuote = afterCreate.indexOf('"', 1);
          if (endQuote > -1) {
            signature = afterCreate.substring(1, endQuote);
          } else {
            signature = afterCreate.substring(1);
          }
        } else if (afterCreate.startsWith("'")) {
          const endQuote = afterCreate.indexOf("'", 1);
          if (endQuote > -1) {
            signature = afterCreate.substring(1, endQuote);
          } else {
            signature = afterCreate.substring(1);
          }
        } else {
          // 따옴표가 없으면 공백으로 분리된 첫 번째 토큰
          signature = afterCreate.split(/\s+/)[0];
        }
      }
    }
    
    // _line에서 추출 실패 시 parsedInput.data 사용
    if (!signature && parsedInput.data.length > 0) {
      signature = parsedInput.data.join(' ').trim();
      // 따옴표 제거
      if ((signature.startsWith('"') && signature.endsWith('"')) ||
          (signature.startsWith("'") && signature.endsWith("'"))) {
        signature = signature.slice(1, -1);
      }
    }
    
    if (!signature || signature.length === 0) {
      throw new Error('Error: Function signature is required for -create option');
    }
    
    try {
      const abi = _parseFunctionSignatureToAbi(signature);
      console.log(JSON.stringify(abi, null, 2));
    } catch (e) {
      throw new Error(e.message);
    }
    return;
  }

  // 기존 -funcSig 옵션 처리
  optionValue = parsedInput.opt['-funcSig'];  
  if( optionValue ) {
    let encodedData     = EncodeFunctionSignature(parsedInput.data[0]);
    printSuccess( encodedData );
  }

}
module.exports._commandAbi = _commandAbi;
module.exports._parseFunctionSignatureToAbi = _parseFunctionSignatureToAbi;